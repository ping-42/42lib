package tastingkit

import (
	"net"
	"time"
)

// MockedICMPConn represents a mocked structure for ICMP connections.
type MockedICMPConn struct {
	WriteToFunc     func(b []byte, addr net.Addr) (int, error)
	ReadFromFunc    func(b []byte) (int, net.Addr, error)
	CloseFunc       func() error
	SetDeadlineFunc func(t time.Time) error
}

// WriteTo implements the WriteTo method of ICMPConn interface.
func (m MockedICMPConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	if m.WriteToFunc != nil {
		return m.WriteToFunc(b, addr)
	}
	return 0, nil
}

// ReadFrom implements the ReadFrom method of ICMPConn interface.
func (m MockedICMPConn) ReadFrom(b []byte) (int, net.Addr, error) {
	if m.ReadFromFunc != nil {
		return m.ReadFromFunc(b)
	}
	return 0, nil, nil
}

// Close implements the Close method of ICMPConn interface.
func (m MockedICMPConn) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// SetDeadline implements the SetDeadline method of ICMPConn interface.
func (m MockedICMPConn) SetDeadline(t time.Time) error {
	if m.SetDeadlineFunc != nil {
		return m.SetDeadlineFunc(t)
	}
	return nil
}
