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
	"time"

	"github.com/ping-42/42lib/testingkit"
	"github.com/stretchr/testify/assert"
)

// will ping one v4 & one v6 IPs(two total) for each will call 3 times(total 6 pings)
// on each  will wait for 200 milliseconds in the WriteToFunc
// and for 300 milliseconds in the ReadFromFunc
// total should wait ~ 1.5 seconds per IP(not calculating the execution time)
func Test_pingHost_general(t *testing.T) {

	ipV4 := "127.0.0.1"
	ipV6 := "::1"
	readBytesPerCall := 101
	callsPerIP := 3
	writeTime := time.Millisecond * 200
	readTime := time.Millisecond * 300
	expectedTotalRtt := (writeTime + readTime) * time.Duration(callsPerIP)
	// thresholdPerRtt - this should cover the time for the execution
	thresholdPerRtt := 200 * time.Millisecond

	// mock message with default win payload
	receivedMessage := []byte(fmt.Sprintf(`{"Id":"b3a74791-4e4f-4457-a601-fbd685d8e389","Name":"ICMP_TASK","SensorId":"b9dc3d20-256b-4ac7-8cae-2f6dc962e183","Opts":{"TargetDomain":"","TargetIPs":["%v", "%v"],"Count":%v,"Payload":"MDgwOTBhMGIwYzBkMGUwZjEwMTExMjEzMTQxNTE2MTcxODE5MWExYjFjMWQxZTFmMjAyMTIyMjMyNDI1MjYyNzI4MjkyYTJiMmMyZDJlMmYzMDMxMzIzMzM0MzUzNjM3"}}`, ipV4, ipV6, callsPerIP))

	// create a context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create new task
	icmpTask, err := NewTaskFromBytes(receivedMessage)
	if err != nil {
		t.Errorf("Received error from NewTaskFromBytes:%v", err)
		return
	}

	// Mocked instance with custom behavior for WriteToFunc
	var mockedICMPConn = testingkit.MockedICMPConn{
		WriteToFunc: func(b []byte, addr net.Addr) (int, error) {
			// custom behavior for WriteToFunc
			time.Sleep(writeTime)
			return 10, nil
		},
		ReadFromFunc: func(b []byte) (int, net.Addr, error) {
			// custom behavior for WriteToFunc
			time.Sleep(readTime)
			return readBytesPerCall, nil, nil
		},
	}

	// the actual function that will test
	resPerIp, err := icmpTask.pingHost(ctx, mockedICMPConn, mockedICMPConn)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(resPerIp), 2)
	// ---V4----
	resIpV4 := resPerIp[ipV4]
	assert.Equal(t, resIpV4.IPAddr, net.ParseIP(ipV4))
	assert.Equal(t, resIpV4.PacketsSent, callsPerIP)
	assert.Equal(t, resIpV4.PacketsReceived, callsPerIP)
	assert.Equal(t, resIpV4.BytesWritten, 312) // the bytes count of Payload structured as json summer per all callsPerIP(Payload * callsPerIP)
	assert.Equal(t, resIpV4.BytesRead, readBytesPerCall*callsPerIP)
	assert.Equal(t, resIpV4.Loss, float64(0))
	assert.Equal(t, resIpV4.FailureMessages, []string(nil))
	// assert the totalRtt to be between 1.5 and 1.51 seconds
	if resIpV4.TotalRTT < expectedTotalRtt || resIpV4.TotalRTT > expectedTotalRtt+thresholdPerRtt {
		t.Errorf("Expected Total Rtt duration between 1.5 and 1.51 seconds, got %s", resIpV4.TotalRTT)
	}
	if resIpV4.MinRTT < (writeTime+readTime) || resIpV4.MinRTT > (writeTime+readTime)+thresholdPerRtt {
		t.Errorf("Expected MinRTT between %v and %v milliseconds, got %s", (writeTime + readTime), resIpV4.MinRTT, (writeTime+readTime)+thresholdPerRtt)
	}
	if resIpV4.MaxRTT < (writeTime+readTime) || resIpV4.MaxRTT > (writeTime+readTime)+thresholdPerRtt {
		t.Errorf("Expected MaxRTT between %v and %v milliseconds, got %s", (writeTime + readTime), resIpV4.MaxRTT, (writeTime+readTime)+thresholdPerRtt)
	}
	if resIpV4.MaxRTT < resIpV4.MinRTT {
		t.Errorf("Expected MaxRTT > MinRTT received:%v, %v", resIpV4.MaxRTT, resIpV4.MinRTT)
	}
	if resIpV4.AverageRTT > resIpV4.MaxRTT || resIpV4.AverageRTT < resIpV4.MinRTT {
		t.Errorf("Expected AverageRTT to be between MinRTT & MaxRTT, received:%v, %v, %v", resIpV4.AverageRTT, resIpV4.MaxRTT, resIpV4.MinRTT)
	}
	// ---V6----
	resIpV6 := resPerIp[ipV6]
	assert.Equal(t, resIpV6.IPAddr, net.ParseIP(ipV6))
	assert.Equal(t, resIpV6.PacketsSent, callsPerIP)
	assert.Equal(t, resIpV6.PacketsReceived, callsPerIP)
	assert.Equal(t, resIpV6.BytesWritten, 312) // the bytes count of Payload structured as json summer per all callsPerIP(Payload * callsPerIP)
	assert.Equal(t, resIpV6.BytesRead, readBytesPerCall*callsPerIP)
	assert.Equal(t, resIpV6.Loss, float64(0))
	assert.Equal(t, resIpV6.FailureMessages, []string(nil))
	// assert the totalRtt to be between 1.5 and 1.51 seconds
	if resIpV6.TotalRTT < expectedTotalRtt || resIpV6.TotalRTT > expectedTotalRtt+thresholdPerRtt {
		t.Errorf("Expected Total Rtt duration between 1.5 and 1.51 seconds, got %s", resIpV6.TotalRTT)
	}
	if resIpV6.MinRTT < (writeTime+readTime) || resIpV6.MinRTT > (writeTime+readTime)+thresholdPerRtt {
		t.Errorf("Expected MinRTT between %v and %v milliseconds, got %s", (writeTime + readTime), resIpV6.MinRTT, (writeTime+readTime)+thresholdPerRtt)
	}
	if resIpV6.MaxRTT < (writeTime+readTime) || resIpV6.MaxRTT > (writeTime+readTime)+thresholdPerRtt {
		t.Errorf("Expected MaxRTT between %v and %v milliseconds, got %s", (writeTime + readTime), resIpV6.MaxRTT, (writeTime+readTime)+thresholdPerRtt)
	}
	if resIpV6.MaxRTT < resIpV6.MinRTT {
		t.Errorf("Expected MaxRTT > MinRTT received:%v, %v", resIpV6.MaxRTT, resIpV6.MinRTT)
	}
	if resIpV6.AverageRTT > resIpV6.MaxRTT || resIpV6.AverageRTT < resIpV6.MinRTT {
		t.Errorf("Expected AverageRTT to be between MinRTT & MaxRTT, received:%v, %v, %v", resIpV6.AverageRTT, resIpV6.MaxRTT, resIpV6.MinRTT)
	}

	fmt.Printf("%+v", resPerIp)
}

