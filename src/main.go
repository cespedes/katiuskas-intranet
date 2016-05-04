package main

import (
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	session := session_get(w, r)

	p := make(map[string]interface{})

	p["session"] = session

	email, ok := session["email"].(string)
	if (!ok) {
		renderTemplate(w, r, "root-nouser", p)
		return
	}
	p["email"] = email
	id, ok1 := session["id"].(int)
	person_type, ok2 := session["type"].(int)
	if id==0 || !ok1 || !ok2 {
		id, person_type := db_mail_2_id(email)
		session["id"] = id
		session["type"] = person_type
	}
	p["id"] = id
	p["type"] = person_type
	if person_type == NoSocio {
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
	p["session"] = session
	renderTemplate(w, r, "root", p)
}

func main() {
	templates_init()
	db_init()

	r := router()

	http.Handle("/", r)

	http.ListenAndServe(":8081", nil)
}
