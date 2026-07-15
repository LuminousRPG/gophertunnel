package packet

import (
	"bytes"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// clientRequestBytes builds the current protocol 1001 SubChunkRequest layout.
func clientRequestBytes(dimension int32, pos protocol.SubChunkPos, offsets []protocol.SubChunkOffset) []byte {
	buf := bytes.NewBuffer(nil)
	w := protocol.NewWriter(buf, 0)

	w.Varint32(&dimension)
	protocol.Slice(w, &offsets)
	w.Int32(&pos[0])
	w.Int32(&pos[1])
	w.Int32(&pos[2])

	return buf.Bytes()
}

// Reading the fields in any order other than the client's makes coordinate
// bytes land where the offsets' length is expected, producing a huge slice
// length and dropping the connection. Limits are enabled here so the test fails
// in the same way as the live connection.
func TestSubChunkRequestDecodesClientLayout(t *testing.T) {
	const chunkX, subChunkY, chunkZ = 573, 4, -12

	pos := protocol.SubChunkPos{chunkX, subChunkY, chunkZ}
	offsets := []protocol.SubChunkOffset{{0, 0, 0}, {1, -1, 2}, {-3, 4, -5}}
	data := clientRequestBytes(0, pos, offsets)

	pk := &SubChunkRequest{}
	r := protocol.NewReader(bytes.NewReader(data), 0, true)

	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("decoding a client request failed: %v", err)
		}
	}()

	pk.Marshal(r)

	if pk.Dimension != 0 {
		t.Errorf("Dimension = %d, want 0", pk.Dimension)
	}
	if pk.Position != pos {
		t.Errorf("Position = %v, want %v", pk.Position, pos)
	}
	if len(pk.Offsets) != len(offsets) {
		t.Fatalf("read %d offsets, want %d", len(pk.Offsets), len(offsets))
	}
	for i, off := range offsets {
		if pk.Offsets[i] != off {
			t.Errorf("Offsets[%d] = %v, want %v", i, pk.Offsets[i], off)
		}
	}
}

// Chunk coordinates far enough from origin are what pushed the misread length
// over the slice limit, so decoding must hold up across the world's range
// rather than only near spawn.
func TestSubChunkRequestDecodesDistantChunks(t *testing.T) {
	for _, chunkX := range []int32{0, 511, 573, 6282, -6282} {
		pos := protocol.SubChunkPos{chunkX, 4, chunkX}
		data := clientRequestBytes(0, pos, []protocol.SubChunkOffset{{1, 2, 3}})

		pk := &SubChunkRequest{}
		func() {
			defer func() {
				if err := recover(); err != nil {
					t.Errorf("chunk X %d failed to decode: %v", chunkX, err)
				}
			}()
			pk.Marshal(protocol.NewReader(bytes.NewReader(data), 0, true))
		}()

		if pk.Position != pos {
			t.Errorf("chunk X %d: Position = %v, want %v", chunkX, pk.Position, pos)
		}
	}
}
