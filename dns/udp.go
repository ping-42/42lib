package dns

import (
	"context"
	"fmt"
	"strings"

	"github.com/miekg/dns"
	"github.com/ping-42/42lib/constants"
	"github.com/ping-42/42lib/logger"
	log "github.com/sirupsen/logrus"
)

const udpProto = "udp4"

func (t task) dnsQueryUDP4(ctx context.Context, nameserver DnsNameServer) (Result, error) {
	var ret = Result{
		Proto: constants.ProtoUDP,
	}

	dnsLogger = dnsLogger.WithFields(log.Fields{
		"proto":      udpProto,
		"host":       t.Host,
		"nameserver": nameserver,
	})

	hostname := strings.TrimPrefix(t.Host, "www.")

	if t.DnsUdpClient == nil {
		err := fmt.Errorf("nil dns udp client")
		logger.LogError(err.Error(), "Unexpected nil udp client", dnsLogger)
		return Result{}, err
	}

	// craft dns query
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
	qt := dns.TypeA
	qc := uint16(dns.ClassINET)
	m.Question[0] = dns.Question{Name: dns.Fqdn(hostname), Qtype: qt, Qclass: qc}
	m.Id = dns.Id()

	// TODO: ctx r, rtt, err := t.DnsUdpClient.ExchangeContext(ctx, m, nameserver.getAddrPort())
	r, rtt, err := t.DnsUdpClient.Exchange(m, nameserver.getAddrPort())
	if err != nil {
		logger.LogError(err.Error(), "UDP exchange error", dnsLogger)
		return Result{}, err
	}

	ret.QueryRtt = rtt
	ret.RespSize = int64(r.Len())
	for _, n := range r.Ns {
		dnsLogger.Infof("used NS: %v", n.String()) // TODO: this should go in the results
	}

	if len(r.Answer) < 1 {
		// TODO: indicate this as valid result - host is down
		logger.Logger.Error("did not receive any A record answers", dnsLogger)
		return ret, err
	}

	for _, a := range r.Answer {
		switch answ := a.(type) {
		case *dns.A:
			ret.AnswerA = append(ret.AnswerA, answ)
		default:
			dnsLogger.Warnf("Got non-A unhandled answer: %v", answ)
		}
	}

	return ret, nil
}
