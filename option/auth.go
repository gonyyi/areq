package option

import (
	"github.com/gonyyi/areq"
	"net/http"
)

func AuthBasic(id, pwd string) *areq.Option {
	return &areq.Option{
		Req: func(req *http.Request) error {
			req.SetBasicAuth(id, pwd)
			return nil
		},
	}
}
