package main

import (
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	Log(r, LOG_DEBUG, "Page /")

	p := make(map[string]interface{})

	if !Ctx(r).roles["user"] {
		renderTemplate(w, r, "root-nouser", p)
		return
	}
	p["userinfo"] = db_get_userinfo(Ctx(r).id)

	renderTemplate(w, r, "root", p)
}
