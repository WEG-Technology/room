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
	Timeout() time.Duration
	Build() Context
}

type ContextBuilder struct {
	timeout time.Duration
}

func (b ContextBuilder) Build() Context {
	if b.Timeout() == 0 {
		return Context{
			Ctx:    context.Background(),
			Cancel: nil,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout())

	return Context{
		Ctx:    ctx,
		Cancel: cancel,
	}
}

func (b ContextBuilder) Timeout() time.Duration {
	return b.timeout
}

func NewContextBuilder(timeout time.Duration) ContextBuilder {
	return ContextBuilder{timeout}
}
