package provider

import (
	"github.com/WEG-Technology/room"
)

func NewRequest() room.IRequest {
	r, e := room.NewGetRequest(
		room.WithEndPoint("api/Traveler"),
		room.WithMethod(room.GET),
		room.WithHeader(defaultHeader()),
		room.WithDto(&TravelerinformationResponse{}),
	)

	if e != nil {
		panic(e)
	}

	return r
}

func defaultHeader() room.IHeader {
	return room.NewHeader().
		Add("Accept", "text/xml")
}
