package server

import (
	"fmt"
	"github.com/kshvakov/jsonrpc2"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestRegisterFunc(t *testing.T) {

	num := rand.Intn(50) + 5

	server := New()

	for i := 0; i < num; i++ {

		server.RegisterFunc(fmt.Sprintf("MethodWithEmptyParams_%d", i), func(params jsonrpc2.EmptyParams) (interface{}, error) {

			return nil, nil
		})
	}

	if assert.Len(t, server.handlers, num) {

		for i := 0; i < num; i++ {

			_, ok := server.handlers[fmt.Sprintf("MethodWithEmptyParams_%d", i)]

			assert.True(t, ok)
		}
	}
}

func TestMethodExists(t *testing.T) {

	server := New()

	assert.NotPanics(t, func() {
		server.RegisterFunc("MethodWithEmptyParams", func(params jsonrpc2.EmptyParams) (interface{}, error) {

			return nil, nil
		})
	})

	assert.Panics(t, func() {
		server.RegisterFunc("MethodWithEmptyParams", func(params jsonrpc2.EmptyParams) (interface{}, error) {

			return nil, nil
		})
	})
}

type testObject struct{}

func (t *testObject) MethodWithEmptyParams(params jsonrpc2.EmptyParams) (interface{}, error) {

	return nil, nil
}

func TestRegisterObject(t *testing.T) {

	methods := []string{"MethodWithEmptyParams"}

	server := New()
	server.RegisterObject("TestObject", &testObject{})

	if assert.Len(t, server.handlers, len(methods)) {

		for _, method := range methods {

			_, ok := server.handlers[fmt.Sprintf("TestObject.%s", method)]

			assert.True(t, ok)
		}
	}
}

type testNotSuitableMethods struct{}

func (t *testNotSuitableMethods) Exported(params jsonrpc2.EmptyParams) (interface{}, error) {

	return nil, nil
}

func (t *testNotSuitableMethods) unexported(params jsonrpc2.EmptyParams) (interface{}, error) {

	return nil, nil
}

func (t *testNotSuitableMethods) WithoutParams() (interface{}, error) {

	return nil, nil
}

func (t *testNotSuitableMethods) WithoutResult(params jsonrpc2.EmptyParams) error {

	return nil
}

func (t *testNotSuitableMethods) WithoutError(params jsonrpc2.EmptyParams) string {

	return ""
}

func (t *testNotSuitableMethods) WithResultWithoutError(params jsonrpc2.EmptyParams) (string, struct{}) {

	return "", struct{}{}
}

func TestNotSuitableMethods(t *testing.T) {

	notSuitableMethod := []string{"unexported", "WithoutParams", "WithoutResult", "WithoutError", "WithResultWithoutError"}

	server := New()
	server.RegisterObject("TestNotSuitableMethods", &testNotSuitableMethods{})

	if assert.Len(t, server.handlers, 1) {

		_, ok := server.handlers["TestNotSuitableMethods.Exported"]

		if assert.True(t, ok) {

			for _, method := range notSuitableMethod {

				_, ok := server.handlers[fmt.Sprintf("TestNotSuitableMethods.%s", method)]

				assert.False(t, ok)
			}
		}
	}
}
