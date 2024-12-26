package katintranet

import (
	"net/http"
)

func (s *server) itemsHandler(w http.ResponseWriter, r *http.Request) {
	s.Log(r, LOG_DEBUG, "Page /items")

	p := make(map[string]interface{})

	p["items"] = s.DBlistItems()
	renderTemplate(w, r, "items", p)
}
