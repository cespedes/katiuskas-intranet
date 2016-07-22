package main

import (
	"fmt"
	"time"
)

func activitiesHandler(ctx *Context) {
	log(ctx, LOG_DEBUG, "Page /actividades")

	p := make(map[string]interface{})

	p["activities"] = db_list_activities()
	p["id"] = ctx.id
	p["people"] = db_list_socios_activos()
	renderTemplate(ctx, "actividades", p)
}

func ajaxActivityHandler(ctx *Context) {
	log(ctx, LOG_DEBUG, "Page /ajax/activity")

	ctx.r.ParseForm()
	action := ctx.r.FormValue("action")

	if action == "new-activity" {
		var date time.Time
		var organizer int
		var title string

		date, _ = time.Parse("2006-01-02", ctx.r.FormValue("date"))
		fmt.Sscan(ctx.r.FormValue("organizer"), &organizer)
		title = ctx.r.FormValue("title")
		db_new_activity(date, organizer, title)
	}
}
