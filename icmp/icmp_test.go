//go:build ignore
// +build ignore

/*
Until the tests are properly implemented and mocked,
ignore this from the pipeline and use for debugging only.
*/
package icmp

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunWithCancelledContext(t *testing.T) {
	// mock message with default win payload
	receivedMessage := []byte(`{"Id":"123","Name":"ICMP_TEST","IcmpOpts":{"TargetIPs":["8.8.8.8", "2001:4860:4860::8844"],"Count":3}}`)

	// create a context and cancel it
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// create new task
	icmpTask := NewTask()
	loggerIcmp.Printf("New ICMP Task: %v\n\n", icmpTask)

	// call Run with the cancelled context
	_, err := icmpTask.Run(ctx, receivedMessage)
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestCreateConnV4(t *testing.T) {
	// create a new icmpV4 connection
	conn, err := createConnV4()
	if err != nil {
		loggerIcmp.Error(err)
	} else {
		loggerIcmp.Infof("connection: %+v\n", conn)
	}
	assert.Nil(t, err)
	assert.NotNil(t, conn)
}

func TestCreateConnV6(t *testing.T) {
	// create a new icmpV4 connection
	conn, err := createConnV6()
	if err != nil {
		loggerIcmp.Error(err)
	} else {
		loggerIcmp.Info("connection: ", conn)
	}
	assert.Nil(t, err)
	assert.NotNil(t, conn)
}

func TestRunICMP4(t *testing.T) {
	// parse an ipv4 to pass to runICMP
	ip := net.ParseIP("8.8.8.8")
	if ip == nil {
		loggerIcmp.Error("no valid ip")
		return
	}
	// create a new icmpV4 connection
	conn, err := createConnV4()
	if err != nil {
		loggerIcmp.Info(err)
		return
	}
	// get the message type based on ip
	msgType := getICMPType(ip)
	// create a new task and run the ping
	icmpTask := NewTask()
	pingResult, err := icmpTask.runICMP(context.TODO(), conn, msgType, ip)
	if err != nil {
		loggerIcmp.Error(err)
	}
	loggerIcmp.Infof("ping result: %v", pingResult)

	assert.Nil(t, err)
	assert.NotNil(t, pingResult)
}

func TestRunICMP6(t *testing.T) {
	// parse an ipv6 to pass to runICMP
	ip := net.ParseIP("2001:4860:4860::8844")
	if ip == nil {
		loggerIcmp.Error("no valid ip")
		return
	}
	// create a new icmpV6 connection
	conn, err := createConnV6()
	if err != nil {
		loggerIcmp.Info(err)
		return
	}
	// get the message type based on ip
	msgType := getICMPType(ip)
	// create a new task and run the ping
	icmpTask := NewTask()
	pingResult, err := icmpTask.runICMP(context.TODO(), conn, msgType, ip)
	if err != nil {
		loggerIcmp.Error(err)
	}
	loggerIcmp.Infof("ping result: %v", pingResult)

	assert.Nil(t, err)
	assert.NotNil(t, pingResult)
}

func TestIcmpTask(t *testing.T) {

	// mock message with default win payload
	receivedMessage := []byte(`{"Id":"123","Name":"ICMP_TEST","IcmpOpts":{"TargetIPs":["8.8.8.8", "2001:4860:4860::8844"],"Count":3}}`)

	// Create an instance of the ICMP task with default options
	icmpTask := NewTask()

	// Call the PingHost function with test options from message
	result, err := icmpTask.Run(context.TODO(), receivedMessage)
	fmt.Printf("New ICMP Task: %v\n\n", icmpTask)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("unit test icmp task results: %+v\n\n", result)

	assert.Nil(t, err)
	assert.NotNil(t, result)

}
