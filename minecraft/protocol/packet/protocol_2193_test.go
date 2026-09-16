package packet

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// Fixtures follow Mojang's v1.26.50 schemas rather than encoding an expected
// packet with the same codec under test.
func TestProtocol2193Lists(t *testing.T) {
	for _, fixture := range []struct {
		pk   Packet
		want []byte
	}{
		{&PlayerList{}, []byte{0}},
		{&SetScore{}, []byte{0}},
	} {
		pk := fixture.pk
		var buf bytes.Buffer
		pk.Marshal(protocol.NewWriter(&buf, 0))
		if !bytes.Equal(buf.Bytes(), fixture.want) {
			t.Fatalf("%T empty list = %x, want %x", pk, buf.Bytes(), fixture.want)
		}
		pk.Marshal(protocol.NewReader(bytes.NewBuffer([]byte{0}), 0, true))
	}
	id := uuid.MustParse("00112233-4455-6677-8899-aabbccddeeff")
	pk := &PlayerList{Entries: []protocol.PlayerListEntry{{ActionType: protocol.PlayerListActionRemove, UUID: id}}}
	var buf bytes.Buffer
	pk.Marshal(protocol.NewWriter(&buf, 0))
	// Count, variant(Remove), action(Remove).
	if buf.Len() != 19 || !bytes.Equal(buf.Bytes()[:3], []byte{1, 0, 1}) {
		t.Fatalf("PlayerList remove envelope = %x", buf.Bytes())
	}
	var decoded PlayerList
	decoded.Marshal(protocol.NewReader(&buf, 0, true))
	if len(decoded.Entries) != 1 || decoded.Entries[0].UUID != id || buf.Len() != 0 {
		t.Fatal("PlayerList did not consume the complete entry")
	}
}

func TestProtocol2193ScoreAction(t *testing.T) {
	pk := &SetScore{Entries: []protocol.ScoreboardEntry{{EntryID: 42, IdentityType: protocol.ScoreboardIdentityRemove}}}
	var buf bytes.Buffer
	pk.Marshal(protocol.NewWriter(&buf, 0))
	want := []byte{1, 0, 6, 'R', 'e', 'm', 'o', 'v', 'e', 84, 0}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Fatalf("SetScore = %x, want %x", buf.Bytes(), want)
	}
}

func TestProtocol2193SubChunkResults(t *testing.T) {
	results := []byte{protocol.SubChunkResultSuccess, protocol.SubChunkResultChunkNotFound, protocol.SubChunkResultInvalidDimension, protocol.SubChunkResultPlayerNotFound, protocol.SubChunkResultIndexOutOfBounds, protocol.SubChunkResultSuccessAllAir}
	for i, result := range results {
		if result != byte(i+1) {
			t.Fatalf("result %d = %d", i, result)
		}
	}
}
