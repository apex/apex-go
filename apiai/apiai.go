package apiai

import (
	"encoding/json"

	"github.com/apex/go-apex"
)

// Context is...
type Context struct {
	// Context name
	Name string `json:"name"`
	// Object consisting of "parameter_name":"parameter_value" and "parameter_name.original":"original_parameter_value" pairs.
	Parameters map[string]string `json:"parameters"`
	// Number of requests after which the context will expire.
	Lifespan int `json:"lifespan"`
}

// Fullfillment is
type Fullfillment struct {
	// Text to be pronounced to the user / shown on the screen
	Speech string `json:"speech"`
	// Matching score for the intent.
	Score float64 `json:"score,omitempty"`
}

// Event data JSON
// https://docs.api.ai/docs/webhook
// https://docs.api.ai/docs/query#response
type Event struct {

	// Unique identifier of the result.
	ID string `json:"id"`

	// Date and time of the request in UTC timezone using ISO-8601 format.
	Timestamp string `json:"timestamp"`

	// Contains the results of the natual language processing.
	Result struct {
		// Source of the answer. Could be "agent" if the response was retrieved from the agent. Or "domains", if Domains functionality is enabled and the source is one of the domains
		Source string `json:"source"`

		// The query that was used to produce this result.
		ResolvedQuery string `json:"resolvedQuery"`

		// An action to take.
		Action string `json:"action"`

		// true if the triggered intent has required parameters and not all the required parameter values have been collected
		// false if all required parameter values have been collected or if the triggered intent doesn't containt any required parameters
		ActionIncomplete bool `json:"actionIncomplete"`

		// Object consisting of "parameter_name":"parameter_value" pairs.
		Parameters map[string]string `json:"parameters"`

		Contexts []*Context `json:"contexts"`

		Metadata struct {
			// ID of the intent that produced this result.
			IntentID string `json:"intentId"`
			// Name of the intent that produced this result.
			IntentName string `json:"intentName"`
			// Indicates wheather webhook functionaly is enabled in the triggered intent.
			WebHookUsed string `json:"webhookUsed"`
		} `json:"metadata,omitempty"`
		Fulfillment *Fullfillment `json:"fulfillment"`
	} `json:"result"`

	// Contains data on how the request succeeded or failed.
	// https://docs.api.ai/docs/status-object
	Status struct {
		// HTTP status code
		Code int `json:"code"`

		// Text description of error, or "success" if no error.
		ErrorType string `json:"errorType"`

		// ID of the error. Optionally returned if the request failed.
		ErrorID string `json:"errorId"`

		// Text details of the error. Only returned if the request failed.
		ErrorDetails string `json:"errorDetails"`
	} `json:"status,omitempty"`

	// Session ID.
	SessionID string `json:"sessionId,omitempty"`
}

// ResponseMessage is JSON.
// @url https://docs.api.ai/docs/webhook
type ResponseMessage struct {

	// Voice response to the request.
	Speech string `json:"speech"`

	// Text displayed on the user device screen.
	DisplayText string `json:"displayText"`

	// Additional data required for performing the action on the client side.
	// The data is sent to the client in the original form and is not processed by Api.ai.
	Data map[string]string `json:"data"`

	// Data source
	Source string `json:"source"`

	ContextOut []*Context `json:"contextOut"`
}

// Handler handles api.ai Events
type Handler interface {
	HandleApiaiEvent(*Event, *apex.Context) (interface{}, error)
}

// HandlerFunc unmarshals Slack Events before passing control.
type HandlerFunc func(*Event, *apex.Context) (interface{}, error)

// Handle implements apex.Handler.
func (h HandlerFunc) Handle(data json.RawMessage, ctx *apex.Context) (interface{}, error) {
	var event Event

	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	return h(&event, ctx)
}

// HandleFunc handles Api.ai Events with callback function.
func HandleFunc(h HandlerFunc) {
	apex.Handle(h)
}

// Handle Api.ai Events with handler.
func Handle(h Handler) {
	HandleFunc(HandlerFunc(h.HandleApiaiEvent))
}
