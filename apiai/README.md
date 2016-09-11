# Apex Api.ai

Providing [Api.ai Webhook Events](https://docs.api.ai/docs/webhook) for [Apex](https://github.com/apex/go-apex)

## Features
This is intended for use with API Gateway or services that use JSON.

## Example

```go
package main

import (
	apex "github.com/apex/go-apex"
	"github.com/apex/go-apex/apiai"
)

func main() {
	apiai.HandleFunc(func(event *apiai.Event, ctx *apex.Context) (interface{}, error) {

		// Construct response message
		var message apiai.ResponseMessage
		message.Speech = "Hello!"
		message.DisplayText = "Hey, nice to meet you!"

		return message, nil
	})
}
```

---

GitHub [@edoardo849](https://github.com/edoardo849) &nbsp;&middot;&nbsp;
Twitter [@edoardo849](https://twitter.com/edoardo849) &nbsp;&middot;&nbsp;
[Medium](https://medium.com/@edoardo849)
