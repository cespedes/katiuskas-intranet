package main

import (
	"fmt"
	"time"
	"strings"
	"net/http"
)

func ajaxAdminHandler(ctx *Context) {
	if !ctx.admin {
		http.Redirect(ctx.w, ctx.r, "/", http.StatusFound)
		return
	}
	ctx.r.ParseForm()
	action := ctx.r.FormValue("action")

	if action == "new-email" {
		var email string
		var id int
		email = ctx.r.FormValue("email")
		fmt.Sscan(ctx.r.FormValue("id"), &id)
		db_person_add_email(id, email)
	} else if action == "update-person" {
		var id int
		var birth time.Time
		fmt.Sscan(ctx.r.FormValue("id"), &id)
		userinfo := db_get_userinfo(id)

		birth, _ = time.Parse("02-01-2006", ctx.r.FormValue("birth"))
		name := ctx.r.FormValue("name")
		surname := ctx.r.FormValue("surname")
		dni := ctx.r.FormValue("dni")
		address := ctx.r.FormValue("address")
		zip := ctx.r.FormValue("zip")
		city := ctx.r.FormValue("city")
		province := ctx.r.FormValue("province")
		gender := map[string]string{"M":"M","F":"F"}[ctx.r.FormValue("gender")]

		db.Exec("UPDATE person SET name=$2,surname=$3,dni=$4,birth=$5,address=$6,zip=$7,city=$8,province=$9,gender=$10 WHERE id=$1", id, name, surname, dni, birth, address, zip, city, province, gender)
		log_msg := fmt.Sprintf("Updated socio %d (%s %s)", id, userinfo["name"], userinfo["surname"])
		fn := func(name, value string) string {
			if userinfo[name] != value {
				return fmt.Sprintf("\n%s: %s -> %s", strings.Title(name), userinfo[name], value)
			} else {
				return ""
			}
		}
		log_msg += fn("name", name)
		log_msg += fn("surname", surname)
		log_msg += fn("dni", dni)
		log_msg += fn("city", city)
		log_msg += fn("province", province)
		gender2 := map[string]string{"M":"Masculino","F":"Femenino"}[gender]
		if userinfo["gender"] != gender2 {
			log_msg += fmt.Sprintf("\nGender: %s -> %s", userinfo["gender"], gender2)
		}
		log(ctx, LOG_NOTICE, log_msg)
	}
}
