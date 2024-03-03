package elevator

import (
	"encoding/json"
	"fmt"
	"github.com/WEG-Technology/room"
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
	Execute(roomKey, requestKey string) room.IResponse
	ExecuteConcurrent(concurrentKey string) map[string]room.IResponse
	WarmUp() IElevatorEngine
	PutBodyParser(roomKey, requestKey string, bodyParser room.IBodyParser) IElevatorEngine
	PutQuery(roomKey, requestKey string, authStrategy room.IQuery) IElevatorEngine
	PutAuthStrategy(roomKey string, authStrategy room.IAuthStrategy) IElevatorEngine
	PutDTO(roomKey, requestKey string, dto any) IElevatorEngine
	GetElapsedTime() float64
}

type ElevatorEngine struct {
	elevator       Elevator
	Segment        room.ISegment
	RoomContainers map[string]RoomContainer
}

func (e *ElevatorEngine) GetElapsedTime() float64 {
	return e.Segment.GetElapsedTime()
}

type RoomContainer struct {
	Room     room.IRoom
	Requests map[string]room.IRequest
}

func NewElevatorEngine(elevator Elevator) IElevatorEngine {
	return &ElevatorEngine{
		elevator: elevator,
	}
}

// TODO add safety check
func (e *ElevatorEngine) Execute(roomKey, requestKey string) room.IResponse {
	if roomContainerEntry, ok := e.RoomContainers[roomKey]; ok {
		if requestEntry, ok := roomContainerEntry.Requests[requestKey]; ok {
			return roomContainerEntry.Room.Send(requestEntry)
		}
	}

	panic("engine not configured")
}

type RoomResponseContainer struct {
	Responses map[string]room.IResponse
	mu        sync.Mutex
}

func (c *RoomResponseContainer) PutResponse(key string, response room.IResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Responses[key] = response
}

func (e *ElevatorEngine) ExecuteConcurrent(concurrentKey string) map[string]room.IResponse {
	e.Segment = room.StartSegmentNow()
	defer e.Segment.End()

	responseContainer := RoomResponseContainer{
		Responses: map[string]room.IResponse{},
	}

	var wg sync.WaitGroup

	for roomKey, configRoom := range e.elevator.Config.Flat.Rooms {
		for requestKey, request := range configRoom.Requests {
			if concurrentKey == request.ConcurrentKey {
				wg.Add(1)
				go func(a, b string) {
					responseContainer.PutResponse(a, e.Execute(a, b))
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
		c, err := room.NewConnection(
			room.WithBaseUrl(r.Connection.BaseURL),
			room.WithDefaultHeader(room.NewHeader(room.NewMapStore(r.Connection.Headers))),
			room.WithTimeout(time.Duration(r.Connection.Timeout)*time.Second),
		)

		if err != nil {
			panic(err)
		}

		roomContainers[roomKey] = RoomContainer{
			Room:     room.NewRoom(c),
			Requests: map[string]room.IRequest{},
		}

		for requestKey, request := range r.Requests {
			roomContainers[roomKey].Requests[requestKey] = e.CreateRequest(request)

			if r.Connection.Auth.Type == "bearer" {
				roomContainers[roomKey].Requests[requestKey].PutPreRequest(e.CreateRequest(r.Connection.Auth.Request))
			}
		}
	}

	e.RoomContainers = roomContainers

	return e
}

func (e *ElevatorEngine) CreateRequest(request Request) room.IRequest {
	var parser room.IBodyParser
	//TODO expand here
	switch request.Body.Type {
	case "json":
		parser = room.NewJsonBodyParser(request.Body.Content)
	case "form":
		parser = room.NewFormURLEncodedBodyParser(request.Body.Content)
	default:
		parser = nil
	}

	r, err := room.NewRequest(
		//TODO use unmarshall to get the method instead of hardcoding
		room.WithMethod(room.HTTPMethod(request.Method)),
		room.WithEndPoint(request.Path),
		room.WithBody(parser),
	)

	if err != nil {
		panic(err)
	}

	return r
}

func (e *ElevatorEngine) PutBodyParser(roomKey, requestKey string, bodyParser room.IBodyParser) IElevatorEngine {
	if roomContainerEntry, ok := e.RoomContainers[roomKey]; ok {
		if requestEntry, ok := roomContainerEntry.Requests[requestKey]; ok {
			requestEntry.PutBodyParser(bodyParser)
			return e
		}
		panic(fmt.Sprintf("engine for %s on %s not configured", roomKey, requestKey))
	}

	panic(fmt.Sprintf("engine for %s not configured", roomKey))
}

func (e *ElevatorEngine) PutQuery(roomKey, requestKey string, query room.IQuery) IElevatorEngine {
	if roomContainerEntry, ok := e.RoomContainers[roomKey]; ok {
		if requestEntry, ok := roomContainerEntry.Requests[requestKey]; ok {
			requestEntry.PutQuery(query)
			return e
		}
		panic(fmt.Sprintf("engine for %s on %s not configured", roomKey, requestKey))
	}

	panic(fmt.Sprintf("engine for %s not configured", roomKey))
}

func (e *ElevatorEngine) PutAuthStrategy(roomKey string, authStrategy room.IAuthStrategy) IElevatorEngine {
	if roomContainerEntry, ok := e.RoomContainers[roomKey]; ok {
		roomContainerEntry.Room.PutAuthStrategy(authStrategy)
		return e
	}

	panic(fmt.Sprintf("engine for %s not configured", roomKey))
}

func (e *ElevatorEngine) PutDTO(roomKey, requestKey string, dto any) IElevatorEngine {
	if roomContainerEntry, ok := e.RoomContainers[roomKey]; ok {
		if requestEntry, ok := roomContainerEntry.Requests[requestKey]; ok {
			requestEntry.PutDTO(dto)
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
