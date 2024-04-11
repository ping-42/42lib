//go:build linux
// +build linux

package dns

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/ping-42/42lib/logger"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

// timeout is the read & write deadline, set on the dns connection
var timeout = time.Duration(20 * time.Second)

// dnsQueryTCP4 establishes connection to the specified name server
// over tcp and queries the specified host.
func (t task) dnsQueryTCP4(ctx context.Context, nameserver DnsNameServer) (Result, error) {
	// TODO: pass the context down the line
	var ret = Result{
		Proto: "tcp4",
	}
	var err error
	dnsLogger = dnsLogger.WithFields(log.Fields{
		"proto":      "tcp4",
		"host":       t.Host,
		"nameserver": nameserver,
	})

	hostname := strings.TrimPrefix(t.Host, "www.")

	co := nameserver.conn
	if co == nil {
		if t.GetDnsConn == nil {
			errMsg := "dns conn getter nil pointer"
			logger.LogError(errMsg, "Couldn't establish dns connection", dnsLogger)
			return Result{}, fmt.Errorf(errMsg)
		}

		co, err = t.GetDnsConn(nameserver.addr, "tcp", nameserver.port)
		if err != nil {
			logger.LogError(err.Error(), "Couldn't establish dns connection", dnsLogger)
			return Result{}, err
		}
		defer co.Close()
	}

	// Craft a DNS query message
	m := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Authoritative:     false,
			AuthenticatedData: false,
			CheckingDisabled:  false,
			RecursionDesired:  true,
			Opcode:            dns.OpcodeQuery,
			Rcode:             dns.RcodeSuccess,
		},
		Question: make([]dns.Question, 1),
	}

	// Add the host to the question
	qt := dns.TypeA
	qc := uint16(dns.ClassINET)
	m.Question[0] = dns.Question{Name: dns.Fqdn(hostname), Qtype: qt, Qclass: qc}
	m.Id = dns.Id()

	err = co.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		logger.LogError(err.Error(), "Dns connection set deadline error", dnsLogger)
		return Result{}, err
	}

	// Snapshot time before sending the msg
	then := time.Now()
	if err := co.WriteMsg(m); err != nil {
		logger.LogError(err.Error(), "Unable to send message to nameserver", dnsLogger)
		return Result{}, err
	}
	r, err := co.ReadMsg()
	if err != nil {
		logger.LogError(err.Error(), "Unable to read message from nameserver", dnsLogger)
		return Result{}, err
	}
	ret.QueryRtt = time.Since(then)

	if r.Id != m.Id {
		err := fmt.Errorf("DNS Query ID Missmatch")
		logger.LogError(err.Error(), "DNS Query ID Mismatch received from nameserver", dnsLogger)
		return Result{}, err
	}

	// Attach the resource record(s)
	for _, a := range r.Answer {
		// TODO: 4A
		switch answ := a.(type) {
		case *dns.A:
			ret.AnswerA = append(ret.AnswerA, answ)
		}
	}

	// Retrieve socket telemetry information
	if runtime.GOOS == "linux" {
		if t.GetSocketInfo == nil {
			err := fmt.Errorf("socket getter nil pointer")
			logger.LogError(err.Error(), "Failed to gather socket info", dnsLogger)
			return Result{}, err
		}

		socketInfo, err := t.GetSocketInfo(co.Conn.(*net.TCPConn))
		if err != nil {
			logger.LogError(err.Error(), "Unable to retrieve TCP socket information", dnsLogger)
			return Result{}, err
		}
		ret.SockRtt = time.Duration(socketInfo.Rtt)
	}

	return ret, err
}
