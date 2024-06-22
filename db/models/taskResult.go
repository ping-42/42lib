package models

import (
	"net"
	"time"
)

type TsDnsResult struct {
	TsSensorTaskBase
	QueryRtt  int64
	SocketRtt int64
	RespSize  int64
	Proto     int8
}

type TsDnsResultAnswer struct {
	TsSensorTaskBase
	HdrName     string
	HdrRrtype   uint16
	HdrClass    uint16
	HdrTtl      uint32
	HdrRdlength uint16
	A           net.IP `gorm:"type:inet"`
}

func (TsDnsResultAnswer) TableName() string {
	return "ts_dns_results_answer"
}

type TsHttpResult struct {
	TsSensorTaskBase
	ResponseCode     int
	DNSLookup        time.Duration
	TCPConnection    time.Duration
	TLSHandshake     time.Duration
	ServerProcessing time.Duration
	NameLookup       time.Duration
	Connect          time.Duration
	Pretransfer      time.Duration
	StartTransfer    time.Duration
	//
	ResponseBody    string
	ResponseHeaders []byte `gorm:"type:jsonb"`
}

type TsIcmpResult struct {
	TsSensorTaskBase
	IPAddr          net.IP `gorm:"type:inet"`
	PacketsSent     int
	PacketsReceived int
	BytesWritten    int
	BytesRead       int
	TotalRTT        time.Duration
	MinRTT          time.Duration
	MaxRTT          time.Duration
	AverageRTT      time.Duration
	Loss            float64
	FailureMessages string
}

type TsTracerouteResult struct {
	TsSensorTaskBase
	DestinationAdress net.IP `gorm:"type:inet"`
}

type TsTracerouteResultHop struct {
	TsSensorTaskBase
	Success       bool
	Address       net.IP `gorm:"type:inet"`
	Host          string
	BytesReceived int
	ElapsedTime   time.Duration
	TTL           int
	Error         string
}

func (TsTracerouteResultHop) TableName() string {
	return "ts_traceroute_results_hop"
}
