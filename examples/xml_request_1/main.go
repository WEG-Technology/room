package main

import (
	"fmt"
	"github.com/WEG-Technology/room/examples/xml_request_1/provider"
)

func main() {
	room := provider.NewRoom()

	response := room.Search()

	fmt.Println("RequestUri", response.RequestUri())
	fmt.Println("RequestHeader", response.RequestHeader())
	fmt.Println("StatusCode", response.StatusCode())
	fmt.Println("Header", response.Header())
	fmt.Println("Status", response.Status())
	fmt.Println("Ok", response.Ok())
	fmt.Println("Dto", response.Dto())
	fmt.Println("RequestMethod", response.RequestMethod())

	fmt.Println("TimeElapsed", room.GetSegment().GetElapsedTime())
}
