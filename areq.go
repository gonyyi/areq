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

const (
	VERSION                 = "areq/0.2.2"
	DefaultDialTimeout      = 5 * time.Second
	DefaultDialKeepAlive    = 600 * time.Second
	DefaultHandshakeTimeout = 5 * time.Second
	DefaultClientTimeout    = 10 * time.Second
)

type Request struct {
	Client    *http.Client
	Transport *http.Transport
}

func (r Request) Init() Request {
	r.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DialContext: (&net.Dialer{
			Timeout:   DefaultDialTimeout,
			KeepAlive: DefaultDialKeepAlive,
		}).DialContext,
		// ForceAttemptHTTP2:     true,
		TLSHandshakeTimeout: DefaultHandshakeTimeout,
		DisableKeepAlives:   false,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}
	r.Client = &http.Client{
		Timeout:   DefaultClientTimeout,
		Transport: r.Transport,
	}
	return r
}

// SetTimeout will set various timeout setting of request.
// If set to -1, it will use default value.
// If set to 0, it means timeout will not be used.
func (r Request) SetTimeout(dialTimeout, dialKeepAlive, handshakeTimeout, clientTimeout time.Duration) Request {
	if dialTimeout == -1 {
		dialTimeout = DefaultDialTimeout
	}
	if dialKeepAlive == -1 {
		dialKeepAlive = DefaultDialKeepAlive
	}
	if handshakeTimeout == -1 {
		handshakeTimeout = DefaultHandshakeTimeout
	}
	if clientTimeout == -1 {
		clientTimeout = DefaultClientTimeout
	}

	if r.Transport == nil || r.Client == nil {
		r = r.Init()
	}

	r.Transport.DialContext = (&net.Dialer{
		Timeout:   dialTimeout,
		KeepAlive: dialKeepAlive,
	}).DialContext
	r.Transport.TLSHandshakeTimeout = handshakeTimeout

	r.Client.Timeout = clientTimeout

	return r
}

func NewRequest() Request {
	return Request{}.Init()
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
