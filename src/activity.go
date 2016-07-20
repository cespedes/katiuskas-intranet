package main

import (
	"net/http"
)

func activitiesHandler(ctx *Context) {
	log(ctx, LOG_DEBUG, "Page /activities")

	http.Redirect(ctx.w, ctx.r, "/", http.StatusFound)
/*
	if ctx.person_type < SocioJunta {
		http.Redirect(ctx.w, ctx.r, "/", http.StatusFound)
		return
	}

	p := make(map[string]interface{})
	if ctx.admin {
		p["admin"] = true
	}
	renderTemplate(ctx, "query", p)
*/
}
