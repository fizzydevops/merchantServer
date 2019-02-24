package server

import "fmt"

func handler(data map[string]interface{}) {

	fmt.Println("We called handler, the data sent is: ", data)
}
