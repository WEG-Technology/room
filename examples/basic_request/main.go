package main

import (
	"fmt"
	"github.com/WEG-Technology/room"
	"time"
)

type IndexTodoRequest struct {
	room.IRequest
}

func NewIndexTodoRequest() (room.IRequest, error) {
	r, err := room.NewGetRequest(
		room.WithEndPoint("products"),
		room.WithMethod(room.GET),
		room.WithHeader(defaultHeader()),
	)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func NewConnection() (room.IConnection, error) {
	return room.NewConnection(room.WithBaseUrl("https://dummyjson.com"),
		room.WithTimeout(10*time.Second),
	)
}

func defaultHeader() room.IHeader {
	return room.NewHeader()
}

func main() {
	r, err := NewIndexTodoRequest()

	if err != nil {
		panic(err)
	}

	c, err := NewConnection()

	if err != nil {
		panic(err)
	}

	response := c.Send(r)

	fmt.Println("RequestUri", response.RequestUri())
	fmt.Println("RequestHeader", response.RequestHeader().Properties())
	fmt.Println("StatusCode", response.StatusCode())
	fmt.Println("Header", response.Header().Properties())
	fmt.Println("Status", response.Status())
	fmt.Println("Ok", response.Ok())
	fmt.Println("Dto", response.Dto())
	fmt.Println("RequestMethod", response.RequestMethod())
}
