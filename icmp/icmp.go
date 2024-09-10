package icmp

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/ping-42/42lib/db/models"
	"github.com/ping-42/42lib/dns"
	"github.com/ping-42/42lib/helpers"
	logger "github.com/ping-42/42lib/logger"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

const (
	// maxIcmpTimeout is the global timeout for all pings to complete
	maxIcmpTimeout   = 25 * time.Second
	timeoutMessage   = "ping timed out"
	winPayloadString = "08090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f3031323334353637" // length 48
)

var (
	loggerIcmp = logger.WithTestType("icmp")
)

// NewTaskFromBytes used in sensor for building the task from the received bytes
func NewTaskFromBytes(msg []byte) (t task, err error) {

	// build the icmp task from the received msg
	err = json.Unmarshal(msg, &t)
	if err != nil {
		err = fmt.Errorf("icmp.NewTaskFromBytes Unmarshal err task:%v, %v", string(msg), err)
		return
	}

	// set here default opts if we need some
	return
}

// NewTaskFromModel used in scheduler for building the task from the db model task
func NewTaskFromModel(t models.Task) (tRes task, err error) {

	tRes.Id = t.ID
	tRes.SensorId = t.SensorID
	tRes.Name = TaskName

	// build the opts
	var o = Opts{}
	err = json.Unmarshal(t.Opts, &o)
	if err != nil {
		err = fmt.Errorf("icmp NewTaskFromModel Unmarshal Opts err:%v", err)
		return
	}
	tRes.Opts = o
	return
}

// run is the entry point for the icmp task
func (t task) Run(ctx context.Context) ([]byte, error) {
	if ctx.Err() != nil {
		return nil, fmt.Errorf("context done detected in Run:%v", ctx.Err())
	}

	var res Result

	// in case we have TargetDomain this means that we need run also DNS to get the target ips
	if t.TargetDomain != "" {

		dnsRes, err := runDnsTask(ctx, t)
		if err != nil {
			return nil, fmt.Errorf("runDnsTask from icmp err: %v", err)
		}

		if len(dnsRes.AnswerA) == 0 {
			return nil, fmt.Errorf("dns.Result{} in icmp Run returns empty dnsRes.AnswerA")
		}

		// append the IPs
		for _, v := range dnsRes.AnswerA {
			t.TargetIPs = append(t.TargetIPs, v.A)
		}

		// set the results
		res.DnsResult = dnsRes
	}

	// create a new icmpV4 connection
	connV4, err := createConnV4()
	if err != nil {
		err = fmt.Errorf("could not create icmpV4 connection:%v", err)
		return nil, err
	}

	var connV6 *icmp.PacketConn
	hasIPv6 := helpers.HasIPv6(t.TargetIPs)
	if hasIPv6 {
		connV6, err = createConnV6()
		if err != nil {
			err = fmt.Errorf("could not create icmpV6 connection:%v", err)
			return nil, err
		}
	}

	resPerIp, err := t.pingHost(ctx, connV4, connV6)
	if err != nil {
		return nil, err
	}

	res.ResultPerIp = resPerIp

	resJson, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return resJson, nil
}

// RandInt returns a random int to mock uint16 process ID for ICMP message
func randInt() (int, error) {
	var b [2]byte
	if _, err := rand.Read(b[:]); err != nil {
		loggerIcmp.Error("could not generate random icmp message ID: ", err)
		return 0, err
	}
	randomInt := int(binary.LittleEndian.Uint16(b[:])) % 65537 // 65536 + 1

	return randomInt, nil
}

// getICMPType takes an IP returns the ICMP type based on the IP version
func getICMPType(ip net.IP) icmp.Type {
	if ip.To4() == nil {
		return ipv6.ICMPTypeEchoRequest
	}
	return ipv4.ICMPTypeEcho
}

