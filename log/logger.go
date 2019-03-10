package log

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
	"log"
	"os"
)

type Logger struct {
	client *firehose.Firehose
}

func New() *Logger {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("", "default"),
	})

	if err != nil {
		log.Println(map[string]interface{}{
			"status":   "error",
			"file":     "log.go",
			"package":  "log",
			"function": "Log",
			"message":  "Failed to create aws session",
			"error":    err.Error(),
		})
		return nil
	}

	client := firehose.New(session)

	return &Logger{client: client}
}

// Log takes in a map and encodes it to JSON, then puts the log record to the AWS Firehose stream
func (logger *Logger) Log(data map[string]interface{}) {

	logBytes, err := json.Marshal(data)

	if err != nil {
		log.Println(map[string]interface{}{
			"status":   "error",
			"file":     "log.go",
			"package":  "log",
			"function": "Log",
			"message":  "Failed to marshal data.",
			"error":    err.Error(),
		})
		return
	}

	putRecord := &firehose.PutRecordInput{
		DeliveryStreamName: aws.String(os.Getenv("AUTH_LOG_STREAM")),
		Record: &firehose.Record{
			Data: append(logBytes, '\n'),
		},
	}

	_, err = logger.client.PutRecord(putRecord)

	// If error failing to put record we will just log the data that was sent in.
	if err != nil {
		log.Println(map[string]interface{}{
			"status":   "error",
			"file":     "log.go",
			"package":  "log",
			"function": "Log",
			"message":  "Failed to put log record.",
			"error":    err.Error(),
		})

		log.Println(data)
	}
}
