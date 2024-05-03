package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ping-42/42lib/sensor"
)

const TaskName = "HTTP_TASK"

// Task implements the TaskRunner interface
type task struct {
	sensor.Task
	Opts `json:"HttpOpts"`
}

// Opts represents the dbs collection of parameters
type Opts struct {
	TargetDomain   string
	HttpMethod     string
	RequestHeaders http.Header
	RequestBody    []byte
}

// GetId gets the id of the task, as received by the server
func (t task) GetId() uuid.UUID {
	return t.Task.Id
}

func (t task) GetName() sensor.TaskName {
	return TaskName
}

func (t task) GetSensorId() uuid.UUID {
	return t.SensorId
}

// Result stores httpstat info. // TODO: no need to pass around as argument, just return?
type Result struct {
	ResponseCode    int
	ResponseBody    string
	ResponseHeaders http.Header

	// Durations for each phase
	DNSLookup        time.Duration
	TCPConnection    time.Duration
	TLSHandshake     time.Duration
	ServerProcessing time.Duration
	contentTransfer  time.Duration

	// The following represent the timeline of the request
	NameLookup    time.Duration
	Connect       time.Duration
	Pretransfer   time.Duration
	StartTransfer time.Duration
	total         time.Duration

	dnsStart      time.Time
	dnsDone       time.Time
	tcpStart      time.Time
	tcpDone       time.Time
	tlsStart      time.Time
	tlsDone       time.Time
	serverStart   time.Time
	serverDone    time.Time
	transferStart time.Time
	transferDone  time.Time // need to be provided from outside

	// isTLS is true when connection seems to use TLS
	isTLS bool

	// isReused is true when connection is reused (keep-alive)
	isReused bool
}

func (r *Result) durations() map[string]time.Duration {
	return map[string]time.Duration{
		"DNSLookup":        r.DNSLookup,
		"TCPConnection":    r.TCPConnection,
		"TLSHandshake":     r.TLSHandshake,
		"ServerProcessing": r.ServerProcessing,
		"ContentTransfer":  r.contentTransfer,

		"NameLookup":    r.NameLookup,
		"Connect":       r.Connect,
		"Pretransfer":   r.Connect,
		"StartTransfer": r.StartTransfer,
		"Total":         r.total,
	}
}

// PrintStatus prints all collected stats.
func (r *Result) PrintStatus() {
	var builder = strings.Builder{}
	builder.WriteString(fmt.Sprintf("DNS lookup:        %4d ms\n",
		int(r.DNSLookup/time.Millisecond)))
	builder.WriteString(fmt.Sprintf("TCP connection:    %4d ms\n",
		int(r.TCPConnection/time.Millisecond)))
	builder.WriteString(fmt.Sprintf("TLS handshake:     %4d ms\n",
		int(r.TLSHandshake/time.Millisecond)))
	builder.WriteString(fmt.Sprintf("Server processing: %4d ms\n",
		int(r.ServerProcessing/time.Millisecond)))

	if r.total > 0 {
		builder.WriteString(fmt.Sprintf("Content transfer:  %4d ms\n\n",
			int(r.contentTransfer/time.Millisecond)))
	} else {
		builder.WriteString(fmt.Sprintf("Content transfer:  %4s ms\n\n", "-"))
	}

	builder.WriteString(fmt.Sprintf("Name Lookup:    %4d ms\n",
		int(r.NameLookup/time.Millisecond)))
	builder.WriteString(fmt.Sprintf("Connect:        %4d ms\n",
		int(r.Connect/time.Millisecond)))
	builder.WriteString(fmt.Sprintf("Pre Transfer:   %4d ms\n",
		int(r.Pretransfer/time.Millisecond)))
	builder.WriteString(fmt.Sprintf("Start Transfer: %4d ms\n",
		int(r.StartTransfer/time.Millisecond)))

	if r.total > 0 {
		builder.WriteString(fmt.Sprintf("Total:          %4d ms\n",
			int(r.total/time.Millisecond)))
	} else {
		builder.WriteString(fmt.Sprintf("Total:          %4s ms\n", "-"))
	}

}

// WithHTTPStat is a wrapper of http.withClientTrace. It records the
// time of each http trace hooks.
func WithHTTPStat(ctx context.Context, r *Result) context.Context {
	return withClientTrace(ctx, r)
}
