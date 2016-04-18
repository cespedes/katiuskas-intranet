package main

import (
	"net/http"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "session").(*sessions.Session)

	p := make(map[string]interface{})
	email, ok := session.Values["email"].(string)
	if (ok) {
		p["email"] = email
		id, ok := session.Values["id"].(int)
		if (ok) {
			p["id"] = id
			renderTemplate(w, r, "root", p)
		} else {
			renderTemplate(w, r, "root-nosocio", p)
		}
	} else {
		renderTemplate(w, r, "root-nouser", p)
	}
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
