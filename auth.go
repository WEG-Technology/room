package room

import (
	"fmt"
)

type IAuthStrategy interface {
	Authenticate(request IRequest, authResponse IResponse) error
}

type ITokenAuthResponse interface {
	GetBearerToken() string
}

type TokenAuth struct{}

func NewTokenAuthStrategy(request IRequest) IAuthStrategy { return TokenAuth{} }
func (a TokenAuth) Authenticate(request IRequest, authResponse IResponse) error {
	response := authResponse.Dto().(ITokenAuthResponse)

	request.Header().Add("Authorization", fmt.Sprintf("Bearer %s", response.GetBearerToken()))

	return nil
}

type DefaultAuth struct {
	request IRequest
}

func (a DefaultAuth) GetAuthRequest() IRequest {
	return a.request
}

func (a DefaultAuth) Authenticate(request IRequest, authResponse IResponse) error {
	response := authResponse.Dto().(map[string]any)

	request.Header().Add("Authorization", fmt.Sprintf("Bearer %s", response["access_token"]))

	return nil
}

func NewDefaultAuth() IAuthStrategy { return DefaultAuth{} }
