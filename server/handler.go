package server

import "fmt"

func merchantHandler(data map[string]interface{}) {
	fmt.Println("We called merchantHandler, the data sent is: ", data)
}
