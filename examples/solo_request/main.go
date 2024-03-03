package main

import (
	"fmt"
	"github.com/WEG-Technology/room"
)

func main() {
	response := room.NewSoloRequest("https://dummyjson.com").Send()

	fmt.Println("RequestUri", response.RequestUri())
	fmt.Println("RequestHeader", response.RequestHeader().Properties())
	fmt.Println("StatusCode", response.StatusCode())
	fmt.Println("Header", response.Header().Properties())
	fmt.Println("Status", response.Status())
	fmt.Println("Ok", response.Ok())
	fmt.Println("Dto", response.Dto())
	fmt.Println("RequestMethod", response.RequestMethod())
}
