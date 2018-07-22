package jsonrpc

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"bytes"
)

var (
	runnerRunTests = []struct {
		Input    *Request
		Expected *Response
		Message  string
	}{
		{
			Input:    NewRequest("toto", "accumulate", []int{1, 2, 3, 4, 5}),
			Expected: NewResponse("toto", 15),
			Message:  "valid request",
		},
		{
			Input:    NewRequest("tutu", "fizzbuzz", "ola"),
			Expected: NewResponseWithError("tutu", Errors.NotFound),
			Message:  "invalid method",
		},
		{
			Input:    NewRequest("tata", "accumulate", "invalid parameter"),
			Expected: NewResponseWithError("tata", Errors.InvalidParams),
			Message:  "invalid params",
		},
		{
			Input: &Request{JsonRPC:"invalid", Id: "titi"},
			Expected: NewResponseWithError("titi", Errors.InvalidRequest),
		},
	}
	commands = MethodMap{
		"accumulate": func(params []byte) (interface{}, Error) {
			args := make([]int, 0)
			err := json.Unmarshal(params, &args)
			if err != nil {
				return nil, Errors.InvalidParams
			}
			acc := 0
			for _, arg := range args {
				acc += arg
			}
			return acc, Error{}
		},
	}
)

func TestRunner_Run(t *testing.T) {
	a := assert.New(t)
	runner := NewRunner(commands)

	for _, test := range runnerRunTests {
		buffer, err := json.Marshal(test.Input)
		if err != nil {
			a.NoError(err, "failed serialisation for " + test.Message)
			continue
		}
		response := runner.Run(bytes.NewReader(buffer))
		a.Equal(test.Expected, response, test.Message)
	}
}

func TestRunner_Batch(t *testing.T) {
	a := assert.New(t)
	runner := NewRunner(commands)

	requests := make([]*Request, len(runnerRunTests))
	expectedResponses := make([]*Response, len(runnerRunTests))
	for index, test := range runnerRunTests {
		requests[index] = test.Input
		expectedResponses[index] = test.Expected
	}

	buffer, err := json.Marshal(requests)
	if err != nil {
		a.NoError(err, "failed serialisation for batch requests")
		return
	}
	responses := runner.Batch(bytes.NewReader(buffer))

	for index, response := range responses {
		a.Equal(expectedResponses[index], response, runnerRunTests[index].Message)
	}
}