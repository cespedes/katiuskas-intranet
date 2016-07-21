package main

func docHandler(ctx *Context) {
	p := make(map[string]interface{})

	p["board"] = db_list_board()

	renderTemplate(ctx, "doc", p)
}
