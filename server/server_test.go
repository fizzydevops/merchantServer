package server_test

import (
	"github.com/auth/server"
	"log"
	"testing"
)

func TestStart(t *testing.T) {
	server.Start()
}

func TestValidateToken(t *testing.T) {
	validToken, err := server.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTIwNjI2NDMsImlhdCI6MTU1MjA2MjY0MiwiaXNzIjoiYXV0aCJ9.SSQoqVyQck-5BcTsMjUnOsPZuKI_JBX8t_mmqkgdApw", "test")

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	} else if !validToken {
		t.Error("Invalid token")
		t.FailNow()
	}
}

func TestGenerateToken(t *testing.T) {
	tkn, err := server.GenerateToken("test")

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	log.Println(tkn)
}