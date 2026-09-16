package protocol

import (
	"bytes"
	"testing"
)

func TestScoreboardRemoveObjectiveUsesSingleOptional(t *testing.T) {
	entry := ScoreboardEntry{
		EntryID:       42,
		ObjectiveName: "sidebar",
		IdentityType:  ScoreboardIdentityRemove,
	}
	buf := bytes.NewBuffer(nil)
	entry.Marshal(NewWriter(buf, 0))

	// 1.26.45 removed the redundant accessor presence byte. Protocol 2193
	// still has exactly one optional marker for the objective.
	wantSuffix := append([]byte{1, byte(len(entry.ObjectiveName))}, entry.ObjectiveName...)
	if !bytes.HasSuffix(buf.Bytes(), wantSuffix) {
		t.Fatalf("remove entry suffix = %x, want suffix %x", buf.Bytes(), wantSuffix)
	}

	var decoded ScoreboardEntry
	decoded.Marshal(NewReader(bytes.NewBuffer(buf.Bytes()), 0, true))
	if decoded != entry {
		t.Fatalf("decoded entry = %#v, want %#v", decoded, entry)
	}
}

func TestScoreboardRemoveWithoutObjectiveUsesAbsentMarker(t *testing.T) {
	entry := ScoreboardEntry{EntryID: 42, IdentityType: ScoreboardIdentityRemove}
	buf := bytes.NewBuffer(nil)
	entry.Marshal(NewWriter(buf, 0))

	want := []byte{0, 6, 'R', 'e', 'm', 'o', 'v', 'e', 84, 0}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Fatalf("remove entry = %x, want %x", buf.Bytes(), want)
	}
}
