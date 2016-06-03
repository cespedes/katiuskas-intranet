package main

import (
	"os"
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
)

func db_mail_2_id(email string) (id int, person_type int, admin bool) {
	var err error
	err = db.QueryRow("SELECT id_person FROM person_email WHERE email=$1", email).Scan(&id)
	if err != nil {
		person_type = NoSocio
		db.Exec("INSERT INTO new_email (email) VALUES ($1)", email) /* ignore errors */
		return
	}
	db.QueryRow("SELECT type FROM vperson WHERE id=$1", id).Scan(&person_type)
	if db_rowExists("SELECT 1 FROM admin WHERE id_person=$1", id) {
		admin = true
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

	var gender string
	// Personal data
	{
		var name, surname, dni, address, zip, city, province string
		var birth time.Time
		var person_type int
		row = db.QueryRow("SELECT name,surname,dni,COALESCE(birth,'1000-01-01') AS birth,address,zip,city,province,CASE WHEN gender='M' THEN 'Masculino' WHEN gender='F' THEN 'Femenino' ELSE '' END AS gender,type FROM vperson WHERE id=$1", id)
		err = row.Scan(&name, &surname, &dni, &birth, &address, &zip, &city, &province, &gender, &person_type)

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
			result["gender"] = gender
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

	if _, err := os.Stat(fmt.Sprintf("files/people/%d.jpg", id)); err == nil {
		result["pic"] = fmt.Sprintf("/files/people/%d.jpg", id)
	} else if gender=="Femenino" {
		result["pic"] = "/files/people/female.jpg"
	} else {
		result["pic"] = "/files/people/male.jpg"
	}
	rows, err = db.Query(`(SELECT issued AS date,'Licencia ' || federation || ' (' || year || ')' AS text FROM person_federation WHERE id_person=$1 UNION SELECT alta,'Alta en el club' FROM socio WHERE id_person=$1 UNION SELECT baja,'Baja del club' FROM socio WHERE id_person=$1 AND baja IS NOT NULL UNION SELECT start, 'Nuevo cargo: ' || position FROM board WHERE id_person=$1 UNION SELECT "end", 'Deja el cargo de ' || position FROM board WHERE id_person=$1 AND "end" IS NOT NULL UNION SELECT start, 'Inicio de baja temporal' FROM baja_temporal WHERE id_person=$1 UNION SELECT "end", 'Fin de baja temporal' FROM baja_temporal WHERE id_person=$1 AND "end" IS NOT NULL) ORDER BY date`, id)
	if err == nil {
		defer rows.Close()
		result["logs"] = []map[string]interface{}{nil}
		for rows.Next() {
			var date time.Time
			var text string
			err = rows.Scan(&date, &text)
			if err == nil {
				log := make(map[string]interface{})
				log["date"] = date.Format("02-01-2006")
				log["text"] = text
				result["logs"] = append(result["logs"].([]map[string]interface{}), log)
			}
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

func db_list_board() (result []map[string]interface{}) {
	rows, err := db.Query("SELECT id,name,surname,position FROM vboard")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			var name,surname,position string
			err = rows.Scan(&id,&name,&surname,&position)
			if err == nil {
				user := make(map[string]interface{})
				user["id"] = id
				user["name"] = name
				user["surname"] = surname
				user["position"] = position
				result = append(result, user)
			}
		}
	}
	return
}

func db_list_altas_bajas(id int) (result []map[string]interface{}) {
	rows, err := db.Query("SELECT alta,COALESCE(baja,'9999-12-31') FROM socio WHERE id_person=$1", id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var alta,baja time.Time
			err = rows.Scan(&alta, &baja)
			if err == nil {
				alta_baja := make(map[string]interface{})
				alta_baja["alta"] = alta.Format("02-01-2006")
				if baja.Format("02-01-2006") != "31-12-9999" {
					alta_baja["baja"] = baja.Format("02-01-2006")
				}
				result = append(result, alta_baja)
			}
		}
	}
	return
}

func db_person_add_email(id int, email string) {
	db.Exec("INSERT INTO person_email (id_person,email) VALUES ($1,$2)", id,email) /* ignore errors */
	db.Exec("DELETE FROM new_email WHERE email=$1", email) /* ignore errors */
}
