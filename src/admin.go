package main

import (
	"os"
	"fmt"
	"time"
	"strings"
	"strconv"
	"net/http"
	"encoding/base64"
)

func ajaxAdminHandler(ctx *Context) {
	if !ctx.roles["admin"] {
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
		emerg_contact := ctx.r.FormValue("emerg_contact")
		gender := map[string]string{"M":"M","F":"F"}[ctx.r.FormValue("gender")]
		phones := strings.Trim(ctx.r.FormValue("phones"), " ")
		emails := strings.Trim(ctx.r.FormValue("emails"), " ")
		if val, ok := userinfo["phones"]; ok {
			userinfo["phones"] = strings.Join(val.([]string), " ")
		} else {
			userinfo["phones"] = ""
		}
		if val, ok := userinfo["emails"]; ok {
			userinfo["emails"] = strings.Join(val.([]string), " ")
		} else {
			userinfo["emails"] = ""
		}

		db.Exec("UPDATE person SET name=$2,surname=$3,dni=$4,birth=$5,address=$6,zip=$7,city=$8,province=$9,emerg_contact=$10,gender=$11 WHERE id=$1", id, name, surname, dni, birth, address, zip, city, province, emerg_contact, gender)
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
		log_msg += fn("emerg_contact", emerg_contact)
		if phones != userinfo["phones"] {
			db.Exec("DELETE FROM person_phone WHERE id_person=$1", id)
			for i, phone := range strings.Split(phones, " ") {
				if phone == "" {
					continue
				}
				if i==0 {
					db.Exec("INSERT INTO person_phone (id_person,phone,main) VALUES ($1,$2,true)", id, phone)
				} else {
					db.Exec("INSERT INTO person_phone (id_person,phone,main) VALUES ($1,$2,false)", id, phone)
				}
			}
			log_msg += fn("phones", phones)
		}
		if emails != userinfo["emails"] {
			db.Exec("DELETE FROM person_email WHERE id_person=$1", id)
			for i, email := range strings.Split(emails, " ") {
				if email == "" {
					continue
				}
				if i==0 {
					db.Exec("INSERT INTO person_email (id_person,email,main) VALUES ($1,$2,true)", id, email)
				} else {
					db.Exec("INSERT INTO person_email (id_person,email,main) VALUES ($1,$2,false)", id, email)
				}
			}
			log_msg += fn("emails", emails)
		}
		gender2 := map[string]string{"M":"Masculino","F":"Femenino"}[gender]
		if userinfo["gender"] != gender2 {
			log_msg += fmt.Sprintf("\nGender: %s -> %s", userinfo["gender"], gender2)
		}
		log(ctx, LOG_NOTICE, log_msg)
	} else if action == "update-person-pic" {
		var id int
		var file string
		fmt.Sscan(ctx.r.FormValue("id"), &id)
		fmt.Sscan(ctx.r.FormValue("file"), &file)
		f, err := os.OpenFile(fmt.Sprintf("files/people/%d.jpg", id), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return
		}
		defer f.Close()
		decoded, err := base64.StdEncoding.DecodeString(strings.Split(file, ",")[1])
		if err != nil {
			return
		}
		f.Write(decoded)
	} else if action == "add-license" {
		var id_person int
		var year int
		var federation string
		var issued time.Time
		var tecnico bool

		fmt.Sscan(ctx.r.FormValue("id"), &id_person)
		userinfo := db_get_userinfo(id_person)

		if y, err := strconv.Atoi(ctx.r.FormValue("license-year")); err == nil {
			year = y
		} else {
			log_msg := fmt.Sprintf("Adding license for socio %d (%s %s): malformed year (%s)", id_person, userinfo["name"], userinfo["surname"], ctx.r.FormValue("license-year"))
			log(ctx, LOG_NOTICE, log_msg)
			return
		}

		federation = ctx.r.FormValue("license-federation")
		issued, _ = time.Parse("2006-01-02", ctx.r.FormValue("license-issued"))
		tecnico = (ctx.r.FormValue("license-tecnico") != "")

		db.Exec("INSERT INTO person_federation (id_person, year, federation, issued, tecnico) VALUES ($1, $2, $3, $4, $5)", id_person, year, federation, issued, tecnico)
		log_msg := fmt.Sprintf("Added license for socio %d (%s %s)", id_person, userinfo["name"], userinfo["surname"])
		log_msg += fmt.Sprintf("\n%s %s (%d)", issued.Format("02-01-2006"), federation, year)
		if tecnico {
			log_msg += fmt.Sprintf(" (técnico)")
		}
		log(ctx, LOG_NOTICE, log_msg)
	} else if action == "add-alta" {
		var id_person int

		fmt.Sscan(ctx.r.FormValue("id"), &id_person)
		userinfo := db_get_userinfo(id_person)
		date, err := time.Parse("2006-01-02", ctx.r.FormValue("date"))
		if err != nil {
			log_msg := fmt.Sprintf("Adding alta for socio %d (%s %s): malformed date (%s)", id_person, userinfo["name"], userinfo["surname"], ctx.r.FormValue("date"))
			log(ctx, LOG_NOTICE, log_msg)
			return
		}
		db.Exec("INSERT INTO socio (id_person, alta) VALUES ($1, $2)", id_person, date)
		log_msg := fmt.Sprintf("Added alta for socio %d (%s %s) with date %s", id_person, userinfo["name"], userinfo["surname"], date.Format("02-01-2006"))
		log(ctx, LOG_NOTICE, log_msg)
	} else if action == "add-baja-temporal" {
		var id_person int

		fmt.Sscan(ctx.r.FormValue("id"), &id_person)
		userinfo := db_get_userinfo(id_person)
		date, err := time.Parse("2006-01-02", ctx.r.FormValue("date"))
		if err != nil {
			log_msg := fmt.Sprintf("Adding baja temporal for socio %d (%s %s): malformed date (%s)", id_person, userinfo["name"], userinfo["surname"], ctx.r.FormValue("date"))
			log(ctx, LOG_NOTICE, log_msg)
			return
		}
		db.Exec("INSERT INTO baja_temporal (id_person, start) VALUES ($1, $2)", id_person, date)
		log_msg := fmt.Sprintf("Added baja temporal for socio %d (%s %s) with start date %s", id_person, userinfo["name"], userinfo["surname"], date.Format("02-01-2006"))
		log(ctx, LOG_NOTICE, log_msg)
	} else if action == "fin-baja-temporal" {
		var id_person int

		fmt.Sscan(ctx.r.FormValue("id"), &id_person)
		userinfo := db_get_userinfo(id_person)
		date, err := time.Parse("2006-01-02", ctx.r.FormValue("date"))
		if err != nil {
			log_msg := fmt.Sprintf("Fin de baja temporal del socio %d (%s %s): malformed date (%s)", id_person, userinfo["name"], userinfo["surname"], ctx.r.FormValue("date"))
			log(ctx, LOG_NOTICE, log_msg)
			return
		}
		db.Exec(`UPDATE baja_temporal SET "end"=$2 WHERE "end" IS NULL AND id_person=$1`, id_person, date)
		log_msg := fmt.Sprintf("Fin de baja temporal del socio %d (%s %s) con fecha %s", id_person, userinfo["name"], userinfo["surname"], date.Format("02-01-2006"))
		log(ctx, LOG_NOTICE, log_msg)
	}
}
