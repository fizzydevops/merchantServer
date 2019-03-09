package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/auth/logger"
	"log"
	"net"
	"strings"
)

const (
	protocol = "tcp"
	ip       = "0.0.0.0"
	port     = "5000"
)

func Start() {
	listener, err := net.Listen(protocol, ip+":"+port)
	defer listener.Close()

	// If we fail to start server this is fatal.
	if err != nil {
		logger.Log(map[string]interface{}{
			"file":     "server.go",
			"package":  "server",
			"function": "Start",
			"message":  "Failed to start auth server.",
			"error":    err.Error(),
		})
		panic(err.Error())
		return
	}

	fmt.Printf("Successfully started auth server; Listening on %s:%s\n", ip, port)

	for {
		//Accepts connections on 0.0.0.0:5000
		conn, err := listener.Accept()

		log.Println("Got a connection from: ", conn.RemoteAddr())
		//Read incoming bytes
		reader := bufio.NewReader(conn)

		if err != nil {
			logger.Log(map[string]interface{}{
				"file":     "server.go",
				"package":  "server",
				"function": "Start",
				"message":  "Failed to accept incoming connection",
				"error":    err.Error(),
			})
			continue
		}

		writer := bufio.NewWriter(conn)

		requestBytes := make([]byte, 1024)

		n, _ := reader.Read(requestBytes)

		var data map[string]interface{}

		err = json.Unmarshal(requestBytes[:n], &data)

		if err != nil {
			logger.Log(map[string]interface{}{
				"file":     "server.go",
				"package":  "server",
				"function": "Start",
				"message":  "Failed to Unmarshal JSON.",
				"error":    err.Error(),
			})
			writeResponse(writer, map[string]interface{}{"status": "error", "message": "Failed to decode JSON"})
			continue
		}

		// Shoot of go routine.
		go merchantHandler(conn, data)
	} // End of infinite for loop

} //end of Start method.

// writeResponse is a utility function for writing back a response to the client.
func writeResponse(writer *bufio.Writer, data map[string]interface{}) {
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
		logger.Log(map[string]interface{}{
			"file":     "server.go",
			"package":  "server",
			"function": "writeResponse",
			"message":  "Failed to write response.",
			"error":    strings.Join(errMsgs, ", "),
		})
	}
}
