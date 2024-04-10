package http

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"time"
)

func DefaultTransport() *http.Transport {
	// It comes from std transport.go
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			Resolver:  nil,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// To avoid shared transport
func DefaultClient() *http.Client {
	return &http.Client{
		Transport: DefaultTransport(),
	}
}

func NewRequest(opts Opts, result *Result) (reqWithContext *http.Request, err error) {

	var requestBody bytes.Buffer
	if opts.RequestBody != nil {
		_, err = requestBody.WriteString(string(opts.RequestBody))
		if err != nil {
			err = fmt.Errorf("requestBody.WriteString err:%v", err)
			return
		}
	}

	req, err := http.NewRequest(opts.HttpMethod, opts.TargetDomain, &requestBody)
	if err != nil {
		err = fmt.Errorf("NewRequest failed:%v", err)
		return
	}

	ctx := WithHTTPStat(req.Context(), result)
	reqWithContext = req.WithContext(ctx)
	return
}
