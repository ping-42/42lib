/*
Until the tests are properly implemented and mocked,
ignore this from the pipeline and use for debugging only.
*/
package traceroute

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ping-42/42lib/db/models"
	"github.com/ping-42/42lib/sensor"
	testingkit "github.com/ping-42/42lib/testingkit"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/unix"
)

// this can be tested on root vscode
func TestTracerouteTaskFromBytes(t *testing.T) {

	if os.Getuid() != 0 {
		t.Skip("Skipping test when running as non-root")
	}
	// mock message with default win payload
	receivedMessage := []byte(`{"Id":"3b241101-e2bb-4255-8caf-4136c566a964","Name":"TRACEROUTE_TASK","SensorID":"3b241101-e2bb-4255-8caf-4136c566a964","Opts":{"Port":33434,"Dest":"8.8.8.8","FirstHop":1,"MaxHops":64,"Timeout":500,"PacketSize":52,"Retries":3}}`)

	// Create an instance of the traceroute task with default options
	tracerouteTask, err := NewTaskFromBytes(receivedMessage)
	assert.Nil(t, err)
	assert.NotNil(t, tracerouteTask)

	t.Logf("New TRACEROUTE Task: %v\n\n", tracerouteTask)

	// Call the traceroute Run with test options from message
	result, err := tracerouteTask.Run(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, result)

	t.Logf("unit test traceroute task results: %+v\n\n", string(result))
}

