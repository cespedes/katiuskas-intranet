package main

import (
	"fmt"
	"time"
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func db_init() {
	var err error
	db, err = sql.Open("postgres", "host=localhost user=katiuskas dbname=katiuskas password=Ohqu8Get")
	if err != nil {
		fmt.Println("Error")
	}
}

func rowExists(query string, args ...interface{}) bool {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil {
		/* fatal error */
	}
	return exists
}

const (
	NoUser int = iota
	NoSocio
	ExSocio
	SocioBajaTemporal
	SocioActivo
	SocioJunta
	SocioAdmin
)

func db_mail_2_id(email string) (id int, person_type int) {
	var err error
	err = db.QueryRow("SELECT id_person FROM person_email WHERE email=$1", email).Scan(&id)
	if err != nil {
		person_type = NoSocio
		db.Exec("INSERT INTO new_email (email) VALUES ($1)", email) /* ignore errors */
		return
	}
	if rowExists("SELECT 1 FROM admin WHERE id_person=$1", id) {
		person_type = SocioAdmin
	} else if rowExists(`SELECT 1 FROM board WHERE "end" IS NOT NULL AND id_person=$1`, id) {
		person_type = SocioJunta
	} else if rowExists(`SELECT 1 FROM baja_temporal WHERE "end" IS NOT NULL AND id_person=$1`, id) {
		person_type = SocioBajaTemporal
	} else if rowExists(`SELECT 1 FROM socio WHERE "baja" IS NOT NULL AND id_person=$1`, id) {
		person_type = SocioActivo
	} else {
		person_type = ExSocio
	}
	return
}

func db_get_new_email_comment(email string) (comment string) {
	db.QueryRow("SELECT comment FROM new_email WHERE email=$1", email).Scan(&comment)
	return
}

func db_set_new_email_comment(email string, comment string) {
	db.Exec("UPDATE new_email SET comment=$1 WHERE email=$2", comment, email)
}

func db_get_userinfo(id int) (result map[string]interface{}) {
	var err error
	var row *sql.Row
	var rows *sql.Rows

	result = make(map[string]interface{})

	// Personal data
	{
		var name, surname, dni, address, zip, city, province string
		var birth time.Time
		row = db.QueryRow("SELECT name,surname,dni,birth,address,zip,city,province FROM person WHERE id=$1", id)
		err = row.Scan(&name, &surname, &dni, &birth, &address, &zip, &city, &province)

		if err == nil {
			result["id"] = id
			result["name"] = name
			result["surname"] = surname
			result["dni"] = dni
			result["birth"] = birth.Format("02-01-2006")
			result["address"] = address
			result["zip"] = zip
			result["city"] = city
			result["province"] = province
		}
	}

	// Phone(s)
	rows, err = db.Query("SELECT phone FROM person_phone WHERE id_person=$1 ORDER BY NOT main,phone", id)
	if err == nil {
		defer rows.Close()
		result["phones"] = []string(nil)
		for rows.Next() {
			var phone string
			err = rows.Scan(&phone)
			if err == nil {
				result["phones"] = append(result["phones"].([]string), phone)
			}
		}
		if len(result["phones"].([]string)) == 0 {
			delete(result, "phones")
		}
	}

	// E-mail(s)
	rows, err = db.Query("SELECT email FROM person_email WHERE id_person=$1 ORDER BY NOT main,email", id)
	if err == nil {
		defer rows.Close()
		result["emails"] = []string(nil)
		for rows.Next() {
			var email string
			err = rows.Scan(&email)
			if err == nil {
				result["emails"] = append(result["emails"].([]string), email)
			}
		}
		if len(result["emails"].([]string)) == 0 {
			delete(result, "emails")
		}
	}

	rows, err = db.Query(`SELECT position,start,COALESCE("end",'9999-12-31'::date) FROM board WHERE id_person=$1 ORDER BY start`, id)
	if err == nil {
		defer rows.Close()
		result["board"] = []interface{}(nil)
		for rows.Next() {
			var position string
			var start,end time.Time
			err = rows.Scan(&position, &start, &end)
			if err == nil {
				end_t := end.Format("02-01-2006")
				if end_t == "31-12-9999" {
					end_t = "actualidad"
				}
				result["board"] = append(result["board"].([]interface{}), struct {
					Position, Start, End string
				}{
					position,
					start.Format("02-01-2006"),
					end_t,
				})
			}
		}
		if len(result["board"].([]interface{})) == 0 {
			delete(result, "board")
		}
	}

	return
}
