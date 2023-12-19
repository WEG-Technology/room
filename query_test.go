package room

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

// Test cases for Query, IMapQuery, and IUrlQuery
func TestNewQuery(t *testing.T) {
	// Create a new Query with an IMap
	mockMap := NewMapStore()
	query := NewQuery(mockMap)

	// Check if the created query is of type IMapQuery
	assert.NotNil(t, query)
	assert.IsType(t, IMapQuery{}, query)

	// Create a new Query with a different type
	nonMapQuery := NewQuery("someValue")

	// Check if the created query is of type IUrlQuery
	assert.NotNil(t, nonMapQuery)
	assert.IsType(t, IUrlQuery{}, nonMapQuery)
}

func TestIMapQuery_String(t *testing.T) {
	// Create a new IMapQuery with a mock map
	mockMap := NewMapStore(map[string]interface{}{"key1": "value1", "key2": "value2"})
	iMapQuery := IMapQuery{v: mockMap}

	// Get the string representation of the IMapQuery
	result := iMapQuery.String()

	// Check if the result matches the expected URL-encoded string
	expected := url.Values{"key1": {"value1"}, "key2": {"value2"}}.Encode()
	assert.Equal(t, expected, result)
}

func TestIUrlQuery_String(t *testing.T) {
	// Create a new IUrlQuery with a mock value
	mockValue := struct {
		Key1 string `url:"key1"`
		Key2 string `url:"key2"`
	}{Key1: "value1", Key2: "value2"}
	iUrlQuery := IUrlQuery{v: mockValue}

	// Get the string representation of the IUrlQuery
	result := iUrlQuery.String()

	// Check if the result matches the expected URL-encoded string
	expected := url.Values{"key1": {"value1"}, "key2": {"value2"}}.Encode()
	assert.Equal(t, expected, result)
}
