// Package sns provides structs for working with AWS SNS records.
package sns

import (
	"encoding/json"
	"time"

	"github.com/apex/go-apex"
)

// Event represents a Kinesis event with one or more records.
type Event struct {
	Records []*Record `json:"Records"`
}

// Record represents a single Kinesis record.
type Record struct {
	SNS struct {
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
	} `json:"Sns"`
	EventSource  string `json:"EventSource"`
	EventVersion string `json:"EventVersion"`
}

// Data returns the payload.
func (r *Record) Data() []byte {
	return r.SNS.Message
}

// HandlerFunc unmarshals Kinesis events before passing control.
type HandlerFunc func(*Event, *apex.Context) error

// Handle implements apex.Handler.
func (h HandlerFunc) Handle(data json.RawMessage, ctx *apex.Context) (interface{}, error) {
	var event Event

	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	return nil, h(&event, ctx)
}
