package server

import (
	"fmt"
	"strings"
)

type InvalidAuthRequest struct {
	missingItems []string
}

func (e *InvalidAuthRequest) Error() string {
	return fmt.Sprintf("Insufficient data sent in request: %s", strings.Join(e.missingItems, ", "))
}
