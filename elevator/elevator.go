package elevator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/WEG-Technology/room"
	"github.com/WEG-Technology/room/segment"
	"github.com/WEG-Technology/room/store"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"slices"
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
	DynamicExecute(roomKey, requestKey string, v any) (room.Response, error)
	ExecuteConcurrent(concurrentKey string, appliedRooms ...string) map[string]room.Response
	WarmUp() IElevatorEngine
	PutBodyParser(roomKey, requestKey string, bodyParser room.IBodyParser) IElevatorEngine
	PutQuery(roomKey, requestKey string, authStrategy room.IQuery) IElevatorEngine
	GetElapsedTime() float64
	Request(roomKey, requestKey string) (*room.Request, error)
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
	Room     room.IRoom
	Requests map[string]*room.Request
}

func NewElevatorEngine(elevator Elevator) IElevatorEngine {
	return &ElevatorEngine{
		elevator: elevator,
	}
}

func (e *ElevatorEngine) Execute(roomKey, requestKey string) (room.Response, error) {
	if roomContainerEntry, ok := e.RoomContainers[roomKey]; ok {
		if requestEntry, ok := roomContainerEntry.Requests[requestKey]; ok {
			return roomContainerEntry.Room.Send(requestEntry)
		}
		panic(fmt.Sprintf("engine for %s on %s not configured", roomKey, requestKey))
	}

	panic(fmt.Sprintf("engine for %s not configured", roomKey))
}

func (e *ElevatorEngine) DynamicExecute(roomKey, requestKey string, v any) (room.Response, error) {
	if roomContainerEntry, ok := e.RoomContainers[roomKey]; ok {
		if requestEntry, ok := roomContainerEntry.Requests[requestKey]; ok {
			elevatorRequest := e.elevator.GetRequest(roomKey, requestKey)
			requestEntry.BodyParser = e.generateDynamicParser(elevatorRequest.Body.Type, elevatorRequest.Body.DynamicContent, NewDynamicExecutionPayload(v).toMap())
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

func (e *ElevatorEngine) ExecuteConcurrent(concurrentKey string, appliedRooms ...string) map[string]room.Response {
	e.Segment = segment.StartSegmentNow()
	defer e.Segment.End()

	responseContainer := RoomResponseContainer{
		Responses: map[string]room.Response{},
	}

	var wg sync.WaitGroup

	for roomKey, configRoom := range e.elevator.Config.Flat.Rooms {
		for requestKey, req := range configRoom.Requests {
			if ((len(appliedRooms) > 0 && slices.Contains(appliedRooms, roomKey)) || len(appliedRooms) == 0) && concurrentKey == req.ConcurrentKey {
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

		var roomObj room.IRoom

		if r.Connection.Auth.Type == "bearer" {
			roomObj = room.NewAuthRoom(c, e.CreateRequest(r.Connection.Auth.Request), r.Connection.Auth.AccessTokenKey)
		} else {
			roomObj = room.NewAuthRoom(c, nil, "")
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
	parser := e.initParser(req.Body.Type, req.Body.Content)

	optionRequests := []room.OptionRequest{
		room.WithMethod(room.HTTPMethod(req.Method)),
		room.WithBody(parser),
	}

	r := room.NewRequest(
		req.Path,
		optionRequests...,
	)

	return r
}

func (e *ElevatorEngine) initParser(bodyType string, content any) room.IBodyParser {
	var parser room.IBodyParser
	//TODO expand here
	switch bodyType {
	case "json":
		parser = room.NewJsonBodyParser(content)
	case "form":
		parser = room.NewFormURLEncodedBodyParser(content)
	case "multipart-form":
		parser = room.NewMultipartFormDataBodyParser(content)
	default:
		parser = nil
	}

	return parser
}

// TODO should be refactored in v2 for use dynamics as `content` instead of `dynamicContent`
func (e *ElevatorEngine) generateDynamicParser(bodyType string, dynamicContents []DynamicContent, v map[string]any) room.IBodyParser {
	requestPayload := map[string]any{}

	for _, dynamicContent := range dynamicContents {

		if dynamicContent.Value != nil {
			requestPayload[dynamicContent.Key] = dynamicContent.Value
			continue
		}

		if _, ok := v[dynamicContent.Key]; !ok {
			panic(fmt.Sprintf("dynamic content key %s not found in payload", dynamicContent.Key))
		}

		requestPayload[dynamicContent.Key] = v[dynamicContent.Key]
	}

	return e.initParser(bodyType, requestPayload)
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

func (e *ElevatorEngine) Request(roomKey, requestKey string) (*room.Request, error) {
	if roomContainerEntry, ok := e.RoomContainers[roomKey]; ok {
		if requestEntry, ok := roomContainerEntry.Requests[requestKey]; ok {
			return requestEntry, nil
		}
		return nil, errors.New("request not found")
	}

	return nil, errors.New("request not found")
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

func (e Elevator) GetRequest(roomKey, requestKey string) Request {
	if roomEntry, ok := e.Config.Flat.Rooms[roomKey]; ok {
		if requestEntry, ok := roomEntry.Requests[requestKey]; ok {
			return requestEntry
		}
	}

	panic(fmt.Sprintf("request for %s on %s not found", roomKey, requestKey))
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
	Type           string           `yaml:"type"`
	Content        any              `yaml:"content"`
	DynamicContent []DynamicContent `yaml:"dynamicContent"`
}

type DynamicContent struct {
	Key   string `yaml:"key"`
	Value any    `yaml:"value"`
}

type Connection struct {
	BaseURL string         `yaml:"baseUrl"`
	Timeout int            `yaml:"timeout"`
	Headers map[string]any `yaml:"headers"`
	Auth    ConnectionAuth `yaml:"auth"`
}

type ConnectionAuth struct {
	Type           string  `yaml:"type"`
	AccessTokenKey string  `yaml:"accessTokenKey"`
	Request        Request `yaml:"request"`
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

func (r *IntegrationConfig) UnmarshalJSON(data []byte) error {
	var res IntegrationConfig

	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	return nil
}

type DynamicExecutionPayload struct {
	Fields []DynamicExecutionField
}

func (p DynamicExecutionPayload) toMap() map[string]any {
	payload := map[string]any{}

	for _, field := range p.Fields {
		payload[field.Key] = field.Value
	}

	return payload
}

type DynamicExecutionField struct {
	Key   string
	Value any
}

func NewDynamicExecutionPayload(v any) DynamicExecutionPayload {
	payload := DynamicExecutionPayload{}
	switch v.(type) {
	case map[string]any:
		for key, value := range v.(map[string]any) {
			payload.Fields = append(payload.Fields, DynamicExecutionField{Key: key, Value: value})
		}
	case store.MapStore:
		storeMap := v.(store.MapStore)
		for key, value := range storeMap.All() {
			payload.Fields = append(payload.Fields, DynamicExecutionField{Key: key, Value: value})
		}
	default:
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Struct {
			structToMap(val, &payload)
		}
	}

	return payload
}

func structToMap(val reflect.Value, payload *DynamicExecutionPayload) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Tag.Get("json")

		if field.Kind() == reflect.Struct {
			structToMap(field, payload)
		} else {
			payload.Fields = append(payload.Fields, DynamicExecutionField{
				Key:   fieldName,
				Value: field.Interface(),
			})
		}
	}
}
