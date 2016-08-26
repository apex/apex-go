package logs

import (
	"io/ioutil"
	"testing"

	"github.com/apex/go-apex"
	"github.com/stretchr/testify/assert"
)

// HandlerFunc apex.Handler assertion.
var _ apex.Handler = HandlerFunc(func(event *Event, ctx *apex.Context) error {
	return nil
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
  "awslogs": {
		"data": "H4sIAAAAAAAA/zWPwQ6CMBBEf4U0HknaIgJ6IxG56AluhpgCKzYBStpFYoz/botxT5uZt5PZNxnAGNFB+ZqAHDxyTMv0dsmKIs0z4hO1jKCdzoNtuIviZM+s2qsu12qenEHFYmgvhroVtOnV3C4CmweCwR9YoAYxODJgPKIsoUFIr5tzWmZFWYm6WZMta+baNFpOKNV4kj2CNvbqSu6/nZNqzcueMOLqvIls12pubABK+wuKwdXiYRywhHEexfut/3/S0a6ZHDsPH+D95U/1+QKw9SHdCQEAAA=="
	}
}
`)

func TestHandlerFunc_Handle(t *testing.T) {
	called := false

	fn := func(e *Event, c *apex.Context) error {
		called = true
		assert.Equal(t, 1, len(e.LogEvents))

		assert.Equal(t, "DATA_MESSAGE", e.MessageType)
		assert.Equal(t, "1234567890", e.Owner)
		assert.Equal(t, "/aws/lambda/cloudwatchtest", e.LogGroup)
		assert.Equal(t, "2016/08/24/[$LATEST]abc12345", e.LogStream)
		assert.Equal(t, 1, len(e.SubscriptionFilters))
		assert.Equal(t, "filters1", e.SubscriptionFilters[0])

		record := e.LogEvents[0]
		assert.Equal(t, "11111", record.ID)
		assert.Equal(t, "testing the message", record.Message)

		return nil
	}

	_, err := HandlerFunc(fn).Handle(event, nil)
	assert.NoError(t, err)

	assert.True(t, called, "function never called")
}
