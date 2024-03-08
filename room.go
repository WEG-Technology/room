package room

type IRoom interface {
	Send(request *Request) (Response, error)
}

type IAuthRoom interface {
	Send(request *Request) (Response, error)
	SetAuthStrategy(auth IAuth)
}

type Room struct {
	Connector *Connector
}

func NewRoom(connector *Connector) *Room {
	return &Room{Connector: connector}
}

func (r *Room) Send(request *Request) (Response, error) {
	return r.Connector.Do(request)
}

type AuthRoom struct {
	*Room
	AuthRequest  *Request
	AuthStrategy IAuth
}

func NewAuthRoom(connector *Connector, auth IAuth, authRequest *Request) IAuthRoom {
	return &AuthRoom{
		Room:         &Room{Connector: connector},
		AuthRequest:  authRequest,
		AuthStrategy: auth,
	}
}

func (r *AuthRoom) Send(request *Request) (Response, error) {
	if r.AuthRequest != nil {
		response, err := r.Connector.Do(r.AuthRequest)

		if err != nil {
			return response, err
		}

		r.AuthStrategy.Apply(r.Connector, response)

		if !response.OK() {
			return response, err
		}
	}

	return r.Room.Send(request)
}

func (r *AuthRoom) SetAuthStrategy(auth IAuth) {
	r.AuthStrategy = auth
}

type IAuth interface {
	Apply(connector *Connector, response Response)
}

type AccessTokenAuth struct{}

func (a AccessTokenAuth) Apply(connector *Connector, response Response) {
	connector.Header.Add("Authorization", "Bearer "+response.DTO.(map[string]any)["access_token"].(string))
}
