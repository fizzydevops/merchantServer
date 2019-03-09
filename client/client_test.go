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

// Testing 10,000 auths and validates.
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
