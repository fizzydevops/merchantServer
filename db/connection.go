package db

import (
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type db struct {
	conn     *sql.DB
	database string
}

// NewConnection establishes a new connection with the database provided.
func NewConnection(database string) (*db, error) {
	credentials, err := getDatabaseCredentials()

	if err != nil {
		panic(err.Error())
	}

	// Try to connect to database with current credentials.
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", credentials["/db/merchantdb/sql/endpoint"], credentials["/db/merchantdb/sql/password"], database))

	if err != nil {
		return nil, err
	}

	// Ping to make sure we establish a connection
	err = conn.Ping()

	if err != nil {
		return nil, err
	}

	return &db{conn: conn, database: database}, err
}

// getDatabaseCredentials retrieves the endpoint and password from aws parameter store.
func getDatabaseCredentials() (map[string]string, error) {
	session := session.Must(session.NewSessionWithOptions(session.Options{Profile: "merchantdb"}))
	svc := ssm.New(session)

	query := &ssm.GetParametersByPathInput{
		MaxResults:     aws.Int64(2),
		Path:           aws.String("/db/merchantdb/sql"),
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

func BuildParamsString(length int) string {

	paramStr := "("

	for i := 0; i < length; i++ {
		if i + 1 == length {
			paramStr += "? )"
		} else {
			paramStr += "?, "
		}
	}

	return paramStr
}

func (db *db) Close() {
	err := db.conn.Close()

	if  err != nil {
		log.Println(map[string]interface{}{
			"status": "error",
			"message": "Failed to close database connection.",
			"package": "db",
			"function": "Close",
			"error": err.Error(),
		})
	}
}
