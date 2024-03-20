package room

import (
	"strings"
)

type URI struct {
	scheme    string
	authority string
	path      string
	query     string
}

func (u URI) String() string {
	url := u.scheme + "://" + u.authority + u.path

	if u.query != "" {
		url += "?" + u.query
	}

	return url
}

func (u URI) Query() string {
	return u.query
}

func (u URI) Path() string {
	return u.path
}

func (u URI) Authority() string {
	return u.authority
}

func (u URI) Scheme() string {
	return u.scheme
}

func NewURI(fullUrl string) URI {
	uri := URI{}

	if strings.Contains(fullUrl, "?") {
		splittedURL := strings.SplitN(fullUrl, "?", 2)

		fullUrl = splittedURL[0]
		// TODO use IStore instead of string
		uri.query = splittedURL[1]
	}

	var urn string
	if strings.HasPrefix(fullUrl, "https://") {
		uri.scheme = "https"
		urn = strings.TrimPrefix(fullUrl, "https://")
	} else {
		uri.scheme = "http"
		urn = strings.TrimPrefix(fullUrl, "http://")
	}

	if len(urn) > 0 && urn[len(urn)-1] != '/' {
		urn = urn + "/"
	}

	splittedURN := strings.SplitN(urn, "/", 2)

	uri.authority = splittedURN[0]
	uri.path = "/" + splittedURN[1]

	if len(uri.path) > 0 && uri.path[len(uri.path)-1] == '/' {
		uri.path = uri.path[:len(uri.path)-1]
	}

	return uri
}
