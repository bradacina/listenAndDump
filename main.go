package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type bodyModifier struct {
	rw http.ResponseWriter
}

func (bm *bodyModifier) Header() http.Header {
	return bm.rw.Header()
}

func (bm *bodyModifier) WriteHeader(status int) {
	bm.rw.WriteHeader(status)
}

func (bm *bodyModifier) Write(content []byte) (int, error) {
	log.Printf("Response Body: %v\n", string(content))
	if strings.Contains(string(content), "00000000-0000-0000-0000-000000000000") {
		log.Printf("FOUND THE TOKEN\n")
		content = []byte(strings.Replace(string(content), "00000000-0000-0000-0000-000000000000", "656e2145-f118-445b-999d-9b5e8faaaaff", -1))
	}
	return bm.rw.Write(content)
}

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

var p *httputil.ReverseProxy

func main() {
	f, err := os.Create("listenAndDump.log")
	defer f.Close()
	log.SetOutput(f)

	url, err := url.Parse("http://knockknock.readify.net/RedPill.svc")
	if err != nil {
		log.Fatal(err)
	}

	p = NewProxy(url)

	http.HandleFunc("/", editResponse)

	for {
		err := http.ListenAndServe(":80", nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func editResponse(w http.ResponseWriter, r *http.Request) {
	bm := &bodyModifier{rw: w}
	p.ServeHTTP(bm, r)
}
