// Package logs provides structs for working with AWS CloudWatch Logs records.
package logs

import (
	"bytes"
	"compress/gzip"
	"encoding/json"

	"github.com/apex/go-apex"
	"github.com/apex/go-apex/kinesis"
)

// LogEvent represents a single log event.
type LogEvent struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

// Event represents a Kinesis event with one or more records.
type Event struct {
	Records []*Record `json:"Records"`
}

// Record represents a single Kinesis record.
type Record struct {
	kinesis.Record
	Logs struct {
		Owner               string      `json:"owner"`
		LogGroup            string      `json:"logGroup"`
		LogStream           string      `json:"logStream"`
		SubscriptionFilters []string    `json:"subscriptionFilters"`
		MessageType         string      `json:"messageType"`
		LogEvents           []*LogEvent `json:"logEvents"`
	}
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

	for _, record := range event.Records {
		r, err := gzip.NewReader(bytes.NewReader(record.Data()))
		if err != nil {
			return nil, err
		}

		err = json.NewDecoder(r).Decode(&record.Logs)
		if err != nil {
			return nil, err
		}

		r.Close()
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