// const (
// 	protocolICMP = 1
// )

// func startTestServer(ip string, expectedTotalRtt time.Duration, msgChan chan<- string) {
// 	// Listen for incoming ICMP packets
// 	c, err := icmp.ListenPacket("ip4:icmp", ip)
// 	if err != nil {
// 		msgChan <- fmt.Sprintf("Mocked ICMP server: Failed to listen for ICMP packets:%v", err)
// 		os.Exit(1)
// 	}
// 	defer c.Close()

// 	// Buffer to read incoming packets
// 	rb := make([]byte, 1500)
// 	for {
// 		// Read incoming ICMP messages
// 		n, addr, err := c.ReadFrom(rb)
// 		if err != nil {
// 			msgChan <- fmt.Sprintf("Mocked ICMP server: Read error:%v", err)
// 			return
// 		}

// 		// Record the time when the request is received
// 		requestReceivedTime := time.Now()

// 		// Parse the message
// 		rm, err := icmp.ParseMessage(protocolICMP, rb[:n])
// 		if err != nil {
// 			msgChan <- fmt.Sprintf("Mocked ICMP server: Error parsing ICMP message:%v", err)
// 			return
// 		}

// 		switch rm.Type {
// 		case ipv4.ICMPTypeEcho:

// 			// Create an echo reply message
// 			m := icmp.Message{
// 				Type: ipv4.ICMPTypeEchoReply,
// 				Code: 0,
// 				Body: &icmp.Echo{
// 					ID:   rm.Body.(*icmp.Echo).ID,
// 					Seq:  rm.Body.(*icmp.Echo).Seq,
// 					Data: rm.Body.(*icmp.Echo).Data,
// 				},
// 			}

