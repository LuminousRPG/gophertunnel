package packet

import (
	"bytes"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// Pin the category and original TextType byte used by the 26.50 reference codec.
func TestTextCategoryNumbering(t *testing.T) {
	for _, tc := range []struct {
		textType         byte
		category, number byte
	}{
		{TextTypeRaw, TextCategoryMessageOnly, 0},
		{TextTypeTip, TextCategoryMessageOnly, 5},
		{TextTypeSystem, TextCategoryMessageOnly, 6},
		{TextTypeObjectWhisper, TextCategoryMessageOnly, 9},
		{TextTypeObject, TextCategoryMessageOnly, 10},
		{TextTypeObjectAnnouncement, TextCategoryMessageOnly, 11},
		{TextTypeChat, TextCategoryAuthoredMessage, 1},
		{TextTypeWhisper, TextCategoryAuthoredMessage, 7},
		{TextTypeAnnouncement, TextCategoryAuthoredMessage, 8},
		{TextTypeTranslation, TextCategoryMessageWithParameters, 2},
		{TextTypePopup, TextCategoryMessageWithParameters, 3},
		{TextTypeJukeboxPopup, TextCategoryMessageWithParameters, 4},
	} {
		pk := &Text{TextType: tc.textType, SourceName: "src", Message: "hi"}
		var buf bytes.Buffer
		pk.Marshal(protocol.NewWriter(&buf, 0))

		// NeedsTranslation, category, then the original TextType.
		if got := buf.Bytes()[1:3]; !bytes.Equal(got, []byte{tc.category, tc.number}) {
			t.Fatalf("text type %d wrote %v, want category/number %v", tc.textType, got, []byte{tc.category, tc.number})
		}
		var decoded Text
		decoded.Marshal(protocol.NewReader(&buf, 0, true))
		if decoded.TextType != tc.textType {
			t.Fatalf("text type %d decoded as %d", tc.textType, decoded.TextType)
		}
		if buf.Len() != 0 {
			t.Fatalf("text type %d left %d unread bytes", tc.textType, buf.Len())
		}
	}
}
