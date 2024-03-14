package room

import (
	"context"
	"time"
)

type Context struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

type IContextBuilder interface {
	Build() Context
}

type ContextBuilder struct {
	timeout time.Duration
}

func (b ContextBuilder) Build() Context {
	if b.timeout == 0 {
		return Context{
			Ctx:    context.Background(),
			Cancel: nil,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)

	return Context{
		Ctx:    ctx,
		Cancel: cancel,
	}
}

func NewContextBuilder(timeout time.Duration) ContextBuilder {
	return ContextBuilder{timeout}
}
