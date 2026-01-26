package scanner

import (
	"errors"
	"net"
	"testing"
	"time"
)

type errConn struct {
	closeErr error
}

func (c errConn) Read([]byte) (int, error)         { return 0, errors.New("read not implemented") }
func (c errConn) Write([]byte) (int, error)        { return 0, errors.New("write not implemented") }
func (c errConn) Close() error                     { return c.closeErr }
func (c errConn) LocalAddr() net.Addr              { return &net.TCPAddr{} }
func (c errConn) RemoteAddr() net.Addr             { return &net.TCPAddr{} }
func (c errConn) SetDeadline(time.Time) error      { return nil }
func (c errConn) SetReadDeadline(time.Time) error  { return nil }
func (c errConn) SetWriteDeadline(time.Time) error { return nil }

func setDialTimeout(t *testing.T, f func(string, string, time.Duration) (net.Conn, error)) func() {
	t.Helper()
	original := dialTimeout
	dialTimeout = f
	return func() { dialTimeout = original }
}

func TestScanPort(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		restore := setDialTimeout(t, func(network, address string, timeout time.Duration) (net.Conn, error) {
			if network != "tcp" {
				t.Fatalf("unexpected network: %q", network)
			}
			if address != "[1.2.3.4]:443" {
				t.Fatalf("unexpected address: %q", address)
			}
			connA, connB := net.Pipe()
			t.Cleanup(func() { _ = connB.Close() })
			return connA, nil
		})
		defer restore()

		err := ScanPort("[1.2.3.4]", 443, time.Second)
		if err != nil {
			t.Fatalf("ScanPort returned error: %v", err)
		}
	})

	t.Run("dial error", func(t *testing.T) {
		expected := errors.New("dial failed")
		restore := setDialTimeout(t, func(network, address string, timeout time.Duration) (net.Conn, error) {
			return nil, expected
		})
		defer restore()

		err := ScanPort("[1.2.3.4]", 443, time.Second)
		if !errors.Is(err, expected) {
			t.Fatalf("expected error %v, got %v", expected, err)
		}
	})

	t.Run("close error", func(t *testing.T) {
		expected := errors.New("close failed")
		restore := setDialTimeout(t, func(network, address string, timeout time.Duration) (net.Conn, error) {
			return errConn{closeErr: expected}, nil
		})
		defer restore()

		err := ScanPort("[1.2.3.4]", 443, time.Second)
		if !errors.Is(err, expected) {
			t.Fatalf("expected close error %v, got %v", expected, err)
		}
	})
}
