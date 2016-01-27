// Package cognito provides structs for working with AWS Cognito records.
package cognito

import (
	"encoding/json"

	"github.com/apex/go-apex"
)

// Event represents a Cognito event with one or more records.
type Event struct {
	Records []*Record `json:"Records"`
}

// Change represents a single Cognito data record change.
type Change struct {
	Old string `json:"oldValue"`
	New string `json:"newValue"`
	Op  string `json:"op"`
}

// Record represents a single Cognito record.
type Record struct {
	Version        int                `json:"version"`
	EventType      string             `json:"eventType"`
	Region         string             `json:"region"`
	IdentityPoolID string             `json:"identityPoolId"`
	IdentityID     string             `json:"identityId"`
	DatasetName    string             `json:"datasetName"`
	DatasetRecords map[string]*Change `json:"datasetRecords"`
}

// Handler handles Cognito events.
type Handler interface {
	HandleCognito(*Event, *apex.Context) error
}

// HandlerFunc unmarshals Cognito events before passing control.
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

// HandleFunc handles Cognito events with callback function.
func HandleFunc(h HandlerFunc) {
	apex.Handle(h)
}

// Handle Cognito events with handler.
func Handle(h Handler) {
	HandleFunc(HandlerFunc(h.HandleCognito))
}
