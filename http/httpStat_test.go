package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

const (
	TestDomainHTTP  = "http://example.com"
	TestDomainHTTPS = "https://example.com"
)

func TestRun(t *testing.T) {
	// The server/endpoint is not mocked!
	receivedMessage := []byte(`{"Id":"51395b2a-b09d-4a1b-87f5-7ca1a9039aa3","Name":"HTTP_TEST","SensorId":"e26e85a8-aa99-43d0-8e1a-78b6056449a2","HttpOpts":{"URL":"https://google.com","HttpMethod":"GET","RequestHeaders":{"Content-Type":["application/json"]},"RequestBody":"c29tZSB0ZXN0IGJvZHk="}}`)
	httpTask, err := NewTaskFromBytes(receivedMessage)
	if err != nil {
		t.Fatal("NewTaskFromBytes failed:", err)
	}

	result, err := httpTask.Run(context.TODO())
	if err != nil {
		t.Fatal("client.Do failed:", err)
	}

	fmt.Printf("%+v\n", string(result))
}

func TestHTTPStat_HTTPS(t *testing.T) {
	var opts = Opts{
		URL: TestDomainHTTPS,
	}
	var result Result
	req, err := NewRequest(opts, &result)
	if err != nil {
		t.Fatal("client.Do failed:", err)
	}

	client := DefaultClient()
	res, err := client.Do(req)
	if err != nil {
		t.Fatal("client.Do failed:", err)
	}

	if _, err := io.Copy(io.Discard, res.Body); err != nil {
		t.Fatal("io.Copy failed:", err)
	}
	err = res.Body.Close()
	if err != nil {
		t.Fatal("res.Body.Close failed:", err)
	}
	result.End(time.Now())

	if !result.isTLS {
		t.Fatal("isTLS should be true")
	}

	for k, d := range result.durations() {
		if d <= 0*time.Millisecond {
			t.Fatalf("expect %s to be non-zero", k)
		}
	}
}

func TestHTTPStat_HTTP(t *testing.T) {
	var opts = Opts{
		URL: TestDomainHTTP,
	}
	var result Result
	req, err := NewRequest(opts, &result)
	if err != nil {
		t.Fatal("client.Do failed:", err)
	}

	client := DefaultClient()
	res, err := client.Do(req)
	if err != nil {
		t.Fatal("client.Do failed:", err)
	}

	if _, err := io.Copy(io.Discard, res.Body); err != nil {
		t.Fatal("io.Copy failed:", err)
	}
	err = res.Body.Close()
	if err != nil {
		t.Fatal("res.Body.Close failed:", err)
	}
	result.End(time.Now())

	if result.isTLS {
		t.Fatal("isTLS should be false")
	}

	if got, want := result.TLSHandshake, 0*time.Millisecond; got != want {
		t.Fatalf("TLSHandshake time of HTTP = %d, want %d", got, want)
	}

	// Except TLS should be non zero
	durations := result.durations()
	result.PrintStatus()
	delete(durations, "TLSHandshake")

	for k, d := range durations {
		if d <= 0*time.Millisecond {
			t.Fatalf("expect %s to be non-zero", k)
		}
	}
}

func TestHTTPStat_KeepAlive(t *testing.T) {

	req1, err := http.NewRequest("GET", TestDomainHTTPS, nil)
	if err != nil {
		t.Fatal("NewRequest failed:", err)
	}

	client := DefaultClient()
	res1, err := client.Do(req1)
	if err != nil {
		t.Fatal("Request failed:", err)
	}

	if _, err := io.Copy(io.Discard, res1.Body); err != nil {
		t.Fatal("Copy body failed:", err)
	}
	err = res1.Body.Close()
	if err != nil {
		t.Fatal("res.Body.Close failed:", err)
	}

	var opts = Opts{
		URL: TestDomainHTTPS,
	}

	var result Result
	req2, err := NewRequest(opts, &result)
	if err != nil {
		t.Fatal("client.Do failed:", err)
	}

	// When second request, connection should be re-used.
	res2, err := client.Do(req2)
	if err != nil {
		t.Fatal("Request failed:", err)
	}

	if _, err := io.Copy(io.Discard, res2.Body); err != nil {
		t.Fatal("Copy body failed:", err)
	}
	err = res2.Body.Close()
	if err != nil {
		t.Fatal("res.Body.Close failed:", err)
	}
	result.End(time.Now())

	// The following values should be zero.
	// Because connection is reused.
	durations := []time.Duration{
		result.DNSLookup,
		result.TCPConnection,
		result.TLSHandshake,
	}

	for i, d := range durations {
		if got, want := d, 0*time.Millisecond; got != want {
			t.Fatalf("#%d expect %d to be eq %d", i, got, want)
		}
	}
}

func TestTotal_Zero(t *testing.T) {
	result := &Result{}
	result.End(time.Now())

	zero := 0 * time.Millisecond
	if result.total != zero {
		t.Fatalf("Total time is %d, want %d", result.total, zero)
	}

	if result.contentTransfer != zero {
		t.Fatalf("Total time is %d, want %d", result.contentTransfer, zero)
	}
}
