// Package s3 provides structs for working with AWS S3 records.
package s3

import (
	"encoding/json"
	"time"

	"github.com/apex/go-apex"
)

// Event represents a S3 event with one or more records.
type Event struct {
	Records []*Record `json:"Records"`
}

// Object is an S3 object.
type Object struct {
	Key       string `json:"key"`
	Size      int    `json:"size"`
	ETag      string `json:"eTag"`
	VersionID string `json:"versionId"`
	Sequencer string `json:"sequencer"`
}

// Bucket is an S3 bucket.
type Bucket struct {
	Name          string `json:"name"`
	OwnerIdentity struct {
		PrincipalID string `json:"principalId"`
	} `json:"ownerIdentity"`
	ARN string `json:"arn"`
}

// Record represents a single S3 record.
type Record struct {
	EventVersion string    `json:"eventVersion"`
	EventSource  string    `json:"eventSource"`
	AWSRegion    string    `json:"awsRegion"`
	EventTime    time.Time `json:"eventTime"`
	EventName    string    `json:"eventName"`
	UserIdentity struct {
		PrincipalID string `json:"principalId"`
	} `json:"userIdentity"`
	RequestParameters struct {
		SourceIPAddress string `json:"sourceIPAddress"`
	} `json:"requestParameters"`
	S3 struct {
		SchemaVersion   string  `json:"s3SchemaVersion"`
		ConfigurationID string  `json:"configurationId"`
		Bucket          *Bucket `json:"bucket"`
		Object          *Object `json:"object"`
	}
}

// Handler handles S3 events.
type Handler interface {
	HandleS3(*Event, *apex.Context) error
}

// HandlerFunc unmarshals S3 events before passing control.
type HandlerFunc func(*Event, *apex.Context) error

// Handle implements apex.Handler.
func (h HandlerFunc) Handle(data json.RawMessage, ctx *apex.Context) (interface{}, error) {
	var event Event

	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	return nil, h(&event, ctx)
}

// HandleFunc handles S3 events with callback function.
func HandleFunc(h HandlerFunc) {
	apex.Handle(h)
}

// Handle S3 events with handler.
func Handle(h Handler) {
	HandleFunc(HandlerFunc(h.HandleS3))
}
