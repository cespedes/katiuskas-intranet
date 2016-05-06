package main

import (
        "fmt"
        "net/http"
)

func ajaxAdminHandler(w http.ResponseWriter, r *http.Request) {
	_, _, person_type := get_id_email_type(w, r)

	if person_type != SocioAdmin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	action := r.URL.Query().Get("action")

	if action == "new-email" {
		var email string
		var id int
		email = r.URL.Query().Get("email")
		fmt.Sscan(r.URL.Query().Get("id"), &id)
		db_person_add_email(id, email)
	}
}
