// (c) 2021 Gon Y Yi. 
// https://gonyyi.com

package areq

// Err is a simple error format compatible as a string
type Err string

func (e Err) Error() string {
	return string(e)
}

var globalReq = NewRequest()

// Req is a global request short cut option
func Req(method string, url string, do ...*DoFn) error {
	return globalReq.Req(method, url, do...)
}
