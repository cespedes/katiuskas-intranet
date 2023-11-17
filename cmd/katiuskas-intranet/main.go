package main

import (
	"log"
	"net/http"
)

func (s *server) MyHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create new context from HTTP session:
		r = s.NewContext(r)
		h.ServeHTTP(w, r)
	})
}

func main() {
	s := NewServer()

	err := http.ListenAndServe(s.config["http_listen_addr"], s.MyHandler(s))
	log.Fatal(err)
}
