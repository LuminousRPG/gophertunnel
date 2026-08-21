package minecraft

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func TestFormatPongDataUsesStandardFieldCount(t *testing.T) {
	status := ServerStatus{
		ServerName:    "Luminous RPG",
		ServerSubName: "PlayServer",
		PlayerCount:   8,
		MaxPlayers:    50,
	}
	got := string(formatPongData(status, 1234, 19132))
	want := fmt.Sprintf("MCPE;Luminous RPG;%d;%s;8;50;1234;PlayServer;Creative;1;19132;19132;",
		protocol.CurrentProtocol, protocol.CurrentVersion)
	if got != want {
		t.Fatalf("formatPongData() = %q, want %q", got, want)
	}
	if fields := strings.Split(strings.TrimSuffix(got, ";"), ";"); len(fields) != 12 {
		t.Fatalf("formatPongData() produced %d fields, want 12: %q", len(fields), got)
	}
}
