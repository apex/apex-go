# Apex CloudWatch

Providing CloudWatch Events for [Apex](https://github.com/apex/go-apex)

## Features
Currently only supports the following event details:
- AutoScaling
- EC2
- Console Sign-In
- Schedule

Unknown events will have empty details

## Example

```go
package main

import (
    "log"

	"github.com/apex/go-apex"
	"github.com/apex/go-apex/cloudwatch"
)

func main() {
	cloudwatch.HandleFunc(func(evt *cloudwatch.Event, ctx *apex.Context) error {
        log.Println("Handler called")
        return nil
    }
}
```

---

GitHub [@sthulb](https://github.com/sthulb) &nbsp;&middot;&nbsp;
Twitter [@sthulb](https://twitter.com/sthulb)
