package protocol

import (
	"bytes"
	"testing"
)

func TestScoreboardRemoveObjectiveUsesDoubleOptional(t *testing.T) {
	entry := ScoreboardEntry{
		EntryID:       42,
		ObjectiveName: "sidebar",
		IdentityType:  ScoreboardIdentityRemove,
	}
	buf := bytes.NewBuffer(nil)
	entry.Marshal(NewWriter(buf, 0))

	// The outer and inner optional markers must both be present before the
	// objective string in 1.26.44.
	wantSuffix := append([]byte{1, 1, byte(len(entry.ObjectiveName))}, entry.ObjectiveName...)
	if !bytes.HasSuffix(buf.Bytes(), wantSuffix) {
		t.Fatalf("remove entry suffix = %x, want suffix %x", buf.Bytes(), wantSuffix)
	}

	var decoded ScoreboardEntry
	decoded.Marshal(NewReader(bytes.NewBuffer(buf.Bytes()), 0, true))
	if decoded != entry {
		t.Fatalf("decoded entry = %#v, want %#v", decoded, entry)
	}
}

func TestScoreboardRemoveWithoutObjectiveUsesNestedAbsentMarker(t *testing.T) {
	entry := ScoreboardEntry{EntryID: 42, IdentityType: ScoreboardIdentityRemove}
	buf := bytes.NewBuffer(nil)
	entry.Marshal(NewWriter(buf, 0))

	if !bytes.HasSuffix(buf.Bytes(), []byte{1, 0}) {
		t.Fatalf("remove entry suffix = %x, want outer-present/inner-absent markers 0100", buf.Bytes())
	}
}
