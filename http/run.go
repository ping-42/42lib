package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/ping-42/42lib/db/models"
)

// Run is the entry point for the http.task
func (t task) Run(ctx context.Context) (res []byte, err error) {
	if ctx.Err() != nil {
		return nil, fmt.Errorf("context done detected in Run:%v", ctx.Err())
	}

	var result Result
	req, err := NewRequest(t.Opts, &result)
	if err != nil {
		err = fmt.Errorf("NewRequest failed:%v", err)
		return
	}

	client := DefaultClient()
	resDo, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("client.Do failed:%v", err)
		return
	}

	respBody, err := io.ReadAll(resDo.Body)
	if err != nil {
		err = fmt.Errorf("reading response body failed: %v", err)
		return
	}
	result.ResponseBody = string(respBody)
	result.ResponseCode = resDo.StatusCode
	result.ResponseHeaders = resDo.Header

	err = resDo.Body.Close()
	if err != nil {
		err = fmt.Errorf("res.Body.Close failed:%v", err)
		return
	}
	result.End(time.Now())

	res, err = json.Marshal(result)
	if err != nil {
		err = fmt.Errorf("Marshal failed:%v", err)
		return
	}
	return
}

// NewTaskFromBytes used in sensor for building the task from the received bytes
func NewTaskFromBytes(msg []byte) (t task, err error) {

	// build the http task from the received msg
	err = json.Unmarshal(msg, &t)
	if err != nil {
		err = fmt.Errorf("http.NewTask Unmarshal err task:%v, %v", string(msg), err)
		return
	}

	return t, nil
}

// NewTaskFromModel used in scheduler for building the task from the db model task
func NewTaskFromModel(t models.Task) (tRes task, err error) {

	var o = Opts{}
	err = json.Unmarshal(t.Opts, &o)
	if err != nil {
		err = fmt.Errorf("http NewTaskFromModel Unmarshal Opts err:%v", err)
		return
	}

	tRes.Id = t.ID
	tRes.SensorId = t.SensorID
	tRes.Name = TaskName
	tRes.Opts = o
	return
}
