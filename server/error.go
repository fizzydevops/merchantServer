package server

import (
	"fmt"
	"strings"
)

type InvalidAuthRequestError struct {
	missingItems []string
}

func (e *InvalidAuthRequestError) Error() string {
	return fmt.Sprintf("Insufficient data sent in request: %s", strings.Join(e.missingItems, ","))
}
