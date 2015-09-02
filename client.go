package jsonrpc2

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Client interface {
	Send(request Request, result interface{}) error
}

func NewClient() Client {

	return &client{
		httpClient: &http.Client{
			Timeout: time.Second * 2,
		},
	}
}

type client struct {
	httpClient *http.Client
	balancer   balancer
}

func (c *client) Send(method string, params Params, result interface{}) error {

	data, _ := json.Marshal(Request{
		Version: "2.0",
		Id:      string(rand.Int63()),
		Method:  method,
		Params:  params,
	})

	for {

		if url, stop := balancer.Next(); !stop {

			if response, err := c.httpClient.Post(url, "application/x-www-form-urlencoded", bytes.NewReader(data)); err == nil {

			}
		}
	}
}
