package main

func infoHandler(ctx *Context) {
	p := make(map[string]interface{})

	p["board"] = db_list_board()

	renderTemplate(ctx, "info", p)
}
