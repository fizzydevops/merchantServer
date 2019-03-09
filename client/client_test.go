package client_test

import (
	"github.com/auth/client"
	"golang.org/x/crypto/bcrypt"
	"log"
	"testing"
	"time"
)

func TestNewMerchantClient(t *testing.T) {
	_, err := client.New()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}

func TestMerchantClient_Send(t *testing.T) {
	c, err := client.New()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	username := "rfoxinc"
	password, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)

	err = c.Send(map[string]interface{}{
		"username": username,
		"password": password,
	})

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}

func TestMerchantClient_Read(t *testing.T) {
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

// This test is to test the server see if we can send 1000 request to it
func TestMerchantClient_Read2(t *testing.T) {
	responseStream := make(chan map[string]interface{})



	// Only doing 100 don't want to open to many connections.
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
			responseStream <- response
		}()
	}

TEST:
	for {
		select {
		case val := <-responseStream:
			log.Printf("Response from server: %v", val)
		case <-time.After(time.Second * 3):
			break TEST
		}
	}
}
