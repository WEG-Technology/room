package main

import (
	"fmt"
	"github.com/WEG-Technology/room"
	"github.com/WEG-Technology/room/elevator"
	"os"
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

func main() {
	_ = os.Setenv("BASE_URL", "https://dummyjson.com")

	el := elevator.NewElevator("examples/yml_concurrent_example/integration.yml")
	engine := elevator.NewElevatorEngine(el).WarmUp()

	responses := engine.
		PutBodyParser("todoRoom", "addTodo", room.NewJsonBodyParser(payload)).
		PutBodyParser("todo1Room", "addTodo", room.NewJsonBodyParser(payload)).
		PutBodyParser("todo2Room", "addTodo", room.NewJsonBodyParser(payload)).
		ExecuteConcurrent("add")

	fmt.Println(engine.GetElapsedTime())

	for i, res := range responses {
		fmt.Println(fmt.Sprintf("response for %s", i))
		fmt.Println("Header", res.Header)
		fmt.Println("Ok", res.OK())
		fmt.Println("-----------------")
	}
}
