package main

import (
	"fmt"
	"time"
	"strconv"
	"net/http"
	"github.com/gorilla/mux"
)

func activitiesHandler(w http.ResponseWriter, r *http.Request) {
	if Ctx(r).person_type == NoUser {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	Log(r, LOG_DEBUG, "Page /actividades")

	p := make(map[string]interface{})

	p["activities"] = db_list_activities()
	p["people"] = db_list_socios_activos()
	renderTemplate(w, r, "actividades", p)
}

func activityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	Log(r, LOG_DEBUG, fmt.Sprintf("Page /actividad/id=%d", id))

	p := make(map[string]interface{})

	p["activity"] = db_one_activity(id)
	p["people"] = db_list_socios_activos()
	renderTemplate(w, r, "actividad", p)
}

func ajaxActivityHandler(w http.ResponseWriter, r *http.Request) {
	Log(r, LOG_DEBUG, "Page /ajax/activity")

	r.ParseForm()
	action := r.FormValue("action")

	if action == "new-activity" {
		var date1, date2 time.Time
		var organizer int
		var title string

		date1, _ = time.Parse("2006-01-02", r.FormValue("date1"))
		date2, _ = time.Parse("2006-01-02", r.FormValue("date2"))
		fmt.Sscan(r.FormValue("organizer"), &organizer)
		title = r.FormValue("title")
		db_new_activity(date1, date2, organizer, title)
	}
}