// createConnV4 returns a new icmpV4 connection
func createConnV4() (conn *icmp.PacketConn, err error) {
	// create a new ipv4 icmp connection
	conn, err = icmp.ListenPacket("ip4:icmp", "0.0.0.0")

	if err != nil {
		return nil, fmt.Errorf("could not create icmpV4 connection: %w", err)
	}

	return conn, nil
}

// createConnV6 returns a new icmpV6 connection
func createConnV6() (conn *icmp.PacketConn, err error) {
	// create a new ipv6 icmp connection
	conn, err = icmp.ListenPacket("ip6:ipv6-icmp", "::")

	if err != nil {
		return nil, fmt.Errorf("could not create icmpV6 connection: %w", err)
	}

	return conn, nil
}

// createICMPMessage takes the task context ICMP type and returns an icmp message with the given type
func (t task) createICMPMessage(ctx context.Context, msgType icmp.Type) (*icmp.Message, error) {
	// rand int64 for the icmp message id
	randInt, err := randInt()
	if err != nil {
		loggerIcmp.Error("could not generate random int64: ", err)
		return nil, err
	}
	return &icmp.Message{
		Type: msgType,
		Code: 0,
		Body: &icmp.Echo{
			ID:   randInt,
			Seq:  1,
			Data: t.Payload,
		},
	}, nil
}

// pingIteration takes the task context, connection, ICMP type, IP, pingStats and iteration number and returns the pingStats
func (t task) pingIteration(ctx context.Context, conn ICMPConn, icmpType icmp.Type, targetIP net.IP, pingStats *IcmpPingStats, i int) error {
	// check if the context is cancelled
	select {
	case <-ctx.Done():
		loggerIcmp.Error("context cancelled")
		pingStats.FailureMessages = append(pingStats.FailureMessages, ctx.Err().Error())
		return ctx.Err()
	default:
		// run the ping
		pingResult, err := t.runICMP(ctx, conn, icmpType, targetIP)
		pingStats.PacketsSent++
		pingStats.BytesWritten += pingResult.BytesWritten
		pingStats.BytesRead += pingResult.BytesRead
		pingStats.TotalRTT += pingResult.RTT
		if err != nil {
			pingStats.FailureMessages = append(pingStats.FailureMessages, err.Error())
			loggerIcmp.Error("error running ping: ", err)
			return err
		}
		pingResult.Sequence = i
		pingStats.PacketsReceived++
		if pingResult.RTT < pingStats.MinRTT || pingStats.MinRTT == 0 {
			pingStats.MinRTT = pingResult.RTT
		}
		if pingResult.RTT > pingStats.MaxRTT || pingStats.MaxRTT == 0 {
			pingStats.MaxRTT = pingResult.RTT
		}

		loggerIcmp.Infof("ping result: %v", pingResult)
	}
	return nil
}

// runICMP takes the task context, connection, ICMP type and IP and returns the ping result
func (t *task) runICMP(ctx context.Context, conn ICMPConn, msgType icmp.Type, ip net.IP) (IcmpPingResult, error) {
	result := IcmpPingResult{
		IPAddr: ip,
	}
	// create an ICMP message
	msg, err := t.createICMPMessage(ctx, msgType)
	if err != nil {
		loggerIcmp.Error("could not create an ICMP message: ", err)
		result.Error = err
		return result, err
	}

	// marshal the message into bytes
	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		loggerIcmp.Error("could not marshal the ICMP message: ", err)
		result.Error = err
		return result, err
	}

	// set the read and write deadline
	err = conn.SetDeadline(time.Now().Add(maxIcmpTimeout))
	if err != nil {
		loggerIcmp.Error("could not set the read and write deadline: ", err)
		result.Error = err
		return result, err
	}

	// write the message to destination
	start := time.Now()
	_, err = conn.WriteTo(msgBytes, &net.IPAddr{IP: ip})
	if err != nil {
		loggerIcmp.Error("could not write message: ", err)
		result.Error = err
		return result, err
	}
	result.BytesWritten = len(msgBytes)

	// read the response
	resp := make([]byte, 1500)
	n, _, err := conn.ReadFrom(resp)
	if err != nil {
		loggerIcmp.Error("could not read message: ", err)
		result.Error = err
		return result, err
	}
	result.BytesRead = n
	elapsed := time.Since(start)
	result.RTT = elapsed

	return result, nil
}

