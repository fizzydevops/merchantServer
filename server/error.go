package server

import (
	"fmt"
	"strings"
)

type InsufficientDataError struct {
	missingItems []string
}

func (e *InsufficientDataError) Error() string {
	return fmt.Sprintf("Insufficient data sent in request: %s", strings.Join(e.missingItems, ","))
}

type InvalidRequestTypeError struct {
	reqType string
}

func (e *InvalidRequestTypeError) Error() string {
	return fmt.Sprintf("Invalid request type. Request type: %s", e.reqType)
}
