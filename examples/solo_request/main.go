package main

import (
	"fmt"
	"github.com/WEG-Technology/room"
)

func main() {
	response, err := room.NewRequest("https://jsonplaceholder.typicode.com/posts/1").Send()

	if err != nil {
		panic(err)
	}

	fmt.Println("Response OK():", response.OK())
	fmt.Println("Response Header:", response.Header)
	fmt.Println("Response StatusCode:", response.StatusCode)
	fmt.Println("Response ResponseBody():", response.ResponseBody())

	fmt.Println("Request Header:", response.Request.Header)
	fmt.Println("Request Full URL:", response.Request.URI.String())
	fmt.Println("Request RequestBody():", response.RequestBody())
}
