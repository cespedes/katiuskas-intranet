package main

import (
	"net/http"
)

func infoHandler(ctx *Context) {
	if ctx.person_type == NoUser {
		http.Redirect(ctx.w, ctx.r, "/", http.StatusSeeOther)
		return
	}
	p := make(map[string]interface{})

	p["board"] = db_list_board()

	renderTemplate(ctx, "info", p)
}
