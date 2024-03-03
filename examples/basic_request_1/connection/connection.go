package connection

import (
	"github.com/WEG-Technology/room"
	"time"
)

var ApiConnection room.IConnection

type DummyJsonApiConnection struct {
	room.IConnection
}

func GetTodos() room.IResponse {
	return ApiConnection.Send(NewIndexTodoRequest())
}

func InitApiConnection() {
	c, err := room.NewConnection(room.WithBaseUrl("https://dummyjson.com"),
		room.WithTimeout(10*time.Second),
	)

	if err != nil {
		panic(err)
	}

	ApiConnection = c
}
