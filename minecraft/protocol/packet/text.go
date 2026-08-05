package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	TextTypeRaw = iota
	TextTypeChat
	TextTypeTranslation
	TextTypePopup
	TextTypeJukeboxPopup
	TextTypeTip
	TextTypeSystem
	TextTypeWhisper
	TextTypeAnnouncement
	TextTypeObjectWhisper
	TextTypeObject
	TextTypeObjectAnnouncement
)

const (
	TextCategoryMessageOnly = iota
	TextCategoryAuthoredMessage
	TextCategoryMessageWithParameters
)

// Text is sent by the client to the server to send chat messages, and by the server to the client to forward
// or send messages, which may be chat, popups, tips etc.
type Text struct {
	// TextType is the type of the text sent. When a client sends this to the server, it should always be
	// TextTypeChat. If the server sends it, it may be one of the other text types above.
	TextType byte
	// NeedsTranslation specifies if any of the messages need to be translated. It seems that where % is found
	// in translatable text types, these are translated regardless of this bool. Translatable text types
	// include TextTypeTranslation, TextTypeTip, TextTypePopup and TextTypeJukeboxPopup.
	NeedsTranslation bool
	// SourceName is the name of the source of the messages. This source is displayed in text types such as
	// the TextTypeChat and TextTypeWhisper, where typically the username is shown.
	SourceName string
	// Message is the message of the packet. This field is set for each TextType and is the main component of
	// the packet.
	Message string
	// Parameters is a list of parameters that should be filled into the message. These parameters are only
	// written if the type of the packet is TextTypeTranslation, TextTypeTip, TextTypePopup or TextTypeJukeboxPopup.
	Parameters []string
	// XUID is the XBOX Live user ID of the player that sent the message. It is only set for packets of
	// TextTypeChat. When sent to a player, the player will only be shown the chat message if a player with
	// this XUID is present in the player list and not muted, or if the XUID is empty.
	XUID string
	// PlatformChatID is an identifier only set for particular platforms when chatting (presumably only for
	// Nintendo Switch). It is otherwise an empty string, and is used to decide which players are able to
	// chat with each other.
	PlatformChatID string
	// FilteredMessage is a filtered version of Message with all the profanity removed. The client will use
	// this over Message if this field is not empty and they have the "Filter Profanity" setting enabled.
	FilteredMessage protocol.Optional[string]
}

// ID ...
func (*Text) ID() uint32 {
	return IDText
}

// textTypeCategories lists, for each category, the text types it carries in the order Mojang's 1.26.40 schema
// gives them (MessageOnly.json, AuthorAndMessage.json and MessageAndParams.json). The byte on the wire is the
// index of the type *within its category*, not the TextType constant: chat is 0 rather than 1, and a type
// like TextTypeSystem would otherwise be written as 6, past the end of the six entry list it belongs to.
var textTypeCategories = [...][]byte{
	TextCategoryMessageOnly:           {TextTypeRaw, TextTypeTip, TextTypeSystem, TextTypeObjectWhisper, TextTypeObject, TextTypeObjectAnnouncement},
	TextCategoryAuthoredMessage:       {TextTypeChat, TextTypeWhisper, TextTypeAnnouncement},
	TextCategoryMessageWithParameters: {TextTypeTranslation, TextTypePopup, TextTypeJukeboxPopup},
}

// textTypeIndex returns the category a text type belongs to and its index within that category.
func textTypeIndex(textType byte) (category, index uint8) {
	for c, types := range textTypeCategories {
		for i, t := range types {
			if t == textType {
				return uint8(c), uint8(i)
			}
		}
	}
	return TextCategoryMessageOnly, 0
}

func (pk *Text) Marshal(io protocol.IO) {
	io.Bool(&pk.NeedsTranslation)

	// Both values are derived from the text type when writing and resolved back into it when reading, so the
	// switch below sees the same type in either direction.
	categoryType, index := textTypeIndex(pk.TextType)
	io.Uint8(&categoryType)
	io.Uint8(&index)
	if int(categoryType) >= len(textTypeCategories) {
		io.UnknownEnumOption(categoryType, "text category")
		return
	}
	types := textTypeCategories[categoryType]
	if int(index) >= len(types) {
		io.UnknownEnumOption(index, "text type")
		return
	}
	pk.TextType = types[index]

	switch pk.TextType {
	case TextTypeChat, TextTypeWhisper, TextTypeAnnouncement:
		io.String(&pk.SourceName)
		io.String(&pk.Message)
	case TextTypeRaw, TextTypeTip, TextTypeSystem, TextTypeObject, TextTypeObjectWhisper, TextTypeObjectAnnouncement:
		io.String(&pk.Message)
	case TextTypeTranslation, TextTypePopup, TextTypeJukeboxPopup:
		io.String(&pk.Message)
		protocol.FuncSlice(io, &pk.Parameters, io.String)
	}

	if len(pk.Message) == 0 {
		io.InvalidValue(pk.Message, "message", "string cannot be empty")
	}
	io.String(&pk.XUID)
	io.String(&pk.PlatformChatID)
	protocol.OptionalFunc(io, &pk.FilteredMessage, io.String)
}
