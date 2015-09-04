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

func TestEnd2End(t *testing.T) {

	app := server.New()
	app.RegisterObject("End2End", &testService{})

	testServer := httptest.NewServer(app)

	defer testServer.Close()

	client := jsonrpc2.NewClient(&testDiscovery{addresses: []string{testServer.URL}})

	var (
		testSumResult     TestSumResult
		emptyParamsResult string
	)

	if err := client.Send("End2End.EmptyParams", &jsonrpc2.EmptyParams{}, &emptyParamsResult); assert.NoError(t, err) {

		assert.Equal(t, "EmptyParams", emptyParamsResult)
	}

	if err := client.Send("End2End.Sum", &TestSumParam{A: 2, B: 3}, &testSumResult); assert.NoError(t, err) {

		assert.Equal(t, 5, testSumResult.Result)
	}
}
