package client

import "github.com/WEG-Technology/room"

type Connection struct {
	*room.Connection
}

func NewConnection() room.IConnection {
	c, err := room.NewConnection(
		room.WithBaseUrl("https://dummyjson.com"),
		room.WithDefaultHeader(defaultHeader()),
	)

	if err != nil {
		panic(err)
	}

	return Connection{c}
}

func defaultHeader() room.IHeader {
	return room.NewHeader().
		Add("Accept", "application/json").
		Add("Content-Type", "application/json")
}
