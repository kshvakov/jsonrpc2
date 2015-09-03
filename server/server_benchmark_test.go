package server

import (
	"fmt"
	"github.com/kshvakov/jsonrpc2"
	"testing"
)

func BenchmarkRegisterFunc(b *testing.B) {

	server := New()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {

		server.RegisterFunc(fmt.Sprintf("MethodWithEmptyParams_%d", i), func(params jsonrpc2.EmptyParams) (interface{}, error) {

			return nil, nil
		})
	}
}
