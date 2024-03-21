package room

import (
	"encoding/json"
	"encoding/xml"
	"github.com/WEG-Technology/room/store"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	RequestURI    URI
	StatusCode    int
	Method        string
	Header        IHeader
	RequestHeader IHeader
	RequestBody   map[string]any
	DTO           any
	Data          []byte
}

func NewResponse(r *http.Response, dto any) (Response, error) {
	response := Response{
		StatusCode:  r.StatusCode,
		Method:      r.Request.Method,
		RequestBody: map[string]any{},
	}.
		setHeader(r.Header).
		setRequestHeader(r.Request.Header).
		setRequestBodyData(r.Request).
		setRequestURI(r.Request)

	var err error

	response, err = response.setData(r)

	if dto != nil {
		response, err = response.setDTO(dto)
	}

	return response, err
}

func (r Response) OK() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

func (r Response) setHeader(header http.Header) Response {
	m := store.NewMapStore()

	for key, values := range header {
		m.Add(key, strings.Join(values, " "))
	}

	r.Header = NewHeader(m)

	return r
}

func (r Response) setRequestHeader(header http.Header) Response {
	m := store.NewMapStore()

	for key, values := range header {
		m.Add(key, strings.Join(values, " "))
	}

	r.RequestHeader = NewHeader(m)

	return r
}

func (r Response) setRequestBodyData(request *http.Request) Response {
	if request.Body != nil {
		decoder := json.NewDecoder(request.Body)

		for decoder.More() {
			if err := decoder.Decode(&r.RequestBody); err != nil {
				panic(err)
			}
		}
	}

	return r
}

func (r Response) setRequestURI(request *http.Request) Response {
	r.RequestURI = NewURI(request.URL.String())

	return r
}

func (r Response) setData(response *http.Response) (Response, error) {
	var err error

	r.Data, err = io.ReadAll(response.Body)

	return r, err
}

func (r Response) setDTO(dto any) (Response, error) {
	r.DTO = dto

	err := NewDTOFactory(r.Header.Get(headerKeyAccept)).marshall(r.Data, &r.DTO)

	return r, err
}

// IDTOFactory declares the interface for creating DTOs.
type IDTOFactory interface {
	marshall(data []byte, v any) error
}

// NewDTOFactory creates a concrete factory based on content type.
func NewDTOFactory(contentType ...string) IDTOFactory {
	var ct string
	if len(contentType) > 0 {
		ct = contentType[0]
	} else {
		ct = ""
	}

	switch ct {
	case headerValueApplicationJson:
		return JsonDTOFactory{}
	case headerValueTextXML:
		return XMLDTOFactory{}
	default:
		return JsonDTOFactory{}
	}
}

// JsonDTOFactory creates JSON DTOs.
type JsonDTOFactory struct{}

func (r JsonDTOFactory) marshall(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// XMLDTOFactory creates XML DTOs.
type XMLDTOFactory struct{}

func (r XMLDTOFactory) marshall(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}
