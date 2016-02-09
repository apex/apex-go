package cloudformation

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apex/go-apex"
	"github.com/stretchr/testify/assert"
)

// HandlerFunc apex.Handler assertion.
var _ apex.Handler = HandlerFunc(func(req *Request, ctx *apex.Context) (interface{}, error) {
	return nil, nil
})

var sampleRequest = Request{
	RequestType:       "Create",
	ResponseURL:       "http://pre-signed-S3-url-for-response",
	StackID:           "arn:aws:cloudformation:us-west-2:EXAMPLE/stack-name/guid",
	RequestID:         "unique id for this create request",
	ResourceType:      "Custom::TestResource",
	LogicalResourceID: "MyTestResource",
	ResourceProperties: map[string]interface{}{
		"Name": "Value",
		"List": []string{"1", "2", "3"},
	},
}

func TestResponseToResponseURL(t *testing.T) {
	responses := []Response{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var resp Response
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		responses = append(responses, resp)
	}))
	defer ts.Close()

	h := HandlerFunc(func(r *Request, ctx *apex.Context) (interface{}, error) {
		return map[string]interface{}{"TestData": "Yes"}, nil
	})

	r := sampleRequest
	r.ResponseURL = ts.URL

	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}

	_, err = h.Handle(json.RawMessage(b), &apex.Context{})
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, responses, 1, "Cloudformation server expected a response")
	assert.Equal(t, responses[0].Status, "SUCCESS")
	assert.Equal(t, responses[0].Data.(map[string]interface{})["TestData"], "Yes")
}
