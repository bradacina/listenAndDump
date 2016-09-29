package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

const (
	// ProxyURL is the url for which we're proxying requests
	ProxyURL = "http://knockknock.readify.net/RedPill.svc"
)

func singleJoiningSlash(a, b string) string {
	if len(b) == 0 || b == "/" {
		return a
	}

	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func NewProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {

		dump, _ := httputil.DumpRequest(req, true)
		log.Printf("%v", string(dump))

		// rewrite the incoming request
		req.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		req.Header.Set("Accept-Encoding","")
		dump, _ = httputil.DumpRequestOut(req, true)
		log.Printf("%v", string(dump))
	}

	return &httputil.ReverseProxy{Director: director}
}

var proxy *httputil.ReverseProxy

func editResponse(w http.ResponseWriter, r *http.Request) {
	rbm := &ResponseBodyModifier{rw: w}
	proxy.ServeHTTP(rbm, r)
}

func main() {
	f, err := os.Create("listenAndDump.log")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	log.SetOutput(f)

	url, err := url.Parse(ProxyURL)
	if err != nil {
		log.Fatal(err)
	}

	proxy = NewProxy(url)

	http.HandleFunc("/", editResponse)

	for {
		err := http.ListenAndServe(":80", nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
