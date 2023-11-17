package main

import (
	"fmt"
	"net/http"
)

func (s *server) apiGetUser(w http.ResponseWriter, r *http.Request) {
	ctx := Ctx(r)
	fmt.Fprintf(w, "This is the API.  Ctx.id=%d.  Request URL is %s\n", ctx.id, r.URL)
}
