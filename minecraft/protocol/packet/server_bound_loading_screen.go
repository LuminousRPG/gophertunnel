package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// These match ServerboundLoadingScreenPacketType in Mojang's 1.26.40 schema, which has exactly two values and
// starts at zero. An earlier definition here led with an Unknown constant, putting both real values one too
// high, so a client reporting that a loading screen had closed was read as reporting that one had opened.
const (
	LoadingScreenTypeStart = iota
	LoadingScreenTypeEnd
)

// ServerBoundLoadingScreen is sent by the client to tell the server about the state of the loading
// screen that the client is currently displaying.
type ServerBoundLoadingScreen struct {
	// Type is the type of the loading screen event. It is one of the constants that may be found above.
	Type int32
	// LoadingScreenID is the ID of the screen that was previously sent by the server in the ChangeDimension
	// packet. The server should validate that the ID matches the last one it sent.
	LoadingScreenID protocol.Optional[uint32]
}

// ID ...
func (*ServerBoundLoadingScreen) ID() uint32 {
	return IDServerBoundLoadingScreen
}

func (pk *ServerBoundLoadingScreen) Marshal(io protocol.IO) {
	io.Varint32(&pk.Type)
	protocol.OptionalFunc(io, &pk.LoadingScreenID, io.Uint32)
}
