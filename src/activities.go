package main

import (
	"fmt"
	"time"
	"strconv"
	"net/http"
	"github.com/gorilla/mux"
)

func activitiesHandler(ctx *Context) {
	if ctx.person_type == NoUser {
		http.Redirect(ctx.w, ctx.r, "/", http.StatusFound)
		return
	}
	log(ctx, LOG_DEBUG, "Page /actividades")

	p := make(map[string]interface{})

	p["activities"] = db_list_activities()
	p["people"] = db_list_socios_activos()
	renderTemplate(ctx, "actividades", p)
}

func activityHandler(ctx *Context) {
	vars := mux.Vars(ctx.r)
	id, _ := strconv.Atoi(vars["id"])
	log(ctx, LOG_DEBUG, fmt.Sprintf("Page /actividad/id=%d", id))

	p := make(map[string]interface{})

	p["activity"] = db_one_activity(id)
	p["people"] = db_list_socios_activos()
	renderTemplate(ctx, "actividad", p)
}

func ajaxActivityHandler(ctx *Context) {
	log(ctx, LOG_DEBUG, "Page /ajax/activity")

	ctx.r.ParseForm()
	action := ctx.r.FormValue("action")

	if action == "new-activity" {
		var date1, date2 time.Time
		var organizer int
		var title string

		date1, _ = time.Parse("2006-01-02", ctx.r.FormValue("date1"))
		date2, _ = time.Parse("2006-01-02", ctx.r.FormValue("date2"))
		fmt.Sscan(ctx.r.FormValue("organizer"), &organizer)
		title = ctx.r.FormValue("title")
		db_new_activity(date1, date2, organizer, title)
	}
}
