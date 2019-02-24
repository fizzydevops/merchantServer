package client

import (
	"log"
	"net"
	"os"
	"strconv"
)

var (
	host = os.Getenv("MERCHANTHOST")
	port = os.Getenv("MERCHANTPORT")
	protocol = "tcp"
)

type merchantClient struct {
	message string
	conn net.Conn
}

func NewMerchantClient() (*merchantClient, error) {
	connPort, err := strconv.Atoi(port)

	if err != nil {
		log.Println(map[string]interface{}{
			"status": "Error",
			"message": "Failed to convert port(str) to int",
			"package": "client",
			"function": "NewClient",
			"error": err.Error(),
		})
		return nil, err
	}

	// create a new connection struct
	conn := NewConnection(host, connPort, protocol)

	// Try's establishes connection to server
	merchantConnection, err := conn.connect()

	if err != nil {
		log.Println(map[string]interface{}{
			"status": "Error",
			"message": "Failed to connect to merchant server",
			"function": "NewClient",
			"package": "client",
			"err": err.Error(),
		})
		return nil, err
	}

	return &merchantClient{conn: merchantConnection}, nil
}

func (mc *merchantClient) SetMessage(msg string) {
	mc.message = msg
}

func (mc *merchantClient) Message() string {
	return mc.message
}

func (mc *merchantClient) SendMessage(msg string) (err error) {
	msgBytes := toBytes(msg)

	_, err = mc.conn.Write(msgBytes)

	if err != nil {
		log.Println(map[string]interface{}{
			"status": "Error",
			"message": "Failed to write to server.",
			"function": "SendMessage",
			"package": "client",
			"error": err.Error(),
		})
		return
	}

	log.Printf("Successfully sent message : %s\n", msg)

	return nil
}


func (mc *merchantClient) ReadMessage() (response string, err error) {

	var responseBytes []byte

	_, err = mc.conn.Read(responseBytes)

	if err != nil {
		log.Println(map[string]interface{}{
			"status": "Error",
			"message": "Failed to read from server",
			"function": "ReadMessage",
			"package": "client",
			"error": err.Error(),
		})
		return
	}

	response = toString(responseBytes)
	log.Printf("Successfully read from server : %s\n", response)

	return
}