package room

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

const (
	connectionErrorStr       = "connection not found"
	protocolNotFoundErrorStr = "protocol not found"
)

type IConnectionConfig interface {
	TimeoutDuration() time.Duration
	Domain() string
	Protocol() HTTPProtocol
	BaseUrl() string
	DefaultHeader() IHeader
}

type IConnection interface {
	Send(request IRequest) IResponse
	GetDomain() string
	InitUrl()
	setProtocol(protocol HTTPProtocol)
	setBaseUrl(baseUrl string)
	setDomain(domain string)
}

type IConnectionObserver interface {
	OnDo(request IRequest)
	OnRequestCreated(request IRequest)
	OnResponseCompleted(request IRequest, response IResponse, err error)
}

type Option func(info *Connection)

func WithProtocol(protocol HTTPProtocol) Option {
	return func(connection *Connection) {
		connection.protocol = protocol
	}
}

func WithDomain(domain string) Option {
	return func(connection *Connection) {
		connection.domain = domain
	}
}

func WithBaseUrl(baseUrl string) Option {
	return func(connection *Connection) {
		connection.baseUrl = baseUrl
	}
}

func WithTimeout(duration time.Duration) Option {
	return func(connection *Connection) {
		connection.timeoutDuration = duration
	}
}

func WithDefaultHeader(header IHeader) Option {
	return func(connection *Connection) {
		connection.defaultHeader = header
	}
}

func NewConnection(opts ...Option) (*Connection, error) {
	c := new(Connection)

	c.client = new(http.Client)

	for _, opt := range opts {
		opt(c)
	}

	if c.observer == nil {
		c.observer = connectionObserver{}
	}

	if c.timeoutDuration == 0 {
		c.timeoutDuration = 15 * time.Second
	}

	c.InitUrl()

	if c.baseUrl == "" {
		return c, errors.New(connectionErrorStr)
	}

	if c.protocol == 0 {
		return c, errors.New(protocolNotFoundErrorStr)
	}

	return c, nil
}

func withProtocol(host string) bool {
	return strings.Contains(host, Http.String()) || strings.Contains(host, Https.String())
}

type ConnectionConfig struct {
	timeoutDuration time.Duration
	protocol        HTTPProtocol
	domain          string
	baseUrl         string
	defaultHeader   IHeader
}

func DefaultConnectionConfig() IConnectionConfig {
	return ConnectionConfig{
		defaultHeader: NewHeader(),
	}
}

func (c ConnectionConfig) TimeoutDuration() time.Duration {
	return c.timeoutDuration
}

func (c ConnectionConfig) Domain() string {
	return c.domain
}

func (c ConnectionConfig) Protocol() HTTPProtocol {
	return c.protocol
}

func (c ConnectionConfig) BaseUrl() string {
	return c.baseUrl
}

func (c ConnectionConfig) DefaultHeader() IHeader {
	return c.defaultHeader
}

type Connection struct {
	client          *http.Client
	domain          string
	protocol        HTTPProtocol
	baseUrl         string
	defaultHeader   IHeader
	timeoutDuration time.Duration
	observer        IConnectionObserver
}

func (c *Connection) do(request IRequest) IResponse {
	c.observer.OnDo(request)

	request.SetConnectionConfig(ConnectionConfig{
		domain:          c.domain,
		protocol:        c.protocol,
		baseUrl:         c.baseUrl,
		timeoutDuration: c.timeoutDuration,
		defaultHeader:   c.defaultHeader,
	})

	defer request.Create().Cancel()

	c.observer.OnRequestCreated(request)

	res, err := c.client.Do(request.Request())

	response := NewResponseDirector().Response(res, request, err)

	defer c.observer.OnResponseCompleted(request, response, err)

	if err != nil {
		return response
	}

	defer res.Body.Close()

	return response
}

func (c *Connection) Send(request IRequest) IResponse {
	return c.do(request)
}

func (c *Connection) GetDomain() string {
	return c.domain
}

func (c *Connection) setProtocol(protocol HTTPProtocol) {
	c.protocol = protocol
}

func (c *Connection) setBaseUrl(baseUrl string) {
	c.baseUrl = baseUrl
}

func (c *Connection) setDomain(domain string) {
	c.domain = domain
}

func (c *Connection) InitUrl() {
	// TODO make string parse great again
	if c.baseUrl != "" {
		if strings.Contains(c.baseUrl, Https.String()) {
			c.protocol = Https
			c.domain = strings.Replace(c.baseUrl, Https.String(), "", 1)
		} else {
			c.protocol = Http
			c.domain = strings.Replace(c.baseUrl, Http.String(), "", 1)
		}
	}

	c.baseUrl = c.protocol.String() + c.domain + "/"
}
