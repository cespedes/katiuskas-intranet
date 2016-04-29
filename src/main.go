package main

import (
	"net/http"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "session").(*sessions.Session)

	p := make(map[string]interface{})

	p["session"] = session.Values

	email, ok := session.Values["email"].(string)
	if (!ok) {
		renderTemplate(w, r, "root-nouser", p)
		return
	}
	p["email"] = email
	id := session.Values["id"].(int)
	person_type, ok := session.Values["type"].(int)
	if id==0 || !ok {
		id, person_type, ok = db_mail_2_id(email)
		session.Values["id"] = id
		session.Values["type"] = person_type
	}
	p["id"] = id
	p["type"] = person_type
	if (!ok) {
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
	p["userinfo"] = db_get_userinfo(100)
	p["session"] = session.Values
	renderTemplate(w, r, "root", p)
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := session_init(w, r)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	templates_init()
	db_init()

	r := router()

	http.Handle("/", middleware(r))

	http.ListenAndServe(":8081", nil)
}
