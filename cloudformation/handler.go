// Package cloudformation provides structs for working with AWS CloudFormation custom resources.
package cloudformation

// See https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-custom-resources.html

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/apex/go-apex"
)

// Request is a request to a CloudFormation Custom Resource
type Request struct {
	RequestType        string                 `json:"RequestType"`
	ResponseURL        string                 `json:"ResponseURL"`
	StackID            string                 `json:"StackId"`
	RequestID          string                 `json:"RequestId"`
	ResourceType       string                 `json:"ResourceType"`
	LogicalResourceID  string                 `json:"LogicalResourceId"`
	ResourceProperties map[string]interface{} `json:"ResourceProperties"`
}

// Response is the response sent to the ResponseURL of the CloudFormation service
type Response struct {
	Status             string `json:"Status"`
	Reason             string `json:"Reason"`
	PhysicalResourceID string `json:"PhysicalResourceId"`
	StackID            string `json:"StackId"`
	RequestID          string `json:"RequestId"`
	LogicalResourceID  string `json:"LogicalResourceId"`
	Data               interface{}
}

func buildResponse(req Request, ctx *apex.Context) Response {
	return Response{
		RequestID:          req.RequestID,
		Status:             "SUCCESS",
		Reason:             "See the details in CloudWatch Log Stream: " + ctx.LogStreamName,
		PhysicalResourceID: ctx.LogStreamName,
		StackID:            req.StackID,
		LogicalResourceID:  req.LogicalResourceID,
	}
}

func sendResponse(resp Response, url string) error {
	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	cfnResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// TODO: Check status of return
	defer cfnResp.Body.Close()
	return nil
}

// Handler handles CloudFormation Custom Resource events.
type Handler interface {
	HandleCloudFormation(*Request, *apex.Context) (interface{}, error)
}

// HandlerFunc unmarshals CloudFormation Requests before passing control.
type HandlerFunc func(*Request, *apex.Context) (interface{}, error)

// Handle implements apex.Handler.
func (h HandlerFunc) Handle(rawReq json.RawMessage, ctx *apex.Context) (interface{}, error) {
	var req Request

	if err := json.Unmarshal(rawReq, &req); err != nil {
		return nil, err
	}

	log.Printf("Request %#v", req)

	data, err := h(&req, ctx)
	resp := buildResponse(req, ctx)
	resp.Data = data

	if err != nil {
		resp.Status = "FAILURE"
		resp.Reason = err.Error()
	}

	log.Printf("Response %#v", resp)
	log.Printf("Data %#v", data)

	return data, sendResponse(resp, req.ResponseURL)
}

// HandleFunc handles CloudFormation Custom Resource events with a callback function.
func HandleFunc(h HandlerFunc) {
	apex.Handle(h)
}

// Handle CloudFormation Custom Resource events with a handler.
func Handle(h Handler) {
	HandleFunc(HandlerFunc(h.HandleCloudFormation))
}
