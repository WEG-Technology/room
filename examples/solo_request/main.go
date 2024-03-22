package main

import (
	"fmt"
	"github.com/WEG-Technology/room"
)

func main() {
	response, err := room.NewRequest("https://jsonplaceholder.typicode.com/posts/1").Send()

	fmt.Println(err)

	fmt.Println("Response:", response.ResponseBody())
}
