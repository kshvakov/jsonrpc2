package server

import (
	"encoding/json"
	"github.com/kshvakov/jsonrpc2"
	"reflect"
)

type handler struct {
	method reflect.Value
	params jsonrpc2.Params
}

func (h *handler) DecodeParams(message json.RawMessage) (jsonrpc2.Params, error) {

	params := reflect.New(reflect.TypeOf(h.params).Elem()).Interface()

	if err := json.Unmarshal(message, &params); err != nil {

		return nil, err
	}

	return params.(jsonrpc2.Params), nil
}

func (h *handler) Call(params jsonrpc2.Params) (interface{}, error) {

	result := h.method.Call([]reflect.Value{reflect.ValueOf(params)})

	if result[1].IsNil() {

		if result[0].IsNil() {

			return nil, nil
		}

		return result[0].Interface(), nil
	}

	return nil, result[1].Interface().(error)
}
