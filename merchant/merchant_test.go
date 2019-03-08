package merchant_test

import (
	"github.com/auth/merchant"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestMerchant_InsertLogin(t *testing.T) {
	username := "test"
	password, err := bcrypt.GenerateFromPassword([]byte("testing123"), bcrypt.MinCost)

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	err = merchant.InsertLogin(username, password)

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}

func TestMerchant_Authenticate(t *testing.T) {
	username := "test"
	authenticated, err := merchant.Authenticate(username, []byte("testing123"))

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	} else if !authenticated {
		t.Error("Failure, passwords do not match")
		t.FailNow()
	}

}
