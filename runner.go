package jsonrpc

import (
	"io"
	"encoding/json"
)

// MethodFunc defines a handler function which will executes a method with serialized params from
// a Request
type MethodFunc func(params []byte) (interface{}, Error)

// MethodMap stores MethodFunc indexed by their method name
type MethodMap map[string]MethodFunc

// Runner contains a user defined MethodMap
type Runner struct {
	methods MethodMap
}

// NewRunner returns a runner with a user defined CommandMap
func NewRunner(methods MethodMap) *Runner {
	return &Runner{methods}
}

// NewResponse is an helper to create Response
func NewResponse(id string, result interface{}) *Response {
	return &Response{JsonRPC: "2.0", Id: id, Result: result}
}

// NewResponseWithError is an helper to create Response with an Error
func NewResponseWithError(id string, error Error) *Response {
	return &Response{JsonRPC: "2.0", Id: id, Error: error}
}

// NewRequest is an helper to create Request
func NewRequest(id string, method string, params interface{}) *Request {
	buffer, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	return &Request{JsonRPC: "2.0", Id: id, Method: method, Params: buffer}
}

// Run takes as input a io.Reader (like http.Request.Body() for example) and returns a Response
func (runner *Runner) Run(body io.Reader) *Response {
	request := &Request{}
	err := json.NewDecoder(body).Decode(&request)
	if err != nil || request.JsonRPC != "2.0" {
		return NewResponseWithError(request.Id, Errors.InvalidRequest)
	}

	fn, ok := runner.methods[request.Method]
	if !ok {
		return NewResponseWithError(request.Id, Errors.NotFound)
	}

	result, status := fn(request.Params)
	if status.Code != 0 {
		return NewResponseWithError(request.Id, status)
	}
	return NewResponse(request.Id, result)
}
