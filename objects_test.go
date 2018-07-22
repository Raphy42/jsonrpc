package jsonrpc

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestError_Error(t *testing.T) {
	a := assert.New(t)

	a.Equal(Errors.NotFound.Error(), "-32601: Method not found")
}