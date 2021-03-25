package main

import (
	"github.com/gonyyi/alog"
	"github.com/gonyyi/alog/ext"
	"net/http"
	"net/http/httputil"
	"os"
)

func main() {
	log := alog.New(os.Stderr).SetFormatter(ext.NewFormatterTerminalColor())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		out, err := httputil.DumpRequest(r, false)
		log.Info().Str("url", r.RequestURI).Write("")
		if err!=nil {
			log.Error().Err("error", err).Write("")
		} else {
			println("------DUMP BEGIN------")
			println(string(out)+"\n")
		}
		w.Write([]byte("ok"))
	})

	http.ListenAndServe(":8080", nil)
}
