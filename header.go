package room

type IHeader interface {
	Properties() IMap
	Add(key string, value string) IHeader
	Get(key string) string
	Merge(header IHeader) IHeader
}

type Header struct {
	properties IMap
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

func (h *Header) Properties() IMap {
	return h.properties
}

func NewHeader(properties ...IMap) IHeader {
	if len(properties) == 0 {
		return &Header{NewMapStore()}
	}

	return &Header{properties[0]}
}
