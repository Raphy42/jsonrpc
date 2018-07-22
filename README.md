# JSONRPC helpers
_currently in a WIP state_

# Usage
```golang
methods := jsonrpc.MethodMap{
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

runner := jsonrpc.NewRunner(methods)

// somewhere in an http handler
response := runner.Run(r.Body())
json.NewEncoder(w).Encode(response)
```
The following json request
```json
{
    "jsonrpc": "2.0",
    "id": "0xdeadbeef",
    "method": "accumulate",
    "params": [1, 2, 3, 4, 5]
}
```
Will give the following JSON response
```json
{
    "jsonrpc": "2.0",
    "id": "0xdeadbeef",
    "result": 15
}
```
Common errors are treated accordingly to `RFC 4627`

For actual usages please have a look at `runner_test.go`