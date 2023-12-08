package client

import "github.com/WEG-Technology/room"

type Connection struct {
	*room.Connection
}

func NewConnection() room.IConnection {
	c, e := room.NewConnection(
		room.WithBaseUrl("https://dummyjson.com"),
		room.WithDefaultHeader(defaultHeader()),
	)

	if e != nil {
		panic(e)
	}

	return Connection{c}
}

func defaultHeader() room.IHeader {
	return room.NewHeader().
		Add("Accept", "application/json").
		Add("Content-Type", "application/json")
}
