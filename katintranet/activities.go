package katintranet

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (s *server) activitiesHandler(w http.ResponseWriter, r *http.Request) {
	s.Log(r, LOG_DEBUG, "Page /actividades")

	p := make(map[string]interface{})

	p["activities"] = s.DBlistActivities()
	p["people"] = s.DBlistSociosActivos()
	renderTemplate(w, r, "actividades", p)
}

func (s *server) activityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	s.Log(r, LOG_DEBUG, fmt.Sprintf("Page /actividad/id=%d", id))

	p := make(map[string]interface{})

	p["activity"] = s.DBoneActivity(id)
	p["people"] = s.DBlistSociosActivos()
	renderTemplate(w, r, "actividad", p)
}

func (s *server) ajaxActivityHandler(w http.ResponseWriter, r *http.Request) {
	s.Log(r, LOG_DEBUG, "Page /ajax/activity")

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
		s.DBnewActivity(date1, date2, organizer, title)
	}
}
