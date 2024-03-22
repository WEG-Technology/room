package main

import (
	"fmt"
	"github.com/WEG-Technology/room"
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

	response, err := engine.
		PutBodyParser("todoRoom", "addTodo", room.NewJsonBodyParser(payload)).
		Execute("todoRoom", "addTodo")

	fmt.Println(err)

	fmt.Println("Response OK: ", response.OK())
	fmt.Println("Response DTO: ", response.DTO(&ResponseSchema{}))
}
