// Package sns provides structs for working with AWS SNS records.
package sns

import (
	"encoding/json"
	"time"

	"github.com/apex/go-apex"
)

// Event represents a SNS event. It is safe to assume a single
// record will be present, as AWS will not send more than one.
type Event struct {
	Records []*Record `json:"Records"`
}

// Record represents a single SNS record.
type Record struct {
	EventSource  string `json:"EventSource"`
	EventVersion string `json:"EventVersion"`
	SNS          struct {
		Type              string                 `json:"Type"`
		MessageID         string                 `json:"MessageID"`
		TopicARN          string                 `json:"TopicArn"`
		Subject           string                 `json:"Subject"`
		Message           string                 `json:"Message"`
		Timestamp         time.Time              `json:"Timestamp"`
		SignatureVersion  string                 `json:"SignatureVersion"`
		Signature         string                 `json:"Signature"`
		SignatureCertURL  string                 `json:"SignatureCertURL"`
		UnsubscribeURL    string                 `json:"UnsubscribeURL"`
		MessageAttributes map[string]interface{} `json:"MessageAttributes"`
	} `json:"Sns"`
}

// Handler handles SNS events.
type Handler interface {
	HandleSNS(*Event, *apex.Context) error
}

// HandlerFunc unmarshals SNS events before passing control.
type HandlerFunc func(*Event, *apex.Context) error

// Handle implements apex.Handler.
func (h HandlerFunc) Handle(data json.RawMessage, ctx *apex.Context) (interface{}, error) {
	var event Event

	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	return nil, h(&event, ctx)
}

// HandleFunc handles SNS events with callback function.
func HandleFunc(h HandlerFunc) {
	apex.Handle(h)
}

// Handle SNS events with handler.
func Handle(h Handler) {
	HandleFunc(HandlerFunc(h.HandleSNS))
}
