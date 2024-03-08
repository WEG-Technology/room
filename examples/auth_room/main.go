package main

import (
	"fmt"
	"github.com/WEG-Technology/room"
	"github.com/WEG-Technology/room/store"
)

type AuthRequestPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthRequest() *room.Request {
	return room.NewRequest("auth/login", room.WithMethod(room.POST), room.WithBody(room.NewJsonBodyParser(
		AuthRequestPayload{
			Username: "kminchelle",
			Password: "0lelplR",
		})), room.ForceDTO())
}

type Auth struct{}

func (a Auth) Apply(connector *room.Connector, response room.Response) {
	connector.Header.Add("Authorization", "Bearer "+response.DTO.(map[string]any)["token"].(string))
}

func NewMeRequest() *room.Request {
	return room.NewRequest("auth/me")
}

func NewConnector() *room.Connector {
	return room.NewConnector("https://dummyjson.com", room.WithHeaderConnector(defaultHeaders()))
}

func defaultHeaders() room.IHeader {
	return room.NewHeader(
		store.NewMapStore().Add("Content-Type", "application/json"),
	)
}

type JsonplaceholderRoom struct {
	room.IAuthRoom
}

func NewRoom() *JsonplaceholderRoom {
	return &JsonplaceholderRoom{
		IAuthRoom: room.NewAuthRoom(NewConnector(), Auth{}, NewAuthRequest()),
	}
}

func main() {
	response, err := NewRoom().Send(NewMeRequest())

	fmt.Println("err: ", err)
	fmt.Println(response.OK())
}
