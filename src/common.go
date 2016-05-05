package main

import (
	"net/http"
)

func get_id_email_type(w http.ResponseWriter, r *http.Request) (id int, email string, person_type int) {
	var ok, ok1, ok2 bool

	session := session_get(w, r)

	email, ok = session["email"].(string)
	if !ok {
		id = 0
		person_type = NoUser
		return
	}
	id, ok1 = session["id"].(int)
	person_type, ok2 = session["type"].(int)
	if id==0 || !ok1 || !ok2 {
		id, person_type = db_mail_2_id(email)
		session["id"] = id
		session["type"] = person_type
	}
	return
}

