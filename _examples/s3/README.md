# S3 Example

The example program receives a S3 event and prints a specific object's value to stdout.

## Setup

Edit event.json file to reference your real S3 object. (awsRegion, bucket.name & object.key)
Run the program locally:

```
go run main.go < event.json
```

## Read S3 Object

In the example the [AWS Go SDK](http://aws.amazon.com/sdk-for-go/) is used to read the S3 Object.
