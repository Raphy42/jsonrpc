# JSONRPC helpers
Based on [jsonrpc 2.0 specification](https://www.jsonrpc.org/specification)

# Usage
A basic example
```golang
package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/raphy42/jsonrpc"
)

type Context struct {
	Accumulator int
}

func accumulate(ctx interface{}, params []byte) (interface{}, jsonrpc.Error) {
	args := make([]int, 0)
	err := json.Unmarshal(params, &args)
	if err != nil {
		return nil, jsonrpc.Errors.InvalidParams
	}
	acc := ctx.(*Context).Accumulator
	for _, arg := range args {
		acc += arg
	}
	return acc, ok
}

var (
	ok      = jsonrpc.Error{}
	methods = jsonrpc.MethodMap{
		"accumulate": accumulate,
	}
	runner = jsonrpc.NewRunner(methods)
)

func main() {
	context := Context{Accumulator: 10}

	request := []byte(`{
		"jsonrpc":"2.0",
		"id": "foobar",
		"method":"accumulate",
		"params":[1, 2, 3, 4, 5]
	}`)
	response := runner.Run(&context, bytes.NewReader(request))

	buffer, _ := json.Marshal(response)
	fmt.Printf("%s", buffer)
}

```
