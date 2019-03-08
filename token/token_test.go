package token_test

import (
	"github.com/auth/token"
	"testing"
)

func TestToken_ValidateToken(t *testing.T) {
	tkn := token.New("test")
	validTkn, err := tkn.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.RH3JYcduYpAH9hPXsSKPKC6FYn39vJ90VblQgY5z0Hw")

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	} else if !validTkn {
		t.Error("Invalid token")
		t.FailNow()
	}
}
