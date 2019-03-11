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

// Testing 10,000 auths and validates.
func TestMerchantAuthAndValidates(t *testing.T) {
	validateStream := make(chan map[string]interface{})
	authenticationStream := make(chan map[string]interface{})

	//// 10,000 validates
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
				"token":    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTIyNDkzMzAsImlhdCI6MTU1MjI0NTczMCwiaXNzIjoiYXV0aCJ9.63yGYgMZD2OAG4WFU8gcSR1Hqsg3vk3tx88pJaWgwVQ",
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

	var validateSuccessCount int
	var authenticationSuccessCount int

TEST:
	for {
		select {
		case validateResponse := <-validateStream:
			if validateResponse["status"] == "success" {
				validateSuccessCount++
			}

		case authenticatonResponse := <-authenticationStream:
			if authenticatonResponse["status"] == "success" {
				authenticationSuccessCount++
			}
		case <-time.After(time.Second * 3):
			break TEST
		}
	}

	t.Log("Amount of validate successes: ", validateSuccessCount)
	t.Log("Amount of authentication successes: ", authenticationSuccessCount)
}
