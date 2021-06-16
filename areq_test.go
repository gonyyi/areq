package areq_test

import (
	"bytes"
	"github.com/gonyyi/areq"
	"strings"
	"testing"
)

func TestRequest_ReqBodyBytes(t *testing.T) {
	b := "hello world"
	out := bytes.Buffer{}
	// var out []byte
	err := areq.Req("POST", "https://httpbin.org/post", areq.Do.ReqBodyBytes([]byte(b)), areq.Do.ResTo(&out))
	if err != nil {println(err.Error())}
	println("ok")
	println(out.String())
}

func TestRequest_ReqBody(t *testing.T) {
	b := "hello world"
	out := bytes.Buffer{}
	// var out []byte
	br := strings.NewReader(b)
	err := areq.Req("POST", "https://httpbin.org/post", areq.Do.ReqBody(br), areq.Do.ResTo(&out))
	if err != nil {println(err.Error())}
	println("ok")
	println(out.String())
}

func TestRequest_DoJoin(t *testing.T) {
	b := "hello world"
	out := bytes.Buffer{}
	// var out []byte
	br := strings.NewReader(b)
	dofns := areq.DoJoin( areq.Do.ReqBody(br), areq.Do.ResTo(&out))
	println(dofns.Name)
	err := areq.Req("POST", "https://httpbin.org/post", dofns)
	if err != nil {println(err.Error())}
	println("ok")
	println(out.String())
}
