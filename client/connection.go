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

	conn.connection = connection

	return nil
}
