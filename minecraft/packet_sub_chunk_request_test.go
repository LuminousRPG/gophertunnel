package minecraft

import (
	"bytes"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func TestDecodeSubChunkRequestLayouts(t *testing.T) {
	position := protocol.SubChunkPos{573, 3, -12}
	offsets := []protocol.SubChunkOffset{{0, 3, 0}, {1, -4, 0}, {-1, 2, 3}}

	tests := []struct {
		name   string
		encode func(protocol.IO)
	}{
		{
			name: "current",
			encode: func(w protocol.IO) {
				dimension := int32(0)
				w.Varint32(&dimension)
				protocol.Slice(w, &offsets)
				w.Int32(&position[0])
				w.Int32(&position[1])
				w.Int32(&position[2])
			},
		},
		{
			name: "legacy",
			encode: func(w protocol.IO) {
				dimension := int32(0)
				w.Varint32(&dimension)
				w.SubChunkPos(&position)
				protocol.SliceUint32Length(w, &offsets)
			},
		},
	}

	conn := &Conn{proto: proto{}, readerLimits: true}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			tt.encode(protocol.NewWriter(buf, 0))
			request, err := decodeSubChunkRequest(buf.Bytes(), conn)
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if request.Position != position {
				t.Fatalf("position = %v, want %v", request.Position, position)
			}
			if len(request.Offsets) != len(offsets) {
				t.Fatalf("offset count = %d, want %d", len(request.Offsets), len(offsets))
			}
			for i := range offsets {
				if request.Offsets[i] != offsets[i] {
					t.Errorf("offset[%d] = %v, want %v", i, request.Offsets[i], offsets[i])
				}
			}
		})
	}
}
