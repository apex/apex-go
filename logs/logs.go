// Package logs provides structs for working with AWS CloudWatch Logs records.
package logs

import (
	"encoding/json"

	"github.com/apex/go-apex"
)

// LogEvent represents a single log event.
type LogEvent struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

// Event represents a Logs event with one or more records.
type Event struct {
	Owner               string      `json:"owner"`
	LogGroup            string      `json:"logGroup"`
	LogStream           string      `json:"logStream"`
	SubscriptionFilters []string    `json:"subscriptionFilters"`
	MessageType         string      `json:"messageType"`
	LogEvents           []*LogEvent `json:"logEvents"`
}

// Handler handles Logs events.
type Handler interface {
	HandleLogs(*Event, *apex.Context) error
}

// HandlerFunc unmarshals Logs events before passing control.
type HandlerFunc func(*Event, *apex.Context) error

// Handle implements apex.Handler.
func (h HandlerFunc) Handle(data json.RawMessage, ctx *apex.Context) (interface{}, error) {
	var event Event

	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	if err := h(&event, ctx); err != nil {
		return nil, err
	}

	return event, nil
}

// HandleFunc handles Logs events with callback function.
func HandleFunc(h HandlerFunc) {
	apex.Handle(h)
}

// Handle Logs events with handler.
func Handle(h Handler) {
	HandleFunc(HandlerFunc(h.HandleLogs))
}
