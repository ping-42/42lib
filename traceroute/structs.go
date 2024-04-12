package traceroute

import (
	"time"

	"github.com/google/uuid"
	"github.com/ping-42/42lib/sensorTask"
)

const TaskName = "TRACEROUTE_TASK"

// task extends base Task struct, that implements TaskRunner interface
type task struct {
	sensorTask.Task
	Opts `json:"Opts"`
}

// Opts
type Opts struct {
	Port       int     `json:"Port"`
	Dest       [4]byte `json:"Dest"`
	FirstHop   int     `json:"FirstHop"`
	MaxHops    int     `json:"MaxHops"`
	Timeout    int     `json:"Timeout"`
	Packetsize int     `json:"PacketSize"`
	Retries    int     `json:"Retries"`
}

// GetId gets the id of the task, as received by the server
func (t task) GetId() uuid.UUID {
	return t.Task.Id
}

func (t task) GetSensorId() uuid.UUID {
	return t.Task.SensorId
}

func (t task) GetName() sensorTask.TaskName {
	return TaskName
}

// TracerouteHop type
type Hop struct {
	Success       bool
	Address       [4]byte
	Host          string
	BytesReceived int
	ElapsedTime   time.Duration
	TTL           int
}

type Result struct {
	DestinationAdress [4]byte
	Hops              []Hop
}
