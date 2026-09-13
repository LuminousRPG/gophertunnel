package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// ClientboundUpdateSoundData is sent by the server to update a sound that is currently playing, identified by
// the handle that the server sent in the PlaySound packet that started it. Each optional field is a Cereal union
// slot that may hold any SoundDataUpdate variant; its name does not constrain the variant on the wire.
type ClientboundUpdateSoundData struct {
	// ServerSoundHandle is the server-side handle of the sound to update.
	ServerSoundHandle uint64
	Stop              protocol.Optional[protocol.SoundDataUpdate]
	SetVolume         protocol.Optional[protocol.SoundDataUpdate]
	SetPitch          protocol.Optional[protocol.SoundDataUpdate]
	Fade              protocol.Optional[protocol.SoundDataUpdate]
	SeekTo            protocol.Optional[protocol.SoundDataUpdate]
	Pause             protocol.Optional[protocol.SoundDataUpdate]
	Resume            protocol.Optional[protocol.SoundDataUpdate]
}

// ID ...
func (*ClientboundUpdateSoundData) ID() uint32 {
	return IDClientboundUpdateSoundData
}

func (pk *ClientboundUpdateSoundData) Marshal(io protocol.IO) {
	io.Uint64(&pk.ServerSoundHandle)
	marshalSoundDataUpdate(io, &pk.Stop, protocol.SoundDataUpdateStop)
	marshalSoundDataUpdate(io, &pk.SetVolume, protocol.SoundDataUpdateSetVolume)
	marshalSoundDataUpdate(io, &pk.SetPitch, protocol.SoundDataUpdateSetPitch)
	marshalSoundDataUpdate(io, &pk.Fade, protocol.SoundDataUpdateFade)
	marshalSoundDataUpdate(io, &pk.SeekTo, protocol.SoundDataUpdateSeekTo)
	marshalSoundDataUpdate(io, &pk.Pause, protocol.SoundDataUpdatePause)
	marshalSoundDataUpdate(io, &pk.Resume, protocol.SoundDataUpdateResume)
}

// marshalSoundDataUpdate writes a Cereal union slot. These fields have a
// default Stop variant, not an optional-presence marker: Writing a bool before
// each union shifts the remainder of the packet and makes Bedrock reject it.
func marshalSoundDataUpdate(io protocol.IO, optional *protocol.Optional[protocol.SoundDataUpdate], expectedType uint32) {
	if _, reading := io.(*protocol.Reader); reading {
		var update protocol.SoundDataUpdate
		update.Marshal(io)
		if update.Type == expectedType {
			*optional = protocol.Option(update)
		} else {
			*optional = protocol.Optional[protocol.SoundDataUpdate]{}
		}
		return
	}

	update, ok := optional.Value()
	if !ok {
		update.Type = protocol.SoundDataUpdateStop
	}
	update.Marshal(io)
}
