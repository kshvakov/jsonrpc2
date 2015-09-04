package end2end_test

import (
	"github.com/kshvakov/jsonrpc2"
	"github.com/kshvakov/jsonrpc2/server"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

type testDiscovery struct {
	addresses []string
	err       error
}

func (t *testDiscovery) Get() ([]string, error) {

	if t.err != nil {

		return nil, t.err
	}

	return t.addresses, nil
}

type testService struct{}

func (t *testService) EmptyParams(_ *jsonrpc2.EmptyParams) (interface{}, error) {

	return "EmptyParams", nil
}

type TestSumParam struct {
	A, B int
}

func (t *TestSumParam) IsValid() bool {

	return true
}

type TestSumResult struct {
	Result int
}

func (t *testService) Sum(param *TestSumParam) (*TestSumResult, error) {

	return &TestSumResult{
		Result: param.A + param.B,
	}, nil
}

func ExampleClient(url string) *exampleClient {

	return &exampleClient{
		rpc: jsonrpc2.NewClient(&testDiscovery{addresses: []string{url}}),
	}
}

type exampleClient struct {
	rpc jsonrpc2.Client
}

func (e *exampleClient) EmptyParams() (string, error) {

	var result string

	err := e.rpc.Send("End2End.EmptyParams", &jsonrpc2.EmptyParams{}, &result)

	if err != nil {

		return "", err
	}

	return result, nil
}

func (e *exampleClient) Sum(a, b int) (int, error) {

	var result TestSumResult

	err := e.rpc.Send("End2End.Sum", &TestSumParam{A: a, B: b}, &result)

	if err != nil {

		return 0, err
	}

	return result.Result, nil
}

func TestEnd2End(t *testing.T) {

	app := server.New()
	app.RegisterObject("End2End", &testService{})

	testServer := httptest.NewServer(app)

	defer testServer.Close()

	client := ExampleClient(testServer.URL)

	if result, err := client.EmptyParams(); assert.NoError(t, err) {

		assert.Equal(t, "EmptyParams", result)
	}

	if result, err := client.Sum(2, 3); assert.NoError(t, err) {

		assert.Equal(t, 5, result)
	}
}
