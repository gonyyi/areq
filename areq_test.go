package areq_test

import (
	"bytes"
	"github.com/gonyyi/areq"
	"github.com/gonyyi/areq/option"
	"testing"
)

func TestRequest_Req(t *testing.T) {
	var out bytes.Buffer

	if err := areq.Req("GET", "https://gonyyi.com/copyright", option.GetResBody(&out)); err != nil {
		println(err)
	}

	println(out.String())
}
