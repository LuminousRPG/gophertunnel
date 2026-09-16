package protocol

import (
	"bytes"
	"reflect"
	"testing"
)

func TestSubChunkEntryHeightMaps(t *testing.T) {
	for _, tc := range []struct {
		name           string
		height, render bool
	}{
		{"absent", false, false},
		{"height", true, false},
		{"render", false, true},
		{"both", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			entry := SubChunkEntry{
				Offset:   SubChunkOffset{-1, 2, 3},
				Result:   SubChunkResultSuccess,
				BlobHash: Option(uint64(0x123456789abcdef0)),
			}
			if tc.height {
				data := make([]int8, SubChunkHeightMapLen)
				for i := range data {
					data[i] = int8(i%18 - 1)
				}
				for z := 0; z < 16; z++ {
					data[z*17] = 16
				}
				entry.HeightMapType = HeightMapDataHasData
				entry.HeightMapData = Option(data)
			}
			if tc.render {
				data := make([]int8, SubChunkHeightMapLen)
				for i := range data {
					data[i] = int8(16 - i%18)
				}
				for z := 0; z < 16; z++ {
					data[z*17] = 16
				}
				entry.RenderHeightMapType = HeightMapDataHasData
				entry.RenderHeightMapData = Option(data)
			}
			var buf bytes.Buffer
			entry.Marshal(NewWriter(&buf, 0))
			// Offset, result, four optional-presence flags, two map types and blob hash.
			wantLen := 18
			// Protocol 2193 includes a one-byte length before each 16-value row.
			if tc.height {
				wantLen += SubChunkHeightMapLen
			}
			if tc.render {
				wantLen += SubChunkHeightMapLen
			}
			if buf.Len() != wantLen {
				t.Fatalf("encoded length = %d, want %d", buf.Len(), wantLen)
			}
			var decoded SubChunkEntry
			decoded.Marshal(NewReader(&buf, 0, true))
			if !reflect.DeepEqual(entry, decoded) {
				t.Fatalf("round trip differs: got %#v, want %#v", decoded, entry)
			}
			if buf.Len() != 0 {
				t.Fatalf("%d unread bytes", buf.Len())
			}
		})
	}
}
