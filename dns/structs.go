//go:build linux
// +build linux

package dns

import (
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/ping-42/42lib/sensor"
	"golang.org/x/sys/unix"

	"github.com/miekg/dns"
)

const TaskName = "DNS_TASK"

// DnsNameServer represents a domain name server.
// nolint:all TODO
type DnsNameServer struct {
	addr     string
	port     int
	asn      string
	asnOrg   string
	location string
	usesTcp  bool
	usesIpv6 bool
	conn     *dns.Conn
}

// Task implements the TaskRunner interface
type task struct {
	sensor.Task
	Opts `json:"DnsOpts"`
}

// Opts represents the dbs collection of parameters
type Opts struct {
	Host          string `json:"Host"`
	Proto         string
	GetSocketInfo func(conn *net.TCPConn) (*unix.TCPInfo, error)        `json:"-"`
	GetDnsConn    func(addr, proto string, port int) (*dns.Conn, error) `json:"-"`
	DnsUdpClient  *dns.Client
}

// GetId gets the id of the task, as received by the server
func (t task) GetId() uuid.UUID {
	return t.Task.Id
}

func (t task) GetName() sensor.TaskName {
	return TaskName
}

func (t task) GetSensorId() uuid.UUID {
	return t.Task.SensorId
}

// // SetId sets the id of the task. Should be done once upon receiving a message from the server
// func (t task) SetId(id string) {
// 	t.Task.Id = id
// }

// func (t task) SetOpts(o Opts) {
// 	t.Opts = o
// }

// buildOpts parses and validates the message from the server and updates the default opts of the task
// func (t task) buildOpts(msg []byte) error {
// 	var ret task
// 	err := json.Unmarshal(msg, &ret)
// 	if err != nil {
// 		return err
// 	}

// 	// validate
// 	_, err = url.Parse(ret.Host)
// 	if err != nil {
// 		return err
// 	}
// 	t.Host = ret.Host

// 	// use UDP by default
// 	t.Proto = ret.Proto
// 	if ret.Proto == "" {
// 		ret.Proto = "dns"
// 	}

// 	return nil
// }

// Result implements the TaskResult interface
// Represents the specific DNS metrics
type Result struct {
	QueryRtt time.Duration
	SockRtt  time.Duration
	RespSize int64
	Proto    int8
	AnswerA  []*dns.A
}

// GetIpSlice returns all IPs, both TCP and UDP, in the DNS answer as a slice
func (r Result) GetIpSlice() []net.IP {
	ipSlice := make([]net.IP, len(r.AnswerA))
	for i, record := range r.AnswerA {
		ipSlice[i] = record.A
	}
	return ipSlice
}

// getAddrPort returns addr:port as a string
func (dns DnsNameServer) getAddrPort() string {
	return fmt.Sprintf("%v:%v", dns.addr, dns.port)
}
