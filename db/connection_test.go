package db_test

import (
	"github.com/merchantServer/db"
	"testing"
)

func TestNewConnection(t *testing.T) {
	db.NewConnection("merchantdb")
}
