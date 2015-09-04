package jsonrpc2

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_sendHttpMethod(t *testing.T) {

	client := &client{
		httpClient: &http.Client{},
	}

	var request *http.Request

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		request = r

		fmt.Fprintln(w, "{}")
	}))

	defer testServer.Close()

	if assert.NoError(t, client.send(testServer.URL, []byte{}, nil)) {

		assert.Equal(t, "POST", request.Method)
	}
}

func TestClient_sendLogicError(t *testing.T) {

	client := &client{
		httpClient: &http.Client{},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		json.NewEncoder(w).Encode(&Response{
			Jsonrpc:   "2.0",
			RequestID: 42,
			Error: &Error{
				Code:    LogicErr,
				Message: "LogicErrror",
			},
		})
	}))

	defer testServer.Close()

	if err := client.send(testServer.URL, []byte{}, nil); assert.Error(t, err) {

		_, ok := err.(*LogicError)

		if assert.True(t, ok) {

			assert.Equal(t, "LogicErrror", err.Error())
		}
	}
}

func TestClient_sendError(t *testing.T) {

	client := &client{
		httpClient: &http.Client{
			Timeout: time.Second * 2,
		},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		json.NewEncoder(w).Encode(&Response{
			Jsonrpc:   "2.0",
			RequestID: 42,
			Error: &Error{
				Code:    InternalError,
				Message: Errors[InternalError],
			},
		})
	}))

	defer testServer.Close()

	if err := client.send(testServer.URL, []byte{}, nil); assert.Error(t, err) {

		assert.Equal(t, errorFmt(InternalError, Errors[InternalError]).Error(), err.Error())
	}

	if err := client.send("http://dev.null", []byte{}, nil); assert.Error(t, err) {

		assert.Contains(t, err.Error(), "dev.null")
	}
}

func TestClient_sendDecodeError(t *testing.T) {

	client := &client{
		httpClient: &http.Client{},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		json.NewEncoder(w).Encode(&Response{
			Jsonrpc:   "2.0",
			RequestID: 42,
			Result:    "42",
		})
	}))

	defer testServer.Close()

	var result []int

	if err := client.send(testServer.URL, []byte{}, &result); assert.Error(t, err) {

		assert.Contains(t, err.Error(), "json: cannot unmarshal")
	}
}

func TestClientBalancer(t *testing.T) {

	var request *http.Request

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		request = r

		json.NewEncoder(w).Encode(&Response{
			Jsonrpc:   "2.0",
			RequestID: 42,
			Result:    "42",
		})
	}))

	defer testServer.Close()

	client := NewClient(&testDiscovery{addresses: []string{"http://dev.null", "http://dev.null2", testServer.URL}})

	if err := client.Send("", &EmptyParams{}, nil); assert.NoError(t, err) {

		assert.Equal(t, "POST", request.Method)
	}
}

func TestClientLogicError(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		json.NewEncoder(w).Encode(&Response{
			Jsonrpc:   "2.0",
			RequestID: 42,
			Error: &Error{
				Code:    LogicErr,
				Message: "LogicErrror",
			},
		})
	}))

	defer testServer.Close()

	client := NewClient(&testDiscovery{addresses: []string{testServer.URL}})

	if err := client.Send("", &EmptyParams{}, nil); assert.Error(t, err) {

		_, ok := err.(*LogicError)

		if assert.True(t, ok) {

			assert.Equal(t, "LogicErrror", err.Error())
		}
	}
}

func TestClientLastError(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		json.NewEncoder(w).Encode(&Response{
			Jsonrpc:   "2.0",
			RequestID: 42,
			Error: &Error{
				Code:    InvalidParams,
				Message: Errors[InvalidParams],
			},
		})
	}))

	defer testServer.Close()

	client := NewClient(&testDiscovery{addresses: []string{"http://dev.null", testServer.URL}})

	if err := client.Send("", &EmptyParams{}, nil); assert.Error(t, err) {

		assert.Equal(t, errorFmt(InvalidParams, Errors[InvalidParams]).Error(), err.Error())
	}
}

func TestClientErrorNoLiveUpstreams(t *testing.T) {

	client := NewClient(&testDiscovery{})

	if err := client.Send("", &EmptyParams{}, nil); assert.Error(t, err) {

		_, ok := err.(*ErrorNoLiveUpstreams)

		assert.True(t, ok)
	}
}

func TestMisc(t *testing.T) {

	empty := EmptyParams{}
	logic := LogicError{message: "Logic"}

	assert.True(t, empty.IsValid())
	assert.Equal(t, "Logic", logic.Error())
}
