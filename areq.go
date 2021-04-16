// (c) 2021 Gon Y Yi. 
// https://gonyyi.com

package areq

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const VERSION = "areq/0.2.0"

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

func (r Request) Req(method string, url string, options ...*DoFn) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", VERSION)
	// if close the connection after the request such when there are tons of diff servers to send few req.
	// req.Close = true

	for _, v := range options {
		if v != nil {
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
		if v != nil {
			if v.Res != nil {
				if err := v.Res(res); err != nil {
					return err
				}
			}
		}
	}

	_, err = io.Copy(ioutil.Discard, res.Body)
	return err
}
