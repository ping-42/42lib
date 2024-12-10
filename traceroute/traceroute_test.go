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
	"time"

	"github.com/google/uuid"
	"github.com/ping-42/42lib/db/models"
	testingkit "github.com/ping-42/42lib/testingkit"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/unix"
)

// this can be tested on root vscode
func TestTracerouteTaskFromBytes(t *testing.T) {

	// mock message with default win payload
	receivedMessage := []byte(`{"Id":"3b241101-e2bb-4255-8caf-4136c566a964","Name":"TRACEROUTE_TASK","SensorID":"3b241101-e2bb-4255-8caf-4136c566a964","Opts":{"Port":33434,"Dest":"8.8.8.8","FirstHop":1,"MaxHops":64,"Timeout":500,"PacketSize":52,"Retries":3}}`)

	// Create an instance of the traceroute task with default options
	tracerouteTask, err := NewTaskFromBytes(receivedMessage)
	if err != nil {
		fmt.Println("eror creating task:", err)
	}
	fmt.Printf("New TRACEROUTE Task: %v\n\n", tracerouteTask)
	// Call the traceroute Run with test options from message
	result, err := tracerouteTask.Run(context.TODO())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("unit test traceroute task results: %+v\n\n", string(result))

	assert.Nil(t, err)
	assert.NotNil(t, result)

}

// this can be tested on root vscode
func TestTracerouteTaskFromModel(t *testing.T) {

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
	if err != nil {
		fmt.Println("eror creating task:", err)
	}

	fmt.Printf("New TRACEROUTE Task: %v\n\n", tracerouteTask)

	// Call the traceoute Run with test options from message
	result, err := tracerouteTask.Run(context.TODO())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("unit test traceroute task results: %+v\n\n", string(result))

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
	if err != nil {
		fmt.Println("eror creating task:", err)
	}

	tracerouteTask.SysUnix = mockSysUnix
	fmt.Printf("New TRACEROUTE Task: %v\n\n", tracerouteTask)

	// Call the traceoute Run with test options from message
	result, err := tracerouteTask.Run(context.TODO())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("unit test traceroute task results: %+v\n\n", string(result))

	assert.Nil(t, err)
	assert.NotNil(t, result)
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
	if err != nil {
		fmt.Println("eror creating task:", err)
	}
	tracerouteTask.SysUnix = mockSysUnix
	fmt.Printf("New TRACEROUTE Task: %v\n\n", tracerouteTask)

	// Call the traceoute Run with test options from message
	result, err := tracerouteTask.Run(context.TODO())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("unit test traceroute task results: %+v\n\n", string(result))

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
