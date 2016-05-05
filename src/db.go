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

func db_rowExists(query string, args ...interface{}) bool {
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

func db_id_2_type(id int) (person_type int) {
	if db_rowExists("SELECT 1 FROM admin WHERE id_person=$1", id) {
		person_type = SocioAdmin
	} else if db_rowExists(`SELECT 1 FROM board WHERE "end" IS NULL AND id_person=$1`, id) {
		person_type = SocioJunta
	} else if db_rowExists(`SELECT 1 FROM baja_temporal b LEFT JOIN socio s ON b.id_socio=s.id WHERE b."end" IS NULL AND s.id_person=$1`, id) {
		person_type = SocioBajaTemporal
	} else if db_rowExists(`SELECT 1 FROM socio WHERE "baja" IS NULL AND id_person=$1`, id) {
		person_type = SocioActivo
	} else {
		person_type = ExSocio
	}
	return
}

func db_mail_2_id(email string) (id int, person_type int) {
	var err error
	err = db.QueryRow("SELECT id_person FROM person_email WHERE email=$1", email).Scan(&id)
	if err != nil {
		person_type = NoSocio
		db.Exec("INSERT INTO new_email (email) VALUES ($1)", email) /* ignore errors */
		return
	}
	err = db.QueryRow("SELECT type FROM person WHERE id=$1", id).Scan(&person_type)
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
		var person_type int
		row = db.QueryRow("SELECT name,surname,dni,COALESCE(birth,'1000-01-01') AS birth,address,zip,city,province,type FROM vperson WHERE id=$1", id)
		err = row.Scan(&name, &surname, &dni, &birth, &address, &zip, &city, &province, &person_type)

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
			result["type"] = person_type
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

	// Board
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

func db_get_new_emails() (result []map[string]interface{}) {
	rows, err := db.Query("SELECT email,comment,date FROM new_email ORDER BY date")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var email,comment string
			var date time.Time
			err = rows.Scan(&email,&comment,&date)
			if err == nil {
				user := make(map[string]interface{})
				user["email"] = email
				user["comment"] = comment
				user["date"] = date.Format("02-01-2006")
				result = append(result, user)
			}
		}
	}
	return
}

func db_list_people() (result []map[string]interface{}) {
	rows, err := db.Query("SELECT id,name,surname,type FROM vperson ORDER BY type<$1,name,surname",SocioBajaTemporal)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			var name,surname string
			var person_type int
			err = rows.Scan(&id,&name,&surname,&person_type)
			if err == nil {
				user := make(map[string]interface{})
				user["id"] = id
				user["name"] = name
				user["surname"] = surname
				user["type"] = person_type
				result = append(result, user)
			}
		}
	}
	return
}
