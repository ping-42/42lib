package testingkit

import "golang.org/x/sys/unix"

type MockedSysUnix struct {
	SocketFunc            func(domain int, typ int, proto int) (fd int, err error)
	CloseFunc             func(fd int) (err error)
	BindFunc              func(fd int, sa unix.Sockaddr) (err error)
	SetsockoptIntFunc     func(fd int, level int, opt int, value int) (err error)
	SetsockoptTimevalFunc func(fd int, level int, opt int, tv *unix.Timeval) (err error)
	SendtoFunc            func(fd int, p []byte, flags int, to unix.Sockaddr) (err error)
	RecvfromFunc          func(fd int, p []byte, flags int) (n int, from unix.Sockaddr, err error)
	NsecToTimevalFunc     func(nsec int64) unix.Timeval
}

func (m MockedSysUnix) Socket(domain int, typ int, proto int) (fd int, err error) {
	if m.SocketFunc != nil {
		return m.SocketFunc(domain, typ, proto)
	}
	return 0, nil
}
func (m MockedSysUnix) Close(fd int) (err error) {
	if m.CloseFunc != nil {
		return m.CloseFunc(fd)
	}
	return nil
}
func (m MockedSysUnix) Bind(fd int, sa unix.Sockaddr) (err error) {
	if m.BindFunc != nil {
		return m.BindFunc(fd, sa)
	}
	return nil
}
func (m MockedSysUnix) SetsockoptInt(fd int, level int, opt int, value int) (err error) {
	if m.SetsockoptTimevalFunc != nil {
		return m.SetsockoptInt(fd, level, opt, value)
	}
	return nil
}
func (m MockedSysUnix) SetsockoptTimeval(fd int, level int, opt int, tv *unix.Timeval) (err error) {
	if m.SetsockoptTimevalFunc != nil {
		return m.SetsockoptTimevalFunc(fd, level, opt, tv)
	}
	return nil
}
func (m MockedSysUnix) Sendto(fd int, p []byte, flags int, to unix.Sockaddr) (err error) {
	if m.SendtoFunc != nil {
		return m.SendtoFunc(fd, p, flags, to)
	}
	return nil
}
func (m MockedSysUnix) Recvfrom(fd int, p []byte, flags int) (n int, from unix.Sockaddr, err error) {
	if m.RecvfromFunc != nil {
		return m.RecvfromFunc(fd, p, flags)
	}
	return 0, nil, nil
}
func (m MockedSysUnix) NsecToTimeval(nsec int64) unix.Timeval {
	if m.NsecToTimevalFunc != nil {
		return m.NsecToTimevalFunc(nsec)
	}
	return unix.Timeval{}
}
