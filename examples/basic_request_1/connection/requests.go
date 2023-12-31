package connection

import (
	"github.com/WEG-Technology/room"
)

type IndexTodoRequest struct {
	room.IRequest
}

func NewIndexTodoRequest() room.IRequest {
	r, e := room.NewGetRequest(
		room.WithEndPoint("todos"),
		room.WithMethod(room.GET),
		room.WithHeader(defaultHeader()),
	)

	if e != nil {
		panic(e)
	}

	return r
}

func defaultHeader() room.IHeader {
	return room.NewHeader().
		Add("Content-Type", "application/json").
		Add("Accept", "application/json")
}
