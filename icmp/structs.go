package icmp

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/ping-42/42lib/dns"
	"github.com/ping-42/42lib/sensor"
)

const TaskName = "ICMP_TASK"

// task extends base Task struct, that implements TaskRunner interface
type task struct {
	sensor.Task
	Opts `json:"Opts"`
}

// Opts defines the parameter payload to call PingHost
type Opts struct {
	// TargetDomain is not required
	TargetDomain string   `json:"TargetDomain"`
	TargetIPs    []net.IP `json:"TargetIPs"`
	Count        int      `json:"Count"`
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

// IcmpPingResult represents a single completed ping
type IcmpPingResult struct {
	IPAddr       net.IP
	Sequence     int
	BytesWritten int
	BytesRead    int
	RTT          time.Duration
	Error        error
}

// IcmpPingResult represents a completed ping with stats
// FailureMessages represents all errors that occurred during the ICMP task.
// TODO add mdev
type IcmpPingStats struct {
	IPAddr          net.IP
	PacketsSent     int
	PacketsReceived int
	BytesWritten    int
	BytesRead       int
	TotalRTT        time.Duration
	MinRTT          time.Duration
	MaxRTT          time.Duration
	AverageRTT      time.Duration
	Loss            float64
	FailureMessages []string
}

type Result struct {
	ResultPerIp map[string]IcmpPingStats
	DnsResult   dns.Result
}

// ICMPConn represents an interface for an ICMP connection.
type ICMPConn interface {
	WriteTo(b []byte, addr net.Addr) (int, error)
	ReadFrom(b []byte) (int, net.Addr, error)
	Close() error
	SetDeadline(t time.Time) error
}
