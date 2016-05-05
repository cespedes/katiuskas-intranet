package main

import (
	"strconv"
	"net/http"
	"github.com/gorilla/mux"
)

func adminHandler(w http.ResponseWriter, r *http.Request) {
	id, email, person_type := get_id_email_type(w, r)

	if person_type != SocioAdmin {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	p := make(map[string]interface{})

	p["id"], p["email"], p["type"] = id, email, person_type

	p["new_emails"] = db_get_new_emails()
	p["people"] = db_list_people()
	for i,v := range p["people"].([]map[string]interface{}) {
		if v["type"].(int) <= ExSocio {
			p["people"].([]map[string]interface{})[i]["first_ex"] = true
			break
		}
	}

	renderTemplate(w, r, "admin", p)
}

func adminPersonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p := make(map[string]interface{})
	p["userinfo"] = db_get_userinfo(id)

	renderTemplate(w, r, "admin-person", p)
}
