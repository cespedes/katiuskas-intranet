package main

import (
	"net/http"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	p := make(map[string]interface{})

	p["session"] = Ctx(r).session.Values
	p["ipaddr"] = Ctx(r).ipaddr
	p["userinfo"] = db_get_userinfo(Ctx(r).id)

	renderTemplate(w, r, "my", p)
}
