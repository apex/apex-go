# Apex Slack

Providing Slack Events for [Apex](https://github.com/apex/go-apex)

## Features
This is intended for use with API Gateway or services that use JSON.

## Example

```go
package main

import (
	apex "github.com/apex/go-apex"
	"github.com/apex/go-apex/slack"
)

func main() {
	slack.HandleFunc(func(event *slack.Event, ctx *apex.Context) (interface{}, error) {
		// useEventData(event)

		// Construct response message
		var message slack.ResponseMessage
		message.ResponseType = "in_channel"
		message.Text = "response text"

		return message, nil
	})
}
```

---

GitHub [@DaveBlooman](https://github.com/DaveBlooman) &nbsp;&middot;&nbsp;
Twitter [@dblooman](https://twitter.com/dblooman)
