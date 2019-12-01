package broker

import (
	"crypto/subtle"
	"net/http"
)

func (b *Broker) BasicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if !b.verifyBasicAuth(rw, req) {
			return
		}
		handler(rw, req)
	}
}

func (b *Broker) verifyBasicAuth(rw http.ResponseWriter, req *http.Request) bool {
	username, password, ok := req.BasicAuth()
	if !ok || subtle.ConstantTimeCompare([]byte(username), []byte(b.Username)) != 1 || subtle.ConstantTimeCompare([]byte(password), []byte(b.Password)) != 1 {
		rw.Header().Set("WWW-Authenticate", `Basic realm="compose.io"`)
		b.Error(rw, req, 401, "Unauthorized", "You are not authorized to access this service broker")
		return false
	}
	return true
}
