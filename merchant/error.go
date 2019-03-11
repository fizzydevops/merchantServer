package merchant

import "fmt"

type UsernameNotFoundError struct {
	username string
}

func (e *UsernameNotFoundError) Error() string {
	return fmt.Sprintf("The username %s was not found.", e.username)
}
