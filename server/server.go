package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func Start() {
	fmt.Println("Starting merchant server...")

	ln, err := net.Listen("tcp", "0.0.0.0:5000")

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

		conn, err := ln.Accept()

		fmt.Println("Accepted connection from : ", conn.RemoteAddr())

		if err != nil {
			log.Println(map[string]interface{}{
				"status": "Error",
				"message": "Failed to accept..",
				"function": "Start",
				"package": "server",
				"error": err.Error(),
			})
			continue
		}
		

		reader := bufio.NewReader(conn)

		n, _ := reader.Read(requestBytes)

		fmt.Println("Got a message : ", string(requestBytes[:n]))
	}

}