package sensor

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"github.com/ping-42/42lib/wss"
)

// TaskRunner
// base task interface
type TaskRunner interface {
	GetId() uuid.UUID
	GetName() TaskName
	GetSensorId() uuid.UUID

	// *idea* we may wrap it as (context.WithValue) and implement a struct to do tracing, similar to xray
	Run(ctx context.Context) ([]byte, error) //TaskResult
}

type Task struct {
	Id       uuid.UUID
	Name     TaskName
	SensorId uuid.UUID
}

type TaskName string

// TResult the implementation of TaskResult interface
type TResult struct {
	// MessageGeneralType is required for all messages sent via WSS
	wss.MessageGeneralType

	TaskId   uuid.UUID
	TaskName TaskName
	Result   []byte
	Error    string
}

func (t TResult) SendToServer(ctx context.Context, wsConn net.Conn) (err error) {

	// Check if the context is done
	if ctx.Err() != nil {
		err = fmt.Errorf("context done detected in SendToServer:%v", ctx.Err())
		return
	}

	r, err := json.Marshal(t)
	if err != nil {
		err = fmt.Errorf("failed Marshal the result: %v", t)
		return
	}

	// Send results to the server
	err = wsutil.WriteClientText(wsConn, r)
	if err != nil {
		err = fmt.Errorf("failed to send task response to server: %v", err)
		return
	}
	return
}
