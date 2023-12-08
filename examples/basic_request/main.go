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
	r, e := room.NewGetRequest(
		room.WithEndPoint("products"),
		room.WithMethod(room.GET),
		room.WithHeader(defaultHeader()),
	)

	if e != nil {
		return nil, e
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
	r, e := NewIndexTodoRequest()

	if e != nil {
		panic(e)
	}

	c, e := NewConnection()

	if e != nil {
		panic(e)
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
