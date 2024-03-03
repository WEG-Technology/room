package room

import (
	"github.com/google/go-querystring/query"
	"net/url"
)

type Query struct {
	v any
}

type IQuery interface {
	String() string
}

func NewQuery(v any) IQuery {
	switch v := v.(type) {
	case IMap:
		return IMapQuery{v}
	default:
		return IUrlQuery{v}
	}
}

type IMapQuery struct {
	v IMap
}

func (q IMapQuery) String() string {
	v := url.Values{}

	q.v.Each(func(key string, value any) {
		v.Add(key, value.(string))
	})

	vEncoded := v.Encode()

	return vEncoded
}

type IUrlQuery struct {
	v any
}

func (q IUrlQuery) String() string {
	urlValues, err := query.Values(q.v)

	if err != nil {
		return ""
	}

	return urlValues.Encode()
}
