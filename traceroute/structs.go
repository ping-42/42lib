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
	TTL           int  `json:"Ttl"`
	Retries       int  `json:"Retries"`
	NetCapRaw     bool `json:"NetCapRaw"`
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
