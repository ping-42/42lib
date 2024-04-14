package traceroute

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/ping-42/42lib/db/models"
	"github.com/ping-42/42lib/logger"
)

// TODO
// retry should occur within each hop if needed
// when should we retry?
// we should analyze ICMP response
// do we need to resolve the domain to an IP address?

const (
	DEFAULT_PORT        = 33434
	DEFAULT_MAX_HOPS    = 64
	DEFAULT_FIRST_HOP   = 1
	DEFAULT_TIMEOUT_MS  = 500
	DEFAULT_RETRIES     = 3
	DEFAULT_PACKET_SIZE = 52
)

var (
	loggerTraceroute = logger.WithTestType("traceroute")
)

// NewTaskFromBytes used in sensor for building the task from the received bytes
func NewTaskFromBytes(msg []byte) (t task, err error) {

	// build the traceroute task from the received msg
	err = json.Unmarshal(msg, &t)
	if err != nil {
		err = fmt.Errorf("traceroute.NewTaskFromBytes Unmarshal err task:%v, %v", string(msg), err)
		return
	}

	return t, nil
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
		err = fmt.Errorf("traceroute NewTaskFromModel Unmarshal Opts err:%v", err)
		return
	}
	tRes.Opts = o
	return
}

// run is the entry point for the traceroute task
func (t task) Run(ctx context.Context) ([]byte, error) {
	if ctx.Err() != nil {
		return nil, fmt.Errorf("context done detected in Run:%v", ctx.Err())
	}

	var res Result

	// run the main traceroute func
	res, err := t.traceroute(ctx)
	if err != nil {
		return nil, err
	}

	// marshal the result
	resJson, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return resJson, nil
}

// socketAddr return the first non-loopback address as a 4 byte IP address. This address
// is used for sending packets out.
func (t task) socketAddr() (addr [4]byte, err error) {
	// get the a list of the system's addresses
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		loggerTraceroute.Error("error retreiving addresses", err)
		return
	}
	// look for an ipv4 address that will be used to send packets
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if len(ipnet.IP.To4()) == net.IPv4len {
				copy(addr[:], ipnet.IP.To4())
				loggerTraceroute.Info("socketAddr: ", addr)
				return addr, nil
			}
		}
	}
	return
}

// runHop takes the task context, does the actual hop operations and returns a completed hop with stats or an error
func (t task) runHop(ctx context.Context) (hop Hop, err error) {
	start := time.Now()

	// set the current hop TTL
	err = syscall.SetsockoptInt(t.SendSocket, 0x0, syscall.IP_TTL, t.TTL)
	if err != nil {
		hop.Error = err
		loggerTraceroute.Infof("hop: %+v", hop)
		return hop, err
	}

	// send a single null byte to the destination
	err = syscall.Sendto(t.SendSocket, []byte{0}, 0, &syscall.SockaddrInet4{Port: t.Port, Addr: t.Dest})
	if err != nil {
		hop.Error = err
		loggerTraceroute.Errorf("hop: %+v", hop)
		return hop, err
	}
	// receive the ICMP response
	bReceived, from, err := syscall.Recvfrom(t.ReceiveSocket, t.Packet, 0)
	//capture time
	elapsed := time.Since(start)
	if err != nil {
		hop.Error = err
		loggerTraceroute.Errorf("hop: %+v", hop)
		return hop, err
	}

	// grab current addr
	t.CurrentAddr = from.(*syscall.SockaddrInet4).Addr
	// reverse lookup
	currAddrStr := fmt.Sprintf("%d.%d.%d.%d", t.CurrentAddr[0], t.CurrentAddr[1], t.CurrentAddr[2], t.CurrentAddr[3])
	t.CurrentHost, err = net.LookupAddr(currAddrStr)
	if err != nil {
		loggerTraceroute.Warn("reverse lookup", err)
	} else {
		hop.Host = t.CurrentHost[0]
	}
	hop.Success = true
	hop.Address = t.CurrentAddr
	hop.BytesReceived = bReceived
	hop.ElapsedTime = elapsed
	hop.TTL = t.TTL
	loggerTraceroute.Infof("hop: %+v", hop)
	return hop, nil

}

// traceroute is the main function that performs the traceroute
func (t task) traceroute(ctx context.Context) (res Result, err error) {
	// set up the result
	res.DestinationAdress = t.Dest
	res.Hops = make([]Hop, 0)
	// initialize the function with options from the task
	//
	timeoutMs := (int64)(t.Timeout)
	maxTracerouteTimeout := 60 * time.Second // arbitrary timeout
	tv := syscall.NsecToTimeval(1000 * 1000 * timeoutMs)

	// get the socket address that packets will be sent from
	socketAddr, err := t.socketAddr()
	if err != nil {
		loggerTraceroute.Error("no valid ip found:", err)
		return
	}
	t.TTL = t.FirstHop
	t.Packet = make([]byte, t.Packetsize)
	retry := 0

	// create a context with a timeout for the entire traceroute operation
	ctx, cancel := context.WithTimeout(ctx, maxTracerouteTimeout)
	defer cancel()

	// set up ICMP receive socket
	t.ReceiveSocket, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		loggerTraceroute.Error("error creating socket: ", err)
		return res, err
	}
	defer syscall.Close(t.ReceiveSocket)

	// bind the receive socket
	err = syscall.Bind(t.ReceiveSocket, &syscall.SockaddrInet4{Port: t.Port, Addr: socketAddr})
	if err != nil {
		loggerTraceroute.Error("error binding socket", err)
		return res, err
	}

	// set the timeout for the socket
	err = syscall.SetsockoptTimeval(t.ReceiveSocket, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv)
	if err != nil {
		loggerTraceroute.Error("error setting timeout", err)
		return res, err
	}

	// set up the UDP send socket
	t.SendSocket, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		loggerTraceroute.Error("error creating socket: ", err)
		return res, err
	}
	defer syscall.Close(t.SendSocket)

	// start the main loop
	//
	for {
		if ctx.Err() != nil {
			return
		}

		// call hop operation and append the hops to the result
		hop, err := t.runHop(ctx)
		if err != nil {
			retry += 1
			// if retries exeeded
			if retry > t.Retries {
				retryErr := errors.New("retries exeeded")
				res.Hops = append(res.Hops, Hop{Success: false, TTL: t.TTL, Error: retryErr})
				t.TTL += 1
				retry = 0
			}
		}

		res.Hops = append(res.Hops, hop)
		t.TTL += 1
		retry = 0
		if t.TTL > t.MaxHops || t.CurrentAddr == t.Dest {
			loggerTraceroute.Infof("res: %+v", res)
			return res, nil
		}
	}
}
