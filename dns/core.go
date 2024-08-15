//go:build linux
// +build linux

package dns

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"time"

	"github.com/ping-42/42lib/db/models"
	"github.com/ping-42/42lib/helpers"
	"github.com/ping-42/42lib/logger"
	"golang.org/x/sys/unix"

	"github.com/docker/docker/libnetwork/resolvconf"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

/*
	Notes & TODOs
	- how do we know what is the DNS that the host uses and get info about it? docker run `--network=host`
	- cross-check default resolver with a trustworthy one, e.g. 1.1.1.1?
	- check for a timeout in the opts from the server message
	- if tcp is requested but fails - fallback to udp
	- add NS info to the result
*/

// defaults
var (
	defaultTcpSocketTimeout = 5 * time.Second
	nsList                  = []DnsNameServer{
		{
			addr:    "1.1.1.1",
			port:    53,
			asnOrg:  "Cloudflare Inc",
			usesTcp: true,
		},
		{
			addr:    "8.8.8.8",
			port:    53,
			asnOrg:  "Google Inc",
			usesTcp: true,
		},
		{
			addr:    "151.80.222.79",
			port:    53,
			asnOrg:  "OVH SAS",
			usesTcp: true,
		},
	}
	defaultOpts = Opts{
		GetSocketInfo: getTcpSocketInfo,
		GetDnsConn:    getDnsTcpConn,
		Proto:         "udp",
		DnsUdpClient:  getDnsUdpClient(),
	}
	dnsLogger = logger.WithTestType("dns")
)

// Run is the entry point for the dns.task
// It iterates the default NS list and returns the first successful resolve
func (t task) Run(ctx context.Context) ([]byte, error) {

	if ctx.Err() != nil {
		return nil, fmt.Errorf("context done detected in Run:%v", ctx.Err())
	}

	//  use default predefined until resolved
	defaultNsList, err := getDefaultNs()
	if err != nil {
		dnsLogger.Warnf("couldn't get any nameservers in resolv.conf: %v", err)
		defaultNsList = nsList
	}

	// TODO: use goroutines here?
	for _, ns := range defaultNsList {
		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(defaultTcpSocketTimeout))
		defer cancel()
		res, err := t.runWithProto(ctx, ns)
		if err != nil {
			// TODO: failed to connect to the NS (or another error)
			// should this be included in the results as info?
			dnsLogger.WithError(err)
		} else {
			res, err := json.Marshal(res)
			if err != nil {
				return nil, err
			}

			return res, nil
		}
	}

	return nil, fmt.Errorf("unknown protocol: %v", t.Proto)
}

// getDefaultNs gets the default nameserver set on the host
func getDefaultNs() ([]DnsNameServer, error) {
	confFile, err := resolvconf.Get()
	if err != nil {
		dnsLogger.Errorf("Couldn't read resolv.conf: %v", err)
		return []DnsNameServer{}, err
	}

	// take the resolv.conf matches without localhost defined ones;
	// TODO: take care of systemd managed DNS, as this will replace 127.0.0.53
	//  with a default one if only local addresses are found
	f, err := resolvconf.FilterResolvDNS(confFile.Content, false)
	if err != nil {
		dnsLogger.Errorf("Couldn't clean resolv.conf: %v", err)
		return []DnsNameServer{}, err
	}

	ns := resolvconf.GetNameservers(f.Content, resolvconf.IPv4)
	if len(ns) < 1 {
		msg := "no default nameservers"
		err := fmt.Errorf("%v", msg)
		dnsLogger.Error(msg)
		return []DnsNameServer{}, err
	}

	ret := make([]DnsNameServer, 0)
	for _, ip := range ns {
		// TODO: "168.63.129.16" means we are in Azure context, might as well ban the sensor?
		dnsLogger.Infof("Parsed nameserver from resolv.conf: %v", ip)
		ret = append(ret, DnsNameServer{
			addr: ip,
			port: 53,
		})
	}

	return ret, nil
}

