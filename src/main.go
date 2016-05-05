package main

import (
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	p := make(map[string]interface{})

	id, email, person_type := get_id_email_type(w, r)
	p["id"], p["email"], p["type"] = id, email, person_type

	if person_type == NoUser {
		renderTemplate(w, r, "root-nouser", p)
		return
	} else if person_type == NoSocio {
		form := r.FormValue("comment")
		if form != "" {
			db_set_new_email_comment(email, form)
			p["comment"] = form
			p["comment_set"] = true
		} else {
			p["comment"] = db_get_new_email_comment(email)
		}
		renderTemplate(w, r, "root-nosocio", p)
		return
	}
	p["userinfo"] = db_get_userinfo(id)
	renderTemplate(w, r, "root", p)
}

func main() {
	templates_init()
	db_init()

	r := router()

	http.Handle("/", r)

	http.ListenAndServe(":8081", nil)
}
