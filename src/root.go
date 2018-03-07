package main

func rootHandler(ctx *Context) {
	log(ctx, LOG_DEBUG, "Page /")

	p := make(map[string]interface{})

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

	if ctx.roles["admin"] {
		p["admin_new_emails"] = db_get_new_emails()
		p["people"] = db_list_people()
		for i,v := range p["people"].([]map[string]interface{}) {
			if v["type"].(int) <= ExSocio {
				p["people"].([]map[string]interface{})[i]["first_ex"] = true
				break
			}
		}
	}
	renderTemplate(ctx, "root", p)
}
