package packet

import (
	"bytes"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func TestClientboundUpdateSoundDataWireFormat(t *testing.T) {
	tests := []struct {
		name   string
		update protocol.SoundDataUpdate
		want   []byte
	}{
		{"stop", protocol.SoundDataUpdate{Type: protocol.SoundDataUpdateStop}, []byte{1, 0, 0, 0, 0, 0, 0, 0, 0}},
		{"volume", protocol.SoundDataUpdate{Type: protocol.SoundDataUpdateSetVolume, Volume: 0.5}, []byte{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0x3f}},
		{"pause", protocol.SoundDataUpdate{Type: protocol.SoundDataUpdatePause}, []byte{1, 0, 0, 0, 0, 0, 0, 0, 5}},
		{"resume", protocol.SoundDataUpdate{Type: protocol.SoundDataUpdateResume}, []byte{1, 0, 0, 0, 0, 0, 0, 0, 6}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pk := ClientboundUpdateSoundData{ServerSoundHandle: 1, Update: test.update}
			buf := bytes.NewBuffer(nil)
			pk.Marshal(protocol.NewWriter(buf, 0))
			// Since 1.26.60 the packet holds a single update: the handle, the update type and its payload.
			if !bytes.Equal(buf.Bytes(), test.want) {
				t.Fatalf("payload = %x, want %x", buf.Bytes(), test.want)
			}

			var decoded ClientboundUpdateSoundData
			decoded.Marshal(protocol.NewReader(bytes.NewBuffer(buf.Bytes()), 0, true))
			if decoded != pk {
				t.Fatalf("decoded = %#v, want %#v", decoded, pk)
			}
		})
	}
}
