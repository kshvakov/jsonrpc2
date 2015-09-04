package server

import (
	"encoding/json"
	"fmt"
	"github.com/kshvakov/jsonrpc2"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type testDecodeParams struct {
	Int         int
	String      string
	SliceInt    []int
	SliceString []string
}

func (t *testDecodeParams) IsValid() bool {

	return t.Int == 42
}

func TestHandlerDecodeParams(t *testing.T) {

	params := testDecodeParams{
		Int:         42,
		String:      "42",
		SliceInt:    []int{1, 2, 3, 4, 5},
		SliceString: []string{"1", "2", "3", "4", "5"},
	}

	handler := handler{
		params: reflect.ValueOf(&testDecodeParams{}).Interface().(jsonrpc2.Params),
	}

	data, _ := json.Marshal(&params)

	p, err := handler.DecodeParams(data)

	if assert.NoError(t, err) && assert.True(t, p.IsValid()) {

		if ps, ok := p.(*testDecodeParams); assert.True(t, ok) {

			assert.Equal(t, ps.Int, params.Int)
			assert.Equal(t, ps.String, params.String)
			assert.Equal(t, ps.SliceInt, params.SliceInt)
			assert.Equal(t, ps.SliceString, params.SliceString)
		}
	}

	data, _ = json.Marshal(&testDecodeParams{Int: 42})

	if p, err = handler.DecodeParams(data); assert.NoError(t, err) {

		if assert.True(t, p.IsValid()) {

			data, _ = json.Marshal(&testDecodeParams{Int: 1})

			if p, err = handler.DecodeParams(data); assert.NoError(t, err) {

				assert.False(t, p.IsValid())
			}
		}
	}

	if _, err := handler.DecodeParams([]byte("message")); assert.Error(t, err) {

		assert.Contains(t, err.Error(), "invalid character")
	}
}

type testCallParams struct {
	Message string
}

func (t *testCallParams) IsValid() bool {
	return true
}

func TestHandlerCall(t *testing.T) {

	fn := func(params *jsonrpc2.EmptyParams) (interface{}, error) {

		return "EmptyParams", nil
	}

	h := handler{
		method: reflect.ValueOf(fn),
		params: reflect.New(reflect.ValueOf(fn).Type().In(0).Elem()).Interface().(jsonrpc2.Params),
	}

	data, _ := json.Marshal(&jsonrpc2.EmptyParams{})

	if p, err := h.DecodeParams(data); assert.NoError(t, err) && assert.True(t, p.IsValid()) {

		if result, err := h.Call(p); assert.NoError(t, err) {

			str, ok := result.(string)

			if assert.True(t, ok) {

				assert.Equal(t, "EmptyParams", str)
			}
		}
	}

	fn2 := func(params *testCallParams) (interface{}, error) {

		return params.Message, nil
	}

	h = handler{
		method: reflect.ValueOf(fn2),
		params: reflect.New(reflect.ValueOf(fn2).Type().In(0).Elem()).Interface().(jsonrpc2.Params),
	}

	data, _ = json.Marshal(&testCallParams{Message: "Message"})

	if p, err := h.DecodeParams(data); assert.NoError(t, err) && assert.True(t, p.IsValid()) {

		if result, err := h.Call(p); assert.NoError(t, err) {

			str, ok := result.(string)

			if assert.True(t, ok) {

				assert.Equal(t, "Message", str)
			}
		}
	}

	fn3 := func(params *testCallParams) (interface{}, error) {

		return nil, fmt.Errorf("error message")
	}

	h = handler{
		method: reflect.ValueOf(fn3),
		params: reflect.New(reflect.ValueOf(fn3).Type().In(0).Elem()).Interface().(jsonrpc2.Params),
	}

	data, _ = json.Marshal(&testCallParams{Message: "Message"})

	if p, err := h.DecodeParams(data); assert.NoError(t, err) && assert.True(t, p.IsValid()) {

		result, err := h.Call(p)

		if assert.Nil(t, result) && assert.Error(t, err) {

			assert.Equal(t, "error message", err.Error())
		}
	}

	fn4 := func(params *testCallParams) (interface{}, error) {

		return nil, nil
	}

	h = handler{
		method: reflect.ValueOf(fn4),
		params: reflect.New(reflect.ValueOf(fn4).Type().In(0).Elem()).Interface().(jsonrpc2.Params),
	}

	data, _ = json.Marshal(&testCallParams{Message: "Message"})

	if p, err := h.DecodeParams(data); assert.NoError(t, err) && assert.True(t, p.IsValid()) {

		result, err := h.Call(p)

		if assert.NoError(t, err) {

			assert.Nil(t, result)
		}
	}
}
