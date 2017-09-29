package apex

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
	"testing"

	"github.com/tj/assert"
)

var eventInput = `{
	"event": {
		"foo": "bar"
	},
	"context": {
		"invokeid": "test"
	}
}`

func TestHandler(t *testing.T) {
	var wg sync.WaitGroup

	n := 50
	wg.Add(n)

	var events []string
	var contexts []string
	h := HandlerFunc(func(event json.RawMessage, ctx *Context) (interface{}, error) {
		events = append(events, string(event))
		contexts = append(contexts, ctx.InvokeID)
		wg.Done()
		return nil, nil
	})

	var buf bytes.Buffer
	pr, pw := io.Pipe()

	m := &manager{
		Reader:  pr,
		Writer:  &buf,
		Handler: h,
	}

	go m.Start()

	go func() {
		for i := 0; i < n; i++ {
			pw.Write([]byte(eventInput))
		}
		pw.Close()
	}()

	wg.Wait()

	for i, e := range events {
		assert.Equal(t, "{\n\t\t\"foo\": \"bar\"\n\t}", e)
		assert.Equal(t, "test", contexts[i])
	}
}

func BenchmarkHandler(b *testing.B) {
	h := HandlerFunc(func(event json.RawMessage, ctx *Context) (interface{}, error) {
		return nil, nil
	})

	var buf bytes.Buffer
	pr, pw := io.Pipe()

	m := &manager{
		Reader:  pr,
		Writer:  &buf,
		Handler: h,
	}

	go m.Start()

	for i := 0; i < b.N; i++ {
		pw.Write([]byte(eventInput))
	}
}
