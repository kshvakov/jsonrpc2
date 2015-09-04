package server

import (
	"encoding/json"
	"fmt"
	"github.com/kshvakov/jsonrpc2"
	"net/http"
)

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	defer func() {

		r.Body.Close()

		if message := recover(); message != nil {

			json.NewEncoder(w).Encode(&jsonrpc2.Response{
				Jsonrpc: "2.0",
				Error: &jsonrpc2.Error{
					Code:    jsonrpc2.InternalError,
					Message: jsonrpc2.Errors[jsonrpc2.InternalError],
					Data:    fmt.Sprint(message),
				},
			})
		}
	}()

	if r.Method != "POST" {

		json.NewEncoder(w).Encode(&jsonrpc2.Response{
			Jsonrpc: "2.0",
			Error: &jsonrpc2.Error{
				Code:    jsonrpc2.InvalidRequest,
				Message: jsonrpc2.Errors[jsonrpc2.InvalidRequest],
			},
		})

		return
	}

	var request jsonrpc2.ServerRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {

		json.NewEncoder(w).Encode(&jsonrpc2.Response{
			Jsonrpc:   "2.0",
			RequestID: request.RequestID,
			Error: &jsonrpc2.Error{
				Code:    jsonrpc2.ParseError,
				Message: jsonrpc2.Errors[jsonrpc2.ParseError],
				Data:    err.Error(),
			},
		})

		return
	}

	handler, found := s.handlers[request.Method]

	if !found {

		json.NewEncoder(w).Encode(&jsonrpc2.Response{
			Jsonrpc:   "2.0",
			RequestID: request.RequestID,
			Error: &jsonrpc2.Error{
				Code:    jsonrpc2.MethodNotFound,
				Message: jsonrpc2.Errors[jsonrpc2.MethodNotFound],
			},
		})

		return
	}

	if params, err := handler.DecodeParams(request.Params); err == nil {

		if !params.IsValid() {

			json.NewEncoder(w).Encode(&jsonrpc2.Response{
				Jsonrpc:   "2.0",
				RequestID: request.RequestID,
				Error: &jsonrpc2.Error{
					Code:    jsonrpc2.InvalidParams,
					Message: jsonrpc2.Errors[jsonrpc2.InvalidParams],
				},
			})

			return
		}

		if result, err := handler.Call(params); err == nil {

			json.NewEncoder(w).Encode(&jsonrpc2.Response{
				Jsonrpc:   "2.0",
				RequestID: request.RequestID,
				Result:    result,
			})

		} else {

			json.NewEncoder(w).Encode(&jsonrpc2.Response{
				Jsonrpc:   "2.0",
				RequestID: request.RequestID,
				Error: &jsonrpc2.Error{
					Code:    jsonrpc2.LogicErr,
					Message: err.Error(),
				},
			})
		}

	} else {

		json.NewEncoder(w).Encode(&jsonrpc2.Response{
			Jsonrpc:   "2.0",
			RequestID: request.RequestID,
			Error: &jsonrpc2.Error{
				Code:    jsonrpc2.ParseError,
				Message: jsonrpc2.Errors[jsonrpc2.ParseError],
				Data:    err.Error(),
			},
		})
	}
}
