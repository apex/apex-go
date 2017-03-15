package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apex/go-apex"
)

// Serve adaptes an http.Handler to the apex.Handler interface
func Serve(handler http.Handler) *Handler {
	return &Handler{handler}
}

// Handler implements the apex.Handler interface and adapts it to an
// http.Handler by converting the incoming event to an http.Request object
type Handler struct {
	Handler http.Handler
}

// Handle accepts a request from the apex shim and dispatches it to an http.Handler
func (p *Handler) Handle(event json.RawMessage, ctx *apex.Context) (interface{}, error) {

	// Parse the incoming proxy event from API Gateway
	proxyEvent := &Event{}
	err := json.Unmarshal(event, proxyEvent)
	if err != nil {
		return nil, fmt.Errorf("Parse proxy event: %s", err)
	}

	// Build an http.Request from the evnt
	request, err := buildRequest(proxyEvent, ctx)
	if err != nil {
		return nil, fmt.Errorf("Build request: %s", err)
	}
	responseWriter := &ResponseWriter{}

	// Handle the request
	p.Handler.ServeHTTP(responseWriter, request)

	// Finish writing the response
	responseWriter.finish()

	return &responseWriter.response, nil
}
