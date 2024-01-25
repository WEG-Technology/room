package room

import (
	"context"
	"net/http"
)

const (
	headerKeyContentType       = "Content-Type"
	headerKeyAccept            = "Accept"
	headerValueFormEncoded     = "application/x-www-form-urlencoded"
	headerValueApplicationJson = "application/json"
	headerValueTextXML         = "text/xml"
)

type IRequest interface {
	Create() IRequest
	Cancel()
	Request() *http.Request
	Dto() any
	Header() IHeader
	SetConnectionConfig(connectionConfig IConnectionConfig) IRequest
	SetHeaders(headers IMap) IRequest
	PutHeaderProperties(header IHeader) IRequest
	InjectHeader()
	NewRequestWithContext() (err error)
	InitContext()
	QueryString() string
	url() string
}

type IRequestObserver interface {
	OnCreate(request IRequest)
	OnCreated(request IRequest)
}

type requestObserver struct{}

func (o requestObserver) OnCreated(request IRequest) {}
func (o requestObserver) OnCreate(request IRequest)  {}

type BaseRequest struct {
	method           HTTPMethod
	R                *http.Request
	ctx              context.Context
	cancel           context.CancelFunc
	header           IHeader
	query            IQuery
	connectionConfig IConnectionConfig
	observer         IRequestObserver
	contextBuilder   IContextBuilder
	bodyParser       IBodyParser
	endPoint         string
	dto              any
}

func (r *BaseRequest) Cancel() {
	if r.cancel != nil {
		r.cancel()
	}
}

func newRequest(opts ...OptionRequest) (IRequest, error) {
	r := new(BaseRequest)

	for _, opt := range opts {
		opt(r)
	}

	if r.header == nil {
		r.header = NewHeader(NewMapStore())
	}

	if r.query == nil {
		r.query = NewQuery(NewMapStore())
	}

	if r.observer == nil {
		r.observer = requestObserver{}
	}

	if r.bodyParser == nil {
		r.bodyParser = dumpBody{}
	}

	//TODO handle endpoint when its nil
	if r.connectionConfig == nil {
		r.connectionConfig = DefaultConnectionConfig()
	}

	if r.contextBuilder == nil {
		r.contextBuilder = NewContextBuilder(r.connectionConfig.TimeoutDuration())
	}

	return r, nil
}

func NewGetRequest(opts ...OptionRequest) (IRequest, error) {
	return newRequest(opts...)
}

func NewPostRequest(opts ...OptionRequest) (IRequest, error) {
	return newRequest(append(opts, WithMethod(POST))...)
}

func (r *BaseRequest) Create() IRequest {
	r.observer.OnCreate(r)

	r.InitContext()

	r.header.Merge(r.connectionConfig.DefaultHeader())

	err := r.NewRequestWithContext()

	if err != nil {
		panic(err)
	}

	r.InjectHeader()

	r.observer.OnCreated(r)

	return r
}

func (r *BaseRequest) InitContext() {
	r.ctx, r.cancel = r.contextBuilder.Build()
}

func (r *BaseRequest) InjectHeader() {
	r.header.Properties().Each(func(key string, value any) {
		r.R.Header.Add(key, value.(string))
	})
}

func (r *BaseRequest) NewRequestWithContext() (err error) {
	r.R, err = http.NewRequestWithContext(r.ctx, r.method.String(), r.url(), r.bodyParser.Parse())

	return err
}

func (r *BaseRequest) Request() *http.Request {
	return r.R
}

func (r *BaseRequest) Dto() any {
	return r.dto
}

func (r *BaseRequest) Header() IHeader {
	return r.header
}

func (r *BaseRequest) SetConnectionConfig(connectionConfig IConnectionConfig) IRequest {
	r.connectionConfig = connectionConfig
	return r
}

func (r *BaseRequest) SetHeaders(headers IMap) IRequest {
	headers.Each(func(key string, value any) {
		r.R.Header.Add(key, value.(string))
	})

	return r
}

func (r *BaseRequest) PutHeaderProperties(header IHeader) IRequest {
	r.header.Merge(header)
	return r
}

func (r *BaseRequest) url() string {
	return r.connectionConfig.BaseUrl() + r.endPoint + r.QueryString()
}

func (r *BaseRequest) QueryString() string {
	str := r.query.String()

	if str == "" {
		return ""
	}

	return "?" + r.query.String()
}

type OptionRequest func(request *BaseRequest)

func WithEndPoint(endPoint string) OptionRequest {
	return func(request *BaseRequest) {
		request.endPoint = endPoint
	}
}

func WithMethod(method HTTPMethod) OptionRequest {
	return func(request *BaseRequest) {
		request.method = method
	}
}

func WithBody(bodyParser IBodyParser) OptionRequest {
	return func(request *BaseRequest) {
		request.bodyParser = bodyParser
	}
}

func WithQuery(query IQuery) OptionRequest {
	return func(request *BaseRequest) {
		request.query = query
	}
}

func WithHeader(header IHeader) OptionRequest {
	return func(request *BaseRequest) {
		request.header = header
	}
}

func WithDto(dto any) OptionRequest {
	return func(request *BaseRequest) {
		request.dto = dto
	}
}

func WithContextBuilder(contextBuilder IContextBuilder) OptionRequest {
	return func(request *BaseRequest) {
		request.contextBuilder = contextBuilder
	}
}
