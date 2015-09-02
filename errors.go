package jsonrpc2

const (
	ParseError     int16 = -32700
	InvalidRequest       = -32600
	MethodNotFound       = -32601
	InvalidParams        = -32602
	InternalError        = -32603
	ServerError          = -32000
)

var Errors map[int16]string
