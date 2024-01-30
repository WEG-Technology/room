package room

import (
	"fmt"
)

type IAuthStrategy interface {
	GetAuthRequest() IRequest
	Authenticate(request IRequest, authResponse IResponse) error
}

type ITokenAuthResponse interface {
	GetBearerToken() string
}

type TokenAuth struct {
	request IRequest
}

func NewTokenAuthStrategy(request IRequest) IAuthStrategy {
	return TokenAuth{request}
}

func (a TokenAuth) GetAuthRequest() IRequest {
	return a.request
}

func (a TokenAuth) Authenticate(request IRequest, authResponse IResponse) error {
	response := authResponse.Dto().(ITokenAuthResponse)

	request.Header().Add("Authorization", fmt.Sprintf("Bearer %s", response.GetBearerToken()))

	return nil
}
