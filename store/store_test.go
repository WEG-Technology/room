package store

import (
	"reflect"
	"testing"
)

func TestMapStore_All(t *testing.T) {
	ms := NewMapStore(map[string]any{"key1": 1, "key2": "value2"})
	allData := ms.All()
	expected := map[string]any{"key1": 1, "key2": "value2"}
	if !reflect.DeepEqual(allData, expected) {
		t.Errorf("All() returned %+v, expected %+v", allData, expected)
	}
}

func TestMapStore_Integer(t *testing.T) {
	ms := NewMapStore(map[string]any{"key1": 42})
	result := ms.Integer("key1")
	expected := 42
	if result != expected {
		t.Errorf("Integer() returned %d, expected %d", result, expected)
	}

	// Test for invalid key
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Integer() did not panic for invalid key")
		}
	}()
	ms.Integer("invalidKey")
}

func TestMapStore_String(t *testing.T) {
	ms := NewMapStore(map[string]any{"key1": "value1"})
	result := ms.String("key1")
	expected := "value1"
	if result != expected {
		t.Errorf("String() returned %s, expected %s", result, expected)
	}

	// Test for invalid key
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("String() did not panic for invalid key")
		}
	}()
	ms.String("invalidKey")
}

func TestMapStore_StringList(t *testing.T) {
	ms := NewMapStore(map[string]any{"key1": []string{"value1", "value2"}})
	result := ms.StringList("key1")
	expected := []string{"value1", "value2"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("StringList() returned %+v, expected %+v", result, expected)
	}

	// Test for invalid key
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("StringList() did not panic for invalid key")
		}
	}()
	ms.StringList("invalidKey")
}

func TestMapStore_Float(t *testing.T) {
	ms := NewMapStore(map[string]any{"key1": 3.14})
	result := ms.Float("key1")
	expected := 3.14
	if result != expected {
		t.Errorf("Float() returned %f, expected %f", result, expected)
	}

	// Test for invalid key
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Float() did not panic for invalid key")
		}
	}()
	ms.Float("invalidKey")
}

func TestMapStore_GetItem(t *testing.T) {
	ms := NewMapStore(map[string]any{"key1": "value1"})
	result, ok := ms.GetItem("key1")
	expectedValue := "value1"
	if !ok || result != expectedValue {
		t.Errorf("GetItem() returned (%v, %t), expected (%v, %t)", result, ok, expectedValue, true)
	}

	// Test for invalid key
	result, ok = ms.GetItem("invalidKey")
	if ok || result != nil {
		t.Errorf("GetItem() returned (%v, %t), expected (%v, %t)", result, ok, nil, false)
	}
}

func TestMapStore_Add(t *testing.T) {
	ms := NewMapStore()
	ms.Add("key1", "value1").Add("key2", 42)
	expected := map[string]any{"key1": "value1", "key2": 42}
	if !reflect.DeepEqual(ms.All(), expected) {
		t.Errorf("Add() did not add expected items")
	}
}

func TestMapStore_IsEmpty(t *testing.T) {
	ms := NewMapStore()
	if !ms.IsEmpty() {
		t.Errorf("IsEmpty() returned false, expected true for empty MapStore")
	}

	ms.Add("key1", "value1")
	if ms.IsEmpty() {
		t.Errorf("IsEmpty() returned true, expected false for non-empty MapStore")
	}
}

func TestMapStore_Remove(t *testing.T) {
	ms := NewMapStore(map[string]any{"key1": "value1", "key2": 42})
	ms.Remove("key1")
	expected := map[string]any{"key2": 42}
	if !reflect.DeepEqual(ms.All(), expected) {
		t.Errorf("Remove() did not remove the expected item")
	}

	// Test removing non-existing key
	ms.Remove("key3")
	if !reflect.DeepEqual(ms.All(), expected) {
		t.Errorf("Remove() removed an item for a non-existing key")
	}
}

func TestMapStore_Set(t *testing.T) {
	ms := NewMapStore()
	ms.Set(map[string]any{"key1": "value1", "key2": 42})
	expected := map[string]any{"key1": "value1", "key2": 42}
	if !reflect.DeepEqual(ms.All(), expected) {
		t.Errorf("Set() did not set the expected items")
	}
}

func TestMapStore_SetMultiple(t *testing.T) {
	ms := NewMapStore()
	ms.SetMultiple(map[string]any{"key1": "value1"}, map[string]any{"key2": 42})
	expected := map[string]any{"key1": "value1", "key2": 42}
	if !reflect.DeepEqual(ms.All(), expected) {
		t.Errorf("SetMultiple() did not set the expected items")
	}
}

func TestMapStore_Each(t *testing.T) {
	ms := NewMapStore(map[string]any{"key1": "value1", "key2": 42})
	expected := map[string]any{"key1": "value1", "key2": 42}
	ms.Each(func(key string, value any) {
		if expected[key] != value {
			t.Errorf("Each() returned unexpected key-value pair")
		}
	})
}

func TestMapStore_Merge(t *testing.T) {
	ms := NewMapStore(map[string]any{"key1": "value1", "key2": 42})
	ms.Merge(map[string]any{"key3": 3.14, "key4": []string{"a", "b"}})
	expected := map[string]any{"key1": "value1", "key2": 42, "key3": 3.14, "key4": []string{"a", "b"}}
	if !reflect.DeepEqual(ms.All(), expected) {
		t.Errorf("Merge() did not merge the expected items")
	}
}

func TestMapStore_MergeIMap(t *testing.T) {
	ms := NewMapStore(map[string]any{"key1": "value1", "key2": 42})
	ms2 := NewMapStore(map[string]any{"key3": 3.14, "key4": []string{"a", "b"}})
	ms.MergeIMap(ms2)
	expected := map[string]any{"key1": "value1", "key2": 42, "key3": 3.14, "key4": []string{"a", "b"}}
	if !reflect.DeepEqual(ms.All(), expected) {
		t.Errorf("MergeIMap() did not merge the expected items")
	}
}

func TestMapStore_StringAll(t *testing.T) {
	ms := NewMapStore(map[string]any{"key1": "value1", "key2": 42})
	result := ms.StringAll()
	expected := "key1: value1, key2: 42"
	if result != expected {
		t.Errorf("StringAll() returned %s, expected %s", result, expected)
	}
}
