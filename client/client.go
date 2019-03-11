package client

import (
	"encoding/json"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	host     = os.Getenv("AUTH_HOST")
	port     = os.Getenv("AUTH_PORT")
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

func (c *client) SetResponseData(response map[string]interface{}) {
	c.response = response
}

func (c *client) Response() map[string]interface{} {
	return c.response
}

func (c *client) SetRequestData(request map[string]interface{}) {
	c.request = request
}

func (c *client) Request() map[string]interface{} {
	return c.request
}

// Send will send a request to the authentication server. If it fails we close the connection.
func (c *client) Send(data map[string]interface{}) error {
	err := validateRequest(data)

	if err != nil {
		c.conn.Close()
		return err
	}

	requestBytes, err := json.Marshal(data)

	if err != nil {
		c.conn.Close()
		return err
	}

	_, err = c.conn.Write(requestBytes)

	if err != nil {
		c.conn.Close()
		return err
	}

	return nil
}

func (c *client) Read() (map[string]interface{}, error) {
	responseBytes := make([]byte, 1024)
	defer c.conn.Close()

	len, err := c.conn.Read(responseBytes)

	if err != nil {
		return nil, err
	}

	var response map[string]interface{}

	err = json.Unmarshal(responseBytes[:len], &response)

	return response, err
}

func validateRequest(data map[string]interface{}) error {
	var reqType, username, token string
	var password []byte
	var ok bool

	// check to make sure type was sent in.
	if reqType, ok = data["type"].(string); !ok || reqType == "" {
		return &InvalidTypeError{"Type cannot be null or zero value."}
	} else if strings.ToLower(reqType) !=  "auth" && strings.ToLower(reqType) != "validate" {
		return &InvalidTypeError{"Type must be \"auth\" or \"validate\"."}
	}

	if username, ok = data["username"].(string); !ok || username == "" {
		return &InsufficientDataError{"Username cannot be null or zero value."}
	}

	if reqType == "auth" {
		if password, ok = data["password"].([]byte); !ok || password == nil {
			return &InsufficientDataError{"Password cannot be null or zero value."}
		}
	} else {
		if token, ok = data["token"].(string); !ok || token == "" {
			return &InsufficientDataError{"Token cannot be null or zero value."}
		}
	}

	return nil
}
