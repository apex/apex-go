
# API Gateway Proxy request handling support

This package provides an Apex-compatible handler for [proxy requests](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-set-up-simple-proxy.html).

Each proxy request matches a wildcard path in AWS API Gateway, which is then converted to a standard `http.Request` before dispatching using the `http.Handler` interface.


## Usage

Any router or middleware framework supporting the `http.Handler` interface should be compatible with this adapter.

```go
package main

import (
	"net/http"

	"github.com/apex/go-apex"
	"github.com/apex/go-apex/proxy"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/foo", foo)
	mux.HandleFunc("/bar", bar)
	apex.Handle(proxy.Serve(mux))
}

...
```


## Notes

### Stdout vs Stderr

As with any Apex handler, you must make sure that your handler doesn't write anything to stdout.
If your web framework logs to stdout by default, such as Martini, you need to change the logger output to use stderr.

### Content-Types passed through as plain text output:

Any output with a Content-Type that doesn't match one of those listed below will be Base64 encoded in the output record.
In order for this to be returned from API Gateway correctly, you will need to enable binary support and
map the content types containing binary data.

In practice you can usually map all types as binary using the `*/*` pattern for binary support if you aren't using other
API Gateway resources which conflict with this.

The text-mode Content-Type regular expressions used by default are:

	* text/.*
	* application/json
	* application/.*\+json
	* application/xml
	* application/.*\+xml

You can override this by calling `proxy.SetTextContentTypes` with a list of regular expressions matching the types that should
not be Base64 encoded.

### Output encoding
API gateway will automatically gzip-encode the output of your API, so it's not necessary to gzip the output of your webapp.

If you use your own gzip encoding, it's likely to interfere with the Base64 output for text content types - this hasn't been tested.

## Differences from eawsy

This implementation reuses a large portion of the event definitions from the eawsy AWS Lambda projects:

 * https://github.com/eawsy/aws-lambda-go-event
 * https://github.com/eawsy/aws-lambda-go-net

However, it wraps a web application in an apex-compatible adapter which makes direct calls to an http.Handler instance
rather than creating a fake `net.Conn` and marshalling/unmarshalling the request data.

A ResponseWriter implementation captures the response of the handler and constructs an API Gateway Proxy response.
