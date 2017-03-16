package cognitouserpools

import (
	"encoding/json"

	"github.com/apex/go-apex"
)

// CustomMessage represents a Cognito event.
type CustomMessage struct {
	CommonEvent
	Request  CustomMessageRequest  `json:"request"`
	Response CustomMessageResponse `json:"response"`
}

// CustomMessageRequest represents a Cognito user's status and the verification code.
type CustomMessageRequest struct {
	CodeParameter  string `json:"codeParameter"`
	UserAttributes struct {
		EmailVerified       string `json:"email_verified"`
		PhoneNumberVerified string `json:"phone_number_verified"`
	} `json:"userAttributes"`
}

// CustomMessageResponse represents a response to the cognito user.
type CustomMessageResponse struct {
	EmailSubject string `json:"emailSubject"`
	EmailMessage string `json:"emailMessage"`
	SmsMessage   string `json:"smsMessage"`
}

// CustomMessageHandler handles Cognito events.
type CustomMessageHandler interface {
	HandleCustomMessage(*CustomMessage, *apex.Context) error
}

// CustomMessageHandlerFunc unmarshals Cognito CustomMessage events before passing control.
type CustomMessageHandlerFunc func(*CustomMessage, *apex.Context) error

// Handle implements apex.Handler.
func (h CustomMessageHandlerFunc) Handle(data json.RawMessage, ctx *apex.Context) (interface{}, error) {
	var event CustomMessage

	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	if err := h(&event, ctx); err != nil {
		return nil, err
	}

	return event, nil
}

// CustomMessageHandleFunc handles Cognito events with callback function.
func CustomMessageHandleFunc(h CustomMessageHandlerFunc) {
	apex.Handle(h)
}

// CustomMessageHandle Cognito events with handler.
func CustomMessageHandle(h CustomMessageHandler) {
	CustomMessageHandleFunc(CustomMessageHandlerFunc(h.HandleCustomMessage))
}
