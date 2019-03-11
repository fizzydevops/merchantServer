package client_test

import (
	"github.com/auth/client"
	"log"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, err := client.New()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}

func TestAuthSendAndRead(t *testing.T) {
	c, err := client.New()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

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

	log.Println(response)

}

// Testing 10,000 auths and validates.
func TestAuthAndValidates(t *testing.T) {
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
				"token":    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTIzNDkwMTQsImlhdCI6MTU1MjM0NTQxNCwiaXNzIjoiYXV0aCJ9.j-WmvILcvT3sKbn7fKdnQ-cd-4keL8tgOd6mXHDA3to",
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
