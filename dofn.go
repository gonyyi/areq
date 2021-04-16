// (c) 2021 Gon Y Yi. 
// https://gonyyi.com

package areq

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
)

type DoFn struct {
	Req func(*http.Request) error
	Cli func(cli *http.Client) error
	Res func(r *http.Response) error
}

func DoIf(condition bool, do *DoFn) *DoFn {
	if condition {
		return do
	}
	return nil
}

var Do dofns

type dofns struct{}

func (dofns) AuthBasic(id, pwd string) *DoFn {
	return &DoFn{
		Req: func(req *http.Request) error {
			req.SetBasicAuth(id, pwd)
			return nil
		},
	}
}

func (dofns) UseCookieJar(jar http.CookieJar) *DoFn {
	return &DoFn{
		Cli: func(cli *http.Client) error {
			cli.Jar = jar
			return nil
		},
	}
}

func (dofns) AddCookie(k, v string) *DoFn {
	return &DoFn{
		Req: func(req *http.Request) error {
			req.AddCookie(&http.Cookie{Name: k, Value: v})
			return nil
		},
	}
}

func (dofns) AddCookies(cookies []*http.Cookie) *DoFn {
	return &DoFn{
		Req: func(req *http.Request) error {
			for i := 0; i < len(cookies); i++ {
				req.AddCookie(cookies[i])
			}
			return nil
		},
	}
}

func (dofns) ReqFunc(f func(r *http.Request) error) *DoFn {
	return &DoFn{
		Req: f,
	}
}

func (dofns) ReqHeader(h http.Header) *DoFn {
	return &DoFn{
		Req: func(r *http.Request) error {
			r.Header = h
			return nil
		},
	}
}

func (dofns) ReqHeaderAdd(k, v string) *DoFn {
	return &DoFn{
		Req: func(r *http.Request) error {
			r.Header.Add(k, v)
			return nil
		},
	}
}

func (dofns) ReqHeaderSet(k, v string) *DoFn {
	return &DoFn{
		Req: func(r *http.Request) error {
			r.Header.Set(k, v)
			return nil
		},
	}
}

func (dofns) ResFunc(f func(*http.Response) error) *DoFn {
	return &DoFn{
		Res: f,
	}
}

func (dofns) GetCookie(name string, dst *string) *DoFn {
	return &DoFn{
		Res: func(r *http.Response) error {
			for _, v := range r.Cookies() {
				if v.Name == name {
					*dst = v.Value
				}
			}
			return nil
		},
	}
}

func (dofns) GetCookies(dst []*http.Cookie) *DoFn {
	return &DoFn{
		Res: func(r *http.Response) error {
			dst = r.Cookies()
			return nil
		},
	}
}

func (dofns) ResJSONTo(dst interface{}) *DoFn {
	return &DoFn{
		Res: func(res *http.Response) error {
			var dec *json.Decoder
			switch res.Header.Get("Content-Encoding") {
			case "gzip":
				tmp, err := gzip.NewReader(res.Body)
				if err != nil {
					return err
				}
				defer tmp.Close()
				dec = json.NewDecoder(tmp)
			default:
				dec = json.NewDecoder(res.Body)
			}

			// maybe need dec.Token() before dec.More().. not sure when..
			for dec.More() {
				err := dec.Decode(dst)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func (dofns) ResTo(dst io.Writer) *DoFn {
	return &DoFn{
		Res: func(res *http.Response) error {
			_, err := io.Copy(dst, res.Body)
			return err
		},
	}
}

func (dofns) ResGzipTo(dst io.Writer) *DoFn {
	return &DoFn{
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
