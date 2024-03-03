package main

import (
	"fmt"
	"github.com/WEG-Technology/room"
)

type CustomAuth struct{}

func (a CustomAuth) Authenticate(request room.IRequest, authResponse room.IResponse) error {
	response := authResponse.Dto().(map[string]any)

	request.Header().Add("Authorization", fmt.Sprintf("Bearer %s", response["token"]))

	return nil
}

func NewCustomAuth() CustomAuth { return CustomAuth{} }
