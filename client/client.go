package client

import (
	"encoding/json"
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

type client struct {
	conn     net.Conn
	request  map[string]interface{}
	response map[string]interface{}
}

func New() (*client, error) {
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

	return &client{conn: merchantConnection}, nil
}

func (c *client) SetMerchantResponse(response map[string]interface{}) {
	c.response = response
}

func (c *client) Response() map[string]interface{} {
	return c.response
}

func (c *client) SetMerchantRequest(request map[string]interface{}) {
	c.request = request
}

func (c *client) Request() map[string]interface{} {
	return c.request
}

func (c *client) Send(data map[string]interface{}) error {

	requestBytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	_, err = c.conn.Write(requestBytes)

	if err != nil {
		return err
	}

	return nil
}

func (c *client) Read() (map[string]interface{}, error) {
	responseBytes := make([]byte, 1024)

	len, err := c.conn.Read(responseBytes)

	if err != nil {
		return nil, err
	}

	var response map[string]interface{}

	err = json.Unmarshal(responseBytes[:len], &response)

	return response, err
}