// pingHost takes the task context and returns the statistics of the pings
func (t task) pingHost(ctx context.Context, connV4 ICMPConn, connV6 ICMPConn) (resultsPerIp map[string]IcmpPingStats, err error) {

	// init the results
	resultsPerIp = make(map[string]IcmpPingStats)

	// hasIPv6 is used to check if there are any IPv6 addresses
	hasIPv6 := helpers.HasIPv6(t.TargetIPs)

	// Create a context with a timeout for the entire PingHost operation
	ctx, cancel := context.WithTimeout(ctx, maxIcmpTimeout)
	defer cancel()

	// defer close v4 connection and return posible error
	defer func() {
		err := connV4.Close()
		if err != nil {
			loggerIcmp.Error("could not close connection: ", err)
			return
		}
	}()

	// create a new icmpV6 connection if available and return error if not
	if hasIPv6 {
		// defer close connection and return posible error
		defer func() {
			err := connV6.Close()
			if err != nil {
				loggerIcmp.Error("could not close connection: ", err)
				return
			}
		}()
	}

	// ping each ip
	for _, targetIP := range t.TargetIPs {
		var pingStats IcmpPingStats
		var conn ICMPConn
		// check if the ip is ipv4 or ipv6 and use the corresponding connection
		if targetIP.To4() != nil {
			conn = connV4
		} else if targetIP.To16() != nil {
			conn = connV6
		}
		// run the icmp ping per opts.Count and populate the pingStats &
		// collect errors from each ping if they occur
		for i := 0; i < t.Count; i++ {
			err := t.pingIteration(ctx, conn, getICMPType(targetIP), targetIP, &pingStats, i)
			if err != nil {
				loggerIcmp.Error("error running ping: ", err)
				pingStats.FailureMessages = append(pingStats.FailureMessages, err.Error())
				continue
			}
		}

		// calculate the average RTT
		if pingStats.PacketsReceived > 0 {
			totalRTTMicroseconds := pingStats.TotalRTT.Microseconds()
			averageRTTMicroseconds := totalRTTMicroseconds / int64(pingStats.PacketsReceived)
			pingStats.AverageRTT = time.Duration(averageRTTMicroseconds) * time.Microsecond
		}

		// calculate the loss
		pingStats.Loss = float64(pingStats.PacketsSent-pingStats.PacketsReceived) / float64(pingStats.PacketsSent) * 100

		pingStats.IPAddr = targetIP

		// collect the results and send them to the channel
		resultsPerIp[targetIP.String()] = pingStats
	}

	loggerIcmp.Infof("all results collected (%d)", len(resultsPerIp))
	return
}

func runDnsTask(ctx context.Context, t task) (dnsRes dns.Result, err error) {
	dnsTask := dns.NewEmptyTask()
	dnsTask.Id = t.Id
	dnsTask.SensorId = t.SensorId
	dnsTask.Name = dns.TaskName
	h, err := helpers.ExtractDomainFromUrl(t.TargetDomain)
	if err != nil {
		err = fmt.Errorf("ExtractDomainFromUrl err:%v, %v", err, t.TargetDomain)
		return
	}
	dnsTask.Opts.Host = h

	dnsResb, err := dnsTask.Run(ctx)
	if err != nil {
		err = fmt.Errorf("error in icmp/dns Run err:%v", err)
		return
	}
	dnsRes = dns.Result{}
	err = json.Unmarshal(dnsResb, &dnsRes)
	if err != nil {
		err = fmt.Errorf("Unmarshal dns.Result{} in icmp Run err:%v", err)
		return
	}
	return
}
