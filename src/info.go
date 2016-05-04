package main

import (
	"net/http"
)

func infoHandler(w http.ResponseWriter, r *http.Request) {
	session := session_get(w, r)
	person_type, ok := session["type"].(int)
	if !ok || person_type != SocioAdmin {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	id := session["id"].(int)

	p := make(map[string]interface{})

	p["session"] = session
	p["email"] = session["email"].(string)
	p["id"] = id
	p["type"] = person_type
	p["userinfo"] = db_get_userinfo(id)
	renderTemplate(w, r, "info", p)
}
