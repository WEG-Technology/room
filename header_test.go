package room

import (
	"github.com/WEG-Technology/room/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test cases for the Header implementation
func TestHeader_Add(t *testing.T) {
	// Create a new Header
	header := NewHeader()

	// Add a key-value pair
	key := "Content-Type"
	value := "application/json"
	resultHeader := header.Add(key, value)

	// Check if the resultHeader is the same as the original header
	assert.Equal(t, header, resultHeader)

	// Check if the added key-value pair exists in the properties
	assert.Equal(t, value, header.Get(key))
}

func TestHeader_Get(t *testing.T) {
	// Create a new Header
	header := NewHeader()

	// Attempt to get a value for a non-existing key
	result := header.Get("NonExistingKey")

	// Check if the result is an empty string
	assert.Equal(t, "", result)

	// Add a key-value pair
	key := "Content-Type"
	value := "application/json"
	header.Add(key, value)

	// Attempt to get the added value using the key
	result = header.Get(key)

	// Check if the result is the same as the added value
	assert.Equal(t, value, result)
}

func TestHeader_Merge(t *testing.T) {
	// Create two headers with different key-value pairs
	header1 := NewHeader()
	header1.Add("Key1", "Value1")
	header1.Add("Key2", "Value2")

	header2 := NewHeader()
	header2.Add("Key3", "Value3")
	header2.Add("Key4", "Value4")

	// Merge header2 into header1
	mergedHeader := header1.Merge(header2)

	// Check if the mergedHeader contains all key-value pairs from both headers
	assert.Equal(t, "Value1", mergedHeader.Get("Key1"))
	assert.Equal(t, "Value2", mergedHeader.Get("Key2"))
	assert.Equal(t, "Value3", mergedHeader.Get("Key3"))
	assert.Equal(t, "Value4", mergedHeader.Get("Key4"))
}

func TestNewHeader(t *testing.T) {
	// Create a new Header without specifying properties
	header := NewHeader()

	// Check if the created header has an empty properties map
	assert.NotNil(t, header)
	assert.NotNil(t, header.Properties())

	// Create a new Header with a specified properties map
	properties := store.NewMapStore()
	properties.Add("Key1", "Value1")
	properties.Add("Key2", "Value2")
	headerWithProperties := NewHeader(properties)

	// Check if the created header has the specified properties map
	assert.NotNil(t, headerWithProperties)
	assert.Equal(t, properties, headerWithProperties.Properties())
}
