package room

import (
	"github.com/WEG-Technology/room/store"
	"testing"
)

func TestHeader_Add(t *testing.T) {
	h := NewHeader()
	h.Add("key", "value")
	val := h.Get("key")
	if val != "value" {
		t.Errorf("Header Add() did not add the value correctly")
	}
}

func TestHeader_Get(t *testing.T) {
	h := NewHeader()
	h.Add("key", "value")
	val := h.Get("key")
	if val != "value" {
		t.Errorf("Header Get() did not return the correct value")
	}

	// Test case for non-existing key
	val = h.Get("nonExistingKey")
	if val != "" {
		t.Errorf("Header Get() did not return an empty string for non-existing key")
	}
}

func TestHeader_Merge(t *testing.T) {
	h1 := NewHeader()
	h2 := NewHeader()
	h2.Add("key", "value")
	h1.Merge(h2)
	val := h1.Get("key")
	if val != "value" {
		t.Errorf("Header Merge() did not merge the headers correctly")
	}
}

func TestHeader_Properties(t *testing.T) {
	h := NewHeader()
	properties := h.Properties()
	if properties == nil {
		t.Error("Header Properties() returned nil")
	}
}

func TestHeader_String(t *testing.T) {
	h := NewHeader()
	h.Add("key1", "value1")
	h.Add("key2", "value2")
	expected := "key1: value1, key2: value2"
	if h.String() != expected {
		t.Errorf("Header String() returned %s, expected %s", h.String(), expected)
	}
}

func TestNewHeader(t *testing.T) {
	h := NewHeader()
	if h == nil {
		t.Error("NewHeader() returned nil")
	}
}

func TestNewHeaderWithIMap(t *testing.T) {
	h := NewHeader(store.NewMapStore(map[string]any{"key": "value"}))
	if h == nil {
		t.Error("NewHeader() returned nil")
	}
}
