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
			Timeout: time.Second,
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

	if err := client.send("http://dev.nul", []byte{}, nil); assert.Error(t, err) {

		assert.Contains(t, err.Error(), "no such host")
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

func TestMisc(t *testing.T) {

	empty := EmptyParams{}
	logic := LogicError{message: "Logic"}

	assert.True(t, empty.IsValid())
	assert.Equal(t, "Logic", logic.Error())
}
