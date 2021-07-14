// (c) 2021 Gon Y Yi. 
// https://gonyyi.com
// 07/14/21 Wed 16:59

package areq

import "time"

var globalReq = NewRequest()

// Req is a global request short cut option
func Req(method string, url string, do ...*DoFn) error {
	return globalReq.Req(method, url, do...)
}

// Err is a simple error format compatible as a string
type Err string

func (e Err) Error() string {
	return string(e)
}

// DefaultSetting will let the user change the default behavior
var DefaultSetting = defaults{
	DialTimeout:         5 * time.Second,
	DialKeepAlive:       600 * time.Second,
	HandshakeTimeout:    5 * time.Second,
	ClientTimeout:       10 * time.Second,
	DisableKeepAlives:   false,
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 100,
}

type defaults struct {
	DialTimeout           time.Duration
	DialKeepAlive         time.Duration
	HandshakeTimeout      time.Duration
	ClientTimeout         time.Duration
	TLSInsecureSkipVerify bool
	ForceAttemptHTTP2     bool
	DisableKeepAlives     bool
	MaxIdleConns          int
	MaxIdleConnsPerHost   int
}

// UpdateGlobal will update the current default setting into global Req method
func (defaults) UpdateGlobal() {
	globalReq = globalReq.Init()
}
