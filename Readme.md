
# Apex Golang

Golang runtime support for Apex/Lambda – providing handlers for Lambda sources, and runtime requirements such as implementing the Node.js shim stdio interface.

## Features

Currently supports:

- Node.js shim
- Environment variable population
- Arbitrary JSON
- CloudWatch Logs
- Cognito
- Kinesis
- Dynamo
- S3
- SNS
- SES

## Example

```go
package main

import (
  "encoding/json"
  "strings"

  "github.com/apex/go-apex"
)

type message struct {
  Value string `json:"value"`
}

func main() {
  apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
    var m message

    if err := json.Unmarshal(event, &m); err != nil {
      return nil, err
    }

    m.Value = strings.ToUpper(m.Value)

    return m, nil
  })
}
```

Run the program:

```
echo '{"event":{"value":"Hello World!"}}' | go run main.go
{"value":{"value":"HELLO WORLD!"}}
```

## Notes

 Due to the Node.js [shim](http://apex.run/#understanding-the-shim) required to run Go in Lambda, you __must__ use stderr for logging – stdout is reserved for the shim.

## Badges

[![Build Status](https://semaphoreci.com/api/v1/projects/66c27cb2-5e00-469e-bfa0-b577cac48053/675168/badge.svg)](https://semaphoreci.com/tj/go-apex)
[![GoDoc](https://godoc.org/github.com/apex/go-apex?status.svg)](https://godoc.org/github.com/apex/go-apex)
![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)
[![](http://apex.sh/images/badge.svg)](https://apex.sh/ping/)

---

> [tjholowaychuk.com](http://tjholowaychuk.com) &nbsp;&middot;&nbsp;
> GitHub [@tj](https://github.com/tj) &nbsp;&middot;&nbsp;
> Twitter [@tjholowaychuk](https://twitter.com/tjholowaychuk)