// 			// Marshal the message into bytes
// 			b, err := m.Marshal(nil)
// 			if err != nil {
// 				msgChan <- fmt.Sprintf("Mocked ICMP server: Error marshalling reply:%v", err)
// 				return
// 			}

// 			// // Calculate the remaining time needed to meet the expected TotalRTT
// 			processingTime := time.Since(requestReceivedTime)
// 			remainingDelay := expectedTotalRtt - processingTime

// 			// Introduce additional delay if necessary to meet the expected TotalRTT
// 			if remainingDelay > 0 {
// 				time.Sleep(remainingDelay)
// 			}

// 			// Send the reply
// 			if _, err := c.WriteTo(b, addr); err != nil {
// 				msgChan <- fmt.Sprintf("Mocked ICMP server: Error sending reply:%v", err)
// 				return
// 			} else {
// 				msgChan <- fmt.Sprintf("Mocked ICMP server: Sent ICMP Echo reply to %s\n", addr)
// 				return
// 			}
// 		default:
// 			// Ignore other ICMP types
// 			msgChan <- fmt.Sprintf("Mocked ICMP server: Got non-echo request from %s: type %v\n", addr, rm.Type)
// 			return
// 		}
// 	}
// }

// func Test_WithMockedServer(t *testing.T) {

// 	fmt.Sprintln("----->Test_Run_Mocked")

// 	mockedIpEndpoint := "127.0.0.1"

// 	// Channel to receive server messages
// 	msgChan := make(chan string)
// 	// Start the server in a goroutine
// 	// go startTestServer(mockedIpEndpoint, 2*time.Second, msgChan)

// 	time.Sleep(1 * time.Second) // Wait a bit for the server to be ready

// 	// mock message with default win payload
// 	receivedMessage := []byte(fmt.Sprintf(`{"Id":"b3a74791-4e4f-4457-a601-fbd685d8e389","Name":"ICMP_TASK","SensorId":"b9dc3d20-256b-4ac7-8cae-2f6dc962e183","Opts":{"TargetDomain":"","TargetIPs":["%v"],"Count":1,"Payload":"MDgwOTBhMGIwYzBkMGUwZjEwMTExMjEzMTQxNTE2MTcxODE5MWExYjFjMWQxZTFmMjAyMTIyMjMyNDI1MjYyNzI4MjkyYTJiMmMyZDJlMmYzMDMxMzIzMzM0MzUzNjM3"}}`, mockedIpEndpoint))

// 	// create a context and cancel it
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	// create new task
// 	icmpTask, err := NewTaskFromBytes(receivedMessage)
// 	if err != nil {
// 		t.Errorf("Received error from NewTaskFromBytes:%v", err)
// 		return
// 	}

// 	res, err := icmpTask.Run(ctx)
// 	if err != nil {
// 		t.Errorf("Received error from Run:%v", err)
// 	}

// 	var icmpRes Result
// 	err = json.Unmarshal(res, &icmpRes)
// 	if err != nil {
// 		t.Errorf("Unmarshal icmp.Result{} err:%v", err)
// 	}

// 	fmt.Printf("%+v", icmpRes)

// 	// Print the first message received from the server
// 	msg := <-msgChan
// 	fmt.Println("\n" + msg + "\n")
// }

// func Test_Run(t *testing.T) {
// 	// mock message with default win payload
// 	receivedMessage := []byte(`{"Id":"b3a74791-4e4f-4457-a601-fbd685d8e389","Name":"ICMP_TASK","SensorId":"b9dc3d20-256b-4ac7-8cae-2f6dc962e183","Opts":{"TargetDomain":"","TargetIPs":["127.0.0.1"],"Count":3,"Payload":"MDgwOTBhMGIwYzBkMGUwZjEwMTExMjEzMTQxNTE2MTcxODE5MWExYjFjMWQxZTFmMjAyMTIyMjMyNDI1MjYyNzI4MjkyYTJiMmMyZDJlMmYzMDMxMzIzMzM0MzUzNjM3"}}`)

