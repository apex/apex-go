package cloudwatch

import (
	"encoding/json"
	"time"

	"github.com/apex/go-apex"
)

// Event represents a CloudWatch Event
type Event struct {
	ID         string          `json:"id"`
	DetailType string          `json:"detail-type"`
	Source     string          `json:"source"`
	Account    string          `json:"account"`
	Time       time.Time       `json:"time"`
	Region     string          `json:"region"`
	Resources  []string        `json:"resources"`
	Detail     json.RawMessage `json:"detail"`
}

// AutoScalingGroupDetail of the triggered event
type AutoScalingGroupDetail struct {
	ActivityID           string            `json:"ActivityId"`
	AutoScalingGroupName string            `json:"AutoScalingGroupName"`
	Cause                string            `json:"Cause"`
	Details              map[string]string `json:"Details"`
	EC2InstanceID        string            `json:"EC2InstanceId"`
	RequestID            string            `json:"RequestId"`
	StatusCode           string            `json:"StatusCode"`

	StartTime time.Time `json:"StartTime"`
	EndTime   time.Time `json:"EndTime"`
}

// EC2Detail of the triggered event
type EC2Detail struct {
	InstanceID string `json:"instance-id"`
	State      string `json:"state"`
}

// APIDetail of the triggered event
// This is useful for API or Console events
type APIDetail struct {
	EventID      string    `json:"eventID"`
	EventName    string    `json:"eventName"`
	EventSource  string    `json:"eventSource"`
	EventTime    time.Time `json:"eventTime"`
	EventType    string    `json:"eventType"`
	EventVersion string    `json:"eventVersion"`

	AWSRegion           string            `json:"awsRegion"`
	AdditionalEventData map[string]string `json:"additionalEventData,omitempty"`
	RequestParams       interface{}       `json:"requestParameters"`
	ResponseElements    map[string]string `json:"responseElements,omitempty"`
	SourceIPAddress     string            `json:"sourceIPAddress"`
	UserAgent           string            `json:"userAgent"`
	UserIdentity        UserIdentity      `json:"userIdentity,omitempty"`
}

type UserIdentity struct {
	Type           string            `json:"type,omitempty"`
	PrincipleID    string            `json:"principalId,omitempty"`
	ARN            string            `json:"arn,omitempty"`
	AccountID      string            `json:"accountId,omitempty"`
	SessionContext map[string]string `json:"sessionContext,omitempty"`
}

// Handler handles CloudWatch Events
type Handler interface {
	HandleCloudWatcEvent(*Event, *apex.Context) error
}

// HandlerFunc unmarshals CloudWatch Events before passing control.
type HandlerFunc func(*Event, *apex.Context) error

// Handle implements apex.Handler.
func (h HandlerFunc) Handle(data json.RawMessage, ctx *apex.Context) (interface{}, error) {
	var event Event

	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	return nil, h(&event, ctx)
}

// HandleFunc handles CloudWatch Events with callback function.
func HandleFunc(h HandlerFunc) {
	apex.Handle(h)
}

// Handle CloudWatch Events with handler.
func Handle(h Handler) {
	HandleFunc(HandlerFunc(h.HandleCloudWatcEvent))
}
