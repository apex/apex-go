// Package logs provides structs for working with AWS CloudWatch Logs records.
package logs

import (
	"bytes"
	"compress/gzip"
	"encoding/json"

	"github.com/apex/go-apex"
)

// LogEvent represents a single log event.
type LogEvent struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

// Record represents a Cloudwatch logs event with one or more records.
type Record struct {
	AWSLogs struct {
		Data []byte `json:"data"`
	} `json:"awslogs"`
}

// Event represents a single log record.
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
	event, record := new(Event), new(Record)

	if err := json.Unmarshal(data, record); err != nil {
		return nil, err
	}

	if err := decode(record, event); err != nil {
		return nil, err
	}

	if err := h(event, ctx); err != nil {
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

// decode decodes the log payload which is gzipped.
func decode(record *Record, event *Event) error {
	r, err := gzip.NewReader(bytes.NewReader(record.AWSLogs.Data))
	if err != nil {
		return err
	}

	if err = json.NewDecoder(r).Decode(&event); err != nil {
		return err
	}

	if err := r.Close(); err != nil {
		return err
	}

	return nil
}
