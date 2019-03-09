package server_test

import (
	"github.com/auth/client"
	"github.com/auth/server"
	"log"
	"testing"
	"time"
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

// This test is to test the server see if we can send 1000 request to it
func TestMerchantClient_Read2(t *testing.T) {
	validateStream := make(chan map[string]interface{})
	authenticationStream := make(chan map[string]interface{})

	// 10,000 validates
	for i := 0; i < 10000; i++ {
		c, err := client.New()

		if err != nil {
			t.Error(err.Error())
			t.FailNow()
		}

		go func() {
			username := "test"

			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}

			err = c.Send(map[string]interface{}{
				"type":     "validate",
				"username": username,
				"token":    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTIxNTcxNzIsImlhdCI6MTU1MjE1MzU3MiwiaXNzIjoiYXV0aCJ9.cc5PZuiISimoAOnNWigFRWwVVKCiFCQ_j6WrJC4TK9I",
			})

			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			// Read response from server
			response, err := c.Read()

			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			validateStream <- response
		}()
	}

	// 10,000 authentication
	for i := 0; i < 10000; i++ {
		c, err := client.New()

		if err != nil {
			t.Error(err.Error())
			t.FailNow()
		}

		go func() {
			username := "test"

			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}

			err = c.Send(map[string]interface{}{
				"type":     "auth",
				"username": username,
				"password": []byte("testing123"),
			})

			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			// Read response from server
			response, err := c.Read()

			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			authenticationStream <- response
		}()
	}

TEST:
	for {
		select {
		case validateResponse := <-validateStream:
			log.Printf("Validation response from server: %v", validateResponse)

		case authenticatonResponse := <-authenticationStream:
			log.Printf("Validation response from server: %v", authenticatonResponse)
		case <-time.After(time.Second * 3):
			break TEST
		}
	}
}
