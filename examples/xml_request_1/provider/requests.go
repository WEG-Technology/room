package provider

import (
	"github.com/WEG-Technology/room"
)

func NewRequest() room.IRequest {
	r, err := room.NewGetRequest(
		room.WithEndPoint("api/Traveler"),
		room.WithMethod(room.GET),
		room.WithHeader(defaultHeader()),
		room.WithDto(&TravelerinformationResponse{}),
	)

	if err != nil {
		panic(err)
	}

	return r
}

func defaultHeader() room.IHeader {
	return room.NewHeader().
		Add("Accept", "text/xml")
}
