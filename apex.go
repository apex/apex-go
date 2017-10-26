// Package apex provides Lambda support for Go via a
// Node.js shim and this package for operating over
// stdio.
package apex

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
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

// Identity as defined in: http://docs.aws.amazon.com/mobile/sdkforandroid/developerguide/lambda.html#identity-context
type Identity struct {
	CognitoIdentityID       string `json:"cognitoIdentityId"`
	CognitoIdentityIDPoolID string `json:"cognitoIdentityPoolId"`
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

// input from the node shim.
type input struct {
	// ID is an itentifier that is boomeranged back to the called,
	// to allow for concurrent commands
	ID      string          `json:"id,omitempty"`
	Event   json.RawMessage `json:"event"`
	Context *Context        `json:"context"`
}

// output for the node shim.
type output struct {
	// The boomeranged ID from the caller
	ID    string      `json:"id,omitempty"`
	Error interface{} `json:"error,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// manager for operating over stdio.
type manager struct {
	Reader  io.Reader
	Writer  io.Writer
	Handler Handler
}

// Error allows apex functions to return structured node errors
// http://docs.aws.amazon.com/lambda/latest/dg/nodejs-prog-mode-exceptions.html
//
// Unfortunately, AWS doesn't provide a mechanism to override the stack trace.
// Consequently, when a custom error is returned the node stack trace will be
// used which isn't particularly helpful.
//
type Error interface {
	// Error implements error
	Error() string

	// ErrorType provides a specific error type useful for switching in Step Functions
	ErrorType() string

	// ErrorMessage contains the human readable string message
	ErrorMessage() string
}

// customError holds custom error message
type customError struct {
	errorType    string
	errorMessage string
}

func (c customError) Error() string {
	return fmt.Sprintf("%v: %v", c.errorType, c.errorMessage)
}

func (c customError) ErrorType() string {
	return c.errorType
}

func (c customError) ErrorMessage() string {
	return c.errorMessage
}

// NewError creates Error message with custom error type; useful for AWS Step Functions
func NewError(errorType, errorMessage string) Error {
	return customError{
		errorType:    errorType,
		errorMessage: errorMessage,
	}
}

// Start the manager.
func (m *manager) Start() {
	dec := json.NewDecoder(m.Reader)
	enc := json.NewEncoder(m.Writer)

	for {
		var msg input
		err := dec.Decode(&msg)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("error decoding input: %s", err)
			break
		}

		v, err := m.Handler.Handle(msg.Event, msg.Context)
		out := output{ID: msg.ID, Value: v}

		if err != nil {
			if ae, ok := err.(Error); ok {
				out.Error = struct {
					ErrorType    string `json:"errorType,omitempty"`
					ErrorMessage string `json:"errorMessage,omitempty"`
				}{
					ErrorType:    ae.ErrorType(),
					ErrorMessage: ae.ErrorMessage(),
				}

			} else {
				out.Error = err.Error()
			}
		}

		if err := enc.Encode(out); err != nil {
			log.Printf("error encoding output: %s", err)
		}
	}
}
