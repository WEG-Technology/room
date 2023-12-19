package room

import "testing"

func TestJsonBody_Parse(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	body := NewJsonBodyParser(data).Parse()

	if body == nil {
		t.Error("Expected a non-nil buffer, but got nil")
	}

	// Optionally, you can check the content of the buffer or other expectations.
}

type testStruct struct {
	Key string `url:"key"`
}

func TestFormURLEncodedBody_Parse(t *testing.T) {
	data := testStruct{"value"}
	body := NewFormURLEncodedBodyParser(data).Parse()

	if body == nil {
		t.Error("Expected a non-nil buffer, but got nil")
	}

	// Optionally, you can check the content of the buffer or other expectations.
}

func TestFormURLEncodedBody_Parse_Error(t *testing.T) {
	// Test case where query.Values returns an error
	data := map[string]interface{}{"key": make(chan int)} // invalid type to trigger an error
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected a panic, but no panic occurred")
		}
	}()

	_ = NewFormURLEncodedBodyParser(data).Parse()
}
