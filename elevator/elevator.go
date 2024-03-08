package elevator

import (
	"encoding/json"
	"fmt"
	"github.com/WEG-Technology/room"
	"github.com/WEG-Technology/room/segment"
	"github.com/WEG-Technology/room/store"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
	"time"
)

func readYml(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func injectEnv(ymlFile []byte) string {
	return os.ExpandEnv(string(ymlFile))
}

func unmarshalYmlContent(ymlContent string) (config IntegrationConfig, err error) {
	err = yaml.Unmarshal([]byte(ymlContent), &config)

	return config, err
}

func marshalJson(config IntegrationConfig) (jsonData []byte, err error) {
	return json.MarshalIndent(config, "", "  ")
}

type IElevatorEngine interface {
	Execute(roomKey, requestKey string) (room.Response, error)
	ExecuteConcurrent(concurrentKey string) map[string]room.Response
	WarmUp() IElevatorEngine
	PutBodyParser(roomKey, requestKey string, bodyParser room.IBodyParser) IElevatorEngine
	PutQuery(roomKey, requestKey string, authStrategy room.IQuery) IElevatorEngine
	PutDTO(roomKey, requestKey string, dto any) IElevatorEngine
	GetElapsedTime() float64
	PutAuthStrategy(roomKey string, authStrategy room.IAuth) IElevatorEngine
}

type ElevatorEngine struct {
	elevator       Elevator
	Segment        segment.ISegment
	RoomContainers map[string]RoomContainer
}

func (e *ElevatorEngine) GetElapsedTime() float64 {
	return e.Segment.GetElapsedTime()
}

type RoomContainer struct {
	Room     room.IAuthRoom
	Requests map[string]*room.Request
}

func NewElevatorEngine(elevator Elevator) IElevatorEngine {
	return &ElevatorEngine{
		elevator: elevator,
	}
}

// TODO add safety check
func (e *ElevatorEngine) Execute(roomKey, requestKey string) (room.Response, error) {
	if roomContainerEntry, ok := e.RoomContainers[roomKey]; ok {
		if requestEntry, ok := roomContainerEntry.Requests[requestKey]; ok {
			return roomContainerEntry.Room.Send(requestEntry)
		}
		panic(fmt.Sprintf("engine for %s on %s not configured", roomKey, requestKey))
	}

	panic(fmt.Sprintf("engine for %s not configured", roomKey))
}

type RoomResponseContainer struct {
	Responses map[string]room.Response
	mu        sync.Mutex
}

func (c *RoomResponseContainer) PutResponse(key string, response room.Response) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Responses[key] = response
}

func (e *ElevatorEngine) ExecuteConcurrent(concurrentKey string) map[string]room.Response {
	e.Segment = segment.StartSegmentNow()
	defer e.Segment.End()

	responseContainer := RoomResponseContainer{
		Responses: map[string]room.Response{},
	}

	var wg sync.WaitGroup

	for roomKey, configRoom := range e.elevator.Config.Flat.Rooms {
		for requestKey, req := range configRoom.Requests {
			if concurrentKey == req.ConcurrentKey {
				wg.Add(1)
				go func(a, b string) {
					//TODO handle errors
					res, _ := e.Execute(a, b)

					responseContainer.PutResponse(a, res)
					wg.Done()
				}(roomKey, requestKey)
			}
		}
	}

	wg.Wait()

	return responseContainer.Responses
}

// TODO remove panics by a parameter
func (e *ElevatorEngine) WarmUp() IElevatorEngine {
	roomContainers := map[string]RoomContainer{}

	for roomKey, r := range e.elevator.Config.Flat.Rooms {
		c := room.NewConnector(
			r.Connection.BaseURL,
			room.WithHeaderConnector(room.NewHeader(store.NewMapStore(r.Connection.Headers))),
			room.WithHeaderContextBuilder(room.NewContextBuilder(time.Duration(r.Connection.Timeout)*time.Second)),
		)

		var roomObj room.IAuthRoom

		if r.Connection.Auth.Type == "bearer" {
			roomObj = room.NewAuthRoom(c, nil, e.CreateRequest(r.Connection.Auth.Request))
		} else {
			roomObj = room.NewAuthRoom(c, nil, nil)
		}

		roomContainers[roomKey] = RoomContainer{
			Room:     roomObj,
			Requests: map[string]*room.Request{},
		}

		for requestKey, req := range r.Requests {
			roomContainers[roomKey].Requests[requestKey] = e.CreateRequest(req)
		}
	}

	e.RoomContainers = roomContainers

	return e
}

