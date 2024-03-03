package provider

import (
	"github.com/WEG-Technology/room"
	"time"
)

func NewConnection() room.IConnection {
	c, err := room.NewConnection(
		room.WithProtocol(room.Http),
		room.WithDomain("restapi.adequateshop.com"),
		room.WithTimeout(15*time.Second),
	)

	if err != nil {
		panic(err)
	}

	return c
}
