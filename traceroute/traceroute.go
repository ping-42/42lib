package traceroute

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/ping-42/42lib/db/models"
	"github.com/ping-42/42lib/logger"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/sys/unix"
)

// // these are typical default values for unix/linux traceroute operations
// const (
// 	DEFAULT_PORT        = 33434 // we target an unreachable port on the destination
// 	DEFAULT_MAX_HOPS    = 64
// 	DEFAULT_FIRST_HOP   = 1 // can be set to a higher value if we want to start hoping from from a certain router
// 	DEFAULT_TIMEOUT_MS  = 500
// 	DEFAULT_RETRIES     = 3
// 	DEFAULT_PACKET_SIZE = 52 // ipv4header (20) + udpheader (8) + payload (24)
// )

var (
	loggerTraceroute = logger.WithTestType("traceroute")
)

// NewTaskFromBytes used in sensor for building the task from the received bytes
func NewTaskFromBytes(msg []byte, sysUnix SysUnix) (t task, err error) {
	t.sysUnix = sysUnix
	// build the traceroute task from the received msg
	err = json.Unmarshal(msg, &t)
	if err != nil {
		err = fmt.Errorf("traceroute.NewTaskFromBytes Unmarshal err task:%v, %v", string(msg), err)
		return
	}

	return t, nil
}

