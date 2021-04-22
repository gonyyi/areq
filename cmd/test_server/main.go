package main

import (
	"github.com/gonyyi/alog"
	"github.com/gonyyi/alog/ext"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
)

func main() {
	log := alog.New(os.Stderr).SetFormatter(ext.NewFormatterTerminalColor())
	var result string
	http.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(result))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		out, err := httputil.DumpRequest(r, false)
		log.Info().Str("url", r.RequestURI).Write("")
		if err != nil {
			w.WriteHeader(400)
			log.Error().Err("error", err).Write("")
			result = err.Error()
			w.Write([]byte(`{"status":"err","err":` + strconv.Quote(err.Error()) + `}`))
		} else {
			result = string(out)
			w.WriteHeader(200)
			w.Write([]byte(`{"status":"ok"}`))
			println(result+"\n")
		}
	})

	http.ListenAndServe(":8080", nil)
}
