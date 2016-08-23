package main

import (
	"io/ioutil"
	"log"

	apex "github.com/apex/go-apex"
	apexS3 "github.com/apex/go-apex/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	apexS3.HandleFunc(func(event *apexS3.Event, ctx *apex.Context) error {
		for _, record := range event.Records {
			svc := s3.New(session.New(&aws.Config{Region: aws.String(record.AWSRegion)}))
			out, err := svc.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(record.S3.Bucket.Name),
				Key:    aws.String(record.S3.Object.Key),
			})
			if err != nil {
				log.Fatal(err)
			}
			bytes, err := ioutil.ReadAll(out.Body)
			if err != nil {
				log.Fatal(err)
			}
			log.Print(string(bytes))
		}
		return nil
	})
}
