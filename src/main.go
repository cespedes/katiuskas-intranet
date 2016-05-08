package main

import (
	"net/http"
)

func rootHandler(ctx *Context) {
	log(ctx, LOG_DEBUG, "rootHandler()")

	p := make(map[string]interface{})

	p["id"], p["email"], p["type"] = ctx.id, ctx.email, ctx.person_type

	if ctx.person_type == NoUser {
		renderTemplate(ctx, "root-nouser", p)
		return
	} else if ctx.person_type == NoSocio {
		if form := ctx.r.FormValue("comment"); form != "" {
			db_set_new_email_comment(ctx.email, form)
			p["comment"] = form
			p["comment_set"] = true
		} else {
			p["comment"] = db_get_new_email_comment(ctx.email)
		}
		renderTemplate(ctx, "root-nosocio", p)
		return
	}
	p["userinfo"] = db_get_userinfo(ctx.id)
	renderTemplate(ctx, "root", p)
}

func main() {
	templates_init()
	db_init()

	r := router()

	http.Handle("/", r)

	http.ListenAndServe(":8081", nil)
}
