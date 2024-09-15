package main

import (
	"fmt"
	"github.com/WEG-Technology/room/elevator"
)

type AddTODORequest struct {
	Todo      string `json:"todo"`
	Completed bool   `json:"completed"`
	UserId    int    `json:"userId"`
}

var payload = AddTODORequest{
	Todo:      "lorem",
	Completed: true,
	UserId:    1,
}

type ResponseSchema struct {
	Completed bool   `json:"completed"`
	Id        int    `json:"id"`
	Todo      string `json:"todo"`
	UserId    int    `json:"userId"`
}

func main() {
	el := elevator.NewElevator("examples/yml_example/integration.yml")
	engine := elevator.NewElevatorEngine(el).WarmUp()

	response, err := engine.DynamicExecute("todoRoom", "addTodo", payload)

	if err != nil {
		panic(err)
	}

	fmt.Println("Response OK: ", response.OK())
	fmt.Println("Request URI: ", response.Request.URI.String())
	fmt.Println("Response Body: ", response.ResponseBody())
	fmt.Println("Request Header: ", response.Request.Header)
	fmt.Println("Request Body: ", response.RequestBody())
}
