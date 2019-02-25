package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

const (
	protocol = "tcp"
	ip       = "0.0.0.0"
	port     = "5000"
)

var (
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	response map[string]interface{}
	listener net.Listener
)

func Start() {
	fmt.Println("Starting merchantServer server...")

	listener, err := net.Listen(protocol, ip+":"+port)
	defer listener.Close()

	// If we fail to start server this is fatal.
	if err != nil {
		logMerchantError("Failed to start merchant server.", "Start", err.Error())
		panic(err.Error())
		return
	}

	fmt.Printf("Successfully started merchantServer server; Listening on %s:%s\n", ip, port)

	requestBytes := make([]byte, 1024)

	for {
		//Accepts connections on 0.0.0.0:5000
		conn, err := listener.Accept()

		//Read incoming bytes
		reader = bufio.NewReader(conn)

		fmt.Println("Accepted connection from : ", conn.RemoteAddr())

		if err != nil {
			logMerchantError("Failed to accept incoming connection", "Start", err.Error())
		}

		writer = bufio.NewWriter(conn)

		n, _ := reader.Read(requestBytes)

		var data map[string]interface{}

		err = json.Unmarshal(requestBytes[:n], &data)

		if err != nil {
			logMerchantError("Failed to decode JSON", "Start", err.Error())
			jsonBytes, _ := json.Marshal(response)
			writer.Write(jsonBytes)
			writer.Flush()
			continue
		}

		merchantHandler(data)

	} // End of infinite for loop

} //end of Start method.

func logMerchantError(message, function, err string) {
	log.Println(map[string]string{
		"status":   "Error",
		"message":  message,
		"function": function,
		"package":  "server",
		"error":    err,
	})
}
