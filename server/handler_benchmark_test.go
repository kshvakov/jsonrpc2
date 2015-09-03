package server

import (
	"encoding/json"
	"github.com/kshvakov/jsonrpc2"
	"reflect"
	"testing"
)

func BenchmarkHandlerDecodeParams(b *testing.B) {

	handler := handler{
		params: reflect.ValueOf(&testDecodeParams{}).Interface().(jsonrpc2.Params),
	}

	data, _ := json.Marshal(&jsonrpc2.EmptyParams{})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {

		handler.DecodeParams(data)
	}
}

func BenchmarkHandlerCall(b *testing.B) {

	fn := func(params *jsonrpc2.EmptyParams) (interface{}, error) {

		return nil, nil
	}

	handler := handler{
		method: reflect.ValueOf(fn),
		params: reflect.New(reflect.ValueOf(fn).Type().In(0).Elem()).Interface().(jsonrpc2.Params),
	}
	data, _ := json.Marshal(&jsonrpc2.EmptyParams{})

	params, err := handler.DecodeParams(data)

	if err != nil {

		b.Fatal(err.Error())
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {

		handler.Call(params)
	}
}
