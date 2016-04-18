package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "session").(*sessions.Session)

	email, ok := session.Values["email"].(string)
	if (ok) {
		fmt.Fprintln(w, "Ya estás autenticado.  Tu email es", email)
		id, ok := session.Values["id"].(int)
		if (ok) {
			fmt.Fprintln(w, "Tu id es", id)
		}
	} else {
		fmt.Fprintln(w, "No te has autenticado.")
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
