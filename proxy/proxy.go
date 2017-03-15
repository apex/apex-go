package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apex/go-apex"
)

// Serve adaptes an http.Handler to the apex.Handler interface.
func Serve(h http.Handler) apex.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	return &handler{h}
}

// Handler implements the apex.Handler interface and adapts it to an
// http.Handler by converting the incoming event to an http.Request object
type handler struct {
	Handler http.Handler
}

// Handle accepts a request from the apex shim and dispatches it to an http.Handler
func (p *handler) Handle(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
	proxyEvent := &Event{}

	err := json.Unmarshal(event, proxyEvent)
	if err != nil {
		return nil, fmt.Errorf("Parse proxy event: %s", err)
	}

	req, err := buildRequest(proxyEvent, ctx)
	if err != nil {
		return nil, fmt.Errorf("Build request: %s", err)
	}

	res := &ResponseWriter{}
	p.Handler.ServeHTTP(res, req)
	res.finish()

	return &res.response, nil
}
