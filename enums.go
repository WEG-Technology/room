package room

type HTTPMethod string

const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	PATCH  HTTPMethod = "PATCH"
	DELETE HTTPMethod = "DELETE"
	HEAD   HTTPMethod = "HEAD"
)

func (s HTTPMethod) String() string {
	switch s {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case PATCH:
		return "PATCH"
	case DELETE:
		return "DELETE"
	case HEAD:
		return "HEAD"
	}
	return "GET"
}

type HTTPProtocol int

const (
	Http HTTPProtocol = iota + 1
	Https
)

func (h HTTPProtocol) String() string {
	switch h {
	case Http:
		return "http://"
	case Https:
		return "https://"
	}
	return "https://"
}
