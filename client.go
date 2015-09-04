package jsonrpc2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func init() {

	rand.Seed(time.Now().UnixNano())
}

type Client interface {
	Send(method string, params Params, result interface{}) error
}

func NewClient(discovery Discovery) Client {

	return &client{
		balancer: newBalancer(discovery),
		httpClient: &http.Client{
			Timeout: time.Second,
		},
	}
}

type client struct {
	balancer   *balancer
	httpClient *http.Client
}

func (c *client) Send(method string, params Params, result interface{}) error {

	data, _ := json.Marshal(Request{
		Jsonrpc:   "2.0",
		RequestID: rand.Int(),
		Method:    method,
		Params:    params,
	})

	var lastError error

	if c.balancer.len() == 0 {

		return &ErrorNoLiveUpstreams{}
	}

	for i := 0; i < c.balancer.len(); i++ {

		url, err := c.balancer.next()

		if err != nil {

			return err
		}

		lastError = c.send(url, data, result)

		if lastError == nil {

			return nil
		}

		if _, ok := lastError.(*LogicError); ok {

			return lastError
		}
	}

	return lastError
}

func (c *client) send(url string, data []byte, result interface{}) error {

	response, err := c.httpClient.Post(url, "application/x-www-form-urlencoded", bytes.NewReader(data))

	if err != nil {

		return err
	}

	r := Response{
		Result: struct{}{},
	}

	if result != nil {

		r.Result = result
	}

	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {

		return err
	}

	if r.Error == nil {

		return nil
	}

	if r.Error.Code == LogicErr {

		return &LogicError{message: r.Error.Message}
	}

	return errorFmt(r.Error.Code, r.Error.Message)
}

func errorFmt(code int16, message string) error {

	return fmt.Errorf("%d:%s", code, message)
}
