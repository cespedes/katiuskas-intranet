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

func db_mail_2_id(email string) (id int, ok bool) {
	var err error
	err = db.QueryRow("SELECT id_person FROM person_email WHERE email=$1", email).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		break
	case err != nil:
		log_error(fmt.Sprintf("SQL Error: %s", err))
		/* deal with error */
		ok = false
		return
	default:
		ok = true
		return
	}
	db.Exec("INSERT INTO new_email (email) VALUES ($1)", email) /* ignore errors */
	ok = false
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
	var name, surname, dni, address, zip, city, province string
	var birth time.Time
	var err error
	var row *sql.Row
	var rows *sql.Rows

	result = make(map[string]interface{})

	// Personal data
	row = db.QueryRow("SELECT name,surname,dni,birth,address,zip,city,province FROM person WHERE id=$1", id)
	err = row.Scan(&name, &surname, &dni, &birth, &address, &zip, &city, &province)

	if err == nil {
		result["name"] = name
		result["surname"] = surname
		result["dni"] = dni
		result["birth"] = birth.Format("02-01-2006")
		result["address"] = address
		result["zip"] = zip
		result["city"] = city
		result["province"] = province
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
	}
	return
}
