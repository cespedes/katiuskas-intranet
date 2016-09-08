package main

import (
	"fmt"
	"time"
	"strconv"
	"github.com/gorilla/mux"
)

func activitiesHandler(ctx *Context) {
	log(ctx, LOG_DEBUG, "Page /actividades")

	p := make(map[string]interface{})

	p["activities"] = db_list_activities()
	p["id"] = ctx.id
	p["people"] = db_list_socios_activos()
	renderTemplate(ctx, "actividades", p)
}

func activityHandler(ctx *Context) {
	vars := mux.Vars(ctx.r)
	id, _ := strconv.Atoi(vars["id"])
	log(ctx, LOG_DEBUG, fmt.Sprintf("Page /actividad/id=%d", id))

	p := make(map[string]interface{})

	p["id"] = ctx.id
	p["activity"] = db_one_activitiy(id)
	p["people"] = db_list_socios_activos()
	renderTemplate(ctx, "actividad", p)
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
