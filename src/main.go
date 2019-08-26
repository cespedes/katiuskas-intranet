package main

import (
	"log"
	"net/http"
)

func MyHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create new context from HTTP session:
		r = NewContext(r)
		h.ServeHTTP(w, r)
	})
}

func main() {
	r := router()

	http.Handle("/", MyHandler(r))

	err := http.ListenAndServe(config("http_listen_addr"), nil)
	log.Fatal(err)
}
