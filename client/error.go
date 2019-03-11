package client

import "fmt"

type InvalidTypeError struct {
	message string
}

func (e * InvalidTypeError) Error() string {
	return fmt.Sprintf("Invalid type error: %s", e.message)
}

type InsufficientDataError struct {
	message string
}

func (e *InsufficientDataError) Error() string {
	return fmt.Sprintf("Insufficient data error: %s", e.message)
}