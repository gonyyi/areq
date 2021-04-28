package main

import (
	"flag"
	"github.com/gonyyi/alog"
	"github.com/gonyyi/alog/ext"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
)

var NAME = []byte("aReq.Tester 0.2.0 (github.com/gonyyi/areq)\n\n")
const VERSION = "0.2.0"

func main() {

	//
	// CLI FLAGS
	//
	var addr string
	var respUri string
	var toIgnore string
	var logFile string
	var verbose bool
	flag.StringVar(&addr, "addr", ":8080", "addr to serve (eg. :8080)")
	flag.StringVar(&respUri, "res", "res", "response URI")
	flag.StringVar(&toIgnore, "ignore", "/abc.ico", "comma delimited URI to be ignored")
	flag.StringVar(&logFile, "log", "", "log file name if to store a log into a file")
	flag.BoolVar(&verbose, "verbose", false, "verbose log")
	flag.Parse()

	//
	// LOGGER
	//
	log := alog.New(nil)
	if verbose {
		log.Control.Level = alog.TraceLevel
	}
	if logFile != "" {
		out, err := os.Create(logFile)
		if err != nil {
			println(err.Error())
			return
		}
		log = log.SetOutput(out)
	} else {
		log = log.SetOutput(os.Stderr).SetFormatter(ext.NewFormatterTerminalColor())
	}
	log.Info().Str("library", "github.com/gonyyi/areq").Str("version", VERSION).Write("Starting aReq.Tester")

	//
	// IGNORE CERTAIN URL
	//
	if !strings.HasPrefix(respUri, "/") {
		respUri = "/" + respUri
	}
	ignores := make(map[string]struct{})
	for idx, v := range strings.Split(toIgnore, ",") {
		log.Debug().Int("id", idx).Str("ignoreURI", v).Write()
		ignores[strings.TrimSpace(v)] = struct{}{}
	}
	ignores[respUri] = struct{}{} // to make sure response won't be recorded..
	log.Info().Str("ignoreURI", toIgnore).Str("respURI", respUri).Int("totalCount", len(ignores)).Write("IgnoreURI")

	shouldIgnore := func(r *http.Request) bool {
		// r.RequestURI
		if _, ok := ignores[r.RequestURI]; ok {
			return true
		}
		return false
	}

	//
	// STORE MOST RECENT RESULT
	//
	var lastURI string
	var lastIP string
	var lastResp string

	//
	// HTTPS
	//
	http.HandleFunc(respUri, func(w http.ResponseWriter, r *http.Request) {
		// For response page, if to be ignored, just not writing it to a log.
		if !shouldIgnore(r) {
			log.Info().Str("ip", r.RemoteAddr).Str("uri", r.RequestURI).Write()
		}
		w.WriteHeader(200)
		w.Write(NAME)
		w.Write([]byte("URI: " + lastURI + "\nIP:  " + lastIP + "\n\n"))
		w.Write([]byte(lastResp))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if shouldIgnore(r) {
			log.Trace().Str("uri", r.RequestURI).Str("ip", r.RemoteAddr).Write("IGNORED")
			return
		}
		lastURI = r.RequestURI
		lastIP = r.RemoteAddr

		out, err := httputil.DumpRequest(r, false)
		if err != nil {
			w.WriteHeader(400)
			log.Error().Str("uri", r.RequestURI).Str("ip", r.RemoteAddr).Err(err).Write()

			w.Write(NAME)
			w.Write([]byte("URI: " + r.RequestURI + "\nIP:  " + r.RemoteAddr + "\n\n"))
			w.Write([]byte(strconv.Quote(err.Error())))

			lastResp = err.Error()
		} else {
			w.WriteHeader(200)
			log.Info().Str("uri", r.RequestURI).Str("ip", r.RemoteAddr).Write("OK")

			w.Write(NAME)
			w.Write([]byte("URI: " + r.RequestURI + "\nIP:  " + r.RemoteAddr + "\n\n"))
			w.Write(out)

			lastResp = string(out)
		}
	})

	log.Info().Str("addr", addr).Write("Starting")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error().Str("addr", addr).Err(err).Write()
	}
}
