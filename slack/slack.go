package slack

import (
	"encoding/json"

	"github.com/apex/go-apex"
)

// Event data JSON
type Event struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	Command     string `json:"command"`
	ResponseURL string `json:"response_url"`
	TeamDomain  string `json:"team_domain"`
	TeamID      string `json:"team_id"`
	Text        string `json:"text"`
	Token       string `json:"token"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
}

// ResponseMessage JSON
type ResponseMessage struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

// Handler handles Slack Events
type Handler interface {
	HandleSlackEvent(*Event, *apex.Context) (interface{}, error)
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

// HandleFunc handles Slack Events with callback function.
func HandleFunc(h HandlerFunc) {
	apex.Handle(h)
}

// Handle Slack Events with handler.
func Handle(h Handler) {
	HandleFunc(HandlerFunc(h.HandleSlackEvent))
}
