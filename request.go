package room

import (
	"net/http"
	"strings"
	"time"
)

const (
	headerKeyContentType       = "Content-Type"
	headerKeyAccept            = "Accept"
	headerValueFormEncoded     = "application/x-www-form-urlencoded"
	headerValueApplicationJson = "application/json"
	headerValueTextXML         = "text/xml"
)

type ISend interface {
	Send() Response
}

type Request struct {
	path           string
	uri            URI
	method         HTTPMethod
	Header         IHeader
	Query          IQuery
	BodyParser     IBodyParser
	contextBuilder IContextBuilder
	DTO            any
	// ForceDTO is a flag to force the parse http response to DTO which is map[string]any
	ForceDTO bool
}

// NewRequest creates a new request
// baseUrl: the base url ex: http://localhost:8080, http://localhost/path, path/, path?lorem=ipsum
// opts: options to configure the request
func NewRequest(path string, opts ...OptionRequest) *Request {
	r := &Request{
		path: path,
	}

	for _, opt := range opts {
		opt(r)
	}

	if r.BodyParser == nil {
		r.BodyParser = dumpBody{}
	}

	if r.method == "" {
		r.method = GET
	}

	return r
}

func (r *Request) Send() (Response, error) {
	c := new(http.Client)

	response, e := c.Do(r.request())

	if e != nil {
		return Response{}, e
	}

	return NewResponse(response, r.ForceDTO)
}

func (r *Request) request() *http.Request {
	var context Context

	if r.contextBuilder != nil {
		context = r.contextBuilder.Build()
	} else {
		context = NewContextBuilder(30 * time.Second).Build()
	}

	if r.Query != nil && r.Query.String() != "" {
		r.uri = NewURI(r.path + "?" + r.Query.String())
	} else {
		r.uri = NewURI(r.path)
	}

	req, _ := http.NewRequestWithContext(context.ctx, r.method.String(), r.uri.String(), r.BodyParser.Parse())

	if r.Header != nil {
		r.Header.Properties().Each(func(k string, v any) {
			req.Header.Add(k, v.(string))
		})
	}

	return req
}

func (r *Request) initBaseUrl(baseUrl string) *Request {
	if strings.HasPrefix(r.path, "/") {
		r.path = r.path[1:]
	}

	if strings.HasSuffix(baseUrl, "/") {
		baseUrl = baseUrl[:len(baseUrl)-1]
	}

	r.path = baseUrl + "/" + r.path

	return r
}

func (r *Request) mergeHeader(header IHeader) *Request {
	if header != nil {
		if r.Header == nil {
			r.Header = header
		} else {
			r.Header.Merge(header)
		}
	}

	return r
}

type OptionRequest func(request *Request)

func WithMethod(method HTTPMethod) OptionRequest {
	return func(request *Request) {
		request.method = method
	}
}

func WithBody(bodyParser IBodyParser) OptionRequest {
	return func(request *Request) {
		request.BodyParser = bodyParser
	}
}

func WithQuery(query IQuery) OptionRequest {
	return func(request *Request) {
		request.Query = query
	}
}

func WithHeader(header IHeader) OptionRequest {
	return func(request *Request) {
		request.Header = header
	}
}

func WithDTO(dto any) OptionRequest {
	return func(request *Request) {
		request.DTO = dto
	}
}

func ForceDTO() OptionRequest {
	return func(request *Request) {
		request.ForceDTO = true
	}
}

func WithContextBuilder(contextBuilder IContextBuilder) OptionRequest {
	return func(request *Request) {
		request.contextBuilder = contextBuilder
	}
}
