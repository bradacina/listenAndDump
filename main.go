package main

import (
	"os"
	"bufio"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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
		log.Printf("%v",string(dump))
		req.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		dump, _ = httputil.DumpRequestOut(req, true)
		log.Printf("%v",string(dump))
	}

	return &httputil.ReverseProxy{Director: director}
}

func main() {
	f, err := os.Create("listenAndDump.log")
	defer f.Close()
	log.SetOutput(f)

	url, err := url.Parse("http://knockknock.readify.net/RedPill.svc")
	if err != nil {
		log.Fatal(err)
	}


	proxy := NewProxy(url)

	http.Handle("/", proxy)
	//http.HandleFunc("/", dump)

	for {
		err := http.ListenAndServe(":80", nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func dump(w http.ResponseWriter, r *http.Request) {

	reader := bufio.NewReader(r.Body)
	log.Println("Method", r.Method)
	log.Println("Header:")
	for k, v := range r.Header {
		log.Printf("%s : %v\n", k, v)
	}
	log.Println("RequestUri", r.RequestURI)
	log.Println("Body:")
	for {
		s, err := reader.ReadString('\n')
		s = strings.Trim(s, "\r\n")
		if err != nil {
			log.Println(s)
			break
		}

		log.Println(s)
	}
}
