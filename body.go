package room

import (
	"bytes"
	"encoding/json"
	"github.com/google/go-querystring/query"
	"mime/multipart"
	"net/url"
)

type IBodyParser interface {
	Parse() *bytes.Buffer
	ContentType() string
}

type JsonBody struct {
	v any
}

func (f *JsonBody) Parse() *bytes.Buffer {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(f.v)

	if err != nil {
		panic(err)
	}

	return &buf
}

func (f *JsonBody) ContentType() string {
	return "application/json"
}

func NewJsonBodyParser(v any) IBodyParser {
	return &JsonBody{v}
}

func NewFormURLEncodedBodyParser(v any) IBodyParser {
	return &FormURLEncodedBody{v}
}

type FormURLEncodedBody struct {
	v any
}

func (f *FormURLEncodedBody) ContentType() string {
	return "application/x-www-form-urlencoded"
}

func (f *FormURLEncodedBody) Parse() *bytes.Buffer {
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

// MultipartFormDataBody handles multipart/form-data encoding
type MultipartFormDataBody struct {
	formData    map[string]string
	contentType string
}

func (f *MultipartFormDataBody) ContentType() string {
	return f.contentType
}

func (f *MultipartFormDataBody) Parse() *bytes.Buffer {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for key, value := range f.formData {
		_ = writer.WriteField(key, value)
	}

	//TODO cover file fields

	_ = writer.Close()

	f.contentType = writer.FormDataContentType()

	return &body
}

func NewMultipartFormDataBodyParser(v any) IBodyParser {
	if _, ok := v.(map[string]string); !ok {
		panic("form data must be a map[string]string")
	}

	return &MultipartFormDataBody{
		formData: v.(map[string]string),
	}
}

type dumpBody struct{}

func (f dumpBody) Parse() *bytes.Buffer { return new(bytes.Buffer) }

func (f dumpBody) ContentType() string { return "" }
