package merchant_test

import (
	"github.com/auth/server/merchant"
	"testing"
)


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
