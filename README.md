# Room - OLD VERSION
⭐️Create your http client wrapper in a few minutes.

## ⚡️ Roadmap

⭐️ This package is still in development. We welcome suggestions for changes that will bring it closer to compliance without overly complicating the code, or useful test cases to add to the test suite.

- [x] init project as mvp
- [ ] add auth example
- [ ] remove panics from `request.go` file
- [ ] create default exceptions for fail status code
- [ ] make better readme
- [ ] add request debugger
- [ ] write unit tests
- [ ] add caching mechanism
- [ ] add concurrent request functionality to rooms
- [ ] add response normalizer
- [ ] make sure struct alignment is in a strict line

## Name
Room

## Description
Create your clients in a minute

## Installation
`go get -u github.com/WEG-Technology/room`

## Usage
There are many examples in examples folder go look for it.

Firstly, we need three files to start wrapping our client. Which are `connection.go`, `requests.go`, `responses.go`.
Finally, `room.go` which is the core wrap of your client.

#### Connection
We are going to wrap our base connection struct to get its skills.

#### Requests
All of your request should implement IRequest.

#### Responses
The responses literarily is a DTO of your response of your endpoint.

#### Room
Your client wrapper should implement IRoom. You can implement your custom interfaces by wrapping base Room struct.

There is an example of this explanation in `./examples/room_example`

#### Step By Step Usage
Let's say u have `https://dummyjson.com` api and you want to wrap it with Room package.

1- Create your `connection.go` file
```go
func NewConnection() room.IConnection {
    c, err := room.NewConnection( // you can use room.NewConnectionWithClient to use your custom http client
            room.WithBaseUrl("https://dummyjson.com"), // base url of your api
	    ),
    )

    if err != nil {
        panic(err) // you can handle error here
    }

    return Connection{c}
}
```

2- Create your `requests.go` file
```go
// AddTODORequest is a request to add to do
func NewAddTODORequest(params AddTODORequest) room.IRequest {
	r, err := room.NewPostRequest( // you can use room.NewGetRequest, room.NewPutRequest, room.NewDeleteRequest
		room.WithEndPoint("todos/add"), // endpoint of your request
		room.WithDto(&AddTODOResponse{}), // if you want to get response as a struct
		room.WithBody(room.NewJsonBodyParser(params)), // if you want to send request as a json
	)

	if err != nil {
		panic(err) // you can handle error here
	}

	return r
}
```

3- Create your `responses.go` file
```go
// AddTODOResponse is a response of add to do
type AddTODOResponse struct {
    ID int `json:"id"`
}
```

4- Create your `room.go` file
```go
// Room is a wrapper of dummyjson api
type DummyJsonClient struct {
	*room.Room // you can implement your custom interfaces by wrapping base Room struct
}
// AddTODO is a request to add to do
func (r DummyJsonClient) AddTODO(parameters AddTODORequest) room.IResponse {
	return r.Send(NewAddTODORequest(parameters)) // you can use r.Send to send your request
}
// NewRoom is a constructor of DummyJsonClient
func NewRoom() DummyJsonClient {
	return DummyJsonClient{room.NewRoom( // you can use room.NewRoom() to use your custom http client
		NewConnection(), // your connection
	)}
}
```

5- Use your client
```go
func main() {
    client := NewRoom() // create your client
    response := client.AddTODO(AddTODORequest{ // send your request
        Title: "Hello World", // your request parameters
    })
    fmt.Println(response.GetDto().(*AddTODOResponse).ID)
}
```


## 📫 Contributing
```TODO```

### 🧬 Clients Developed with Room
```TODO```
Please, feel free to add your client to this list.
- [Parasut](https://github.com/yahya077/parasut)
- ...