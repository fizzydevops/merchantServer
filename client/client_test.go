package client_test

import (
	"encoding/json"
	"fmt"
	"github.com/merchantServer/client"
	"testing"
)

func TestNewMerchantClient(t *testing.T) {
	_, err := client.NewMerchantClient()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	fmt.Println("Successfully connected to server...")
}

func TestMerchantClient_SendMessage(t *testing.T) {
	c, err := client.NewMerchantClient()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	err = c.SendMessage("Hello from the client")

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	fmt.Println("Successfully sent message to server...")
}

func TestMerchantClient_ReadMessage(t *testing.T) {
	c, err := client.NewMerchantClient()

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	err = c.SendMessage("Hello from the client")

	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	fmt.Println("Successfully sent message to server...")

	responseBytes, err := c.ReadMessage()


	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	fmt.Println("Successfully got a response back")

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

	fmt.Println("Sent a successful request.")
}


