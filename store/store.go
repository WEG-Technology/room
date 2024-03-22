package store

import (
	"fmt"
	"reflect"
	"strings"
)

type IMap interface {
	All() map[string]any
	Integer(key string) int
	String(key string) string
	StringList(key string) []string
	Float(key string) float64
	GetItem(key string) (any, bool)
	Add(key string, value any) IMap
	IsEmpty() bool
	Remove(key string) IMap
	Set(data map[string]any) IMap
	Merge(data map[string]any) IMap
	Each(callback func(key string, value any)) IMap
	MergeIMap(m IMap) IMap
	SetMultiple(data ...map[string]any) IMap
	StringAll() string
}

func NewMapStore(defData ...map[string]any) IMap {
	mapStore := new(MapStore)
	mapStore.data = map[string]any{}
	if len(defData) > 0 {
		mapStore.SetMultiple(defData...)
	}
	return mapStore
}

type MapStore struct {
	data map[string]any
}

func (s *MapStore) getValidData(key string) any {
	if v, ok := s.data[key]; ok {
		return v
	}
	return nil
}

func (s *MapStore) All() map[string]any {
	return s.data
}

func (s *MapStore) Integer(key string) int {
	return s.getValidData(key).(int)
}

func (s *MapStore) String(key string) string {
	return s.getValidData(key).(string)
}

func (s *MapStore) StringAll() string {
	v := reflect.ValueOf(s.All())

	keys := v.MapKeys()

	parts := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		value := v.MapIndex(keys[i]).Interface()
		parts[i] = fmt.Sprintf("%v: %v", keys[i].Interface(), value)
	}

	return strings.Join(parts, ", ")
}

func (s *MapStore) StringList(key string) []string {
	return s.getValidData(key).([]string)
}

func (s *MapStore) Float(key string) float64 {
	return s.getValidData(key).(float64)
}

func (s *MapStore) GetItem(key string) (any, bool) {
	v, ok := s.data[key]
	return v, ok
}

func (s *MapStore) Add(key string, value any) IMap {
	s.data[key] = value
	return s
}

func (s *MapStore) IsEmpty() bool {
	return s.data == nil || len(s.data) == 0
}

func (s *MapStore) Remove(key string) IMap {
	delete(s.data, key)
	return s
}

func (s *MapStore) Set(data map[string]any) IMap {
	s.data = data
	return s
}

func (s *MapStore) SetMultiple(data ...map[string]any) IMap {
	for i := 0; i < len(data); i++ {
		s.Merge(data[i])
	}
	return s
}

func (s *MapStore) Each(callback func(key string, value any)) IMap {
	for key, value := range s.data {
		callback(key, value)
	}

	return s
}

func (s *MapStore) Merge(data map[string]any) IMap {
	merged := make(map[string]interface{})
	for k, v := range s.data {
		if _, ok := s.data[k]; ok {
			merged[k] = v
		}
	}

	for k, v := range data {
		if _, ok := data[k]; ok {
			merged[k] = v
		}
	}

	return s.Set(merged)
}

func (s *MapStore) MergeIMap(m IMap) IMap {
	return s.Merge(m.All())
}
