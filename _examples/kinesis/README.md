# Kinesis Example

The example program receives a Kinesis event with an array of records and forwards the `message` from each record to a [RequestBin](http://requestb.in/) endpoint.

## Setup

Run the program locally:

```
go run main.go < event.json
```

View the the message in the RequestBin bucket:

<img alt="request_bin" src="https://cloud.githubusercontent.com/assets/739782/14099243/a0f4fde8-f538-11e5-9f5d-fd1600b832d4.png">

### Base64 Encoding

The record data on the Kinesis stream is Base64 encoded. Therefore in order to emulate a live message from Kinesis we'll manually encode the JSON payload.

Example JSON payload:

```
{"message": "Hello from Apex!"}
```

After Base64 encoding:

```
eyJtZXNzYWdlIjogIkhlbGxvIGZyb20gQXBleCEifQ==
```

The encoded data withing the `event.json` payload:

```json
{
  "event": {
    "Records": [
      {
        "eventID": "shardId-000000000000:49545115243490985018280067714973144582180062593244200961",
        "eventVersion": "1.0",
        "Kinesis": {
          "partitionKey": "partitionKey-3",
          "data": "eyJtZXNzYWdlIjogIkhlbGxvIGZyb20gQXBleCEifQ==",   // Base64 encoded data
          "kinesisSchemaVersion": "1.0",
          "sequenceNumber": "49545115243490985018280067714973144582180062593244200961"
        },
        "invokeIdentityArn": "arn:aws:iam::EXAMPLE",
        "eventName": "aws:kinesis:record",
        "eventSourceARN": "arn:aws:kinesis:EXAMPLE",
        "eventSource": "aws:kinesis",
        "awsRegion": "us-east-1"
      }
    ]
  }
}
```
