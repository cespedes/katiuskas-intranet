package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "session").(*sessions.Session)

	if session.Values["email"] == nil {
		fmt.Fprintln(w, "No te has autenticado.")
	} else {
		var email = session.Values["email"].(string)
		fmt.Fprintln(w, "Ya estás autenticado.  Tu email es", email)
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
	r := router()

	http.Handle("/", middleware(r))

	http.ListenAndServe(":8081", nil)
}
