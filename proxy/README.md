
# API Gateway Proxy request handling support

This package provides an Apex-compatible handler for [proxy requests](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-set-up-simple-proxy.html).

Each proxy request matches a wildcard path in AWS API Gateway, which is then converted to a standard `http.Request` before dispatching using the `http.Handler` interface.


## Usage
Any router or middleware framework supporting the `http.Handler` interface should be compatible with this adapter.

~~~ go
package main

import (
	"net/http"

	"github.com/apex/go-apex"
	"github.com/apex/go-apex/proxy"
)

func main() {
	// Adapts a single function to the http.Handler interface
	hf := http.HandlerFunc(handler)

	// Handles incoming apex requests using the proxy adapter
	apex.Handle(proxy.Serve(hf))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
~~~


## Notes
As with any Apex handler, you must make sure that your handler doesn't write anything to stdout. If your web framework logs to stdout by default, such as Martini, you need to change the logger output to use stderr.


## Differences from eawsy
This implementation reuses a large portion of the event definitions from the eawsy AWS Lambda projects:

 * https://github.com/eawsy/aws-lambda-go-event
 * https://github.com/eawsy/aws-lambda-go-net

However, it wraps a web application in an apex-compatible adapter which makes direct calls to an http.Handler instance rather than creating a fake `net.Conn` and marshalling/unmarshalling the request data.

A ResponseWriter implementation captures the response of the handler and constructs an API Gateway Proxy response.