package room

import (
	"strings"
	"testing"
)

func TestJsonBody_Parse(t *testing.T) {
	// Test case where JSON encoding is successful
	data := map[string]interface{}{"key": "value"}
	body := JsonBody{v: data}
	buffer := body.Parse()
	expected := `{"key":"value"}`
	bufferString := strings.TrimRight(buffer.String(), "\n") // Remove trailing newline
	if bufferString != expected {
		t.Errorf("JsonBody Parse() returned %s, expected %s", bufferString, expected)
	}

	// Test case where JSON encoding fails
	defer func() {
		if r := recover(); r == nil {
			t.Error("JsonBody Parse() did not panic when JSON encoding failed")
		}
	}()
	invalidData := make(chan int) // Invalid data for JSON encoding
	invalidBody := JsonBody{v: invalidData}
	_ = invalidBody.Parse()
}

func TestFormURLEncodedBodyAsStruct_Parse(t *testing.T) {
	mapData := struct {
		Key1 string `url:"key1"`
		Key2 int    `url:"key2"`
	}{"value1", 42}
	body := FormURLEncodedBody{v: mapData}
	buffer := body.Parse()
	expected := "key1=value1&key2=42"
	if buffer.String() != expected {
		t.Errorf("FormURLEncodedBody Parse() returned %s, expected %s", buffer.String(), expected)
	}
}

func TestFormURLEncodedBodyAsMap_Parse(t *testing.T) {
	mapData := map[string]any{"key1": "value1", "key2": "42"}
	body := FormURLEncodedBody{v: mapData}
	buffer := body.Parse()
	expected := "key1=value1&key2=42"
	if buffer.String() != expected {
		t.Errorf("FormURLEncodedBody Parse() returned %s, expected %s", buffer.String(), expected)
	}
}

func TestDumpBody_Parse(t *testing.T) {
	// Test DumpBody Parse() always returns an empty buffer
	body := dumpBody{}
	buffer := body.Parse()
	expected := ""
	if buffer.String() != expected {
		t.Errorf("DumpBody Parse() returned %s, expected empty", buffer.String())
	}
}

func TestNewJsonBodyParser(t *testing.T) {
	// Test NewJsonBodyParser() creates a JsonBody instance
	data := map[string]interface{}{"key": "value"}
	bodyParser := NewJsonBodyParser(data)
	_, ok := bodyParser.(JsonBody)
	if !ok {
		t.Error("NewJsonBodyParser() did not return a JsonBody instance")
	}
}

func TestNewFormURLEncodedBodyParser(t *testing.T) {
	// Test NewFormURLEncodedBodyParser() creates a FormURLEncodedBody instance
	data := map[string]interface{}{"key": "value"}
	bodyParser := NewFormURLEncodedBodyParser(data)
	_, ok := bodyParser.(FormURLEncodedBody)
	if !ok {
		t.Error("NewFormURLEncodedBodyParser() did not return a FormURLEncodedBody instance")
	}
}
