package room

import "github.com/WEG-Technology/room/store"

type IHeader interface {
	Properties() store.IMap
	Add(key string, value string) IHeader
	Get(key string) string
	Merge(header IHeader) IHeader
	String() string
}

type Header struct {
	properties store.IMap
}

func (h *Header) Add(key string, value string) IHeader {
	h.properties.Add(key, value)

	return h
}

func (h *Header) Get(key string) string {
	s, ok := h.properties.GetItem(key)

	if !ok {
		return ""
	}

	return s.(string)
}

func (h *Header) Merge(header IHeader) IHeader {
	if header != nil {
		h.properties = h.properties.MergeIMap(header.Properties())
	}

	return h
}

func (h *Header) Properties() store.IMap {
	return h.properties
}

func (h *Header) String() (str string) {
	return h.Properties().StringAll()
}

func NewHeader(properties ...store.IMap) IHeader {
	if len(properties) == 0 {
		return &Header{store.NewMapStore()}
	}

	return &Header{properties[0]}
}
