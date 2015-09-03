package server

import (
	"fmt"
	"github.com/kshvakov/jsonrpc2"
	"reflect"
)

func New() *server {

	return &server{
		handlers: make(map[string]handler),
	}
}

type server struct {
	handlers map[string]handler
}

func (s *server) RegisterFunc(method string, fn interface{}) {

	s.addHandler(method, reflect.ValueOf(fn))
}

func (s *server) RegisterObject(name string, obj interface{}) {

	for i := 0; i < reflect.TypeOf(obj).NumMethod(); i++ {

		method := reflect.TypeOf(obj).Method(i)

		if method.PkgPath != "" {

			continue
		}

		s.addHandler(fmt.Sprintf("%s.%s", name, method.Name), reflect.ValueOf(obj).Method(i))
	}
}

func (s *server) addHandler(method string, fn reflect.Value) {

	if _, found := s.handlers[method]; !found {

		ft := fn.Type()

		if ft.Kind() == reflect.Ptr {

			ft = ft.Elem()
		}

		if ft.NumIn() != 1 || ft.NumOut() != 2 || !ft.In(0).Implements(reflect.TypeOf((*jsonrpc2.Params)(nil)).Elem()) {

			return
		}

		if !ft.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {

			return
		}

		params := ft.In(0)

		if params.Kind() == reflect.Ptr {

			params = params.Elem()
		}

		s.handlers[method] = handler{
			Method: fn,
			Params: reflect.New(params).Interface().(jsonrpc2.Params),
		}

	} else {

		panic(fmt.Sprintf("Method '%s' is exists", method))
	}
}
