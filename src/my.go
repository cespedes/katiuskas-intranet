package main

import (
	"net/http"
)

func (s *server) myHandler(w http.ResponseWriter, r *http.Request) {
	s.Log(r, LOG_DEBUG, "Page /info")

	p := make(map[string]interface{})

	p["session"] = Ctx(r).session.Values
	p["ipaddr"] = Ctx(r).ipaddr
	p["userinfo"] = s.DBgetUserinfo(Ctx(r).id)

	renderTemplate(w, r, "my", p)
}