// 	// create a context and cancel it
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	// create new task
// 	icmpTask, err := NewTaskFromBytes(receivedMessage)
// 	if err != nil {
// 		t.Errorf("Received error from NewTaskFromBytes:%v", err)
// 		return
// 	}

// 	_, err = icmpTask.Run(ctx)
// 	if err == nil {
// 		t.Errorf("Received error from Run:%v", err)
// 	}
// }

// -------------

// func Test_RunWithCancelledContext(t *testing.T) {
// 	// mock message with default win payload
// 	receivedMessage := []byte(`{"Id":"123","Name":"ICMP_TEST","IcmpOpts":{"TargetIPs":["8.8.8.8", "2001:4860:4860::8844"],"Count":3}}`)

// 	// create a context and cancel it
// 	ctx, cancel := context.WithCancel(context.Background())
// 	cancel()

// 	// create new task
// 	icmpTask := NewTask()
// 	loggerIcmp.Printf("New ICMP Task: %v\n\n", icmpTask)

// 	// call Run with the cancelled context
// 	_, err := icmpTask.Run(ctx, receivedMessage)
// 	if err == nil {
// 		t.Errorf("Expected an error")
// 	}
// }

// func TestCreateConnV4(t *testing.T) {
// 	// create a new icmpV4 connection
// 	conn, err := createConnV4()
// 	if err != nil {
// 		loggerIcmp.Error(err)
// 	} else {
// 		loggerIcmp.Infof("connection: %+v\n", conn)
// 	}
// 	assert.Nil(t, err)
// 	assert.NotNil(t, conn)
// }

// func TestCreateConnV6(t *testing.T) {
// 	// create a new icmpV4 connection
// 	conn, err := createConnV6()
// 	if err != nil {
// 		loggerIcmp.Error(err)
// 	} else {
// 		loggerIcmp.Info("connection: ", conn)
// 	}
// 	assert.Nil(t, err)
// 	assert.NotNil(t, conn)
// }

// func TestRunICMP4(t *testing.T) {
// 	// parse an ipv4 to pass to runICMP
// 	ip := net.ParseIP("8.8.8.8")
// 	if ip == nil {
// 		loggerIcmp.Error("no valid ip")
// 		return
// 	}
// 	// create a new icmpV4 connection
// 	conn, err := createConnV4()
// 	if err != nil {
// 		loggerIcmp.Info(err)
// 		return
// 	}
// 	// get the message type based on ip
// 	msgType := getICMPType(ip)
// 	// create a new task and run the ping
// 	icmpTask := NewTask()
// 	pingResult, err := icmpTask.runICMP(context.TODO(), conn, msgType, ip)
// 	if err != nil {
// 		loggerIcmp.Error(err)
// 	}
// 	loggerIcmp.Infof("ping result: %v", pingResult)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, pingResult)
// }

// func TestRunICMP6(t *testing.T) {
// 	// parse an ipv6 to pass to runICMP
// 	ip := net.ParseIP("2001:4860:4860::8844")
// 	if ip == nil {
// 		loggerIcmp.Error("no valid ip")
// 		return
// 	}
// 	// create a new icmpV6 connection
// 	conn, err := createConnV6()
// 	if err != nil {
// 		loggerIcmp.Info(err)
// 		return
// 	}
// 	// get the message type based on ip
// 	msgType := getICMPType(ip)
// 	// create a new task and run the ping
// 	icmpTask := NewTask()
// 	pingResult, err := icmpTask.runICMP(context.TODO(), conn, msgType, ip)
// 	if err != nil {
// 		loggerIcmp.Error(err)
// 	}
// 	loggerIcmp.Infof("ping result: %v", pingResult)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, pingResult)
// }

// func TestIcmpTask(t *testing.T) {

// 	// mock message with default win payload
// 	receivedMessage := []byte(`{"Id":"123","Name":"ICMP_TEST","IcmpOpts":{"TargetIPs":["8.8.8.8", "2001:4860:4860::8844"],"Count":3}}`)

// 	// Create an instance of the ICMP task with default options
// 	icmpTask := NewTask()

// 	// Call the PingHost function with test options from message
// 	result, err := icmpTask.Run(context.TODO(), receivedMessage)
// 	fmt.Printf("New ICMP Task: %v\n\n", icmpTask)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Printf("unit test icmp task results: %+v\n\n", result)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, result)

// }
