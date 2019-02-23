package client

import (
	"log"
	"net"
	"strconv"
)

type connection struct {
	ip string
	port int
	protocol string
	connection net.Conn
}

func New(ip string, port int) *connection {
	return &connection{
		ip: ip,
		port: port,
	}
}

func (conn *connection) SetIP(ip string) {
	conn.ip  = ip
}

func (conn *connection) IP() string {
	return conn.ip
}

func (conn *connection) SetPort(port int) {
	conn.port = port
}

func (conn *connection) Port() int {
	return conn.port
}

func (conn *connection) SetProtocol(protocol string) {
	conn.protocol = protocol
}

func (conn *connection) Protocol() string {
	return conn.protocol
}

func (conn *connection) getConnectionString() string {
	return conn.ip + ":" + strconv.Itoa(conn.port)
}

func toBytes(msg string) []byte {
	return []byte(msg)
}

func toString (msgBytes []byte) string {
	return string(msgBytes)
}

func (conn *connection) Connect() error {
	connection, err := net.Dial(conn.protocol, conn.ip + ":" + conn.getConnectionString())

	if err != nil {
		log.Println(map[string]interface{}{
			"status": "Error",
			"message": "Failed to establish server connection",
			"function": "Connect",
			"package": "client",
			"error": err.Error(),
		})
		return err
	}

	log.Printf("Successfully created connection to : %s\n", conn.getConnectionString())
	conn.connection = connection

	return nil
}

func (conn *connection) SendMessage(msg string) (err error) {
	msgBytes := toBytes(msg)

	_, err = conn.connection.Write(msgBytes)

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


func (conn *connection) ReadMessage() (response string, err error) {

	var responseBytes []byte

	_, err = conn.connection.Read(responseBytes)

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