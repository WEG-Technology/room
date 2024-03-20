package room

import (
	"bytes"
	"encoding/json"
	"github.com/google/go-querystring/query"
	"net/url"
)

type IBodyParser interface {
	Parse() *bytes.Buffer
}

type JsonBody struct {
	v any
}

func (f JsonBody) Parse() *bytes.Buffer {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(f.v)

	if err != nil {
		panic(err)
	}

	return &buf
}

func NewJsonBodyParser(v any) IBodyParser {
	return JsonBody{v}
}

func NewFormURLEncodedBodyParser(v any) IBodyParser {
	return FormURLEncodedBody{v}
}

type FormURLEncodedBody struct {
	v any
}

func (f FormURLEncodedBody) Parse() *bytes.Buffer {
	values := url.Values{}

	switch f.v.(type) {
	case map[string]any:
		for key, value := range f.v.(map[string]any) {
			values.Add(key, value.(string))
		}
	default:
		values, _ = query.Values(f.v)
	}

	return bytes.NewBufferString(values.Encode())
}

type dumpBody struct{}

func (f dumpBody) Parse() *bytes.Buffer { return new(bytes.Buffer) }
