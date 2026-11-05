package packet

import (
	"bytes"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func TestClientCacheBlobStatusManyHashes(t *testing.T) {
	in := &ClientCacheBlobStatus{MissHashes: []uint64{7}, HitHashes: make([]uint64, 1377)}
	for i := range in.HitHashes {
		in.HitHashes[i] = uint64(i)
	}
	buf := new(bytes.Buffer)
	in.Marshal(protocol.NewWriter(buf, 0))
	out := &ClientCacheBlobStatus{}
	r := protocol.NewReader(buf, 0, true)
	out.Marshal(r)
	if len(out.MissHashes) != 1 || out.MissHashes[0] != 7 || len(out.HitHashes) != 1377 || out.HitHashes[1376] != 1376 {
		t.Fatalf("hashes not read back: miss=%v hits=%d", out.MissHashes, len(out.HitHashes))
	}
}
