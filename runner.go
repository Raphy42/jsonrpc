package jsonrpc

import (
	"io"
	"encoding/json"
)

// MethodFunc defines a handler function which will executes a method with serialized a user defined context and params
// - context can be whatever you want
// - params is a JSON serialized Request.Params
// returns a Result and an Error
type MethodFunc func(context interface{}, params []byte) (interface{}, Error)

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
	return &Response{JsonRPC: "2.0", Id: id, Error: &error}
}

// NewRequest is an helper to create Request
func NewRequest(id string, method string, params interface{}) *Request {
	buffer, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	return &Request{JsonRPC: "2.0", Id: id, Method: method, Params: buffer}
}

// Run takes as input a user defined context and an io.Reader (like http.Request.Body() for example)
// returns a Response
func (runner *Runner) Run(context interface{}, body io.Reader) *Response {
	request := Request{}
	err := json.NewDecoder(body).Decode(&request)
	if err != nil || request.JsonRPC != "2.0" {
		return NewResponseWithError(request.Id, Errors.InvalidRequest)
	}

	fn, ok := runner.methods[request.Method]
	if !ok {
		return NewResponseWithError(request.Id, Errors.NotFound)
	}

	result, status := fn(context, request.Params)
	if status.Code != 0 {
		return NewResponseWithError(request.Id, status)
	}
	return NewResponse(request.Id, result)
}

// Batch takes as input a user defined context and an io.Reader containing multiple Requests
// returns an array of Responses
func (runner *Runner) Batch(context interface{}, body io.Reader) []*Response {
	requests := make([]Request, 0)
	err := json.NewDecoder(body).Decode(&requests)
	if err != nil {
		return []*Response{NewResponseWithError("", Errors.InvalidRequest)}
	}

	responses := make([]*Response, len(requests))
	for index, request := range requests {
		if request.JsonRPC != "2.0" {
			responses[index] = NewResponseWithError(request.Id, Errors.InvalidRequest)
			continue
		}

		fn, ok := runner.methods[request.Method]
		if !ok {
			responses[index] = NewResponseWithError(request.Id, Errors.NotFound)
			continue
		}

		result, status := fn(context, request.Params)
		if status.Code != 0 {
			responses[index] = NewResponseWithError(request.Id, status)
			continue
		}
		responses[index] = NewResponse(request.Id, result)
	}
	return responses
}
