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
