package cloudwatch

import (
	"encoding/json"
	"reflect"
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

// ConsoleSignInDetail of the triggered event
// This is useful for unknown or schedule event types
type ConsoleSignInDetail struct {
	EventID      string    `json:"eventID"`
	EventName    string    `json:"eventName"`
	EventSource  string    `json:"eventSource"`
	EventTime    time.Time `json:"eventTime"`
	EventType    string    `json:"eventType"`
	EventVersion string    `json:"eventVersion"`

	AWSRegion           string            `json:"awsRegion"`
	AdditionalEventData map[string]string `json:"additionalEventData"`
	RequestParams       interface{}       `json:"requestParameters"`
	ResponseElements    map[string]string `json:"responseElements"`
	SourceIPAddress     string            `json:"sourceIPAddress"`
	UserAgent           string            `json:"userAgent"`
	UserIdentity        map[string]string `json:"userIdentity"`
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

func listKeys(i interface{}) []string {
	v := reflect.ValueOf(i).Elem()
	t := v.Type()

	values := []string{}
	for j := 0; j < v.NumField(); j++ {
		values = append(values, t.Field(j).Name)
	}

	return values

}

func keyValue(i interface{}, key string) (interface{}, bool) {
	v := reflect.ValueOf(i).Elem()
	t := v.Type()

	for j := 0; j < v.NumField(); j++ {
		k := t.Field(j).Name
		if k == key {
			return v.Field(j).Interface(), true
		}
	}

	return nil, false
}
