//go:build ignore
// +build ignore

/*
Until the tests are properly implemented and mocked,
ignore this from the pipeline and use for debugging only.
*/
package traceroute

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestRunHop(t *testing.T) {
// 	receivedMessage := []byte(`{"Id":"3b241101-e2bb-4255-8caf-4136c566a964","Name":"TRACEROUTE_TASK","SensorID":"3b241101-e2bb-4255-8caf-4136c566a964","Opts":{"Port":33434,"Dest":[8,8,8,8],"FirstHop":1,"MaxHops":30,"Timeout":5000,"PacketSize":52,"Retries":3}}`)

// 	task, err := NewTaskFromBytes(receivedMessage)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	hop, err := task.runHop(context.TODO())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	loggerTraceroute.Infof("hop: %+v\n", hop)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, hop)

// }

func TestTracerouteTask(t *testing.T) {

	// mock message with default win payload
	receivedMessage := []byte(`{"Id":"3b241101-e2bb-4255-8caf-4136c566a964","Name":"TRACEROUTE_TASK","SensorID":"3b241101-e2bb-4255-8caf-4136c566a964","Opts":{"Port":33434,"Dest":[8,8,8,8],"FirstHop":1,"MaxHops":64,"Timeout":500,"PacketSize":52,"Retries":3}}`)

	// Create an instance of the traceroute task with default options
	tracerouteTask, err := NewTaskFromBytes(receivedMessage)
	if err != nil {
		fmt.Println("eror creating task:", err)
	}
	fmt.Printf("New ICMP Task: %v\n\n", tracerouteTask)
	// Call the PingHost function with test options from message
	result, err := tracerouteTask.Run(context.TODO())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("unit test traceroute task results: %+v\n\n", result)

	assert.Nil(t, err)
	assert.NotNil(t, result)

}
