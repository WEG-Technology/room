# Room
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

## 📫 Contributing
```TODO```

### 🧬 Clients Developed with Room
```TODO```
Please, feel free to add your client to this list.
- ...
- ...