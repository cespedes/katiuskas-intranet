package main

import (
	"log"
	"net/http"
)

func main() {
	templates_init()
	db_init()

	r := router()

	http.Handle("/", r)

	err := http.ListenAndServe("localhost:8081", nil)
	log.Fatal(err)
}
