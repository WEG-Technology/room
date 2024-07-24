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
	StatusCode int
	Header     IHeader
	Data       []byte
	Request    RequestDTO
}

type RequestDTO struct {
	Header IHeader
	Data   []byte
	URI    URI
	Method string
}

func NewResponse(response *http.Response, request *http.Request) Response {
	responseDTO := newResponse(request).setHeader(response.Header).setData(response)

	responseDTO.StatusCode = response.StatusCode

	return responseDTO
}

func NewErrorResponse(request *http.Request, err error) (Response, error) {
	responseDTO := newResponse(request)

	errData, _ := json.Marshal(map[string]string{"runtime_error": err.Error()})
	responseDTO.Data = errData

	return responseDTO, err
}

func newResponse(request *http.Request) Response {
	responseDTO := Response{
		Header: NewHeader(),
		Request: RequestDTO{
			Method: request.Method,
			URI:    NewURI(request.URL.String()),
			Header: NewHeader(),
		},
	}.
		setRequestHeader(request.Header).
		setRequestData(request)

	return responseDTO
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

	r.Request.Header = NewHeader(m)

	return r
}

func (r Response) setData(response *http.Response) Response {
	if response.Body != nil {
		r.Data, _ = io.ReadAll(response.Body)
	}

	return r
}

func (r Response) ResponseBodyOrFail() (map[string]any, error) {
	var body map[string]any

	err := NewDTOFactory(r.Header.Get(headerKeyAccept)).marshall(r.Data, &body)

	return body, err
}

func (r Response) ResponseBody() map[string]any {
	var body map[string]any

	if r.Data != nil {
		_ = NewDTOFactory(r.Header.Get(headerKeyAccept)).marshall(r.Data, &body)
	}

	return body
}

func (r Response) DTO(v any) any {
	if r.Data != nil {
		_ = NewDTOFactory(r.Header.Get(headerKeyAccept)).marshall(r.Data, v)
	}

	return v
}

func (r Response) DTOorFail(v any) error {
	return NewDTOFactory(r.Header.Get(headerKeyAccept)).marshall(r.Data, v)
}

func (r Response) setRequestData(request *http.Request) Response {
	if request.Body != nil {
		r.Request.Data, _ = io.ReadAll(request.Body)
	}

	return r
}

func (r Response) RequestBodyOrFail() (map[string]any, error) {
	var body map[string]any

	err := NewDTOFactory(r.Header.Get(headerKeyContentType)).marshall(r.Request.Data, &body)

	return body, err
}

func (r Response) RequestBody() map[string]any {
	var body map[string]any

	if r.Request.Data != nil {
		_ = NewDTOFactory(r.Header.Get(headerKeyAccept)).marshall(r.Request.Data, &body)
	}

	return body
}

func (r Response) RequestDTO(v any) any {
	if r.Data != nil {
		_ = NewDTOFactory(r.Header.Get(headerKeyContentType)).marshall(r.Request.Data, &v)
	}

	return v
}

func (r Response) RequestDTOorFail(v any) any {
	return NewDTOFactory(r.Header.Get(headerKeyContentType)).marshall(r.Request.Data, &v)
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
