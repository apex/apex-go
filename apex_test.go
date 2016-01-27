package apex

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// HandlerFunc Handler assertion.
var _ Handler = HandlerFunc(func(event json.RawMessage, ctx *Context) (interface{}, error) {
	return nil, nil
})

func Test(t *testing.T) {
	assert.Nil(t, nil)
	// TODO: shim tests
}
