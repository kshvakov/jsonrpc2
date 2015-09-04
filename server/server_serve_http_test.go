package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/kshvakov/jsonrpc2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerMethodNotAllowed(t *testing.T) {

	server := New()

	testServer := httptest.NewServer(server)

	defer testServer.Close()

	client := &http.Client{}

	response, err := client.Get(testServer.URL)

	if assert.NoError(t, err) {

		var result jsonrpc2.Response

		err := json.NewDecoder(response.Body).Decode(&result)

		if assert.NoError(t, err) && assert.NotNil(t, result.Error) {

			assert.True(t, jsonrpc2.InvalidRequest == result.Error.Code)
			assert.Equal(t, jsonrpc2.Errors[jsonrpc2.InvalidRequest], result.Error.Message)
		}
	}
}

func TestServerParseError(t *testing.T) {

	server := New()

	testServer := httptest.NewServer(server)

	defer testServer.Close()

	client := &http.Client{}

	response, err := client.Post(testServer.URL, "application/x-www-form-urlencoded", bytes.NewReader([]byte("string")))

	if assert.NoError(t, err) {

		var result jsonrpc2.Response

		err := json.NewDecoder(response.Body).Decode(&result)

		if assert.NoError(t, err) && assert.NotNil(t, result.Error) {

			assert.True(t, jsonrpc2.ParseError == result.Error.Code)
			assert.Equal(t, jsonrpc2.Errors[jsonrpc2.ParseError], result.Error.Message)
			assert.Contains(t, result.Error.Data, "invalid character")
		}
	}
}

func TestServerPanic(t *testing.T) {

	server := New()
	server.RegisterFunc("TestPanic", func(_ *jsonrpc2.EmptyParams) (interface{}, error) {

		panic("Panic")

		return nil, nil
	})

	testServer := httptest.NewServer(server)

	defer testServer.Close()

	client := &http.Client{}

	req, _ := json.Marshal(&jsonrpc2.Request{
		Method: "TestPanic",
		Params: &jsonrpc2.EmptyParams{},
	})

	response, err := client.Post(testServer.URL, "application/x-www-form-urlencoded", bytes.NewReader(req))

	if assert.NoError(t, err) {

		var result jsonrpc2.Response

		err := json.NewDecoder(response.Body).Decode(&result)

		if assert.NoError(t, err) && assert.NotNil(t, result.Error) {

			assert.True(t, jsonrpc2.InternalError == result.Error.Code)
			assert.Equal(t, jsonrpc2.Errors[jsonrpc2.InternalError], result.Error.Message)
			assert.Contains(t, result.Error.Data, "Panic")
		}
	}
}

func TestServerMethodNotFound(t *testing.T) {

	server := New()
	server.RegisterFunc("Method", func(_ *jsonrpc2.EmptyParams) (interface{}, error) {

		return nil, nil
	})

	testServer := httptest.NewServer(server)

	defer testServer.Close()

	client := &http.Client{}

	req, _ := json.Marshal(&jsonrpc2.Request{
		Method: "MethodNotFound",
		Params: &jsonrpc2.EmptyParams{},
	})

	response, err := client.Post(testServer.URL, "application/x-www-form-urlencoded", bytes.NewReader(req))

	if assert.NoError(t, err) {

		var result jsonrpc2.Response

		err := json.NewDecoder(response.Body).Decode(&result)

		if assert.NoError(t, err) && assert.NotNil(t, result.Error) {

			assert.True(t, jsonrpc2.MethodNotFound == result.Error.Code)
			assert.Equal(t, jsonrpc2.Errors[jsonrpc2.MethodNotFound], result.Error.Message)
		}
	}
}

type testInvalidParams struct{}

func (t *testInvalidParams) IsValid() bool {
	return false
}

