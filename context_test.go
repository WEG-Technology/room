package room

import (
	"context"
	"testing"
	"time"
)

func TestContextBuilder_Build(t *testing.T) {
	// Test case where timeout is not set
	builder := NewContextBuilder(0)
	ctx := builder.Build()

	// Ensure that the context is background context and cancel function is nil
	if ctx.Ctx != context.Background() {
		t.Error("Context is not background context when timeout is not set")
	}
	if ctx.Cancel != nil {
		t.Error("Cancel function is not nil when timeout is not set")
	}

	// Test case where timeout is set
	timeout := time.Second
	builder = NewContextBuilder(timeout)
	ctx = builder.Build()

	// Ensure that the context is context with timeout and cancel function is not nil
	if ctx.Ctx.Err() != nil {
		t.Error("Context is not context with timeout when timeout is set")
	}
	if ctx.Cancel == nil {
		t.Error("Cancel function is nil when timeout is set")
	}
}