// runWithProto calls the appropriate function depending on the requested protocol
func (t task) runWithProto(ctx context.Context, ns DnsNameServer) (res Result, err error) {
	switch t.Proto {
	case "udp":
		res, err = t.dnsQueryUDP4(ctx, ns)
		if err != nil {
			return
		}
	case "tcp":
		res, err = t.dnsQueryTCP4(ctx, ns)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("unknown protocol requested: %v", t.Proto)
		return
	}
	return
}

func NewEmptyTask() (t task) {
	// set default opts
	t.Opts.DnsUdpClient = defaultOpts.DnsUdpClient
	t.Opts.GetSocketInfo = defaultOpts.GetSocketInfo
	t.Opts.GetDnsConn = defaultOpts.GetDnsConn
	if t.Opts.Proto == "" {
		t.Opts.Proto = defaultOpts.Proto
	}
	return
}

// NewTaskFromBytes used in sensor for building the task from the received bytes
func NewTaskFromBytes(msg []byte) (t task, err error) {

	// build the dns task from the received msg
	err = json.Unmarshal(msg, &t)
	if err != nil {
		err = fmt.Errorf("dns.NewTask Unmarshal err task:%v, %v", string(msg), err)
		return
	}

	// set default opts
	t.Opts.DnsUdpClient = defaultOpts.DnsUdpClient
	t.Opts.GetSocketInfo = defaultOpts.GetSocketInfo
	t.Opts.GetDnsConn = defaultOpts.GetDnsConn
	t.Opts.Host, err = helpers.ExtractDomainFromUrl(t.Opts.Host)
	if err != nil {
		return t, err
	}
	if t.Opts.Proto == "" {
		t.Opts.Proto = defaultOpts.Proto
	}

	return t, nil
}

// NewTaskFromModel used in scheduler for building the task from the db model task
func NewTaskFromModel(t models.Task) (tRes task, err error) {

	var o = Opts{}
	err = json.Unmarshal(t.Opts, &o)
	if err != nil {
		err = fmt.Errorf("dns NewTaskFromModel Unmarshal Opts err:%v", err)
		return
	}

	tRes.Id = t.ID
	tRes.SensorId = t.SensorID
	tRes.Name = TaskName
	tRes.Opts = o
	return
}

// getSocketRtt Retrieve socket RTT on Linux (with some caveats)
// GetsockoptTCPInfo() does not seem to really return all options outlined in
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/tcp.h#L214
// We only get some of them
func getTcpSocketInfo(conn *net.TCPConn) (*unix.TCPInfo, error) {
	if runtime.GOOS != "linux" {
		return nil, errors.New("unsupported OS - we need Linux")
	}

	raw, err := conn.SyscallConn()
	if err != nil {
		dnsLogger.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("SyscallConn() error")
		return nil, err
	}

	var info *unix.TCPInfo
	ctrlErr := raw.Control(func(fd uintptr) {
		info, err = unix.GetsockoptTCPInfo(int(fd), unix.IPPROTO_TCP, unix.TCP_INFO)
	})

	// Figure out if anything failed
	switch {
	case ctrlErr != nil:
		return nil, ctrlErr
	case err != nil:
		return nil, err
	}
	dnsLogger.WithFields(log.Fields{
		"socketRtt":    time.Duration(info.Rtt).Nanoseconds(),
		"socketRttVar": time.Duration(info.Rttvar).Nanoseconds(),
		"socketLost":   info.Lost,
	}).Debug("Time duration of socket detected")

	return info, nil
}

// getDnsTcpConn returns a new dns connection pointer
func getDnsTcpConn(addr, proto string, port int) (*dns.Conn, error) {
	co := new(dns.Conn)
	var err error
	co.Conn, err = net.DialTimeout(proto, net.JoinHostPort(addr, strconv.Itoa(port)), defaultTcpSocketTimeout)
	return co, err
}

// getDnsUdpClient returns a dns client set to use UDP
func getDnsUdpClient() *dns.Client {
	dnsClient := new(dns.Client)
	dnsClient.Net = "udp4"
	return dnsClient
}
