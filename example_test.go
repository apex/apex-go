package apex_test

import (
	"encoding/json"

	"github.com/apex/go-apex"
)

type Message struct {
	Hello string `json:"hello"`
}

// Example of a Lambda function handling arbitrary JSON input.
func Example() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		return &Message{"world"}, nil
	})
}
