package main

import (
	"fmt"
	"github.com/WEG-Technology/room"
)

func main() {
	c := room.NewConnector("https://jsonplaceholder.typicode.com")

	response, err := c.Send("posts/1")

	fmt.Println(err)

	fmt.Println("Response:", response)
}
