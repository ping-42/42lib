package sensor

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/gobwas/ws/wsutil"
	"github.com/ping-42/42lib/wss"
)

type HostTelemetry struct {
	// MessageGeneralType is required for all messages sent via WSS
	wss.MessageGeneralType

	Cpu        Cpu
	Memory     Memory
	GoRoutines int
	Network    []Network
}

type Cpu struct {
	ModelName string
	Cores     uint16
	CpuUsage  float64
}

type Memory struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

type Network struct {
	Name      string
	BytesSent uint64
	BytesRecv uint64
}

func (ht HostTelemetry) SendToServer(ctx context.Context, wsConn net.Conn) (err error) {
	// Check if the context is done
	if ctx.Err() != nil {
		err = fmt.Errorf("context done detected in SendToServer: %v", ctx.Err())
		return
	}

	r, err := json.Marshal(ht)
	if err != nil {
		err = fmt.Errorf("failed Marshal the hostTelemetry: %v", ht)
		return
	}

	// Send results to the server
	err = wsutil.WriteClientText(wsConn, r)
	if err != nil {
		err = fmt.Errorf("failed to send telemetry response to server: %v", err)
		return
	}
	return
}
