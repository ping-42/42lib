package models

import (
	"net"
	"time"
)

type TsDnsResult struct {
	TsSensorTaskBase
	QueryRtt    int64
	SocketRtt   int64
	RespSize    int64
	Proto       int64
	IPAddresses []net.IP `gorm:"type:inet[]"`
}

type TsHttpResult struct {
	TsSensorTaskBase
	ResponseCode     uint8
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
