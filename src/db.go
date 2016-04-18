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

func db_new_email(email string) (id int, ok bool) {
	var err error
	err = db.QueryRow("SELECT id_person FROM person_email WHERE email=$1", email).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		break
	case err != nil:
		/* deal with error */
		ok = false
		return
	default:
		ok = true
		return
	}
	var comment string
	err = db.QueryRow("SELECT comment FROM new_email WHERE email=$1", email).Scan(&comment)
	switch {
	case err == sql.ErrNoRows:
		break
	case err != nil:
		/* deal with error */
		ok = false
		return
	default:
		id = 0
		ok = true
		return
	}
	_, err = db.Exec("INSERT INTO new_email (email) VALUES ($1)", email)
	if err != nil {
		/* deal with error */
		ok = false
		return
	}
	id = 0
	ok = true
	return
}
