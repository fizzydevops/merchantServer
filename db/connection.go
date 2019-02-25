package db

import (
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

type db struct {
	conn *sql.DB
}

func NewConnection(database string) *db {
	credentials, err := getDatabaseCredentials()

	if err !=  nil {
		panic(err.Error())
	}

	// Try to connect to database with current credentials.
	conn, err := sql.Open("mysql",fmt.Sprintf("%s:%s@/%s", credentials["/db/merchantdb/sql/endpoint"], credentials["db/merchantdb/sql/password"], database))

	if err != nil {
		panic(err.Error())
	}

	return &db{conn:conn}
}

// getDatabaseCredentials retrieves the endpoint and password from aws parameter store.
func getDatabaseCredentials() (map[string]string, error) {
	session := session.Must(session.NewSessionWithOptions(session.Options{Profile: "merchantdb"}))
	svc := ssm.New(session)

	query := &ssm.GetParametersByPathInput{
		MaxResults: aws.Int64(2),
		Path: aws.String("/db/merchantdb/sql"),
		WithDecryption: aws.Bool(true),
	}

	resp, err := svc.GetParametersByPath(query)

	if err != nil {
		return nil, err
	}

	var results = map[string]string{}

	if len(resp.Parameters) > 0 {
		for _, value := range resp.Parameters {
			results[*value.Name] = *value.Value
		}
	} else {
		// Throw some error because nothing came back.
	}

	log.Println(results)

	return results, nil
}

