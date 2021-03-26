package option

import (
	"github.com/gonyyi/areq"
	"net/http"
)

func UseCookieJar(jar http.CookieJar) *areq.Option {
	return &areq.Option{
		Cli: func(cli *http.Client) error {
			cli.Jar = jar
			return nil
		},
	}
}

func AddCookie(k, v string) *areq.Option {
	return &areq.Option{
		Req: func(req *http.Request) error {
			req.AddCookie(&http.Cookie{Name: k, Value: v})
			return nil
		},
	}
}

func AddCookies(cookies []*http.Cookie) *areq.Option {
	return &areq.Option{
		Req: func(req *http.Request) error {
			for i:=0; i<len(cookies); i++ {
				req.AddCookie(cookies[i])
			}
			return nil
		},
	}
}
