package mock

import (
	"github.com/WEG-Technology/room"
)

func NewConnection() room.IConnection {
	return new(mockConnection)
}

type mockConnection struct {
	baseUrl string
}

func (c *mockConnection) Send(request room.IRequest) room.IResponse {
	return &room.Response{}
}

func (c *mockConnection) GetDomain() string {
	return ""
}

func (c *mockConnection) SetProtocol(protocol room.HTTPProtocol) {}

func (c *mockConnection) SetBaseUrl(baseUrl string) {}

func (c *mockConnection) SetDomain(domain string) {}

func (c *mockConnection) InitUrl() {}
