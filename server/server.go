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
	ip = "0.0.0.0"
	port = "5000"
)

var (
	conn net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	response map[string]interface{}
)

func Start() {
	fmt.Println("Starting merchant server...")

	listener, err := net.Listen(protocol, ip + ":" + port)
	defer listener.Close()

	if err != nil {
		log.Fatal(map[string]interface{}{
			"status": "Error",
			"message": "Failed to start server",
			"function": "start",
			"package": "server",
			"error": err.Error(),
		})
		return
	}

	fmt.Println("Successfully started merchant server")

	requestBytes := make([]byte, 1024)

	for {

		fmt.Println("Waiting for connections...")

		conn, err = listener.Accept()
		reader = bufio.NewReader(conn)

		fmt.Println("Accepted connection from : ", conn.RemoteAddr())

		if err != nil {
			log.Println(map[string]interface{}{
				"status": "Error",
				"message": "Failed to accept..",
				"function": "Start",
				"package": "server",
				"error": err.Error(),
			})

			jsonBytes, _ := json.Marshal("{}")
			writer.Write(jsonBytes)
			writer.Flush()
		}

		writer = bufio.NewWriter(conn)

		n, _ := reader.Read(requestBytes)

		requestStr := string(requestBytes[:n])

		fmt.Printf("The incoming request... %s\n", requestStr)

		var data map[string]interface{}

		err  = json.Unmarshal(requestBytes, &data)

		if err != nil {
			response = map[string]interface{}{
				"status": "Error",
				"message": "Failed to decode json.",
				"function": "Start",
				"package": "server",
				"error": err.Error(),
			}

			log.Println(response)
			jsonBytes, _ := json.Marshal(response)
			log.Println("JSON BYTES : ", string(jsonBytes))

			writer.Write(jsonBytes)
			writer.Flush()
			continue
		}

		fmt.Println("Successfully decoded JSON: ", data)

		handler(data)

	} // End of infinite for loop

} //end of Start method.