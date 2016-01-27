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
  "Records": [
		{
			"eventId": "whatever",
			"kinesis": {
				"data": "H4sIAKDlqFYAA91Sy07DMBC89ysqn1vkR+zYvUUQKg5cSMSFoCokVmUpjSvbAVVR/h07DaH0joTwabUzszuzcr9Y+gf0RysN2CwBwiSiLOYCIgxWZ7DR+63R3THgt43u6tyUqrlAM2dkebiW7765u86uZWndGn2pbPdmK6OOTun2XjVOGuv1LyM4Ep60dklVSWvB2HydhAffKfcyPx1lWHiX5MnuMc2yZJteOErfZet+TuznaiSpOsgJEpQgyBhkgjHBCYk5hwJiCmPBYYQRQhHBXDBEIxxhwTFGlPFp1TzNKe/LlYdwo1GAGacUQnjFm9yH1X0BZHD57KP7IxRgUwB0A0kBVgXorDQPtUeVO3nEc50PPHLCYQowgHnwsPqVjOL/Z4zhn8p4/ueL4RNf623ylAMAAA=="
			}
		}
	]
}
`)

func TestHandlerFunc_Handle(t *testing.T) {
	called := false

	fn := func(e *Event, c *apex.Context) error {
		called = true
		assert.Equal(t, 1, len(e.Records))

		record := e.Records[0]
		assert.Equal(t, "whatever", record.EventID)
		assert.Equal(t, "123456789012_CloudTrail_us-east-1", record.Logs.LogStream)
		assert.Equal(t, 3, len(record.Logs.LogEvents))

		log := record.Logs.LogEvents[0]
		assert.Equal(t, "31953106606966983378809025079804211143289615424298221568", log.ID)
		assert.Equal(t, int64(1432826855000), log.Timestamp)
		assert.Equal(t, `{"eventVersion":"1.03","userIdentity":{"type":"Root"}`, log.Message)
		return nil
	}

	_, err := HandlerFunc(fn).Handle(event, nil)
	assert.NoError(t, err)

	assert.True(t, called, "function never called")
}
