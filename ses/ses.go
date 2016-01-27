// Package ses provides structs for working with AWS SES records.
package ses

import (
	"encoding/json"
	"time"

	"github.com/apex/go-apex"
)

// Event represents a SES event with one or more records.
type Event struct {
	Records []*Record `json:"Records"`
}

// Record represents a single SES record.
type Record struct {
	EventSource  string `json:"eventSource"`
	EventVersion string `json:"eventVersion"`
	SES          struct {
		Receipt struct {
			Type              string                 `json:"Type"`
			MessageID         string                 `json:"MessageID"`
			TopicARN          string                 `json:"TopicArn"`
			Subject           string                 `json:"Subject"`
			Message           []byte                 `json:"Message"`
			Timestamp         time.Time              `json:"Timestamp"`
			SignatureVersion  string                 `json:"SignatureVersion"`
			Signature         string                 `json:"Signature"`
			SignatureCertURL  string                 `json:"SignatureCertURL"`
			UnsubscribeURL    string                 `json:"UnsubscribeURL"`
			MessageAttributes map[string]interface{} `json:"MessageAttributes"`
		} `json:"receipt"`
		Mail struct {
			Type              string                 `json:"Type"`
			MessageID         string                 `json:"MessageID"`
			TopicARN          string                 `json:"TopicArn"`
			Subject           string                 `json:"Subject"`
			Message           []byte                 `json:"Message"`
			Timestamp         time.Time              `json:"Timestamp"`
			SignatureVersion  string                 `json:"SignatureVersion"`
			Signature         string                 `json:"Signature"`
			SignatureCertURL  string                 `json:"SignatureCertURL"`
			UnsubscribeURL    string                 `json:"UnsubscribeURL"`
			MessageAttributes map[string]interface{} `json:"MessageAttributes"`
		} `json:"mail"`
	} `json:"ses"`
}

// Handler handles SES events.
type Handler interface {
	HandleSES(*Event, *apex.Context) error
}

// HandlerFunc unmarshals SES events before passing control.
type HandlerFunc func(*Event, *apex.Context) error

// Handle implements apex.Handler.
func (h HandlerFunc) Handle(data json.RawMessage, ctx *apex.Context) (interface{}, error) {
	var event Event

	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	return nil, h(&event, ctx)
}

// HandleFunc handles SES events with callback function.
func HandleFunc(h HandlerFunc) {
	apex.Handle(h)
}

// Handle SES events with handler.
func Handle(h Handler) {
	HandleFunc(HandlerFunc(h.HandleSES))
}
