package main

import (
	"log"
	"net/http"
)

func MyHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/* nothing to do on all petitions, before or after the main handler */
		h.ServeHTTP(w, r)
	})
}

func main() {
	templates_init()
	db_init()

	r := router()

	http.Handle("/", MyHandler(r))

	err := http.ListenAndServe("localhost:8081", nil)
	log.Fatal(err)
}
