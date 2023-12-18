package client

import (
	"github.com/WEG-Technology/room"
)

type AddTODORequest struct {
	Todo      string `json:"todo"`
	Completed bool   `json:"completed"`
	UserId    int    `json:"userId"`
}

func NewAddTODORequest(params AddTODORequest) room.IRequest {
	r, e := room.NewPostRequest(
		room.WithEndPoint("todos/add"),
		room.WithDto(&AddTODOResponse{}),
		room.WithBody(room.NewJsonBodyParser(params)),
	)

	if e != nil {
		panic(e)
	}

	return r
}
