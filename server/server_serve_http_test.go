package server

import (
	"encoding/json"
	"github.com/kshvakov/jsonrpc2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerMethodNotAllowed(t *testing.T) {

	server := New()
	server.RegisterFunc("Test", func() (string, error) {

		return "", nil
	})

	testServer := httptest.NewServer(server)

	defer testServer.Close()

	client := &http.Client{}

	response, err := client.Get(testServer.URL)

	if assert.NoError(t, err) {

		var result jsonrpc2.Response

		err := json.NewDecoder(response.Body).Decode(&result)

		if assert.NoError(t, err) {
			assert.True(t, jsonrpc2.InvalidRequest == result.Error.Code)
			assert.Equal(t, jsonrpc2.Errors[jsonrpc2.InvalidRequest], result.Error.Message)
		}
	}
}