// NewTaskFromModel used in scheduler for building the task from the db model task
func NewTaskFromModel(t models.Task, sysUnix SysUnix) (tRes task, err error) {

	tRes.Id = t.ID
	tRes.SensorId = t.SensorID
	tRes.Name = TaskName
	tRes.sysUnix = sysUnix

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

// // socketAddr return the first non-loopback address as a 4 byte IP address. This address
// // is used for sending packets out.
// func (t task) socketAddr() (addr [4]byte, err error) {
// 	// get the a list of the system's addresses
// 	addrs, err := net.InterfaceAddrs()
// 	if err != nil {
// 		loggerTraceroute.Error("error retreiving addresses", err)
// 		return
// 	}
// 	// look for an ipv4 address that will be used to send packets
// 	for _, a := range addrs {
// 		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
// 			if len(ipnet.IP.To4()) == net.IPv4len {
// 				copy(addr[:], ipnet.IP.To4())
// 				loggerTraceroute.Info("socketAddr: ", addr)
// 				return addr, nil
// 			}
// 		}
// 	}
// 	return
// }

// runHop takes the task context, does the actual hop operation and returns a completed hop with stats or an error
func (t *task) runHop() (hop Hop, err error) {
	hop = Hop{}
	// set the current hop TTL
	err = t.sysUnix.SetsockoptInt(t.SendSocket, 0x0, unix.IP_TTL, t.TTL)
	if err != nil {
		return hop, err
	}
	start := time.Now()

	for retries := 0; retries < t.Retries; retries++ {

		// send empty udp packet
		err = t.sysUnix.Sendto(t.SendSocket, []byte{0}, 0, &unix.SockaddrInet4{Port: t.Port, Addr: t.Dest})
		if err != nil {
			loggerTraceroute.Errorf("Failed to send packet on hop #%d: %v", t.TTL, err)
			continue //retry sending
		}

		// read the ICMP response into the buffer we created
		bReceived, from, err := t.sysUnix.Recvfrom(t.ReceiveSocket, t.Packet, 0)
		if err != nil {
			loggerTraceroute.Errorf("Failed to receive packet on hop #%d: %v", t.TTL, err)
			continue //retry receiving
		}

		// get the current address
		t.CurrentAddr = from.(*unix.SockaddrInet4).Addr
		addrStr := fmt.Sprintf("%d.%d.%d.%d", t.CurrentAddr[0], t.CurrentAddr[1], t.CurrentAddr[2], t.CurrentAddr[3])

		// parse the ICMP message
		parsedIcmpMessage, err := icmp.ParseMessage(1, t.Packet[20:])
		if err != nil {
			loggerTraceroute.Error("error parsing message", err)
		}

		// // do we need the header??
		// // parse the header
		// parsedHeader, err := icmp.ParseIPv4Header(t.Packet)
		// if err != nil {
		// 	loggerTraceroute.Error("error parsing header: ", err)
		// }
		// loggerTraceroute.Infof("parsedHeader: %+v", parsedHeader)

		// switch on the ICMP message type to determine hop result
		switch parsedIcmpMessage.Type {
		case ipv4.ICMPTypeTimeExceeded: // time exeeded. this means that the packet was dropped and the TTL will increment by 1 on next hop.
			loggerTraceroute.Infof("Time Exceeded received from %s", addrStr)
			t.CurrentHost, err = net.LookupAddr(addrStr)
			if err != nil {
				loggerTraceroute.Warn("reverse lookup failed: ", err)
			} else {
				hop.Host = t.CurrentHost[0]
			}
			hop.Success = true
			hop.Address = t.CurrentAddr
			hop.BytesReceived = bReceived
			hop.ElapsedTime = time.Since(start)
			hop.TTL = t.TTL
			loggerTraceroute.Infof("hop: %+v", hop)
			return hop, nil
		case ipv4.ICMPTypeDestinationUnreachable: // port unreachable. this means we reached the dest(yay) but the port is not available (cuz we send to a weird port).
			loggerTraceroute.Warn("Port unreachable")
			t.CurrentHost, err = net.LookupAddr(addrStr)
			if err != nil {
				loggerTraceroute.Warn("reverse lookup failed: ", err)
			} else {
				hop.Host = t.CurrentHost[0]
			}
			hop.Success = true
			hop.Address = t.CurrentAddr
			hop.BytesReceived = bReceived
			hop.ElapsedTime = time.Since(start)
			hop.TTL = t.TTL
			loggerTraceroute.Infof("hop: %+v", hop)
			return hop, nil
		case ipv4.ICMPTypeEchoReply: // we hit the destination address and port. This is possible but very unlikely.
			loggerTraceroute.Infof("Destination reached: %s", addrStr)
			t.CurrentHost, err = net.LookupAddr(addrStr)
			if err != nil {
				loggerTraceroute.Warn("reverse lookup failed: ", err)
			} else {
				hop.Host = t.CurrentHost[0]
			}
			hop.Success = true
			hop.Address = t.CurrentAddr
			hop.BytesReceived = bReceived
			hop.ElapsedTime = time.Since(start)
			hop.TTL = t.TTL
			loggerTraceroute.Infof("hop: %+v", hop)
			return hop, nil
		default:
			loggerTraceroute.Infof("received non-handled ICMP type: %d", parsedIcmpMessage.Type)
			continue
		}
	}
	hop.Success = false
	hop.ElapsedTime = time.Since(start)
	hop.TTL = t.TTL
	hop.Error = fmt.Errorf("max retries exceeded for hop")
	loggerTraceroute.Infof("hop: %+v", hop)
	return hop, hop.Error
}

// traceroute is the main function that performs the traceroute
func (t task) traceroute(ctx context.Context) (res Result, err error) {
	// set up the result
	res.DestinationAdress = t.Dest
	res.Hops = make([]Hop, 0)

	// initialize the function with options from the task
	maxTracerouteTimeout := 60 * time.Second // arbitrary timeout
	timeoutMs := (int64)(t.Timeout)
	timeValue := t.sysUnix.NsecToTimeval(1000 * 1000 * timeoutMs)
	t.TTL = t.FirstHop
	t.Packet = make([]byte, t.Packetsize) // create packet buffer that will store the ICMP response

	// create a context with a timeout for the entire traceroute operation
	ctx, cancel := context.WithTimeout(ctx, maxTracerouteTimeout)
	defer cancel()

	// set up raw socket for receiving ICMP replies
	t.ReceiveSocket, err = t.sysUnix.Socket(unix.AF_INET, unix.SOCK_RAW, unix.IPPROTO_ICMP)
	if err != nil {
		loggerTraceroute.Error("error creating socket: ", err)
		return res, err
	}
	defer t.sysUnix.Close(t.ReceiveSocket)

	// bind the receive socket to 0.0.0.0 to listen on all interfaces
	err = t.sysUnix.Bind(t.ReceiveSocket, &unix.SockaddrInet4{})
	if err != nil {
		loggerTraceroute.Error("error binding socket", err)
		return res, err
	}

	// set the timeout for the socket
	err = t.sysUnix.SetsockoptTimeval(t.ReceiveSocket, unix.SOL_SOCKET, unix.SO_RCVTIMEO, &timeValue)
	if err != nil {
		loggerTraceroute.Error("error setting timeout", err)
		return res, err
	}

	// set up datagram socket for sending UDP packets
	t.SendSocket, err = t.sysUnix.Socket(unix.AF_INET, unix.SOCK_DGRAM, unix.IPPROTO_UDP)
	if err != nil {
		loggerTraceroute.Error("error creating socket: ", err)
		return res, err
	}
	defer t.sysUnix.Close(t.SendSocket)

	// start the main loop
	for {
		if ctx.Err() != nil {
			return
		}

		// call hop operation and append the hops to the result
		hop, err := t.runHop()
		if err != nil {
			loggerTraceroute.Error("error running hop: ", err)
		}

		res.Hops = append(res.Hops, hop)
		t.TTL += 1
		if t.TTL > t.MaxHops || t.CurrentAddr == t.Dest {
			loggerTraceroute.Infof("res: %+v", res)
			return res, nil
		}
	}
}
