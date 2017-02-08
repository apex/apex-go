//
// Copyright 2017 Alsanium, SAS. or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Originally from https://github.com/eawsy/aws-lambda-go-event/blob/master/service/lambda/runtime/event/apigatewayproxyevt/definition.go
// Changes (kothar - 2017-02):
//   Relocated to go-apex/proxy

package proxy

import "encoding/json"

// Identity provides identity information about the API caller.
type Identity struct {
	// The API owner key associated with the API.
	APIKey string

	// The AWS account ID associated with the request.
	AccountID string

	// The User Agent of the API caller.
	UserAgent string

	// The source IP address of the TCP connection making the request to
	// Amazon API Gateway.
	SourceIP string

	// The Amazon Access Key associated with the request.
	AccessKey string

	// The principal identifier of the caller making the request.
	// It is same as the User and interchangeable.
	Caller string

	// The principal identifier of the user making the request.
	// It is same as the Caller and interchangeable.
	User string

	// The Amazon Resource Name (ARN) of the effective user identified after
	// authentication.
	UserARN string

	// The Amazon Cognito identity ID of the caller making the request.
	// Available only if the request was signed with Amazon Cognito credentials.
	CognitoIdentityID string

	// The Amazon Cognito identity pool ID of the caller making the request.
	// Available only if the request was signed with Amazon Cognito credentials.
	CognitoIdentityPoolID string

	// The Amazon Cognito authentication type of the caller making the request.
	// Available only if the request was signed with Amazon Cognito credentials.
	CognitoAuthenticationType string

	// The Amazon Cognito authentication provider used by the caller making the
	// request.
	// Available only if the request was signed with Amazon Cognito credentials.
	CognitoAuthenticationProvider string
}

// RequestContext provides contextual information about an Amazon API Gateway
// Proxy event.
type RequestContext struct {
	// The identifier Amazon API Gateway assigns to the API.
	APIID string

	// The identifier Amazon API Gateway assigns to the resource.
	ResourceID string

	// An automatically generated ID for the API call.
	RequestID string

	// The incoming request HTTP method name.
	// Valid values include: DELETE, GET, HEAD, OPTIONS, PATCH, POST, and PUT.
	HTTPMethod string

	// The resource path as defined in Amazon API Gateway.
	ResourcePath string

	// The AWS account ID associated with the API.
	AccountID string

	// The deployment stage of the API call (for example, Beta or Prod).
	Stage string

	// The API caller identification information.
	Identity *Identity

	// If used with Amazon Cognito, it represents the claims returned from the
	// Amazon Cognito user pool after the method caller is successfully
	// authenticated.
	// If used with Amazon API Gateway custom authorizer, it represents the
	// specified key-value pair of the context map returned from the custom
	// authorizer AWS Lambda function.
	Authorizer map[string]string `json:"-"`
}

// Event represents an Amazon API Gateway Proxy Event.
type Event struct {
	// The incoming request HTTP method name.
	// Valid values include: DELETE, GET, HEAD, OPTIONS, PATCH, POST, and PUT.
	HTTPMethod string

	// The incoming reauest HTTP headers.
	// Duplicate entries are not supported.
	Headers map[string]string

	// The resource path with raw placeholders as defined in Amazon API Gateway.
	Resource string

	// The incoming request path parameters corresponding to the resource path
	// placeholders values as defined in Resource.
	PathParameters map[string]string

	// The real path corresponding to the path parameters injected into the
	// Resource placeholders.
	Path string

	// The incoming request query string parameters.
	// Duplicate entries are not supported.
	QueryStringParameters map[string]string

	// If used with Amazon API Gateway binary support, it represents the Base64
	// encoded binary data from the client.
	// Otherwise it represents the raw data from the client.
	Body string

	// A flag to indicate if the applicable request payload is Base64 encoded.
	IsBase64Encoded bool

	// The name-value pairs defined as configuration attributes associated with
	// the deployment stage of the API.
	StageVariables map[string]string

	// The contextual information associated with the API call.
	RequestContext *RequestContext
}

// String returns the string representation.
func (e *Event) String() string {
	s, _ := json.MarshalIndent(e, "", "  ")
	return string(s)
}

// GoString returns the string representation.
func (e *Event) GoString() string {
	return e.String()
}
