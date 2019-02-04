package main

import (
	"net/http"
)

func infoHandler(w http.ResponseWriter, r *http.Request) {
	if Ctx(r).person_type == NoUser {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	p := make(map[string]interface{})

	p["board"] = db_list_board()

	renderTemplate(w, r, "info", p)
}
