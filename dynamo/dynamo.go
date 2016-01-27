// Package dynamo provides structs for working with AWS Dynamo records.
package dynamo

import (
	"encoding/json"

	"github.com/apex/go-apex"
)

// Event represents a Dynamo event with one or more records.
type Event struct {
	Records []*Record `json:"Records"`
}

// Record represents a single Dynamo record.
type Record struct {
	EventID      string `json:"eventID"`
	EventName    string `json:"eventName"`
	EventSource  string `json:"eventSource"`
	EventVersion string `json:"eventVersion"`
	AWSRegion    string `json:"awsRegion"`
	Dynamodb     struct {
		Keys struct {
			ForumName struct {
				S string `json:"S"`
			} `json:"ForumName"`
			Subject struct {
				S string `json:"S"`
			} `json:"Subject"`
		} `json:"Keys"`
		SequenceNumber string `json:"SequenceNumber"`
		SizeBytes      int    `json:"SizeBytes"`
		StreamViewType string `json:"StreamViewType"`
	} `json:"dynamodb"`
}

// Handler handles Dynamo events.
type Handler interface {
	HandleDynamo(*Event, *apex.Context) error
}

// HandlerFunc unmarshals Dynamo events before passing control.
type HandlerFunc func(*Event, *apex.Context) error

// Handle implements apex.Handler.
func (h HandlerFunc) Handle(data json.RawMessage, ctx *apex.Context) (interface{}, error) {
	var event Event

	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	return nil, h(&event, ctx)
}

// HandleFunc handles Dynamo events with callback function.
func HandleFunc(h HandlerFunc) {
	apex.Handle(h)
}

// Handle Dynamo events with handler.
func Handle(h Handler) {
	HandleFunc(HandlerFunc(h.HandleDynamo))
}
