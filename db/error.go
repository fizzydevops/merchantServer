package db

import "fmt"

type noCredentialsFoundError struct {
	err      string
	database string
}

func (e *noCredentialsFoundError) Error() string {
	return fmt.Sprintf(e.err+"\n database: %s", e.database)
}
