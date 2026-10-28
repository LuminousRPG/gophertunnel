package protocol

import (
	"bytes"
	"reflect"
	"testing"
)

func TestSubChunkEntryHeightMaps(t *testing.T) {
	// Each height map is 16 rows of a one-byte row length (always 16) followed by 16 int8 heights.
	const heightMapLen = 16 * (1 + 16)
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
				var data HeightMap
				for z := range data {
					for x := range data[z] {
						data[z][x] = int8((z*16+x)%18 - 1)
					}
				}
				entry.HeightMapType = HeightMapDataHasData
				entry.HeightMapData = Option(data)
			}
			if tc.render {
				var data HeightMap
				for z := range data {
					for x := range data[z] {
						data[z][x] = int8(16 - (z*16+x)%18)
					}
				}
				entry.RenderHeightMapType = HeightMapDataHasData
				entry.RenderHeightMapData = Option(data)
			}
			var buf bytes.Buffer
			entry.Marshal(NewWriter(&buf, 0))
			// Offset, result, four optional-presence flags, two map types and blob hash.
			wantLen := 18
			if tc.height {
				wantLen += heightMapLen
			}
			if tc.render {
				wantLen += heightMapLen
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
