package client

import (
	"github.com/WEG-Technology/room"
)

type DummyJsonClient struct {
	*room.Room
}

func (r DummyJsonClient) AddTODO(parameters AddTODORequest) room.IResponse {
	return r.Send(NewAddTODORequest(parameters))
}

func NewRoom() DummyJsonClient {
	return DummyJsonClient{room.NewRoom(
		NewConnection(),
	)}
}
