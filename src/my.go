package main

func myHandler(ctx *Context) {
	p := make(map[string]interface{})

	p["session"] = ctx.session.Values
	p["ipaddr"] = ctx.ipaddr
	p["userinfo"] = db_get_userinfo(ctx.id)

	renderTemplate(ctx, "my", p)
}
