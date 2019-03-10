package db

import (
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
)

type db struct {
	conn     *sql.DB
	database string
}

var conn *sql.DB
var once sync.Once

// NewConnection establishes a new connection with the database provided.
func New(database string) *db {

	once.Do(func() {
		credentials, err := getDatabaseCredentials()

		if err != nil {
			panic(err.Error())
		}

		// Try to connect to database with current credentials.
		conn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp([%s]:3306)/%s", credentials["/db/merchantdb/sql/username"], credentials["/db/merchantdb/sql/password"], credentials["/db/merchantdb/sql/endpoint"], database))

		if err != nil {
			panic(err.Error())
		}

		// Ping to make sure we establish a connection
		err = conn.Ping()

		if err != nil {
			panic(err.Error())
		}
		conn.SetMaxOpenConns(20)
		conn.SetMaxIdleConns(0)
	})

	return &db{conn: conn, database: database}
}

// getDatabaseCredentials retrieves the endpoint, username, and password from aws parameter store.
func getDatabaseCredentials() (map[string]string, error) {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("", "default"),
	})

	if err != nil {
		return nil, err
	}

	svc := ssm.New(session)

	query := &ssm.GetParametersByPathInput{
		MaxResults:     aws.Int64(3),
		Path:           aws.String("/db/merchantdb/sql"),
		WithDecryption: aws.Bool(false), //Don't feel like spending money on kms
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
		err = &noCredentialsFoundError{err: "No credentials found.", database: "merchantdb"}
	}

	return results, err
}

// PrepareAndExecute a method to be used for INSERTS, UPDATES, and DELETES
func (db *db) PrepareAndExecute(query string, values []interface{}) (sql.Result, error) {
	stmt, err := db.conn.Prepare(query)
	defer stmt.Close()

	if err != nil {
		return nil, err
	}

	result, err := stmt.Exec(values...)

	if err != nil {
		return nil, err
	}

	return result, err
}

// QueryAndScan is a method to do a select query and return it as a map.
func (db *db) QueryAndScan(query string, values []interface{}) (map[string]interface{}, error) {
	rows, err := db.conn.Query(query, values...)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	cols, _ := rows.Columns()
	results := make(map[string]interface{})

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))

		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		err := rows.Scan(columnPointers...)

		if err != nil {
			return nil, err
		}

		for i, colName := range cols {
			value := columnPointers[i].(*interface{})
			results[colName] = *value
		}
	}

	return results, err
}

func (db *db) Close() {
	err := db.conn.Close()

	if err != nil {
		log.Println(map[string]interface{}{
			"status":   "error",
			"message":  "Failed to close database connection.",
			"package":  "db",
			"function": "Close",
			"error":    err.Error(),
		})
	}
}
