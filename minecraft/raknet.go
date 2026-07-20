package minecraft

import (
	"context"
	"github.com/sandertv/go-raknet"
	"log/slog"
	"net"
)

// RakNet is an implementation of a RakNet v10 Network.
type RakNet struct {
	l *slog.Logger
}

// DialContext ...
func (r RakNet) DialContext(ctx context.Context, address string) (net.Conn, error) {
	return raknet.Dialer{ErrorLog: r.l.With("net origin", "raknet")}.DialContext(ctx, address)
}

// PingContext ...
func (r RakNet) PingContext(ctx context.Context, address string) (response []byte, err error) {
	return raknet.Dialer{ErrorLog: r.l.With("net origin", "raknet")}.PingContext(ctx, address)
}

// Listen ...
func (r RakNet) Listen(address string) (NetworkListener, error) {
	return raknet.ListenConfig{
		ErrorLog: r.l.With("net origin", "raknet"),
		// Capped conservatively below the common real-world path MTU floor
		// (e.g. some mobile carrier NAT/tunnelling paths silently black-hole
		// larger UDP datagrams instead of returning ICMP fragmentation-needed
		// errors, which RakNet's MTU discovery can't detect). This trades a
		// little more fragmentation for connections that would otherwise
		// have data go missing above their real path MTU.
		MaxMTU: 1200,
	}.Listen(address)
}

// init registers the RakNet network.
func init() {
	RegisterNetwork("raknet", func(l *slog.Logger) Network { return RakNet{l: l} })
}
