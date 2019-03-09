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
		"token":    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTIxMDc4OTIsImlhdCI6MTU1MjEwNDI5MiwiaXNzIjoiYXV0aCJ9.LtwpD5-XtlVBqVH9YCTetJLi4BFmZ6DP9vttuU9y9eM",
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

// This test is to penitrate the server see if we can send 1000 request to it
func TestMerchantClient_Read2(t *testing.T) {
	responseStream := make(chan map[string]interface{})

	for i := 0; i < 100; i++ {
		go func() {
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
				"type":     "validate",
				"username": username,
				"token":    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTIxMTE2MDcsImlhdCI6MTU1MjEwODAwNywiaXNzIjoiYXV0aCJ9.7OdUdPObJH6YPCJwodMooeW89ciSnTA4JhdoDdYr40s",
				//"password": []byte("testing123"),
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
		case <-time.After(time.Second * 10):
			break TEST
		}
	}
}
