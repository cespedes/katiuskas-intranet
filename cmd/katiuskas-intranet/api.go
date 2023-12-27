package main

import (
	"fmt"
	"net/http"
)

func (s *server) apiHandler() http.Handler {
	api := http.NewServeMux()

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Called API\n")
		fmt.Fprintf(w, "Path: %v\n", r.URL)
		fmt.Fprintf(w, "Headers: %v\n", r.Header)
	})

	return api
}

func (s *server) apiGetUser(w http.ResponseWriter, r *http.Request) {
	ctx := Ctx(r)
	fmt.Fprintf(w, "This is the API.  Ctx.id=%d.  Request URL is %s\n", ctx.id, r.URL)
}
