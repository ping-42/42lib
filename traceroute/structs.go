package traceroute

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/ping-42/42lib/sensor"
	"golang.org/x/sys/unix"
)

const TaskName = "TRACEROUTE_TASK"

// task extends base Task struct, that implements TaskRunner interface
type task struct {
	sensor.Task
	Opts    `json:"Opts"`
	SysUnix SysUnix
}

// Opts for the task
type Opts struct {
	Port          int      `json:"Port"`
	Dest          net.IP   `json:"Dest"`
	CurrentAddr   net.IP   `json:"CurrentAddr"`
	CurrentHost   []string `json:"CurrentHost"`
	ReceiveSocket int      `json:"ReceiveSocket"`
	SendSocket    int      `json:"SendSocket"`
	FirstHop      int      `json:"FirstHop"`
	MaxHops       int      `json:"MaxHops"`
	Timeout       int      `json:"Timeout"`
	Packetsize    int      `json:"PacketSize"`
	Packet        []byte
	TTL           int `json:"Ttl"`
	Retries       int `json:"Retries"`
	// NetCapRaw     bool `json:"NetCapRaw"`
}

// SysUnix is an interface for interacting with low-level system.
type SysUnix interface {
	Socket(domain int, typ int, proto int) (fd int, err error)
	Close(fd int) (err error)
	Bind(fd int, sa unix.Sockaddr) (err error)
	SetsockoptInt(fd int, level int, opt int, value int) (err error)
	SetsockoptTimeval(fd int, level int, opt int, tv *unix.Timeval) (err error)
	Sendto(fd int, p []byte, flags int, to unix.Sockaddr) (err error)
	Recvfrom(fd int, p []byte, flags int) (n int, from unix.Sockaddr, err error)
	NsecToTimeval(nsec int64) unix.Timeval
}

// GetId gets the id of the task, as received by the server
func (t task) GetId() uuid.UUID {
	return t.Task.Id
}

func (t task) GetSensorId() uuid.UUID {
	return t.Task.SensorId
}

func (t task) GetName() sensor.TaskName {
	return TaskName
}

// Traceroute Hop type
type Hop struct {
	Success       bool
	Address       net.IP
	Host          string
	BytesReceived int
	ElapsedTime   time.Duration
	TTL           int
	Error         error
}

type Result struct {
	DestinationAdress net.IP
	Hops              []Hop
}

// SysUnixReal will be used for the real socket operation methods
// should these methods be defined in structs.go?
type SysUnixReal struct{}

func (s SysUnixReal) Socket(domain int, typ int, proto int) (fd int, err error) {
	return unix.Socket(domain, typ, proto)
}

func (s SysUnixReal) Close(fd int) (err error) {
	return unix.Close(fd)
}

func (s SysUnixReal) Bind(fd int, sa unix.Sockaddr) (err error) {
	return unix.Bind(fd, sa)
}

func (s SysUnixReal) SetsockoptInt(fd int, level int, opt int, value int) error {
	return unix.SetsockoptInt(fd, level, opt, value)
}

func (s SysUnixReal) SetsockoptTimeval(fd int, level int, opt int, tv *unix.Timeval) error {
	return unix.SetsockoptTimeval(fd, level, opt, tv)
}

func (s SysUnixReal) Sendto(fd int, p []byte, flags int, to unix.Sockaddr) (err error) {
	return unix.Sendto(fd, p, flags, to)
}

func (s SysUnixReal) Recvfrom(fd int, p []byte, flags int) (n int, from unix.Sockaddr, err error) {
	return unix.Recvfrom(fd, p, flags)
}

func (s SysUnixReal) NsecToTimeval(nsec int64) unix.Timeval {
	return unix.NsecToTimeval(nsec)
}
