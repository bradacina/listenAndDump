package main

import (
	"bufio"
	"log"
	"net/http"
	"strings"
)

func main() {

	http.HandleFunc("/", dump)

	for {
		err := http.ListenAndServe(":8080", nil)
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
