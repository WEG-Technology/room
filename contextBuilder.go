package room

import (
	"context"
	"time"
)

type IContextBuilder interface {
	Timeout() time.Duration
	Build() (context.Context, context.CancelFunc)
}

type ContextBuilder struct {
	timeout time.Duration
}

func (b ContextBuilder) Build() (context.Context, context.CancelFunc) {
	if b.Timeout() == 0 {
		return context.Background(), nil
	}

	return context.WithTimeout(context.Background(), b.Timeout())
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
