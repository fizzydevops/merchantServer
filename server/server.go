package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	protocol = "tcp"
	ip       = "0.0.0.0"
	port     = "5000"
)

var (
	reader   *bufio.Reader
	writer   *bufio.Writer
)

func Start() {
	fmt.Println("Starting merchantServer server...")

	listener, err := net.Listen(protocol, ip+":"+port)
	defer listener.Close()

	// If we fail to start server this is fatal.
	if err != nil {
		logServerError("Failed to start merchant server.", "Start", err.Error())
		panic(err.Error())
		return
	}

	fmt.Printf("Successfully started auth server; Listening on %s:%s\n", ip, port)

	requestBytes := make([]byte, 1024)

	for {
		//Accepts connections on 0.0.0.0:5000
		conn, err := listener.Accept()

		//Read incoming bytes
		reader = bufio.NewReader(conn)

		fmt.Println("Accepted connection from : ", conn.RemoteAddr())

		if err != nil {
			logServerError("Failed to accept incoming connection", "Start", err.Error())
		}

		writer = bufio.NewWriter(conn)

		n, _ := reader.Read(requestBytes)

		var data map[string]interface{}

		err = json.Unmarshal(requestBytes[:n], &data)

		if err != nil {
			logServerError("Failed to decode JSON", "Start", err.Error())
			writeResponse(map[string]interface{}{"status" : "error", "message": "Failed to decode JSON"})
			continue
		}

		merchantHandler(data)

	} // End of infinite for loop

} //end of Start method.


// logMerchantError is a utility function for logging errors to the server.
func logServerError(message, function, err string) {
	log.Println(map[string]string{
		"status":   "Error",
		"message":  message,
		"function": function,
		"package":  "server",
		"error":    err,
	})
}

// writeResponse is a utility function for writing back a response to the client.
func writeResponse(data map[string]interface{}) {
	var errMsgs []string
	response, err := json.Marshal(data)

	if err != nil {
		errMsgs = append(errMsgs, err.Error())
	}

	_, err = writer.Write(response)

	if err != nil {
		errMsgs = append(errMsgs, err.Error())
	}

	if err = writer.Flush(); err != nil {
		errMsgs = append(errMsgs, err.Error())
	}

	if len(errMsgs) > 0 {
		logServerError("Failed to write reponse data.", "writeResponse", strings.Join(errMsgs,", "))
	}
}
