package provider

import (
	"github.com/WEG-Technology/room"
)

type Room struct {
	*room.Room
}

func (r Room) Search() room.IResponse {
	return r.Send(NewRequest())
}

func NewRoom() Room {
	return Room{room.NewRoom(
		NewConnection(),
	)}
}
