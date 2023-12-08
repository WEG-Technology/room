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
	c, e := room.NewConnection(room.WithBaseUrl("https://dummyjson.com"),
		room.WithTimeout(10*time.Second),
	)

	if e != nil {
		panic(e)
	}

	ApiConnection = c
}
