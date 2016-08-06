// Package apex provides Lambda support for Go via a
// Node.js shim and this package for operating over
// stdio.
package apex

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
)

// Handler handles Lambda events.
type Handler interface {
	Handle(json.RawMessage, *Context) (interface{}, error)
}

// HandlerFunc implements Handler.
type HandlerFunc func(json.RawMessage, *Context) (interface{}, error)

// Handle Lambda event.
func (h HandlerFunc) Handle(event json.RawMessage, ctx *Context) (interface{}, error) {
	return h(event, ctx)
}

// Context represents the context data provided by a Lambda invocation.
type Context struct {
	InvokeID                 string          `json:"invokeid"`
	RequestID                string          `json:"awsRequestId"`
	FunctionName             string          `json:"functionName"`
	FunctionVersion          string          `json:"functionVersion"`
	LogGroupName             string          `json:"logGroupName"`
	LogStreamName            string          `json:"logStreamName"`
	MemoryLimitInMB          string          `json:"memoryLimitInMB"`
	IsDefaultFunctionVersion bool            `json:"isDefaultFunctionVersion"`
	ClientContext            json.RawMessage `json:"clientContext"`
	Identity                 Identity        `json:"identity,omitempty"`
	InvokedFunctionARN       string          `json:"invokedFunctionArn"`
}

// As defined in: http://docs.aws.amazon.com/mobile/sdkforandroid/developerguide/lambda.html#identity-context
type Identity struct {
	CognitoIdentityId       string `json:"cognitoIdentityId"`
	CognitoIdentityIdPoolId string `json:"cognitoIdentityPoolId"`
}

// Handle Lambda events with the given handler.
func Handle(h Handler) {
	m := &manager{
		Reader:  os.Stdin,
		Writer:  os.Stdout,
		Handler: h,
	}

	m.Start()
}

// HandleFunc handles Lambda events with the given handler function.
func HandleFunc(h HandlerFunc) {
	Handle(h)
}

// input for the node shim.
type input struct {
	Event   json.RawMessage `json:"event"`
	Context *Context        `json:"context"`
}

// output from the node shim.
type output struct {
	Error string      `json:"error,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// manager for operating over stdio.
type manager struct {
	Reader  io.Reader
	Writer  io.Writer
	Handler Handler
}

// Start the manager.
func (m *manager) Start() {
	m.output(m.handle(m.input()))
}

// input reads from the Reader and decodes JSON messages.
func (m *manager) input() <-chan *input {
	dec := json.NewDecoder(m.Reader)
	ch := make(chan *input)

	go func() {
		defer close(ch)

		for {
			msg := new(input)
			err := dec.Decode(msg)

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Printf("error decoding input: %s", err)
				break
			}

			ch <- msg
		}
	}()

	return ch
}

// handle invokes the handler and sends the response to the output channel
func (m *manager) handle(in <-chan *input) <-chan *output {
	ch := make(chan *output)
	var wg sync.WaitGroup

	go func() {
		defer close(ch)

		for msg := range in {
			msg := msg
			wg.Add(1)

			go func() {
				defer wg.Done()
				ch <- m.invoke(msg)
			}()
		}

		wg.Wait()
	}()

	return ch
}

// invoke calls the handler with `msg`.
func (m *manager) invoke(msg *input) *output {
	v, err := m.Handler.Handle(msg.Event, msg.Context)

	if err != nil {
		return &output{Error: err.Error()}
	}

	return &output{Value: v}
}

// output encodes the JSON messages and writes to the Writer.
func (m *manager) output(ch <-chan *output) {
	enc := json.NewEncoder(m.Writer)

	for msg := range ch {
		if err := enc.Encode(msg); err != nil {
			log.Printf("error encoding output: %s", err)
		}
	}
}
