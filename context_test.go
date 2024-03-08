package room

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Test cases for ContextBuilder
func TestContextBuilder_Build(t *testing.T) {
	// Test case: Build context without a timeout
	builderWithoutTimeout := NewContextBuilder(0)
	ctxWithoutTimeout := builderWithoutTimeout.Build()

	// Check if the returned context is background and cancel function is nil
	assert.Equal(t, context.Background(), ctxWithoutTimeout.ctx)
	assert.Nil(t, ctxWithoutTimeout.cancel)

	// Test case: Build context with a timeout
	timeoutDuration := time.Second * 5
	builderWithTimeout := NewContextBuilder(timeoutDuration)
	ctxWithTimeout := builderWithTimeout.Build()

	// Check if the returned context is with a timeout and cancel function is not nil
	assert.NotEqual(t, context.Background(), ctxWithTimeout.ctx)
	assert.NotNil(t, ctxWithTimeout.cancel)

	// Check if the context is canceled after the specified timeout
	select {
	case <-ctxWithTimeout.ctx.Done():
		// Context is canceled, as expected
	case <-time.After(timeoutDuration + time.Second): // Allow extra time for the context to be canceled
		t.Error("Context is not canceled within the specified timeout")
	}

	// Ensure cancel function is called after the test
	defer ctxWithTimeout.cancel()
}

func TestContextBuilder_Timeout(t *testing.T) {
	// Test case: Get timeout from the ContextBuilder
	timeoutDuration := time.Second * 10
	builder := NewContextBuilder(timeoutDuration)

	// Check if the returned timeout matches the expected value
	assert.Equal(t, timeoutDuration, builder.Timeout())
}
