package main

func infoHandler(ctx *Context) {
	p := make(map[string]interface{})

	p["session"] = ctx.session.Values
	p["id"] = ctx.id
	p["email"] = ctx.email
	p["type"] = ctx.person_type
	p["ipaddr"] = ctx.ipaddr
	p["userinfo"] = db_get_userinfo(ctx.id)

	renderTemplate(ctx, "info", p)
}
