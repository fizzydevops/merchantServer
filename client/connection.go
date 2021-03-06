package client

import (
	"net"
	"strconv"
)

type connection struct {
	ip       string
	port     int
	protocol string
}

func NewConnection(ip string, port int, protocol string) *connection {
	return &connection{
		ip:       ip,
		port:     port,
		protocol: protocol,
	}
}

func (conn *connection) SetIP(ip string) {
	conn.ip = ip
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

func (conn *connection) connect() (net.Conn, error) {
	connection, err := net.Dial(conn.protocol, net.JoinHostPort(conn.IP(), strconv.Itoa(conn.Port())))

	if err != nil {
		return nil, err
	}

	return connection, nil
}
