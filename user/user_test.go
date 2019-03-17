package user_test

import (
	"auth/user"
	"testing"
)


func TestMerchant_Authenticate(t *testing.T) {
	username := "test"
	authenticated, err := user.Authenticate(username, []byte("testing123"))

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	} else if !authenticated {
		t.Error("Failure, passwords do not match")
		t.FailNow()
	}

}
