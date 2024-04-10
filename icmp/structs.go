package icmp

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/ping-42/42lib/dns"
	"github.com/ping-42/42lib/sensorTask"
)

const TaskName = "ICMP_TASK"

// task extends base Task struct, that implements TaskRunner interface
type task struct {
	sensorTask.Task
	Opts `json:"Opts"`
}

// Opts defines the parameter payload to call PingHost
// Since we have pointer parameters, the struct should be also passed by pointer
// to avoid nasty bugs due to struct value copy
type Opts struct {
	// TargetDomain is not required
	TargetDomain string   `json:"TargetDomain"`
	TargetIPs    []net.IP `json:"TargetIPs"`
	Count        int      `json:"Count"`
	Payload      []byte   `json:"Payload"`
}

// overrides the default Unmarshal cuz []net.IP is not working
// here have custom logic per TargetIPs
// func (o *Opts) UnmarshalJSON(data []byte) error {
// 	type Alias Opts
// 	aux := &struct {
// 		TargetIPs []string `json:"TargetIPs"`
// 		*Alias
// 	}{
// 		Alias: (*Alias)(o),
// 	}
// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}
// 	o.TargetIPs = make([]net.IP, len(aux.TargetIPs))
// 	for i, ipStr := range aux.TargetIPs {
// 		o.TargetIPs[i] = net.ParseIP(ipStr)
// 	}
// 	return nil
// }

// // MarshalJSON overrides the default Marshal for Opts
// func (o *Opts) MarshalJSON() ([]byte, error) {
// 	type Alias Opts
// 	aux := &struct {
// 		TargetIPs []string `json:"TargetIPs"`
// 		*Alias
// 	}{
// 		Alias: (*Alias)(o),
// 	}
// 	aux.TargetIPs = make([]string, len(o.TargetIPs))
// 	for i, ip := range o.TargetIPs {
// 		aux.TargetIPs[i] = ip.String()
// 	}
// 	return json.Marshal(aux)
// }

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

// // buildOpts parses and validates the message from the server and updates the default opts of the task.
// func (t task) buildOpts(msg []byte) error {
// 	var ret task
// 	err := json.Unmarshal(msg, &ret)
// 	if err != nil {
// 		return err
// 	}

// 	// assign task ID from server
// 	if ret.Task.Id == "" {
// 		errMsg := "no id found"
// 		err := fmt.Errorf(errMsg)
// 		loggerIcmp.Error(err)
// 		return err
// 	}
// 	t.Id = ret.Task.Id

// 	// assign task name from server
// 	if ret.Task.Name == "" {
// 		errMsg := "no task name found"
// 		err := fmt.Errorf(errMsg)
// 		loggerIcmp.Error(err)
// 		return err
// 	}
// 	t.Name = ret.Task.Name

// 	// validate ips
// 	TargetIPs := ret.TargetIPs
// 	for _, ip := range TargetIPs {
// 		ipString := ip.String()
// 		validIP := net.ParseIP(ipString)
// 		if validIP == nil {
// 			errMsg := "no valid ip"
// 			err := fmt.Errorf(errMsg)
// 			loggerIcmp.Error(err)
// 			return err
// 		}
// 	}
// 	t.TargetIPs = ret.TargetIPs

// 	// set count
// 	if ret.Count == 0 {
// 		errMsg := "count cannot be 0"
// 		err := fmt.Errorf(errMsg)
// 		loggerIcmp.Error(err)
// 		return err
// 	}
// 	t.Count = ret.Count

// 	// set payload
// 	if ret.Payload != nil {
// 		t.Payload = ret.Payload
// 	}

// 	return nil
// }

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
