package option

import (
	"compress/gzip"
	"github.com/gonyyi/areq"
	"io"
	"net/http"
)

func GetResBody(dst io.Writer) *areq.Option {
	return &areq.Option{
		Res: func(res *http.Response) error {
			_, err := io.Copy(dst, res.Body)
			return err
		},
	}
}

func BasicAuth(id, pwd string) *areq.Option {
	return &areq.Option{
		Req: func(req *http.Request) error {
			req.SetBasicAuth(id, pwd)
			return nil
		},
	}
}

func GetResGzip(dst io.Writer) *areq.Option {
	return &areq.Option{
		Req: func(req *http.Request) error {
			req.Header.Add("Accept-Encoding", "gzip")
			return nil
		},
		Res: func(res *http.Response) error {
			_, err := io.Copy(dst, res.Body)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func GetResGzipBody(dst io.Writer) *areq.Option {
	return &areq.Option{
		Req: func(req *http.Request) error {
			req.Header.Add("Accept-Encoding", "gzip")
			return nil
		},
		Res: func(res *http.Response) error {
			// println("encoding: ", res.Header.Get("Content-Encoding"))
			switch res.Header.Get("Content-Encoding") {
			case "gzip":
				tmp, err := gzip.NewReader(res.Body)
				if err != nil {
					return err
				}
				defer tmp.Close()

				_, err = io.Copy(dst, tmp)
				if err != nil {
					return err
				}
			default:
				_, err := io.Copy(dst, res.Body)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}
