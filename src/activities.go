package main

func activitiesHandler(ctx *Context) {
	log(ctx, LOG_DEBUG, "Page /actividades")

	p := make(map[string]interface{})

	p["activities"] = db_list_activities()
	renderTemplate(ctx, "actividades", p)
}
