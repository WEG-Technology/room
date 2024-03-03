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
	r, err := room.NewPostRequest(
		room.WithEndPoint("todos/add"),
		room.WithDto(new(AddTODOResponse)),
		room.WithBody(room.NewJsonBodyParser(params)),
	)

	if err != nil {
		panic(err)
	}

	return r
}
