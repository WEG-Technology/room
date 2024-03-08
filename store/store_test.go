package store

import (
	"reflect"
	"testing"
)

func TestMapStore(t *testing.T) {
	// Test NewMapStore with default data
	store := NewMapStore(map[string]interface{}{"key1": 42, "key2": "value"})
	expectedAllData := map[string]interface{}{"key1": 42, "key2": "value"}
	checkAllData(t, store, expectedAllData)

	// Test Set method
	store.Set(map[string]interface{}{"key3": 3.14})
	expectedAllData = map[string]interface{}{"key3": 3.14}
	checkAllData(t, store, expectedAllData)

	// Test SetMultiple method
	store.SetMultiple(map[string]interface{}{"key4": true}, map[string]interface{}{"key5": "hello"})
	expectedAllData = map[string]interface{}{"key3": 3.14, "key4": true, "key5": "hello"}
	checkAllData(t, store, expectedAllData)

	// Test GetItem method
	value, ok := store.GetItem("key4")
	if !ok || value.(bool) != true {
		t.Errorf("GetItem() method did not return the expected value. Got: %v, Expected: true", value)
	}

	// Test Add method
	store.Add("key6", "world")
	expectedAllData["key6"] = "world"
	checkAllData(t, store, expectedAllData)

	// Test Integer method
	store.SetMultiple(map[string]interface{}{"key1": 42, "key2": "value"})
	intValue := store.Integer("key1")
	expectedIntValue := 42
	if intValue != expectedIntValue {
		t.Errorf("Integer() method did not return the expected value. Got: %v, Expected: %v", intValue, expectedIntValue)
	}

	// Test String method
	stringValue := store.String("key2")
	expectedStringValue := "value"
	if stringValue != expectedStringValue {
		t.Errorf("String() method did not return the expected value. Got: %v, Expected: %v", stringValue, expectedStringValue)
	}

	// Test StringList method
	store.Add("key7", []string{"one", "two", "three"})
	stringListValue := store.StringList("key7")
	expectedStringListValue := []string{"one", "two", "three"}
	if !reflect.DeepEqual(stringListValue, expectedStringListValue) {
		t.Errorf("StringList() method did not return the expected value. Got: %v, Expected: %v", stringListValue, expectedStringListValue)
	}

	// Test Float method
	floatValue := store.Float("key3")
	expectedFloatValue := 3.14
	if floatValue != expectedFloatValue {
		t.Errorf("Float() method did not return the expected value. Got: %v, Expected: %v", floatValue, expectedFloatValue)
	}

	// Test IsEmpty method
	isEmpty := store.IsEmpty()
	if isEmpty {
		t.Errorf("IsEmpty() method returned true, but the map is not empty")
	}
}

func checkAllData(t *testing.T, store IMap, expected map[string]interface{}) {
	t.Helper()
	allData := store.All()
	if !reflect.DeepEqual(allData, expected) {
		t.Errorf("All() method did not return the expected data. Got: %v, Expected: %v", allData, expected)
	}
}
