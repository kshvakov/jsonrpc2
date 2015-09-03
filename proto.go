package jsonrpc2

import (
	"encoding/json"
)

type Params interface {
	IsValid() bool
}

type EmptyParams struct{}

func (_ EmptyParams) IsValid() bool {

	return true
}

type Request struct {
	Jsonrpc   string `json:"jsonrpc"`
	RequestID string `json:"id"`
	Method    string `json:"method"`
	Params    Params `json:"params"`
}

type ServerRequest struct {
	Request
	Params json.RawMessage `json:"params"`
}

type Response struct {
	Jsonrpc   string      `json:"jsonrpc"`
	RequestID string      `json:"id"`
	Result    interface{} `json:"result,omitempty"`
	Error     *Error      `json:"error,omitempty"`
}

type Error struct {
	Code    int16  `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}
