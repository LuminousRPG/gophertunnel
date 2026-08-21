package minecraft

import (
	"context"
	"log/slog"
	"net"

	"github.com/sandertv/go-raknet"
)

// RakNet is an implementation of a RakNet v10 Network.
type RakNet struct {
	l *slog.Logger

	// MaxSendBytesPerSecond limits each accepted connection's outbound
	// datagram rate. Zero disables send pacing.
	MaxSendBytesPerSecond int
	// SendBurstBytes is the number of bytes that may be sent at once after an
	// idle period. Zero uses go-raknet's default.
	SendBurstBytes int
}

// NewRakNet returns a RakNet network that reports transport errors to log.
func NewRakNet(log *slog.Logger) RakNet {
	return RakNet{l: log}
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
		MaxMTU:                1200,
		MaxSendBytesPerSecond: r.MaxSendBytesPerSecond,
		SendBurstBytes:        r.SendBurstBytes,
	}.Listen(address)
}

// init registers the RakNet network.
func init() {
	RegisterNetwork("raknet", func(l *slog.Logger) Network { return NewRakNet(l) })
}
