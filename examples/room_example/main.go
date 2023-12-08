package main

import (
	"fmt"
	"github.com/WEG-Technology/room/examples/room_example/client"
)

func main() {
	parameters := client.AddTODORequest{
		Todo:      "Use DummyJSON in the project",
		Completed: false,
		UserId:    1,
	}

	dummyJsonClient := client.NewRoom()

	response := dummyJsonClient.AddTODO(parameters)

	fmt.Println("RequestUri", response.RequestUri())
	fmt.Println("RequestHeader", response.RequestHeader())
	fmt.Println("StatusCode", response.StatusCode())
	fmt.Println("Header", response.Header())
	fmt.Println("Status", response.Status())
	fmt.Println("Ok", response.Ok())
	fmt.Println("Dto", response.Dto())
	fmt.Println("RequestMethod", response.RequestMethod())

	fmt.Println("TimeElapsed", dummyJsonClient.GetSegment().GetElapsedTime())
}
