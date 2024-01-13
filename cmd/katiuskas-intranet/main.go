package main

import (
	"log"
	"net/http"
	"os"
)

func (s *server) MyHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create new context from HTTP session and token:
		r = s.NewContext(r)
		h.ServeHTTP(w, r)
	})
}

func main() {
	s := NewServer()

	log.Println("katiuskas-intranet starting")

	err := http.ListenAndServe(s.config["http_listen_addr"], s.MyHandler(s))
	log.Fatal(err)
}

func init() {
	// if running from systemd, do not show date and time in logs:
	if os.Getenv("INVOCATION_ID") != "" {
		log.SetFlags(0)
	}
}
