package db_test

import (
	"github.com/auth/db"
	"testing"
)

func TestNewConnection(t *testing.T) {
	_, err := db.New("merchantdb")

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}
