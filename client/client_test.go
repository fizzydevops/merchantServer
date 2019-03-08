package client_test

import (
	"encoding/json"
	"github.com/auth/client"
	"golang.org/x/crypto/bcrypt"
	"log"
	"testing"
)

func TestNewMerchantClient(t *testing.T) {
	_, err := client.New()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}

func TestMerchantClient_SendMessage(t *testing.T) {
	c, err := client.New()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	username := "rfoxinc"
	password, err := bcrypt.GenerateFromPassword([]byte("password123"),bcrypt.MinCost)

	err = c.Send(map[string]interface{}{
		"username": username,
		"password": password,
	})

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}

func TestMerchantClient_ReadMessage(t *testing.T) {
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
		"username": username,
		"password": []byte("testing123"),
	})

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	// Read response from server
	responseBytes, err := c.Read()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	var response map[string]interface{}

	err = json.Unmarshal(responseBytes, &response)

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	log.Println(response)

}
