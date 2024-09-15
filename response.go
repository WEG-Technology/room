package room

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/WEG-Technology/room/store"
	"io"
	"mime"
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
	responseDTO := newResponse(response.Request).setHeader(response.Header).setData(response)

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

	err := NewDTOFactory(r.Header.Get(headerKeyContentType)).marshall(r.Data, &body)

	return body, err
}

func (r Response) ResponseBody() map[string]any {
	var body map[string]any

	if r.Data != nil {
		_ = NewDTOFactory(r.Header.Get(headerKeyContentType)).marshall(r.Data, &body)
	}

	return body
}

func (r Response) DTO(v any) any {
	if r.Data != nil {
		_ = NewDTOFactory(r.Header.Get(headerKeyContentType)).marshall(r.Data, v)
	}

	return v
}

func (r Response) DTOorFail(v any) error {
	return NewDTOFactory(r.Header.Get(headerKeyContentType)).marshall(r.Data, v)
}

func (r Response) setRequestData(request *http.Request) Response {
	if request.Body != nil {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, request.Body)
		if err != nil {
			fmt.Println("Error reading request body:", err)
		}

		// Converting the body back to a string or further process it
		requestBodyData := buf.String()

		// Optionally, parse the body back into a map-like structure
		requestMap := make(map[string]string)
		for _, pair := range strings.Split(requestBodyData, "&") {
			parts := strings.Split(pair, "=")
			if len(parts) == 2 {
				requestMap[parts[0]] = parts[1]
			}
		}

		// Print the map structure
		fmt.Println("Request body as map:", requestMap)

		r.Request.Data, _ = io.ReadAll(request.Body)

		fmt.Println("DDddddddd", r.Request.Data)
	}

	return r
}

func (r Response) RequestBodyOrFail() (map[string]any, error) {
	var body map[string]any

	err := NewDTOFactory(r.Request.Header.Get(headerKeyContentType)).marshall(r.Request.Data, &body)

	return body, err
}

func (r Response) RequestBody() map[string]any {
	var body map[string]any

	if r.Request.Data != nil {
		_ = NewDTOFactory(r.Request.Header.Get(headerKeyContentType)).marshall(r.Request.Data, &body)
	}

	return body
}

func (r Response) RequestDTO(v any) any {
	if r.Data != nil {
		_ = NewDTOFactory(r.Request.Header.Get(headerKeyContentType)).marshall(r.Request.Data, &v)
	}

	return v
}

func (r Response) RequestDTOorFail(v any) any {
	return NewDTOFactory(r.Request.Header.Get(headerKeyContentType)).marshall(r.Request.Data, &v)
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

	if strings.Contains(ct, headerValueMultipartFormData) {
		return MultipartFormDataDTOFactory{contentType: ct}
	} else {
		switch ct {
		case headerValueApplicationJson:
			return JsonDTOFactory{}
		case headerValueTextXML:
			return XMLDTOFactory{}
		default:
			return JsonDTOFactory{}
		}
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

type MultipartFormDataDTOFactory struct {
	contentType string
}

// marshall parses the multipart/form-data and populates the provided DTO (v).
func (f MultipartFormDataDTOFactory) marshall(data []byte, v any) error {
	//TODO read request data
	return nil
}

func getBoundary(contentType string) string {
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		panic("invalid content type for multipart/form-data")
	}
	return params["boundary"]
}

func DTO[T any](response Response, v T) T {
	if response.Data != nil {
		_ = NewDTOFactory(response.Header.Get(headerKeyContentType)).marshall(response.Data, v)
	}

	return v
}

func DTOorFail[T any](response Response, v T) (T, error) {
	if response.Data != nil {
		err := NewDTOFactory(response.Header.Get(headerKeyContentType)).marshall(response.Data, v)

		if err != nil {
			return v, err
		}
	}

	return v, nil
}
