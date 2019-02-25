package client

import (
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
	message  string
	conn     net.Conn
	request  []byte
	response []byte
}

func (mc *merchantClient) SetMessage(msg string) {
	mc.message = msg
}

func (mc *merchantClient) Message() string {
	return mc.message
}

func (mc *merchantClient) SetMerchantResponse(response []byte) {
	mc.response = response
}

func (mc *merchantClient) Response() []byte {
	return mc.response
}

func (mc *merchantClient) SetMerchantRequest(request []byte) {
	mc.request = request
}

func (mc *merchantClient) Request() []byte {
	return mc.request
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

func (mc *merchantClient) SendMessage(msg []byte) error {
	_, err := mc.conn.Write(msg)

	if err != nil {
		return err
	}

	log.Printf("Successfully sent message : %s\n", msg)

	return nil
}
func (mc *merchantClient) ReadMessage() ([]byte, error) {
	responseBytes := make([]byte, 1024)

	_, err := mc.conn.Read(responseBytes)

	if err != nil {
		return nil, err
	}

	return responseBytes, err
}
