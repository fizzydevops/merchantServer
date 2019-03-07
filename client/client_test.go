package client_test

import (
	"encoding/json"
	"github.com/auth/client"
	"testing"
)

func TestNewMerchantClient(t *testing.T) {
	_, err := client.NewMerchantClient()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}

func TestMerchantClient_SendMessage(t *testing.T) {
	c, err := client.NewMerchantClient()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	data := map[string]interface{}{
		"name":     "Ryan Claude Fox",
		"username": "fizzyFox101",
	}

	jsonBytes, err := json.Marshal(data)

	err = c.SendMessage(jsonBytes)

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
}

func TestMerchantClient_ReadMessage(t *testing.T) {
	c, err := client.NewMerchantClient()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	data := map[string]interface{}{
		"name":     "Ryan Claude Fox",
		"username": "fizzyFox101",
	}

	jsonBytes, err := json.Marshal(data)

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	err = c.SendMessage(jsonBytes)

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	// Read response from server
	responseBytes, err := c.ReadMessage()

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

	if val := response["status"]; val == "Error" {
		t.Error("Failed to send a successful message")
		t.FailNow()
	}

}
