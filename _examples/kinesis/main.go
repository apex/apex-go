package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apex/go-apex"
	"github.com/apex/go-apex/kinesis"
)

const requestBinURL = "http://requestb.in/14cd4uz1"

type data struct {
	Message string `json:"message"`
}

func main() {
	kinesis.HandleFunc(func(event *kinesis.Event, ctx *apex.Context) error {
		var d data

		for _, record := range event.Records {
			// Unmarshal the message data
			err := json.Unmarshal(record.Kinesis.Data, &d)
			if err != nil {
				return err
			}

			// New request
			req, err := http.NewRequest("GET", requestBinURL, nil)
			if err != nil {
				return err
			}

			// Query params
			c := &http.Client{}
			q := req.URL.Query()
			q.Set("message", d.Message)
			req.URL.RawQuery = q.Encode()

			// Make request
			resp, err := c.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// Check response
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("%v", resp)
			}
		}

		return nil
	})
}
