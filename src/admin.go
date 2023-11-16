package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (s *server) adminHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is for admins")
}

func (s *server) ajaxAdminHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	action := r.FormValue("action")

	if action == "update-person" {
		var id int
		var birth time.Time
		fmt.Sscan(r.FormValue("id"), &id)
		userinfo := s.DBgetUserinfo(id)

		birth, _ = time.Parse("02-01-2006", r.FormValue("birth"))
		name := r.FormValue("name")
		surname := r.FormValue("surname")
		dni := r.FormValue("dni")
		address := r.FormValue("address")
		zip := r.FormValue("zip")
		city := r.FormValue("city")
		province := r.FormValue("province")
		emerg_contact := r.FormValue("emerg_contact")
		gender := map[string]string{"M": "M", "F": "F"}[r.FormValue("gender")]
		phones := strings.Trim(r.FormValue("phones"), " ")
		emails := strings.Trim(r.FormValue("emails"), " ")
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

		s.db.Exec("UPDATE person SET name=$2,surname=$3,dni=$4,birth=$5,address=$6,zip=$7,city=$8,province=$9,emerg_contact=$10,gender=$11 WHERE id=$1", id, name, surname, dni, birth, address, zip, city, province, emerg_contact, gender)
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
			s.db.Exec("DELETE FROM person_phone WHERE id_person=$1", id)
			for i, phone := range strings.Split(phones, " ") {
				if phone == "" {
					continue
				}
				if i == 0 {
					s.db.Exec("INSERT INTO person_phone (id_person,phone,main) VALUES ($1,$2,true)", id, phone)
				} else {
					s.db.Exec("INSERT INTO person_phone (id_person,phone,main) VALUES ($1,$2,false)", id, phone)
				}
			}
			log_msg += fn("phones", phones)
		}
		if emails != userinfo["emails"] {
			s.db.Exec("DELETE FROM person_email WHERE id_person=$1", id)
			for i, email := range strings.Split(emails, " ") {
				if email == "" {
					continue
				}
				if i == 0 {
					s.db.Exec("INSERT INTO person_email (id_person,email,main) VALUES ($1,$2,true)", id, email)
				} else {
					s.db.Exec("INSERT INTO person_email (id_person,email,main) VALUES ($1,$2,false)", id, email)
				}
			}
			log_msg += fn("emails", emails)
		}
		gender2 := map[string]string{"M": "Masculino", "F": "Femenino"}[gender]
		if userinfo["gender"] != gender2 {
			log_msg += fmt.Sprintf("\nGender: %s -> %s", userinfo["gender"], gender2)
		}
		s.Log(r, LOG_NOTICE, log_msg)
	} else if action == "update-person-pic" {
		var id int
		var file string
		fmt.Sscan(r.FormValue("id"), &id)
		fmt.Sscan(r.FormValue("file"), &file)
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

		fmt.Sscan(r.FormValue("id"), &id_person)
		userinfo := s.DBgetUserinfo(id_person)

		if y, err := strconv.Atoi(r.FormValue("license-year")); err == nil {
			year = y
		} else {
			log_msg := fmt.Sprintf("Adding license for socio %d (%s %s): malformed year (%s)", id_person, userinfo["name"], userinfo["surname"], r.FormValue("license-year"))
			s.Log(r, LOG_NOTICE, log_msg)
			return
		}

		federation = r.FormValue("license-federation")
		issued, _ = time.Parse("2006-01-02", r.FormValue("license-issued"))
		tecnico = (r.FormValue("license-tecnico") != "")

		s.db.Exec("INSERT INTO person_federation (id_person, year, federation, issued, tecnico) VALUES ($1, $2, $3, $4, $5)", id_person, year, federation, issued, tecnico)
		log_msg := fmt.Sprintf("Added license for socio %d (%s %s)", id_person, userinfo["name"], userinfo["surname"])
		log_msg += fmt.Sprintf("\n%s %s (%d)", issued.Format("02-01-2006"), federation, year)
		if tecnico {
			log_msg += " (técnico)"
		}
		s.Log(r, LOG_NOTICE, log_msg)
	} else if action == "add-alta" {
		var id_person int

		fmt.Sscan(r.FormValue("id"), &id_person)
		userinfo := s.DBgetUserinfo(id_person)
		date, err := time.Parse("2006-01-02", r.FormValue("date"))
		if err != nil {
			log_msg := fmt.Sprintf("Adding alta for socio %d (%s %s): malformed date (%s)", id_person, userinfo["name"], userinfo["surname"], r.FormValue("date"))
			s.Log(r, LOG_NOTICE, log_msg)
			return
		}
		s.db.Exec("INSERT INTO socio (id_person, alta) VALUES ($1, $2)", id_person, date)
		log_msg := fmt.Sprintf("Added alta for socio %d (%s %s) with date %s", id_person, userinfo["name"], userinfo["surname"], date.Format("02-01-2006"))
		s.Log(r, LOG_NOTICE, log_msg)
	} else if action == "add-baja" {
		var id_person int

		fmt.Sscan(r.FormValue("id"), &id_person)
		userinfo := s.DBgetUserinfo(id_person)
		date, err := time.Parse("2006-01-02", r.FormValue("date"))
		if err != nil {
			log_msg := fmt.Sprintf("Adding baja definitiva for socio %d (%s %s): malformed date (%s)", id_person, userinfo["name"], userinfo["surname"], r.FormValue("date"))
			s.Log(r, LOG_NOTICE, log_msg)
			return
		}
		s.db.Exec("UPDATE socio SET baja=$2 WHERE baja IS NULL AND id_person=$1", id_person, date)
		log_msg := fmt.Sprintf("Added baja definitiva for socio %d (%s %s) with date %s", id_person, userinfo["name"], userinfo["surname"], date.Format("02-01-2006"))
		s.Log(r, LOG_NOTICE, log_msg)
	} else if action == "add-baja-temporal" {
		var id_person int

		fmt.Sscan(r.FormValue("id"), &id_person)
		userinfo := s.DBgetUserinfo(id_person)
		date, err := time.Parse("2006-01-02", r.FormValue("date"))
		if err != nil {
			log_msg := fmt.Sprintf("Adding baja temporal for socio %d (%s %s): malformed date (%s)", id_person, userinfo["name"], userinfo["surname"], r.FormValue("date"))
			s.Log(r, LOG_NOTICE, log_msg)
			return
		}
		s.db.Exec("INSERT INTO baja_temporal (id_person, start) VALUES ($1, $2)", id_person, date)
		log_msg := fmt.Sprintf("Added baja temporal for socio %d (%s %s) with start date %s", id_person, userinfo["name"], userinfo["surname"], date.Format("02-01-2006"))
		s.Log(r, LOG_NOTICE, log_msg)
	} else if action == "fin-baja-temporal" {
		var id_person int

		fmt.Sscan(r.FormValue("id"), &id_person)
		userinfo := s.DBgetUserinfo(id_person)
		date, err := time.Parse("2006-01-02", r.FormValue("date"))
		if err != nil {
			log_msg := fmt.Sprintf("Fin de baja temporal del socio %d (%s %s): malformed date (%s)", id_person, userinfo["name"], userinfo["surname"], r.FormValue("date"))
			s.Log(r, LOG_NOTICE, log_msg)
			return
		}
		s.db.Exec(`UPDATE baja_temporal SET "end"=$2 WHERE "end" IS NULL AND id_person=$1`, id_person, date)
		log_msg := fmt.Sprintf("Fin de baja temporal del socio %d (%s %s) con fecha %s", id_person, userinfo["name"], userinfo["surname"], date.Format("02-01-2006"))
		s.Log(r, LOG_NOTICE, log_msg)
	}
}
