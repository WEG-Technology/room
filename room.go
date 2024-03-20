package room

import "errors"

const (
	ErrAuthRoomCanNotFoundKey = "authToken can not found in response"
)

type IRoom interface {
	Send(request *Request) (Response, error)
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
	AuthRequest *Request
	AuthToken   string
}

func NewAuthRoom(connector *Connector, authRequest *Request, authToken string) IRoom {
	return &AuthRoom{
		Room:        &Room{Connector: connector},
		AuthRequest: authRequest,
		AuthToken:   authToken,
	}
}

func (r *AuthRoom) Send(request *Request) (Response, error) {
	if r.AuthRequest != nil {
		response, err := r.Connector.Do(r.AuthRequest)

		if err != nil {
			return response, err
		}

		if !response.OK() {
			return response, err
		}

		if token, found := findToken(response.DTO.(map[string]any), r.AuthToken); found {
			r.Connector.Header.Add("Authorization", "Bearer "+token)
		} else {
			return response, errors.New(ErrAuthRoomCanNotFoundKey)
		}
	}

	return r.Room.Send(request)
}

type IAuth interface {
	Apply(connector *Connector, response Response)
}

type AccessTokenAuth struct{}

func (a AccessTokenAuth) Apply(connector *Connector, response Response) {
	connector.Header.Add("Authorization", "Bearer "+response.DTO.(map[string]any)["access_token"].(string))
}

func NewAccessTokenAuth() IAuth {
	return AccessTokenAuth{}
}

func findToken(responseMap map[string]any, tokenKey string) (string, bool) {
	for k, v := range responseMap {
		if k == tokenKey {
			if strValue, ok := v.(string); ok {
				return strValue, true
			}
			return "", false
		}
		if nestedData, ok := v.(map[string]any); ok {
			if value, found := findToken(nestedData, tokenKey); found {
				return value, true
			}
		}
	}
	return "", false
}
