package areq_test

import "testing"

func TestRequest_ReqBody(t *testing.T) {
	b := "hello world"
	out := bytes.Buffer{}
	// var out []byte
	err := areq.Req("POST", "https://httpbin.org/post", areq.Do.ReqBody([]byte(b)), areq.Do.ResTo(&out))
	if err != nil {println(err.Error())}
	println("ok")
	println(out.String())
}
