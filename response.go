package room

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net"
	"net/http"
	"strings"
)

type ResponseDirector struct{}

func (f ResponseDirector) Response(response *http.Response, request IRequest, err error) IResponse {
	var e error
	responseDTO := NewResponse(request)

	if err != nil {
		return NewErrorResponse(err).Response(responseDTO)
	}

	responseDTO.SetHeader(response.Header).setStatus(response.Status).setStatusCode(response.StatusCode)

	if e = responseDTO.SetRequestBodyData(request); e != nil {
		return responseDTO.Error(e)
	}

	if e = responseDTO.SetData(response); e != nil {
		return responseDTO.Error(e)
	}

	if request.Dto() != nil {
		_ = NewDTOFactory(request.Header().Get(headerKeyAccept)).marshallDTO(responseDTO.Data(), request.Dto())

		responseDTO.setDto(request.Dto())
	}

	return responseDTO
}

type IErrorResponse interface {
	Response(response IResponse) IResponse
}

func NewErrorResponse(err error) IErrorResponse {
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return TimeoutErrorResponse{err}
	}

	return DefaultErrorResponse{err}
}

type DefaultErrorResponse struct{ error }

func (e DefaultErrorResponse) Response(response IResponse) IResponse {
	return response.Error(e.error)
}

type TimeoutErrorResponse struct{ error }

func (e TimeoutErrorResponse) Response(response IResponse) IResponse {
	return response.setStatus(e.error.Error()).setStatusCode(http.StatusRequestTimeout)
}

type IResponseDirector interface {
	Response(response *http.Response, request IRequest, err error) IResponse
}

func NewResponseDirector() IResponseDirector {
	return ResponseDirector{}
}

func NewDTOFactory(contentType string) IDTOFactory {
	switch contentType {
	case headerValueApplicationJson:
		return JsonDTOFactory{}
	case headerValueTextXML:
		return XMLDTOFactory{}
	default:
		return JsonDTOFactory{}
	}
}

type DTOFactory struct {
	contentType string
}

type IDTOFactory interface {
	marshallDTO(data []byte, v any) error
}

type JsonDTOFactory struct{}

func (r JsonDTOFactory) marshallDTO(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

type XMLDTOFactory struct{}

func (r XMLDTOFactory) marshallDTO(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}

type IResponse interface {
	Status() string
	StatusCode() int
	Header() IHeader
	RequestHeader() IHeader
	Data() []byte
	RequestBodyData() []byte
	Dto() any
	Ok() bool
	RequestUri() string
	RequestMethod() string
	setRequestBodyData(requestBodyData []byte) IResponse
	setDto(dto any) IResponse
	setData(data []byte) IResponse
	setStatus(status string) IResponse
	setStatusCode(status int) IResponse
	SetRequestBodyData(request IRequest) error
	SetData(response *http.Response) error
	SetHeader(header http.Header) IResponse
	Error(err error) IResponse
}

func NewResponse(request IRequest) IResponse {
	return &Response{
		requestHeader: request.Header(),
		requestUri:    request.Request().URL.RequestURI(),
		requestMethod: request.Request().Method,
	}
}

type Response struct {
	status          string
	requestMethod   string
	requestUri      string
	statusCode      int
	header          IHeader
	requestHeader   IHeader
	requestBodyData []byte
	data            []byte
	dto             any
}

func (r *Response) SetRequestBodyData(request IRequest) error {
	if request.Request().Body != nil {
		data, err := io.ReadAll(request.Request().Body)

		if err != nil {
			return err
		}

		r.setData(data)
	}

	return nil
}

func (r *Response) SetData(response *http.Response) error {
	d, err := io.ReadAll(response.Body)

	if err != nil {
		return err
	}

	r.setData(d)

	return nil
}

func (r *Response) SetHeader(header http.Header) IResponse {
	m := NewMapStore()
	for key, values := range header {
		m.Add(key, strings.Join(values, " "))
	}

	r.header = NewHeader(m)

	return r
}

func (r *Response) Error(err error) IResponse {
	return r.setStatus(err.Error()).
		setStatusCode(http.StatusInternalServerError)
}

func (r *Response) setStatus(status string) IResponse {
	r.status = status
	return r
}

func (r *Response) setStatusCode(status int) IResponse {
	r.statusCode = status
	return r
}

func (r *Response) setRequestBodyData(requestBodyData []byte) IResponse {
	r.requestBodyData = requestBodyData
	return r
}

func (r *Response) setDto(dto any) IResponse {
	r.dto = dto
	return r
}

func (r *Response) setData(data []byte) IResponse {
	r.data = data
	return r
}

func (r *Response) Status() string {
	return r.status
}

func (r *Response) StatusCode() int {
	return r.statusCode
}

func (r *Response) Header() IHeader {
	return r.header
}

func (r *Response) RequestHeader() IHeader {
	return r.requestHeader
}

func (r *Response) RequestMethod() string {
	return r.requestMethod
}

func (r *Response) RequestUri() string {
	return r.requestUri
}

func (r *Response) Data() []byte {
	return r.data
}

func (r *Response) RequestBodyData() []byte {
	return r.requestBodyData
}

func (r *Response) Dto() any {
	return r.dto
}

func (r *Response) Ok() bool {
	return r.statusCode < 400
}
