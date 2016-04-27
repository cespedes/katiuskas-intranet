package main

import (
	"fmt"
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

func db_get_info(id int) (result map[string]interface{}) {
	var name, surname, dni, birth string
	var err error
	var row *sql.Row

	result = make(map[string]interface{})
	row = db.QueryRow("SELECT name,surname,dni,birth FROM person WHERE id=$1", id)
	err = row.Scan(&name, &surname, &dni, &birth)

	if err == nil {
		result["name"] = name
		result["surname"] = surname
		result["dni"] = dni
		result["birth"] = birth
	}
	return
}