func TestServerInvalidParams(t *testing.T) {

	server := New()
	server.RegisterFunc("InvalidParamsMethod", func(_ *testInvalidParams) (interface{}, error) {

		return nil, nil
	})

	testServer := httptest.NewServer(server)

	defer testServer.Close()

	client := &http.Client{}

	req, _ := json.Marshal(&jsonrpc2.Request{
		Method: "InvalidParamsMethod",
		Params: &testInvalidParams{},
	})

	response, err := client.Post(testServer.URL, "application/x-www-form-urlencoded", bytes.NewReader(req))

	if assert.NoError(t, err) {

		var result jsonrpc2.Response

		err := json.NewDecoder(response.Body).Decode(&result)

		if assert.NoError(t, err) && assert.NotNil(t, result.Error) {

			assert.True(t, jsonrpc2.InvalidParams == result.Error.Code)
			assert.Equal(t, jsonrpc2.Errors[jsonrpc2.InvalidParams], result.Error.Message)
		}
	}
}

type testErrorDecodeParams struct {
	A string
	B int
	C interface{}
}

func (t *testErrorDecodeParams) IsValid() bool {

	return true
}

type errorDecodeParams struct {
	A int
	B interface{}
	C string
}

func (e *errorDecodeParams) IsValid() bool {

	return true
}

func TestServerErrorDecodeParams(t *testing.T) {

	server := New()
	server.RegisterFunc("ErrorDecodeParams", func(_ *testErrorDecodeParams) (interface{}, error) {

		return nil, nil
	})

	testServer := httptest.NewServer(server)

	defer testServer.Close()

	client := &http.Client{}

	req, _ := json.Marshal(&jsonrpc2.Request{
		Method: "ErrorDecodeParams",
		Params: &errorDecodeParams{
			A: 42,
			B: "B",
			C: "C",
		},
	})

	response, err := client.Post(testServer.URL, "application/x-www-form-urlencoded", bytes.NewReader(req))

	if assert.NoError(t, err) {

		var result jsonrpc2.Response

		err := json.NewDecoder(response.Body).Decode(&result)

		if assert.NoError(t, err) && assert.NotNil(t, result.Error) {

			assert.True(t, jsonrpc2.ParseError == result.Error.Code)
			assert.Equal(t, jsonrpc2.Errors[jsonrpc2.ParseError], result.Error.Message)
			assert.Contains(t, result.Error.Data, "json: cannot unmarshal")
		}
	}
}

func TestServerLogicError(t *testing.T) {

	server := New()
	server.RegisterFunc("LogicError", func(_ *jsonrpc2.EmptyParams) (interface{}, error) {

		return nil, errors.New("Logic Error")
	})

	testServer := httptest.NewServer(server)

	defer testServer.Close()

	client := &http.Client{}

	req, _ := json.Marshal(&jsonrpc2.Request{
		Method: "LogicError",
		Params: &jsonrpc2.EmptyParams{},
	})

	response, err := client.Post(testServer.URL, "application/x-www-form-urlencoded", bytes.NewReader(req))

	if assert.NoError(t, err) {

		var result jsonrpc2.Response

		err := json.NewDecoder(response.Body).Decode(&result)

		if assert.NoError(t, err) && assert.NotNil(t, result.Error) {

			assert.True(t, jsonrpc2.LogicErr == result.Error.Code)
			assert.Equal(t, "Logic Error", result.Error.Message)
		}
	}
}

type testOk struct {
	Message string
}

func TestServerOk(t *testing.T) {

	server := New()
	server.RegisterFunc("Ok", func(_ *jsonrpc2.EmptyParams) (interface{}, error) {

		return &testOk{Message: "OK"}, nil
	})

	testServer := httptest.NewServer(server)

	defer testServer.Close()

	client := &http.Client{}

	req, _ := json.Marshal(&jsonrpc2.Request{
		Method: "Ok",
		Params: &jsonrpc2.EmptyParams{},
	})

	response, err := client.Post(testServer.URL, "application/x-www-form-urlencoded", bytes.NewReader(req))

	if assert.NoError(t, err) {

		result := jsonrpc2.Response{
			Result: &testOk{},
		}

		err := json.NewDecoder(response.Body).Decode(&result)

		if assert.NoError(t, err) && assert.Nil(t, result.Error) {

			if r, ok := result.Result.(*testOk); assert.True(t, ok) {

				assert.Equal(t, "OK", r.Message)
			}
		}
	}
}
