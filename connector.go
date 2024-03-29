package room

type Connector struct {
	baseUrl        string
	Header         IHeader
	contextBuilder IContextBuilder
	authRequest    *Request
}

type OptionConnector func(info *Connector)

func WithHeaderConnector(header IHeader) OptionConnector {
	return func(request *Connector) {
		request.Header = header
	}
}

func WithHeaderContextBuilder(ctxBuilder IContextBuilder) OptionConnector {
	return func(request *Connector) {
		request.contextBuilder = ctxBuilder
	}
}

func NewConnector(baseUrl string, opts ...OptionConnector) *Connector {
	c := &Connector{
		baseUrl: NewURI(baseUrl).String(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Connector) Send(path string) (Response, error) {
	return c.Do(NewRequest(path))
}

func (c *Connector) Do(request *Request) (Response, error) {
	request.
		initBaseUrl(c.baseUrl).
		mergeHeader(c.Header)

	if c.contextBuilder != nil {
		request.contextBuilder = c.contextBuilder
	}

	return request.Send()
}
