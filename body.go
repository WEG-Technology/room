package room

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/go-querystring/query"
	"net/url"
)

// TODO make contentType as a custom type
func NewBodyFactory(contentType string) IBody {
	switch contentType {
	case headerValueApplicationJson:
		return JsonBody{}
	case headerValueFormEncoded:
		return FormURLEncodedBody{}
	default:
		panic(errors.New("unsupported type"))
	}
}

type IBody interface {
	Body(v any) *bytes.Buffer
}

type JsonBody struct{}

func (f JsonBody) Body(v any) *bytes.Buffer {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(v)

	if err != nil {
		panic(err)
	}

	return &buf
}

type FormURLEncodedBody struct{}

func (f FormURLEncodedBody) Body(v any) *bytes.Buffer {
	formValues := url.Values{}

	formValues, err := query.Values(v)

	if err != nil {
		panic(err)
	}

	return bytes.NewBufferString(formValues.Encode())
}