// this can be tested on root vscode
func TestTracerouteTaskFromModel(t *testing.T) {

	if os.Getuid() != 0 {
		t.Skip("Skipping test when running as non-root")
	}

	createdAt, err := time.Parse(time.RFC3339Nano, "2024-09-05T15:04:05.999999+00:00")
	if err != nil {
		fmt.Println(err)
		return
	}

	uuid := uuid.New()

	// mock message with default win payload
	lvTaskType := models.LvTaskType{
		ID:   3,
		Type: "TRACEROUTE",
	}

	lvTaskStatus := models.LvTaskStatus{
		ID:     3,
		Status: "pending",
	}

	org := models.Organization{
		ID:   uuid,
		Name: "ping",
	}

	sensor := models.Sensor{
		ID:             uuid,
		OrganizationID: uuid,
		Organization:   org,
		Name:           "pingSensor",
		Location:       "TLV",
		Secret:         "123",
		IsActive:       true,
		CreatedAt:      createdAt,
	}

	subscription := models.Subscription{}

	modelTask := models.Task{
		ID:             uuid,
		TaskTypeID:     3,
		TaskType:       lvTaskType,
		TaskStatusID:   3,
		TaskStatus:     lvTaskStatus,
		SensorID:       uuid,
		Sensor:         sensor,
		SubscriptionID: 3,
		Subscription:   subscription,
		CreatedAt:      createdAt,
		Opts:           []byte(`{"Port":33434,"Dest":"8.8.8.8","FirstHop":1,"MaxHops":64,"Timeout":500,"PacketSize":52,"Retries":3}`),
	}

	// Create an instance of the traceroute task with default options
	tracerouteTask, err := NewTaskFromModel(modelTask)
	assert.Nil(t, err)
	assert.NotNil(t, tracerouteTask)

	t.Logf("New TRACEROUTE Task: %v\n\n", tracerouteTask)

	// Call the traceoute Run with test options from message
	result, err := tracerouteTask.Run(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// uses SysUnix interface's mocked methods
func TestTracerouteTaskFromBytesMocked(t *testing.T) {
	mockSysUnix := &testingkit.MockedSysUnix{
		SocketFunc: func(domain, typ, proto int) (int, error) {
			return 1, nil
		},
		RecvfromFunc: func(fd int, p []byte, flags int) (int, unix.Sockaddr, error) {
			// // mock getting package
			return len(p), &unix.SockaddrInet4{}, nil
		},
	}

	// mock message with default win payload
	receivedMessage := []byte(`{"Id":"3b241101-e2bb-4255-8caf-4136c566a964","Name":"TRACEROUTE_TASK","SensorID":"3b241101-e2bb-4255-8caf-4136c566a964","Opts":{"Port":33434,"Dest":"8.8.8.8","FirstHop":1,"MaxHops":64,"Timeout":500,"PacketSize":52,"Retries":3}}`)

	// Create an instance of the traceroute task with default options
	tracerouteTask, err := NewTaskFromBytes(receivedMessage)
	assert.Nil(t, err)
	assert.NotNil(t, tracerouteTask)

	tracerouteTask.SysUnix = mockSysUnix
	t.Logf("New TRACEROUTE Task: %v\n\n", tracerouteTask)

	// Call the traceoute Run with test options from message
	result, err := tracerouteTask.Run(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, result)

	t.Logf("unit test traceroute task results: %+v\n\n", string(result))

}

// uses SysUnix interface's mocked methods
func TestTracerouteTaskFromModelMocked(t *testing.T) {
	mockSysUnix := &testingkit.MockedSysUnix{
		SocketFunc: func(domain, typ, proto int) (int, error) {
			return 1, nil
		},
		RecvfromFunc: func(fd int, p []byte, flags int) (int, unix.Sockaddr, error) {
			// // mock getting package
			return len(p), &unix.SockaddrInet4{}, nil
		},
	}
	createdAt, err := time.Parse(time.RFC3339Nano, "2024-09-05T15:04:05.999999+00:00")
	assert.Nil(t, err)

	uuid := uuid.New()

	// mock message with default win payload
	lvTaskType := models.LvTaskType{
		ID:   3,
		Type: "TRACEROUTE",
	}

	lvTaskStatus := models.LvTaskStatus{
		ID:     3,
		Status: "pending",
	}

	org := models.Organization{
		ID:   uuid,
		Name: "ping",
	}

	sensor := models.Sensor{
		ID:             uuid,
		OrganizationID: uuid,
		Organization:   org,
		Name:           "pingSensor",
		Location:       "TLV",
		Secret:         "123",
		IsActive:       true,
		CreatedAt:      createdAt,
	}

	subscription := models.Subscription{}

	modelTask := models.Task{
		ID:             uuid,
		TaskTypeID:     3,
		TaskType:       lvTaskType,
		TaskStatusID:   3,
		TaskStatus:     lvTaskStatus,
		SensorID:       uuid,
		Sensor:         sensor,
		SubscriptionID: 3,
		Subscription:   subscription,
		CreatedAt:      createdAt,
		Opts:           []byte(`{"Port":33434,"Dest":"8.8.8.8","FirstHop":1,"MaxHops":64,"Timeout":500,"PacketSize":52,"Retries":3}`),
	}

	// Create an instance of the traceroute task with default options
	tracerouteTask, err := NewTaskFromModel(modelTask)
	assert.Nil(t, err)

	tracerouteTask.SysUnix = mockSysUnix
	t.Logf("New TRACEROUTE Task: %v\n\n", tracerouteTask)

	// Call the traceoute Run with test options from message
	result, err := tracerouteTask.Run(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, result)

	t.Logf("unit test traceroute task results: %+v\n\n", string(result))

	assert.Nil(t, err)
	assert.NotNil(t, result)
}

/*
unmarshal fail from docker log
time="2024-12-04T21:13:22Z" level=info msg="received task:{\"Id\":\"b08206cf-ba52-419d-9c1b-8391eb8a2361\",\"Name\":\"TRACEROUTE_TASK\",\"SensorId\":\"98952434-b6fe-465e-8d9b-8bca56e99873\",\"Opts\":{\"Port\":33434,\"Dest\":\"8.8.8.8\",\"CurrentAddr\":\"\",\"CurrentHost\":null,\"ReceiveSocket\":0,\"SendSocket\":0,\"FirstHop\":1,\"MaxHops\":64,\"Timeout\":500,\"PacketSize\":52,\"Packet\":null,\"Ttl\":0,\"Retries\":3},\"SysUnix\":{}}" testType=sensor
time="2024-12-04T21:13:22Z" level=error msg="error in factoryTask()" error="traceroute.NewTaskFromBytes Unmarshal err task:{\"Id\":\"b08206cf-ba52-419d-9c1b-8391eb8a2361\",\"Name\":\"TRACEROUTE_TASK\",\"SensorId\":\"98952434-b6fe-465e-8d9b-8bca56e99873\",\"Opts\":{\"Port\":33434,\"Dest\":\"8.8.8.8\",\"CurrentAddr\":\"\",\"CurrentHost\":null,\"ReceiveSocket\":0,\"SendSocket\":0,\"FirstHop\":1,\"MaxHops\":64,\"Timeout\":500,\"PacketSize\":52,\"Packet\":null,\"Ttl\":0,\"Retries\":3},\"SysUnix\":{}}, json: cannot unmarshal object into Go struct field task.SysUnix of type traceroute.SysUnix" testType=sensor
time="2024-12-04T21:14:22Z" level=info msg="received task:{\"Id\":\"3a923e95-6991-4d94-817a-ecffd6bed1bb\",\"Name\":\"TRACEROUTE_TASK\",\"SensorId\":\"98952434-b6fe-465e-8d9b-8bca56e99873\",\"Opts\":{\"Port\":33434,\"Dest\":\"8.8.8.8\",\"CurrentAddr\":\"\",\"CurrentHost\":null,\"ReceiveSocket\":0,\"SendSocket\":0,\"FirstHop\":1,\"MaxHops\":64,\"Timeout\":500,\"PacketSize\":52,\"Packet\":null,\"Ttl\":0,\"Retries\":3},\"SysUnix\":{}}" testType=sensor
*/

// func TestRunHop(t *testing.T) {
// 	receivedMessage := []byte(`{"Id":"3b241101-e2bb-4255-8caf-4136c566a964","Name":"TRACEROUTE_TASK","SensorID":"3b241101-e2bb-4255-8caf-4136c566a964","Opts":{"Port":33434,"Dest":"8.8.8.8","FirstHop":1,"MaxHops":64,"Timeout":500,"PacketSize":52,"Retries":3}}`)

// 	task, err := NewTaskFromBytes(receivedMessage)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	hop, err := task.runHop()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	loggerTraceroute.Infof("hop: %+v\n", hop)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, hop)

// }

func TestTracerouteResultSubmission(t *testing.T) {
	taskResult := []byte(`
	{"MessageGeneralType" : 0,
"TaskId" : "b62f5016-40b2-4b9d-a8e3-46a57af7a9e6",
"TaskName" : "TRACEROUTE_TASK",
"Result" : "eyJEZXN0aW5hdGlvbkFkcmVzcyI6IjEuMS4xLjEiLCJIb3BzIjpbeyJTdWNjZXNzIjp0cnVlLCJBZGRyZXNzIjoiMTAuMjEwLjEwLjEiLCJIb3N0IjoidW5pZmkubG9jYWxkb21haW4uIiwiQnl0ZXNSZWNlaXZlZCI6NTIsIkVsYXBzZWRUaW1lIjoxNTQ2MDc5LCJUVEwiOjIsIkVycm9yIjpudWxsfSx7IlN1Y2Nlc3MiOnRydWUsIkFkZHJlc3MiOiIxOTIuMTY4LjEuMjU0IiwiSG9zdCI6IiIsIkJ5dGVzUmVjZWl2ZWQiOjUyLCJFbGFwc2VkVGltZSI6MTY3OTMwMiwiVFRMIjozLCJFcnJvciI6bnVsbH0seyJTdWNjZXNzIjpmYWxzZSwiQWRkcmVzcyI6IiIsIkhvc3QiOiIiLCJCeXRlc1JlY2VpdmVkIjowLCJFbGFwc2VkVGltZSI6MTU1MDQ5Mzk2MSwiVFRMIjo0LCJFcnJvciI6e319LHsiU3VjY2VzcyI6ZmFsc2UsIkFkZHJlc3MiOiIiLCJIb3N0IjoiIiwiQnl0ZXNSZWNlaXZlZCI6MCwiRWxhcHNlZFRpbWUiOjE1MzU5MTQ4OTIsIlRUTCI6NSwiRXJyb3IiOnt9fSx7IlN1Y2Nlc3MiOnRydWUsIkFkZHJlc3MiOiIxMi4yNDIuMTE2LjE3IiwiSG9zdCI6IiIsIkJ5dGVzUmVjZWl2ZWQiOjUyLCJFbGFwc2VkVGltZSI6NDEzODY0NCwiVFRMIjo2LCJFcnJvciI6bnVsbH0seyJTdWNjZXNzIjp0cnVlLCJBZGRyZXNzIjoiMzIuMTQyLjEzMy4xMDYiLCJIb3N0IjoiIiwiQnl0ZXNSZWNlaXZlZCI6NTIsIkVsYXBzZWRUaW1lIjo1NDA4Nzk3LCJUVEwiOjcsIkVycm9yIjpudWxsfSx7IlN1Y2Nlc3MiOnRydWUsIkFkZHJlc3MiOiIxMDguMTYyLjIxMS4yMjgiLCJIb3N0IjoiIiwiQnl0ZXNSZWNlaXZlZCI6NTIsIkVsYXBzZWRUaW1lIjo2NDYxNzA5OSwiVFRMIjo4LCJFcnJvciI6bnVsbH0seyJTdWNjZXNzIjpmYWxzZSwiQWRkcmVzcyI6IiIsIkhvc3QiOiIiLCJCeXRlc1JlY2VpdmVkIjowLCJFbGFwc2VkVGltZSI6MTUyNTU5NTA5NywiVFRMIjo5LCJFcnJvciI6e319LHsiU3VjY2VzcyI6ZmFsc2UsIkFkZHJlc3MiOiIiLCJIb3N0IjoiIiwiQnl0ZXNSZWNlaXZlZCI6MCwiRWxhcHNlZFRpbWUiOjE1MzU5MjE4NTgsIlRUTCI6MTAsIkVycm9yIjp7fX0seyJTdWNjZXNzIjpmYWxzZSwiQWRkcmVzcyI6IiIsIkhvc3QiOiIiLCJCeXRlc1JlY2VpdmVkIjowLCJFbGFwc2VkVGltZSI6MTUzNTk0MDQxMiwiVFRMIjoxMSwiRXJyb3IiOnt9fSx7IlN1Y2Nlc3MiOmZhbHNlLCJBZGRyZXNzIjoiIiwiSG9zdCI6IiIsIkJ5dGVzUmVjZWl2ZWQiOjAsIkVsYXBzZWRUaW1lIjoxNTM1OTQyMTkxLCJUVEwiOjEyLCJFcnJvciI6e319LHsiU3VjY2VzcyI6ZmFsc2UsIkFkZHJlc3MiOiIiLCJIb3N0IjoiIiwiQnl0ZXNSZWNlaXZlZCI6MCwiRWxhcHNlZFRpbWUiOjE1MzU5MjkwOTcsIlRUTCI6MTMsIkVycm9yIjp7fX0seyJTdWNjZXNzIjpmYWxzZSwiQWRkcmVzcyI6IiIsIkhvc3QiOiIiLCJCeXRlc1JlY2VpdmVkIjowLCJFbGFwc2VkVGltZSI6MTUzNTg4NTIyOCwiVFRMIjoxNCwiRXJyb3IiOnt9fSx7IlN1Y2Nlc3MiOmZhbHNlLCJBZGRyZXNzIjoiIiwiSG9zdCI6IiIsIkJ5dGVzUmVjZWl2ZWQiOjAsIkVsYXBzZWRUaW1lIjoxNTM1OTQxNzQ2LCJUVEwiOjE1LCJFcnJvciI6e319LHsiU3VjY2VzcyI6ZmFsc2UsIkFkZHJlc3MiOiIiLCJIb3N0IjoiIiwiQnl0ZXNSZWNlaXZlZCI6MCwiRWxhcHNlZFRpbWUiOjE1MzU5NDQ2OTEsIlRUTCI6MTYsIkVycm9yIjp7fX0seyJTdWNjZXNzIjpmYWxzZSwiQWRkcmVzcyI6IiIsIkhvc3QiOiIiLCJCeXRlc1JlY2VpdmVkIjowLCJFbGFwc2VkVGltZSI6MTUzNTkyOTAyNSwiVFRMIjoxNywiRXJyb3IiOnt9fSx7IlN1Y2Nlc3MiOmZhbHNlLCJBZGRyZXNzIjoiIiwiSG9zdCI6IiIsIkJ5dGVzUmVjZWl2ZWQiOjAsIkVsYXBzZWRUaW1lIjoxNTM1OTkxMDA1LCJUVEwiOjE4LCJFcnJvciI6e319LHsiU3VjY2VzcyI6ZmFsc2UsIkFkZHJlc3MiOiIiLCJIb3N0IjoiIiwiQnl0ZXNSZWNlaXZlZCI6MCwiRWxhcHNlZFRpbWUiOjE1MzU4OTEzOTUsIlRUTCI6MTksIkVycm9yIjp7fX0seyJTdWNjZXNzIjpmYWxzZSwiQWRkcmVzcyI6IiIsIkhvc3QiOiIiLCJCeXRlc1JlY2VpdmVkIjowLCJFbGFwc2VkVGltZSI6MTUzNTkzMDY1MywiVFRMIjoyMCwiRXJyb3IiOnt9fSx7IlN1Y2Nlc3MiOmZhbHNlLCJBZGRyZXNzIjoiIiwiSG9zdCI6IiIsIkJ5dGVzUmVjZWl2ZWQiOjAsIkVsYXBzZWRUaW1lIjoxNTM1OTQwMjI3LCJUVEwiOjIxLCJFcnJvciI6e319LHsiU3VjY2VzcyI6ZmFsc2UsIkFkZHJlc3MiOiIiLCJIb3N0IjoiIiwiQnl0ZXNSZWNlaXZlZCI6MCwiRWxhcHNlZFRpbWUiOjE1MzU5MTUyODIsIlRUTCI6MjIsIkVycm9yIjp7fX0seyJTdWNjZXNzIjpmYWxzZSwiQWRkcmVzcyI6IiIsIkhvc3QiOiIiLCJCeXRlc1JlY2VpdmVkIjowLCJFbGFwc2VkVGltZSI6MTUzNTk0MjM5NCwiVFRMIjoyMywiRXJyb3IiOnt9fSx7IlN1Y2Nlc3MiOmZhbHNlLCJBZGRyZXNzIjoiIiwiSG9zdCI6IiIsIkJ5dGVzUmVjZWl2ZWQiOjAsIkVsYXBzZWRUaW1lIjoxNTM1ODgyODE3LCJUVEwiOjI0LCJFcnJvciI6e319LHsiU3VjY2VzcyI6ZmFsc2UsIkFkZHJlc3MiOiIiLCJIb3N0IjoiIiwiQnl0ZXNSZWNlaXZlZCI6MCwiRWxhcHNlZFRpbWUiOjE1MzU5MzU5NjcsIlRUTCI6MjUsIkVycm9yIjp7fX0seyJTdWNjZXNzIjpmYWxzZSwiQWRkcmVzcyI6IiIsIkhvc3QiOiIiLCJCeXRlc1JlY2VpdmVkIjowLCJFbGFwc2VkVGltZSI6MTUzNTkwMzc2MSwiVFRMIjoyNiwiRXJyb3IiOnt9fSx7IlN1Y2Nlc3MiOmZhbHNlLCJBZGRyZXNzIjoiIiwiSG9zdCI6IiIsIkJ5dGVzUmVjZWl2ZWQiOjAsIkVsYXBzZWRUaW1lIjoxNTM1OTcxNDg0LCJUVEwiOjI3LCJFcnJvciI6e319LHsiU3VjY2VzcyI6ZmFsc2UsIkFkZHJlc3MiOiIiLCJIb3N0IjoiIiwiQnl0ZXNSZWNlaXZlZCI6MCwiRWxhcHNlZFRpbWUiOjE1MzU5MjQwMTgsIlRUTCI6MjgsIkVycm9yIjp7fX0seyJTdWNjZXNzIjpmYWxzZSwiQWRkcmVzcyI6IiIsIkhvc3QiOiIiLCJCeXRlc1JlY2VpdmVkIjowLCJFbGFwc2VkVGltZSI6MTUzNTk2MDMzNCwiVFRMIjoyOSwiRXJyb3IiOnt9fSx7IlN1Y2Nlc3MiOmZhbHNlLCJBZGRyZXNzIjoiIiwiSG9zdCI6IiIsIkJ5dGVzUmVjZWl2ZWQiOjAsIkVsYXBzZWRUaW1lIjoxNTM1ODc3MzUwLCJUVEwiOjMwLCJFcnJvciI6e319LHsiU3VjY2VzcyI6ZmFsc2UsIkFkZHJlc3MiOiIiLCJIb3N0IjoiIiwiQnl0ZXNSZWNlaXZlZCI6MCwiRWxhcHNlZFRpbWUiOjE1MzU5MjU5ODAsIlRUTCI6MzEsIkVycm9yIjp7fX0seyJTdWNjZXNzIjpmYWxzZSwiQWRkcmVzcyI6IiIsIkhvc3QiOiIiLCJCeXRlc1JlY2VpdmVkIjowLCJFbGFwc2VkVGltZSI6MTUzNTkzMDM2OSwiVFRMIjozMiwiRXJyb3IiOnt9fV19",
"Error" : ""}
	`)
	var testTResult sensor.TResult
	err := json.Unmarshal(taskResult, &testTResult)
	assert.NotNil(t, err)

	tResult, err := base64.StdEncoding.DecodeString(string(testTResult.Result))
	assert.Nil(t, err)
	assert.NotNil(t, tResult)


	fmt.Println(string(tResult))
}