func (e *ElevatorEngine) CreateRequest(req Request) *room.Request {
	var parser room.IBodyParser
	//TODO expand here
	switch req.Body.Type {
	case "json":
		parser = room.NewJsonBodyParser(req.Body.Content)
	case "form":
		parser = room.NewFormURLEncodedBodyParser(req.Body.Content)
	default:
		parser = nil
	}

	r := room.NewRequest(
		req.Path,
		room.WithMethod(room.HTTPMethod(req.Method)),
		room.WithBody(parser),
	)

	if req.ForceDTO {
		r.ForceDTO = true
	}

	return r
}

func (e *ElevatorEngine) PutBodyParser(roomKey, requestKey string, bodyParser room.IBodyParser) IElevatorEngine {
	if roomContainerEntry, ok := e.RoomContainers[roomKey]; ok {
		if requestEntry, ok := roomContainerEntry.Requests[requestKey]; ok {
			requestEntry.BodyParser = bodyParser
			return e
		}
		panic(fmt.Sprintf("engine for %s on %s not configured", roomKey, requestKey))
	}

	panic(fmt.Sprintf("engine for %s not configured", roomKey))
}

func (e *ElevatorEngine) PutQuery(roomKey, requestKey string, query room.IQuery) IElevatorEngine {
	if roomContainerEntry, ok := e.RoomContainers[roomKey]; ok {
		if requestEntry, ok := roomContainerEntry.Requests[requestKey]; ok {
			requestEntry.Query = query
			return e
		}
		panic(fmt.Sprintf("engine for %s on %s not configured", roomKey, requestKey))
	}

	panic(fmt.Sprintf("engine for %s not configured", roomKey))
}

func (e *ElevatorEngine) PutAuthStrategy(roomKey string, authStrategy room.IAuth) IElevatorEngine {
	if roomContainerEntry, ok := e.RoomContainers[roomKey]; ok {
		roomContainerEntry.Room.SetAuthStrategy(authStrategy)
		return e
	}

	panic(fmt.Sprintf("engine for %s not configured", roomKey))
}

func (e *ElevatorEngine) PutDTO(roomKey, requestKey string, dto any) IElevatorEngine {
	if roomContainerEntry, ok := e.RoomContainers[roomKey]; ok {
		if requestEntry, ok := roomContainerEntry.Requests[requestKey]; ok {
			requestEntry.DTO = dto
			return e
		}
		panic(fmt.Sprintf("engine for %s on %s not configured", roomKey, requestKey))
	}

	panic(fmt.Sprintf("engine for %s not configured", roomKey))
}

func NewElevator(integrationYmlPath string) Elevator {
	ymlFile, err := readYml(integrationYmlPath)

	if err != nil {
		panic(err)
	}

	ymlContent := injectEnv(ymlFile)

	config, err := unmarshalYmlContent(ymlContent)

	if err != nil {
		panic(err)
	}

	return Elevator{
		config,
	}
}

type IElevator interface {
	AddBody(body Body, room string, request string) Elevator
	AddHeader(header string, value string, room string) Elevator
	JsonData() ([]byte, error)
}

type Elevator struct {
	Config IntegrationConfig
}

func (e Elevator) AddBody(body Body, room string, request string) Elevator {
	if roomEntry, ok := e.Config.Flat.Rooms[room]; ok {
		if requestEntry, ok := roomEntry.Requests[request]; ok {
			requestEntry.Body = body
		}
	}

	return e
}

func (e Elevator) AddHeader(header string, value string, room string) Elevator {
	e.Config.Flat.Rooms[room].Connection.Headers[header] = value

	return e
}

func (e Elevator) JsonData() ([]byte, error) {
	return marshalJson(e.Config)
}

func NewBody(contentType string, content any) Body {
	return Body{
		Type:    contentType,
		Content: content,
	}
}

type Body struct {
	Type    string `yaml:"type"`
	Content any    `yaml:"content"`
}

type Connection struct {
	BaseURL string         `yaml:"baseUrl"`
	Timeout int            `yaml:"timeout"`
	Headers map[string]any `yaml:"headers"`
	Auth    ConnectionAuth `yaml:"auth"`
}

type ConnectionAuth struct {
	Type    string  `yaml:"type"`
	Request Request `yaml:"request"`
}

type Request struct {
	ConcurrentKey string `yaml:"concurrentKey"`
	Method        string `yaml:"method"`
	Path          string `yaml:"path"`
	ForceDTO      bool   `yaml:"forceDTO"`
	Body          Body   `yaml:"body"`
}

type Room struct {
	Connection Connection         `yaml:"connection"`
	Requests   map[string]Request `yaml:"requests"`
}

type Flat struct {
	Rooms map[string]Room `yaml:"rooms"`
}

type IntegrationConfig struct {
	Flat Flat `yaml:"flat"`
}

func (r *IntegrationConfig) UnmarshalJSON(data []byte) error {
	var res IntegrationConfig

	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	return nil
}
