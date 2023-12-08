package main

import (
	"fmt"
	"github.com/WEG-Technology/room/examples/basic_request_1/connection"
)

func main() {
	connection.InitApiConnection()

	resp := connection.GetTodos()

	fmt.Println(resp)
}
