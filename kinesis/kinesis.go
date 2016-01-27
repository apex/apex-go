// Package kinesis provides structs for working with AWS Kinesis records.
package kinesis

import (
	"encoding/json"

	"github.com/apex/go-apex"
)

// Event represents a Kinesis event with one or more records.
type Event struct {
	Records []*Record `json:"Records"`
}

// Record represents a single Kinesis record.
type Record struct {
	Kinesis struct {
		SchemaVersion  string `json:"kinesisSchemaVersion"`
		PartitionKey   string `json:"partitionKey"`
		SequenceNumber string `json:"sequenceNumber"`
		Data           []byte `json:"data"`
	}
	EventSource       string `json:"eventSource"`
	EventVersion      string `json:"eventVersion"`
	EventID           string `json:"eventID"`
	EventName         string `json:"eventName"`
	InvokeIdentityARN string `json:"invokeIdentityArn"`
	AWSRegion         string `json:"awsRegion"`
	EventSourceARN    string `json:"eventSourceARN"`
}

// Data returns the payload.
func (r *Record) Data() []byte {
	return r.Kinesis.Data
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
