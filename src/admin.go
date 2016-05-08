package main

import (
	"strconv"
	"net/http"
	"github.com/gorilla/mux"
)

func adminHandler(ctx *Context) {
	log(ctx, LOG_DEBUG, "adminHandler()")

	if ctx.person_type != SocioAdmin {
		http.Redirect(ctx.w, ctx.r, "/", http.StatusFound)
		return
	}

	p := make(map[string]interface{})

	p["id"], p["email"], p["type"] = ctx.id, ctx.email, ctx.person_type

	p["new_emails"] = db_get_new_emails()
	p["people"] = db_list_people()
	for i,v := range p["people"].([]map[string]interface{}) {
		if v["type"].(int) <= ExSocio {
			p["people"].([]map[string]interface{})[i]["first_ex"] = true
			break
		}
	}

	renderTemplate(ctx, "admin", p)
}

func adminPersonHandler(ctx *Context) {
	vars := mux.Vars(ctx.r)
	id, _ := strconv.Atoi(vars["id"])

	p := make(map[string]interface{})
	p["userinfo"] = db_get_userinfo(id)

	renderTemplate(ctx, "admin-person", p)
}
