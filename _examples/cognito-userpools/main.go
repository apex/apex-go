package main

import (
	apex "github.com/apex/go-apex"
	idp "github.com/apex/go-apex/cognito-userpools"
)

func main() {
	idp.CustomMessageHandleFunc(func(event *idp.CustomMessage, ctx *apex.Context) error {
		code := event.Request.CodeParameter

		switch event.TriggerSource {
		case "CustomMessage_SignUp":
			event.Response.EmailSubject = "Welcome to the service"
			event.Response.EmailMessage = "Thank you for signing up. Your verification code is " + code
			event.Response.SmsMessage = "Welcome to the service. Your verification code is " + code
			break
		case "CustomMessage_ForgotPassword":
			event.Response.EmailSubject = "Verification code"
			event.Response.EmailMessage = "Your verification code is " + code
			event.Response.SmsMessage = "Your verification code is " + code
			break
		}
		return nil
	})
}
