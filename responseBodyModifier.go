package main

import (
	"net/http"
	"strings"
)

const (
	// RedPillToken is the token returned by the RedPill.svc
	RedPillToken = "00000000-0000-0000-0000-000000000000"

	// OurToken is the token that our service is supposed to return
	OurToken = "656e2145-f118-445b-999d-9b5e8faaaaff"
)

// ResponseBodyModifier modifies the body of a Response replacing the token returned by
// the RedPill.svc with our own token
type ResponseBodyModifier struct {
	rw http.ResponseWriter
}

func (m *ResponseBodyModifier) Header() http.Header {
	return m.rw.Header()
}

func (m *ResponseBodyModifier) WriteHeader(status int) {
	m.rw.WriteHeader(status)
}

func (m *ResponseBodyModifier) Write(content []byte) (int, error) {
	contentString := string(content)
	if strings.Contains(contentString, RedPillToken) {
		contentString = strings.Replace(contentString, RedPillToken, OurToken, -1)
		content = []byte(contentString)
	}
	return m.rw.Write(content)
}
