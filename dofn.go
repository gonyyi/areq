// (c) 2021 Gon Y Yi. 
// https://gonyyi.com

package areq

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type DoFn struct {
	Name string
	Req  func(*http.Request) error
	Cli  func(cli *http.Client) error
	Res  func(r *http.Response) error
}

var Do dofns

type dofns struct{}

func (dofns) If(condition bool, do *DoFn) *DoFn {
	if condition {
		return do
	}
	return nil
}

func (dofns) Join(do ...*DoFn) *DoFn {
	var dofn DoFn
	var dofnNames []string
	var req []func(*http.Request) error
	var cli []func(cli *http.Client) error
	var res []func(r *http.Response) error
	for _, f := range do {
		if f.Name != "" {
			dofnNames = append(dofnNames, f.Name)
		} else {
			dofnNames = append(dofnNames, "NONAME-DOFN")
		}
		if f.Req != nil {
			req = append(req, f.Req)
		}
		if f.Cli != nil {
			cli = append(cli, f.Cli)
		}
		if f.Res != nil {
			res = append(res, f.Res)
		}
	}
	dofn.Req = func(r *http.Request) error {
		for _, f := range req {
			if err := f(r); err != nil {
				return err
			}
		}
		return nil
	}
	dofn.Cli = func(c *http.Client) error {
		for _, f := range cli {
			if err := f(c); err != nil {
				return err
			}
		}
		return nil
	}
	dofn.Res = func(r *http.Response) error {
		for _, f := range res {
			if err := f(r); err != nil {
				return err
			}
		}
		return nil
	}
	dofn.Name = "DoJoin(" + strings.Join(dofnNames, ",") + ")"
	return &dofn
}

func (dofns) AuthBasic(id, pwd string) *DoFn {
	return &DoFn{
		Name: "AuthBasic",
		Req: func(req *http.Request) error {
			req.SetBasicAuth(id, pwd)
			return nil
		},
	}
}

func (dofns) SetClientTimeout(duration time.Duration) *DoFn {
	return &DoFn{
		Name: "SetClientTimeout="+duration.String(),
		Cli: func(cli *http.Client) error {
			cli.Timeout = duration
			return nil
		},
	}
}

func (dofns) UseCookieJar(jar http.CookieJar) *DoFn {
	return &DoFn{
		Name: "UseCookieJar",
		Cli: func(cli *http.Client) error {
			cli.Jar = jar
			return nil
		},
	}
}

func (dofns) AddCookie(k, v string) *DoFn {
	return &DoFn{
		Name: "AddCookie",
		Req: func(req *http.Request) error {
			req.AddCookie(&http.Cookie{Name: k, Value: v})
			return nil
		},
	}
}

func (dofns) AddCookies(cookies []*http.Cookie) *DoFn {
	return &DoFn{
		Name: "AddCookies",
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
		Name: "ReqFunc",
		Req:  f,
	}
}

func (dofns) ReqBodyStr(body string) *DoFn {
	return &DoFn{
		Name: "ReqBodyStr",
		Req: func(r *http.Request) error {
			r.Body = ioutil.NopCloser(strings.NewReader(body))
			return nil
		},
	}
}

func (dofns) ReqBody(body io.Reader) *DoFn {
	return &DoFn{
		Name: "ReqBody",
		Req: func(r *http.Request) error {
			r.Body = ioutil.NopCloser(body)
			return nil
		},
	}
}

func (dofns) ReqBodyBytes(body []byte) *DoFn {
	return &DoFn{
		Name: "ReqBodyBytes",
		Req: func(r *http.Request) error {
			r.Body = ioutil.NopCloser(bytes.NewReader(body))
			return nil
		},
	}
}

func (dofns) ReqHeader(h http.Header) *DoFn {
	return &DoFn{
		Name: "ReqHeader",
		Req: func(r *http.Request) error {
			r.Header = h
			return nil
		},
	}
}

func (dofns) ReqHeaderAdd(k, v string) *DoFn {
	return &DoFn{
		Name: "ReqHeaderAdd",
		Req: func(r *http.Request) error {
			r.Header.Add(k, v)
			return nil
		},
	}
}

func (dofns) ReqHeaderSet(k, v string) *DoFn {
	return &DoFn{
		Name: "ReqHeaderSet",
		Req: func(r *http.Request) error {
			r.Header.Set(k, v)
			return nil
		},
	}
}

func (dofns) ResFunc(f func(*http.Response) error) *DoFn {
	return &DoFn{
		Name: "ResFunc",
		Res:  f,
	}
}

func (dofns) GetCookie(name string, dst *string) *DoFn {
	return &DoFn{
		Name: "GetCookie",
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
		Name: "GetCookies",
		Res: func(r *http.Response) error {
			dst = r.Cookies()
			return nil
		},
	}
}

func (dofns) ResJSONTo(dst interface{}) *DoFn {
	return &DoFn{
		Name: "ResJSONTo",
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
		Name: "ResTo",
		Res: func(res *http.Response) error {
			_, err := io.Copy(dst, res.Body)
			return err
		},
	}
}

func (dofns) ResGzipTo(dst io.Writer) *DoFn {
	return &DoFn{
		Name: "ResGzipTo",
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
