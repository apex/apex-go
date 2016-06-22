package main

import (
	"fmt"

	"github.com/apex/go-apex"
	"github.com/apex/go-apex/dynamo"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type newImage struct {
	ExampleKey string
}

func main() {
	dynamo.HandleFunc(func(event *dynamo.Event, ctx *apex.Context) error {
		// Iterate all event records
		for _, record := range event.Records {

			// only act on INSERTs
			if record.EventName == "INSERT" {

				n := newImage{}
				// Unmarshal the data for NewImage
				dynamodbattribute.UnmarshalMap(record.Dynamodb.NewImage, &n)

				// Print the example attribute. (Don't do this in your function!)
				fmt.Println(n.ExampleKey)
			}
		}

		return nil
	})
}
