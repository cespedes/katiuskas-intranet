package main

import (
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	Log(r, LOG_DEBUG, "Page /")

	p := make(map[string]interface{})

	if Ctx(r).person_type == NoUser {
		renderTemplate(w, r, "root-nouser", p)
		return
	} else if Ctx(r).person_type == NoSocio {
		if form := r.FormValue("comment"); form != "" {
			db_set_new_email_comment(Ctx(r).email, form)
			p["comment"] = form
			p["comment_set"] = true
		} else {
			p["comment"] = db_get_new_email_comment(Ctx(r).email)
		}
		renderTemplate(w, r, "root-nosocio", p)
		return
	}
	p["userinfo"] = db_get_userinfo(Ctx(r).id)

	if Ctx(r).roles["admin"] {
		p["admin_new_emails"] = db_get_new_emails()
		p["people"] = db_list_people()
		for i,v := range p["people"].([]map[string]interface{}) {
			if v["type"].(int) <= ExSocio {
				p["people"].([]map[string]interface{})[i]["first_ex"] = true
				break
			}
		}
	}
	renderTemplate(w, r, "root", p)
}
