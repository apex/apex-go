package s3

import (
	"testing"
    
    "github.com/apex/go-apex/apiai"
	"github.com/apex/go-apex"
	"github.com/stretchr/testify/assert"
)

// HandlerFunc apex.Handler assertion.
var _ apex.Handler = apiai.HandlerFunc(func(event *apiai.Event, ctx *apex.Context) error {
	return nil
})

func Test(t *testing.T) {
	assert.Nil(t, nil)
	// TODO: unmarshalling test
}
