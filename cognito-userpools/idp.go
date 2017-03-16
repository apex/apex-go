// Package cognitouserpools provides structs for working with AWS Cognito User Pools records.
package cognitouserpools

// CommonEvent represents a Cognito event.
type CommonEvent struct {
	Version       string `json:"version"`
	Region        string `json:"region"`
	UserPoolID    string `json:"userPoolId"`
	UserName      string `json:"userName"`
	TriggerSource string `json:"triggerSource"`
	CallerContext struct {
		AwsSdkVersion string `json:"awsSdkVersion"`
		ClientID      string `json:"clientId"`
	} `json:"callerContext"`
}
