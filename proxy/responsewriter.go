package proxy

import (
	"bytes"
	"net/http"
)

// Response defines parameters for a well formed response AWS Lambda should
// return to Amazon API Gateway.
// Originally from https://github.com/eawsy/aws-lambda-go-net/blob/master/service/lambda/runtime/net/apigatewayproxy/server.go
type Response struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers,omitempty"`
	Body            string            `json:"body,omitempty"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

// ResponseWriter implements the http.ResponseWriter interface and
// collects the results of an HTTP request in an API Gateway proxy
// response object.
type ResponseWriter struct {
	response       Response
	output         bytes.Buffer
	headers        http.Header
	headersWritten bool
}

// Header returns the header map that will be sent by
// WriteHeader. Changing the header after a call to
// WriteHeader (or Write) has no effect unless the modified
// headers were declared as trailers by setting the
// "Trailer" header before the call to WriteHeader (see example).
// To suppress implicit response headers, set their value to nil.
func (w *ResponseWriter) Header() http.Header {
	if w.headers == nil {
		w.headers = make(http.Header)
	}
	return w.headers
}

// Write writes the data to the connection as part of an HTTP reply.
//
// If WriteHeader has not yet been called, Write calls
// WriteHeader(http.StatusOK) before writing the data. If the Header
// does not contain a Content-Type line, Write adds a Content-Type set
// to the result of passing the initial 512 bytes of written data to
// DetectContentType.
//
// Depending on the HTTP protocol version and the client, calling
// Write or WriteHeader may prevent future reads on the
// Request.Body. For HTTP/1.x requests, handlers should read any
// needed request body data before writing the response. Once the
// headers have been flushed (due to either an explicit Flusher.Flush
// call or writing enough data to trigger a flush), the request body
// may be unavailable. For HTTP/2 requests, the Go HTTP server permits
// handlers to continue to read the request body while concurrently
// writing the response. However, such behavior may not be supported
// by all HTTP/2 clients. Handlers should read before writing if
// possible to maximize compatibility.
func (w *ResponseWriter) Write(bs []byte) (int, error) {
	if !w.headersWritten {
		w.WriteHeader(http.StatusOK)
	}
	return w.output.Write(bs)
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (w *ResponseWriter) WriteHeader(status int) {
	if w.headersWritten {
		return
	}

	w.response.StatusCode = status

	finalHeaders := make(map[string]string)
	for k, v := range w.headers {
		if len(v) > 0 {
			finalHeaders[k] = v[len(v)-1]
		}
	}

	if value, ok := finalHeaders["Content-Type"]; !ok || value == "" {
		finalHeaders["Content-Type"] = "text/html"
	}
	w.response.Headers = finalHeaders

	w.headersWritten = true
}

// finish writes the accumulated output to the response.Body
func (w *ResponseWriter) finish() {
	w.response.Body = w.output.String()
}
