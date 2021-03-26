// (c) Gon Yi

package areq

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var globalReq = NewRequest()

func Req(method string, url string, options ...*Option) error {
	return globalReq.Req(method, url, options...)
}

type Option struct {
	Req   func(*http.Request) error
	Cli   func(cli *http.Client) error
	Res   func(r *http.Response) error
}

type Request struct {
	Client    *http.Client
	Transport *http.Transport
}

func NewRequest() *Request {
	r := Request{}
	r.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DialContext: (&net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 600 * time.Second,
		}).DialContext,
		// ForceAttemptHTTP2:     true,
		DisableKeepAlives:   false,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}
	r.Client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: r.Transport,
	}
	return &r
}

func (r *Request) Req(method string, url string, options ...*Option) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	// if close the connection after the request such when there are tons of diff servers to send few req.
	// req.Close = true

	for _, v := range options {
		if v.Req != nil {
			if err := v.Req(req); err != nil {
				return err
			}
		}
		if v.Cli != nil {
			if err := v.Cli(r.Client); err != nil {
				return err
			}
		}
	}

	res, err := r.Client.Do(req)
	// if body.Close is after the error check, when redirect, res won't be nil.
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return err
	}

	for _, v := range options {
		if v.Res != nil {
			if err := v.Res(res); err != nil {
				return err
			}
		}
	}

	_, err = io.Copy(ioutil.Discard, res.Body)
	return err
}
