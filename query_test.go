package room

import (
	"github.com/WEG-Technology/room/store"
	"testing"
)

func TestNewQuery(t *testing.T) {
	// Test case for store.IMap type
	mapData := map[string]any{"key1": "value1", "key2": "value2"}
	query := NewQuery(store.NewMapStore(mapData))
	_, ok := query.(IMapQuery)
	if !ok {
		t.Error("NewQuery() did not return an IMapQuery instance for store.IMap")
	}

	// Test case for other types
	urlData := map[string]interface{}{"key1": "value1", "key2": "value2"}
	query = NewQuery(urlData)
	_, ok = query.(IUrlQuery)
	if !ok {
		t.Error("NewQuery() did not return an IUrlQuery instance for other types")
	}
}

func TestIMapQuery_String(t *testing.T) {
	// Test case for empty map
	mapData := map[string]any{}
	query := IMapQuery{v: store.NewMapStore(mapData)}
	result := query.String()
	if result != "" {
		t.Errorf("IMapQuery String() returned %s, expected an empty string for empty map", result)
	}

	// Test case for non-empty map
	mapData = map[string]any{"key1": "value1", "key2": "value2"}
	query = IMapQuery{v: store.NewMapStore(mapData)}
	result = query.String()
	expected := "key1=value1&key2=value2"
	if result != expected {
		t.Errorf("IMapQuery String() returned %s, expected %s", result, expected)
	}
}

func TestIUrlQuery_String(t *testing.T) {
	// Test case for valid URL values
	urlData := struct {
		Key1 string `url:"key1"`
		Key2 string `url:"key2"`
	}{"value1", "value2"}
	query := IUrlQuery{v: urlData}
	result := query.String()
	expected := "key1=value1&key2=value2"
	if result != expected {
		t.Errorf("IUrlQuery String() returned %s, expected %s", result, expected)
	}

	// Test case for invalid URL values
	invalidUrlData := map[string]interface{}{"key1": []string{"value1", "value2"}}
	query = IUrlQuery{v: invalidUrlData}
	result = query.String()
	if result != "" {
		t.Errorf("IUrlQuery String() returned %s for invalid URL values, expected an empty string", result)
	}
}
