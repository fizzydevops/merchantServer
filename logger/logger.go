package logger

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
	"log"
	"os"
)

// Log takes in a map and encodes it to JSON, then puts the log record to the AWS Firehose stream
func Log(data map[string]interface{}) {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("", "default"),
	})

	if err != nil {
		log.Println(map[string]interface{}{
			"message": "Failed to create aws session",
			"error": err.Error(),
		})
	}

	svc := firehose.New(session)

	logBytes, err := json.Marshal(data)

	if err != nil {
		log.Println(map[string]interface{}{
			"message": "Failed to marshal data.",
			"error": err.Error(),
		})
	}

	putRecord := &firehose.PutRecordInput{
		DeliveryStreamName: aws.String(os.Getenv("AUTH_LOG_STREAM")),
		Record: &firehose.Record{
			Data: logBytes,
		},
	}

	_, err = svc.PutRecord(putRecord)

	if err != nil {
		log.Println(map[string]interface{}{
			"message": "Failed to put log record.",
			"error": err.Error(),
		})
	}
}
