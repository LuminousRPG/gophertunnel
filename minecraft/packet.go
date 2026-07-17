package minecraft

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// packetData holds the data of a Minecraft packet.
type packetData struct {
	h       *packet.Header
	full    []byte
	payload *bytes.Buffer
}

// parseData parses the packet data slice passed into a packetData struct.
func parseData(data []byte, conn *Conn) (*packetData, error) {
	buf := bytes.NewBuffer(data)
	header := &packet.Header{}
	if err := header.Read(buf); err != nil {
		// We don't return this as an error as it's not in the hand of the user to control this. Instead,
		// we return to reading a new packet.
		return nil, fmt.Errorf("read packet header: %w", err)
	}
	if conn.packetFunc != nil {
		// The packet func was set, so we call it.
		conn.packetFunc(*header, buf.Bytes(), conn.RemoteAddr(), conn.LocalAddr())
	}
	return &packetData{h: header, full: data, payload: buf}, nil
}

type unknownPacketError struct {
	id uint32
}

func (err unknownPacketError) Error() string {
	return fmt.Sprintf("unexpected packet (ID=%v)", err.id)
}

// decode decodes the packet payload held in the packetData and returns the packet.Packet decoded.
func (p *packetData) decode(conn *Conn) (pks []packet.Packet, err error) {
	// Attempt to fetch the packet with the right packet ID from the pool.
	pkFunc, ok := conn.pool[p.h.PacketID]
	var pk packet.Packet
	if !ok {
		// No packet with the ID. This may be a custom packet of some sorts.
		pk = &packet.Unknown{PacketID: p.h.PacketID}
		if conn.disconnectOnUnknownPacket {
			_ = conn.Close()
			return nil, unknownPacketError{id: p.h.PacketID}
		}
	} else {
		pk = pkFunc()
	}

	defer func() {
		if recoveredErr := recover(); recoveredErr != nil {
			err = fmt.Errorf("decode packet %T: %w", pk, recoveredErr.(error))
		}
		if err != nil && !errors.Is(err, unknownPacketError{}) && conn.disconnectOnInvalidPacket {
			_ = conn.Close()
		}
	}()

	if request, ok := pk.(*packet.SubChunkRequest); ok {
		decoded, decodeErr := decodeSubChunkRequest(p.payload.Bytes(), conn)
		if decodeErr != nil {
			err = fmt.Errorf("decode packet %T: %w", pk, decodeErr)
		} else {
			*request = *decoded
			p.payload.Next(p.payload.Len())
		}
	} else {
		r := conn.proto.NewReader(p.payload, conn.shieldID.Load(), conn.readerLimits)
		pk.Marshal(r)
	}
	if err == nil && p.payload.Len() != 0 {
		err = fmt.Errorf("decode packet %T: %v unread bytes left: 0x%x", pk, p.payload.Len(), p.payload.Bytes())
	}
	if conn.disconnectOnInvalidPacket && err != nil {
		return nil, err
	}
	return conn.proto.ConvertToLatest(pk, conn), err
}

// decodeSubChunkRequest accepts both layouts used by protocol 1001 clients.
// The 1.26.30 layout places the offset list before a fixed-width position,
// while some patch clients still send the former varint position first.
func decodeSubChunkRequest(data []byte, conn *Conn) (*packet.SubChunkRequest, error) {
	request, currentErr := tryDecodeSubChunkRequest(data, conn, false)
	if currentErr == nil {
		return request, nil
	}
	request, legacyErr := tryDecodeSubChunkRequest(data, conn, true)
	if legacyErr == nil {
		return request, nil
	}
	return nil, fmt.Errorf("unsupported layout (current: %v; legacy: %v), payload=0x%x", currentErr, legacyErr, data)
}

func tryDecodeSubChunkRequest(data []byte, conn *Conn, legacy bool) (request *packet.SubChunkRequest, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("%v", recovered)
		}
	}()

	buf := bytes.NewBuffer(data)
	io := conn.proto.NewReader(buf, conn.shieldID.Load(), conn.readerLimits)
	request = &packet.SubChunkRequest{}
	io.Varint32(&request.Dimension)

	var count uint32
	if legacy {
		io.SubChunkPos(&request.Position)
		io.Uint32(&count)
	} else {
		io.Varuint32(&count)
	}

	// SubChunkRequest may legitimately contain more offsets than the generic
	// slice limit (1024). Validate the count against the exact payload size
	// before allocating so large, valid requests remain safe to decode.
	trailingBytes := uint64(0)
	if !legacy {
		trailingBytes = 12 // Three fixed-width int32 position components.
	}
	expectedBytes := uint64(count)*3 + trailingBytes
	if expectedBytes != uint64(buf.Len()) {
		return nil, fmt.Errorf("offset count %d requires %d bytes, have %d", count, expectedBytes, buf.Len())
	}

	request.Offsets = make([]protocol.SubChunkOffset, count)
	for i := range request.Offsets {
		request.Offsets[i].Marshal(io)
	}
	if !legacy {
		io.Int32(&request.Position[0])
		io.Int32(&request.Position[1])
		io.Int32(&request.Position[2])
	}
	if buf.Len() != 0 {
		return nil, fmt.Errorf("%d unread bytes", buf.Len())
	}
	return request, nil
}
