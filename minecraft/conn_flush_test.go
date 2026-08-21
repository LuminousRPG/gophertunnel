package minecraft

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"testing"
	"time"
)

func TestFlushReturnsNetworkWriteError(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	conn := newConn(errorWriteConn{err: context.Canceled}, nil, log, proto{}, 0, false)
	conn.bufferedSend = [][]byte{{1}}

	err := conn.Flush()
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if len(conn.bufferedSendSpare) != 0 {
		t.Fatalf("expected flushed buffer to be cleared, got length %d", len(conn.bufferedSendSpare))
	}
}

type errorWriteConn struct {
	err error
}

func (c errorWriteConn) Read([]byte) (int, error)       { return 0, c.err }
func (c errorWriteConn) Write([]byte) (int, error)      { return 0, c.err }
func (errorWriteConn) Close() error                     { return nil }
func (errorWriteConn) LocalAddr() net.Addr              { return testAddr("local") }
func (errorWriteConn) RemoteAddr() net.Addr             { return testAddr("remote") }
func (errorWriteConn) SetDeadline(time.Time) error      { return nil }
func (errorWriteConn) SetReadDeadline(time.Time) error  { return nil }
func (errorWriteConn) SetWriteDeadline(time.Time) error { return nil }

type testAddr string

func (a testAddr) Network() string { return "test" }
func (a testAddr) String() string  { return string(a) }
