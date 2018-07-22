package jsonrpc

import (
	"encoding/json"
	"fmt"
)

var (
	Errors = struct {
		// invalid JSON was received by the server
		Parse          Error

		// the JSON sent is not a valid Request object
		InvalidRequest Error

		// the method does not exist or is not available
		NotFound       Error

		// invalid method parameter(s)
		InvalidParams  Error

		// internal JSON-RPC error
		Internal       Error

		// implementation server error
		Server         Error
	}{
		Parse: Error{
			Code:    -32700,
			Message: "Parse error",
		},
		InvalidRequest: Error{
			Code:    -32600,
			Message: "Invalid request",
		},
		NotFound: Error{
			Code:    -32601,
			Message: "Method not found",
		},
		InvalidParams: Error{
			Code:    -32602,
			Message: "Invalid params",
		},
		Internal: Error{
			Code:    -32603,
			Message: "Internal error",
		},
		Server: Error{
			Code:    -32010,
			Message: "Server error",
		},
	}
)

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// Request represents a rpc call from a client to a server
type Request struct {
	// version of the protocol, should be equal to `"2.0"`
	JsonRPC string `json:"jsonrpc"`

	// a client established identifier
	// if not included or empty, the response is treated like a notification
	Id string `json:"id,omitempty"`

	// contains the name of the method to be invoked
	// methods prefixed by "rpc." are reserved for internals and must not be used
	// for anything else
	Method string `json:"method"`

	// a structured value that holds the parameter values to be used during the invocation of the method
	// this member may be omitted.
	Params json.RawMessage `json:"params,omitempty"`
}

// Error defines an error
type Error struct {
	// an integer indicating the error type that occurred
	Code int32 `json:"code"`

	// a single concise sentence describing the error
	Message string `json:"message"`

	// primitive or structured value containing additional information about the error
	Data interface{} `json:"data,omitempty"`
}

// Response is issued by the server when a rpc call, except for in the case of Notifications
type Response struct {
	// version of the protocol, should be equal to `"2.0"`
	JsonRPC string `json:"jsonrpc"`

	// it must be the same as the value of the id member in the Request object
	// if there was an error in detecting the id in the Request object, it must be equal to ""
	Id string `json:"id,omitempty"`

	// required on success, result of the executed command on server
	Result interface{} `json:"result,omitempty"`

	// required on error or omitted
	Error Error `json:"error,omitempty"`
}
