package main

import (
	"net/http"
)

func main() {
	templates_init()
	db_init()

	r := router()

	http.Handle("/", r)

	http.ListenAndServe("localhost:8081", nil)
}
