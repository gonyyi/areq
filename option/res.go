package option

import (
	"compress/gzip"
	"encoding/json"
	"github.com/gonyyi/areq"
	"io"
	"net/http"
)

func ResFunc(f func(*http.Response)error) *areq.Option {
	return &areq.Option{
		Res: f,
	}
}

func ResJSONTo(dst interface{}) *areq.Option {
	return &areq.Option{
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

func ResTo(dst io.Writer) *areq.Option {
	return &areq.Option{
		Res: func(res *http.Response) error {
			_, err := io.Copy(dst, res.Body)
			return err
		},
	}
}

func ResGzipTo(dst io.Writer) *areq.Option {
	return &areq.Option{
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
