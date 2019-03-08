package merchant_test

import (
	"github.com/auth/merchant"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestMerchant_InsertLogin(t *testing.T) {
	username := "test"
	password, err := bcrypt.GenerateFromPassword([]byte("testing123"),bcrypt.MinCost)

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	m := merchant.New(username, password)
	err = m.InsertLogin()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}


func TestMerchant_Authenticate(t *testing.T) {
	username := "test"

	m := merchant.New(username, []byte("testing123"))
	authenticated, err := m.Authenticate()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	} else if !authenticated {
		t.Error("Failure, passwords do not match")
		t.FailNow()
	}

}