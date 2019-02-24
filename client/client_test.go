package client_test

import (
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


