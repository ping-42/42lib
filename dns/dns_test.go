/*
Until the tests are properly implemented and mocked,
ignore this from the pipeline and use for debugging only.
*/

package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_udp(t *testing.T) {
	receivedMessage := []byte(`{"Id":"4a812264-46f0-4bc9-b49b-ae36164fdaa8","Name":"DNS_TASK","SensorId":"b9dc3d20-256b-4ac7-8cae-2f6dc962e183","DnsOpts":{"Host":"https://example.com","Proto":"udp","DnsUdpClient":null}}`)

	dnsTask, err := NewTaskFromBytes(receivedMessage)
	assert.Nil(t, err)
	result, err := dnsTask.Run(context.TODO())

	assert.Nil(t, err)
	assert.NotNil(t, result)

	var dnsRes = Result{}
	err = json.Unmarshal(result, &dnsRes)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(dnsRes.AnswerA), 1)

	fmt.Printf("dns Result: %+v\n", dnsRes)
}

func Test_defaultNsGetter(t *testing.T) {
	ns, err := getDefaultNs()
	assert.Nil(t, err)
	assert.NotNil(t, ns)
}

func TestDnsOverTcp(t *testing.T) {
	receivedMessage := []byte(`{"Id":"4a812264-46f0-4bc9-b49b-ae36164fdaa8","Name":"DNS_TASK","SensorId":"b9dc3d20-256b-4ac7-8cae-2f6dc962e183","DnsOpts":{"Host":"https://example.com","Proto":"tcp","DnsUdpClient":null}}`)

	dnsTask, err := NewTaskFromBytes(receivedMessage)
	assert.Nil(t, err)
	result, err := dnsTask.Run(context.TODO())

	assert.Nil(t, err)
	assert.NotNil(t, result)

	var dnsRes = Result{}
	err = json.Unmarshal(result, &dnsRes)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(dnsRes.AnswerA), 1)

	fmt.Printf("dns Result: %+v\n", dnsRes)
}
