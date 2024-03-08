package room

import (
	"context"
	"time"
)

type Context struct {
	ctx    context.Context
	cancel context.CancelFunc
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
			ctx:    context.Background(),
			cancel: nil,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout())

	return Context{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (b ContextBuilder) Timeout() time.Duration {
	return b.timeout
}

func NewContextBuilder(timeout time.Duration) ContextBuilder {
	return ContextBuilder{timeout}
}

/*
// If you want to use this, you need to add a dependency to the transaction package in the go.mod file. (eq. newrelic, opentracing, etc.)
type ContextWithTransaction struct {
	txn *Transaction
}

func (b ContextWithTransaction) Build() (context.Context, context.CancelFunc) {
	return NewContext(context.Background(), txn), nil
}
*/

//Example: NewRelicContextBuilder
/*

type NewRelicContextBuilder struct {
	txn *newrelic.Transaction
}

func (b NewRelicContextBuilder) Build() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), b.Timeout())

	return newrelic.NewContext(ctx, b.txn), cancel
}

func (b NewRelicContextBuilder) Timeout() time.Duration {
	return time.Second * 15
}

func NewContextBuilder(txn *newrelic.Transaction) room.IContextBuilder {
	return NewRelicContextBuilder{txn}
}

*/
