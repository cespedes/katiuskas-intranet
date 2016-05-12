package main

import (
	"fmt"
	"net/http"
)

func ajaxAdminHandler(ctx *Context) {
	if ctx.person_type != SocioAdmin {
		http.Redirect(ctx.w, ctx.r, "/", http.StatusFound)
		return
	}
	action := ctx.r.URL.Query().Get("action")

	if action == "new-email" {
		var email string
		var id int
		email = ctx.r.URL.Query().Get("email")
		fmt.Sscan(ctx.r.URL.Query().Get("id"), &id)
		db_person_add_email(id, email)
	}
}
