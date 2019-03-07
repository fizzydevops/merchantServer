package client

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
)

var (
	host     = os.Getenv("MERCHANTHOST")
	port     = os.Getenv("MERCHANTPORT")
	protocol = "tcp"
)

func init() {
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5000"
	}
}

type merchantClient struct {
	conn     net.Conn
	request  map[string]interface{}
	response map[string]interface{}
}

func NewMerchantClient() (*merchantClient, error) {
	connPort, err := strconv.Atoi(port)

	if err != nil {
		return nil, err
	}

	// create a new connection struct
	conn := NewConnection(host, connPort, protocol)

	// Try's establishes connection to server
	merchantConnection, err := conn.connect()

	if err != nil {
		return nil, err
	}

	return &merchantClient{conn: merchantConnection}, nil
}

func (mc *merchantClient) SetMerchantResponse(response map[string]interface{}) {
	mc.response = response
}

func (mc *merchantClient) Response() map[string]interface{} {
	return mc.response
}

func (mc *merchantClient) SetMerchantRequest(request map[string]interface{}) {
	mc.request = request
}

func (mc *merchantClient) Request() map[string]interface{} {
	return mc.request
}

func (mc *merchantClient) Send(data map[string]interface{}) error {

	requestBytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	_, err = mc.conn.Write(requestBytes)

	if err != nil {
		return err
	}

	log.Printf("Successfully sent message : %s\n", requestBytes)

	return nil
}

func (mc *merchantClient) Read() ([]byte, error) {
	responseBytes := make([]byte, 1024)

	_, err := mc.conn.Read(responseBytes)

	if err != nil {
		return nil, err
	}

	return responseBytes, err
}
