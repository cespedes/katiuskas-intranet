package main

import (
	"net/http"
)

func infoHandler(w http.ResponseWriter, r *http.Request) {
	p := make(map[string]interface{})

	session := session_get(w, r)
	p["session"] = session

	id, email, person_type := get_id_email_type(w, r)
	p["id"] = id
	p["email"] = email
	p["type"] = person_type

	p["userinfo"] = db_get_userinfo(id)

	renderTemplate(w, r, "info", p)
}
