package apiai_test

import (
	"io/ioutil"
	"testing"

	"github.com/apex/go-apex"
	"github.com/apex/go-apex/apiai"
	"github.com/stretchr/testify/assert"
)

// HandlerFunc apex.Handler assertion.
var _ apex.Handler = apiai.HandlerFunc(func(event *apiai.Event, ctx *apex.Context) (interface{}, error) {
	return nil, nil
})

func fixture(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return b
}

var event = []byte(`
{
  "id": "9b49f2fb-fdd4-46f1-aa0d-7c4ed2caccdc",
  "timestamp": "2016-09-08T05:34:23.167Z",
  "result": {
    "source": "agent",
    "resolvedQuery": "my name is Sam and I live in Paris",
    "action": "",
    "actionIncomplete": false,
    "parameters": {
      "city": "Paris",
      "user_name": "Sam"
    },
    "contexts": [
      {
        "name": "greetings",
        "parameters": {
          "city": "Paris",
          "user_name": "Sam",
          "city.original": "Paris",
          "user_name.original": "Sam"
        },
        "lifespan": 5
      }
    ],
    "metadata": {
      "intentId": "373a354b-c15a-4a60-ac9d-a9f2aee76cb4",
      "webhookUsed": "false",
      "intentName": "greetings"
    },
    "fulfillment": {
      "speech": "Nice to meet you, Sam!"
    },
    "score": 1
  },
  "status": {
    "code": 200,
    "errorType": "success"
  },
  "sessionId": "7501656c-b86e-496f-ae03-c2c800b851ff"
}
`)

func TestHandlerFunc_Handle(t *testing.T) {
	called := false

	fn := func(e *apiai.Event, c *apex.Context) (interface{}, error) {
		called = true

		assert.Equal(t, "9b49f2fb-fdd4-46f1-aa0d-7c4ed2caccdc", e.ID)
		assert.Equal(t, "7501656c-b86e-496f-ae03-c2c800b851ff", e.SessionID)
		assert.Equal(t, "2016-09-08T05:34:23.167Z", e.Timestamp)

		r := e.Result
		assert.Equal(t, "agent", r.Source)
		assert.Equal(t, "my name is Sam and I live in Paris", r.ResolvedQuery)
		assert.Equal(t, "", r.Action)
		assert.IsType(t, true, r.ActionIncomplete)
		assert.Equal(t, false, r.ActionIncomplete)

		assert.Equal(t, "Paris", r.Parameters["city"])
		assert.Equal(t, "Sam", r.Parameters["user_name"])

		ctx := r.Contexts
		assert.NotEmpty(t, ctx)

		grCtx := ctx[0]
		assert.NotEmpty(t, ctx[0])
		assert.Equal(t, "greetings", grCtx.Name)
		assert.IsType(t, 1, grCtx.Lifespan)
		assert.Equal(t, 5, grCtx.Lifespan)

		assert.Equal(t, "Paris", grCtx.Parameters["city"])
		assert.Equal(t, "Sam", grCtx.Parameters["user_name"])
		assert.Equal(t, "Paris", grCtx.Parameters["city.original"])
		assert.Equal(t, "Sam", grCtx.Parameters["user_name.original"])

		m := r.Metadata
		assert.Equal(t, "373a354b-c15a-4a60-ac9d-a9f2aee76cb4", m.IntentID)
		assert.Equal(t, "false", m.WebHookUsed)
		assert.Equal(t, "greetings", m.IntentName)

		f := r.Fulfillment
		assert.Equal(t, "Nice to meet you, Sam!", f.Speech)

		s := e.Status
		assert.IsType(t, 1, s.Code)
		assert.Equal(t, 200, s.Code)
		assert.Equal(t, "success", s.ErrorType)

		return nil, nil
	}

	_, err := apiai.HandlerFunc(fn).Handle(event, nil)
	assert.NoError(t, err)

	assert.True(t, called, "function never called")
}
