package katintranet

import (
	"net/http"
	"runtime/debug"
)

func (s *server) myHandler(w http.ResponseWriter, r *http.Request) {
	s.Log(r, LOG_DEBUG, "Page /info")

	p := make(map[string]interface{})

	p["session"] = C(r).session.Values
	p["ipaddr"] = C(r).ipaddr
	p["userinfo"] = s.DBgetUserinfo(C(r).id)
	p["buildinfo"], _ = debug.ReadBuildInfo()

	renderTemplate(w, r, "my", p)
}
