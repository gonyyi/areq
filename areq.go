// (c) 2021 Gon Y Yi. 
// https://gonyyi.com
// 07/14/21 Wed 17:00

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
	VERSION = "areq/0.2.2"
)

type Request struct {
	Client    *http.Client
	Transport *http.Transport
}

func (r Request) Init() Request {
	r.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: DefaultSetting.TLSInsecureSkipVerify,
		},
		DialContext: (&net.Dialer{
			Timeout:   DefaultSetting.DialTimeout,
			KeepAlive: DefaultSetting.DialKeepAlive,
		}).DialContext,
		ForceAttemptHTTP2:   DefaultSetting.ForceAttemptHTTP2,
		TLSHandshakeTimeout: DefaultSetting.HandshakeTimeout,
		DisableKeepAlives:   DefaultSetting.DisableKeepAlives,
		MaxIdleConns:        DefaultSetting.MaxIdleConns,
		MaxIdleConnsPerHost: DefaultSetting.MaxIdleConnsPerHost,
	}
	r.Client = &http.Client{
		Timeout:   DefaultSetting.ClientTimeout,
		Transport: r.Transport,
	}
	return r
}

// SetTimeout will set various timeout setting of request.
// If set to 0, it will use default value.
// If set to -1, it means timeout will not be used.
func (r Request) SetTimeout(dialTimeout, dialKeepAlive, handshakeTimeout, clientTimeout time.Duration) Request {
	if dialTimeout == 0 {
		dialTimeout = DefaultSetting.DialTimeout
	}
	if dialKeepAlive == 0 {
		dialKeepAlive = DefaultSetting.DialKeepAlive
	}
	if handshakeTimeout == 0 {
		handshakeTimeout = DefaultSetting.HandshakeTimeout
	}
	if clientTimeout == 0 {
		clientTimeout = DefaultSetting.ClientTimeout
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
