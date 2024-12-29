package katintranet

import (
	"net/http"
)

func (s *server) rootHandler(w http.ResponseWriter, r *http.Request) {
	s.Log(r, LOG_DEBUG, "Page /")

	p := make(map[string]interface{})

	if !HasRole(r, "user") {
		renderTemplate(w, r, "root-nouser", p)
		return
	}
	p["userinfo"] = s.DBgetUserinfo(C(r).id)

	renderTemplate(w, r, "root", p)
}
