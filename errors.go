package jsonrpc2

const (
	ParseError     int16 = -32700
	InvalidRequest       = -32600
	MethodNotFound       = -32601
	InvalidParams        = -32602
	InternalError        = -32603
	ServerError          = -32000
	LogicErr             = -32001
)

var Errors = map[int16]string{
	ParseError:     "Parse Error",
	InvalidRequest: "Invalid Request",
	MethodNotFound: "Method not found",
	InvalidParams:  "Invalid params",
	InternalError:  "Internal error",
	ServerError:    "Server error",
}

type LogicError struct {
	message string
}

func (l *LogicError) Error() string {

	return l.message
}
