package packet

import (
	"bytes"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func TestClientboundUpdateSoundDataStopWireFormat(t *testing.T) {
	pk := ClientboundUpdateSoundData{
		ServerSoundHandle: 1,
		Stop: protocol.Option(protocol.SoundDataUpdate{
			Type: protocol.SoundDataUpdateStop,
		}),
	}
	buf := bytes.NewBuffer(nil)
	pk.Marshal(protocol.NewWriter(buf, 0))

	// uint64 handle followed directly by seven one-byte zero-valued union
	// discriminators. There are no optional-presence bools on these slots.
	want := []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Fatalf("stop payload = %x, want %x", buf.Bytes(), want)
	}

	var decoded ClientboundUpdateSoundData
	decoded.Marshal(protocol.NewReader(bytes.NewBuffer(buf.Bytes()), 0, true))
	if decoded.ServerSoundHandle != 1 {
		t.Fatalf("decoded handle = %d, want 1", decoded.ServerSoundHandle)
	}
	if update, ok := decoded.Stop.Value(); !ok || update.Type != protocol.SoundDataUpdateStop {
		t.Fatalf("decoded stop = %#v, present=%v", update, ok)
	}
}

func TestClientboundUpdateSoundDataVolumeWireFormat(t *testing.T) {
	pk := ClientboundUpdateSoundData{
		ServerSoundHandle: 2,
		SetVolume: protocol.Option(protocol.SoundDataUpdate{
			Type:   protocol.SoundDataUpdateSetVolume,
			Volume: 0.5,
		}),
	}
	buf := bytes.NewBuffer(nil)
	pk.Marshal(protocol.NewWriter(buf, 0))

	// handle + Stop default + SetVolume discriminator/payload + five defaults.
	want := []byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0x00, 0x3f, 0, 0, 0, 0, 0}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Fatalf("volume payload = %x, want %x", buf.Bytes(), want)
	}
}

func TestClientboundUpdateSoundDataPauseResumeWireFormat(t *testing.T) {
	tests := []struct {
		name string
		pk   ClientboundUpdateSoundData
		want []byte
	}{
		{
			name: "pause",
			pk: ClientboundUpdateSoundData{ServerSoundHandle: 3, Pause: protocol.Option(protocol.SoundDataUpdate{
				Type: protocol.SoundDataUpdatePause,
			})},
			want: []byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0},
		},
		{
			name: "resume",
			pk: ClientboundUpdateSoundData{ServerSoundHandle: 4, Resume: protocol.Option(protocol.SoundDataUpdate{
				Type: protocol.SoundDataUpdateResume,
			})},
			want: []byte{4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			test.pk.Marshal(protocol.NewWriter(buf, 0))
			if !bytes.Equal(buf.Bytes(), test.want) {
				t.Fatalf("payload = %x, want %x", buf.Bytes(), test.want)
			}
		})
	}
}
