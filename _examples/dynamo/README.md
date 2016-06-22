# DynamoDB Example

The example program receives a DynamoDB event and prints a specific attribute's value to stdout.

## Setup

Run the program locally:

```
go run main.go < event.json
```

## Unmarshaling DynamoDB Records

DynamoDB uses a somewhat awkward notation to represent Attributes and Values in the event. It's a straight representation of the DynamoDB JSON which needs to be parsed to be useful.

In the example the [AWS Go SDK](http://aws.amazon.com/sdk-for-go/) is used to Unmarshal the Dynamo Attributes. These are stored in the struct `newImage`.

Depending on the event settings you could also Unmarshal `record.Dynamodb.OldImage` as needed.

This method of Unmarshalling is optional you may elect to handle `record.Dynamodb.NewImage` (etc) in another way.
