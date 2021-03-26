package option

import (
	"github.com/gonyyi/areq"
	"net/http"
)

func ReqFunc(f func(r *http.Request)error) *areq.Option {
	return &areq.Option{
		Req: f,
	}
}

func ReqHeader(h http.Header) *areq.Option {
	return &areq.Option{
		Req: func(r *http.Request) error {
			r.Header = h
			return nil
		},
	}
}

func ReqHeaderAdd(k, v string) *areq.Option {
	return &areq.Option{
		Req: func(r *http.Request) error {
			r.Header.Add(k, v)
			return nil
		},
	}
}

func ReqHeaderSet(k, v string) *areq.Option {
	return &areq.Option{
		Req: func(r *http.Request) error {
			r.Header.Set(k, v)
			return nil
		},
	}
}
