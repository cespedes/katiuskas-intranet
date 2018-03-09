package main

import (
	"os"
	"fmt"
	"time"
	"strings"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func db_init() {
	var err error
	db, err = sqlx.Open("postgres", "host=localhost user=katiuskas dbname=katiuskas password=Ohqu8Get")
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
	NoUser int = iota   /* 0 */
	NoSocio             /* 1 */
	ExSocio             /* 2 */
	SocioBajaTemporal   /* 3 */
	SocioActivo         /* 4 */
)

func db_get_roles(id int) (roles map[string]bool) {
	// Roles
	rows, err := db.Query("SELECT role FROM role WHERE person_id=$1", id)
	if err == nil {
		defer rows.Close()
		roles = make(map[string]bool)
		for rows.Next() {
			var role string
			err = rows.Scan(&role)
			if err == nil {
				roles[role] = true
			}
		}
	}
	return roles
}

func db_mail_2_id(email string) (id int, person_type int, board bool) {
	var err error
	email = strings.ToLower(email)
	err = db.QueryRow("SELECT id_person FROM person_email WHERE email=$1", email).Scan(&id)
	if err != nil {
		person_type = NoSocio
		db.Exec("INSERT INTO new_email (email) VALUES ($1)", email) /* ignore errors */
		return
	}
	db.QueryRow("SELECT type FROM vperson WHERE id=$1", id).Scan(&person_type)
	if db_rowExists(`SELECT 1 FROM board WHERE "end" IS NULL AND id_person=$1`, id) {
		board = true
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
	var row *sqlx.Row
	var rows *sqlx.Rows

	result = make(map[string]interface{})

	// Personal data
	row = db.QueryRowx("SELECT name,surname,dni,COALESCE(birth,'1000-01-01') AS birth,address,zip,city,province,CASE WHEN gender='M' THEN 'Masculino' WHEN gender='F' THEN 'Femenino' ELSE '' END AS gender,emerg_contact,type FROM vperson WHERE id=$1", id)
	if err = row.MapScan(result); err==nil {
		result["id"] = id
		result["birth"] = result["birth"].(time.Time).Format("02-01-2006")
	}

	// Phone(s)
	var phones []string
	db.Select(&phones, "SELECT phone FROM person_phone WHERE id_person=$1 ORDER BY NOT main,phone", id)
	result["phones"] = phones

	// E-mail(s)
	var emails []string
	db.Select(&emails, "SELECT email FROM person_email WHERE id_person=$1 ORDER BY NOT main,email", id)
	result["emails"] = emails

	// Board
	rows, err = db.Queryx(`SELECT position,start,COALESCE("end",'9999-12-31'::date) FROM board WHERE id_person=$1 ORDER BY start`, id)
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
	} else if result["gender"].(string)=="Femenino" {
		result["pic"] = "/files/people/female.jpg"
	} else {
		result["pic"] = "/files/people/male.jpg"
	}

	// Logs
	rows, err = db.Queryx(`
		SELECT date,text FROM (
		  SELECT alta AS date, 'Alta en el club' AS text, 1 AS sub FROM socio WHERE id_person=$1
		    UNION
		  SELECT issued, 'Licencia ' || federation || ' (' || year || ')' || CASE WHEN tecnico THEN ' (técnico)' ELSE '' END, 2 FROM person_federation WHERE id_person=$1
		    UNION
		  SELECT "end", 'Deja el cargo de ' || position, 3 FROM board WHERE id_person=$1 AND "end" IS NOT NULL
		    UNION
		  SELECT "end", 'Fin de baja temporal', 4 FROM baja_temporal WHERE id_person=$1 AND "end" IS NOT NULL
		    UNION
		  SELECT start, 'Inicio de baja temporal', 5 FROM baja_temporal WHERE id_person=$1
		    UNION
		  SELECT start, 'Nuevo cargo: ' || position, 6 FROM board WHERE id_person=$1
		    UNION
		  SELECT baja, 'Baja del club', 7 FROM socio WHERE id_person=$1 AND baja IS NOT NULL 
		) a ORDER BY date, sub`, id)

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

	result["roles"] = db_get_roles(id)
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

func db_list_federations() (result []string) {
	rows, err := db.Query("SELECT name FROM federation ORDER BY id,name")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			err = rows.Scan(&name)
			if err == nil {
				result = append(result, name)
			}
		}
	}
	return
}

func db_list_socios_activos() (result []map[string]interface{}) {
	rows, err := db.Query("SELECT id,name,surname,type FROM vperson WHERE type=$1 ORDER BY name,surname",SocioActivo)
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
				result = append(result, user)
			}
		}
	}
	return
}

func db_list_board() (result []map[string]interface{}) {
	rows, err := db.Query("SELECT id,name,surname,position,phone FROM vboard")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id int
			var name,surname,phone,position string
			err = rows.Scan(&id,&name,&surname,&position,&phone)
			if err == nil {
				user := make(map[string]interface{})
				user["id"] = id
				user["name"] = name
				user["surname"] = surname
				user["position"] = position
				user["phone"] = phone
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

func db_new_activity(date1 time.Time, date2 time.Time, organizer int, title string) {
	db.Exec("INSERT INTO activity (organizer,date_begin,date_end,title) VALUES ($1,$2,$3,$4)", organizer,date1,date2,title) /* ignore errors */
}

func db_list_activities() (result map[string][]map[string]interface{}) {
	act_type := map[int]string{0:"active", 1:"finished", 2:"cancelled"}
	result = make(map[string][]map[string]interface{})
	rows, err := db.Query(`
		SELECT
			a.id, a.date_begin, a.date_end, a.title, state,
			p.name || ' ' || p.surname AS organizer,
			COALESCE(pe.persons, 0) AS num_persons,
			COALESCE(eq.items, 0) AS num_items,
			COALESCE(pl.places, 0) AS num_places
		FROM activity a
		LEFT JOIN person p ON a.organizer=p.id
		LEFT JOIN (SELECT activity_id,count(person_id) as persons FROM activity_person GROUP BY activity_id) pe ON a.id=pe.activity_id
		LEFT JOIN (SELECT activity_id,count(id) as items FROM activity_item GROUP BY activity_id) eq ON a.id=eq.activity_id
		LEFT JOIN (SELECT activity_id,count(place_id) as places FROM activity_place GROUP BY activity_id) pl ON a.id=pl.activity_id
		ORDER BY date_begin;
        `)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			activity := db_fill_activity(rows)
			result[act_type[activity["state"].(int)]] = append(result[act_type[activity["state"].(int)]], activity)
		}
	}
	return
}

func db_fill_activity(rows *sql.Rows) (result map[string]interface{}) {
	var date_begin, date_end time.Time
	var title, organizer string
	var id, state, num_persons, num_items, num_places int
	err := rows.Scan(&id, &date_begin, &date_end, &title, &state, &organizer, &num_persons, &num_items, &num_places)
	if err == nil {
		result = make(map[string]interface{})
		result["id"] = id
		result["date_begin"] = date_begin.Format("02-01-2006")
		result["date_end"] = date_end.Format("02-01-2006")
		result["title"] = title
		result["state"] = state
		result["organizer"] = organizer
		result["num_persons"] = num_persons
		result["num_items"] = num_items
		result["num_places"] = num_places
	}
	return
}

func db_one_activity(id int) (result map[string]interface{}) {
	rows, err := db.Query(`
		SELECT
			a.id, a.date_begin, a.date_end, a.title, state,
			p.name || ' ' || p.surname AS organizer,
			COALESCE(pe.persons, 0) AS num_persons,
			COALESCE(eq.items, 0) AS num_items,
			COALESCE(pl.places, 0) AS num_places
		FROM activity a
		LEFT JOIN person p ON a.organizer=p.id
		LEFT JOIN (SELECT activity_id,count(person_id) as persons FROM activity_person GROUP BY activity_id) pe ON a.id=pe.activity_id
		LEFT JOIN (SELECT activity_id,count(id) as items FROM activity_item GROUP BY activity_id) eq ON a.id=eq.activity_id
		LEFT JOIN (SELECT activity_id,count(place_id) as places FROM activity_place GROUP BY activity_id) pl ON a.id=pl.activity_id
		WHERE a.id=$1;
        `, id)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		result = db_fill_activity(rows)
	}

/*
	// Persons
	rows, err = db.Query("SELECT pe.name FROM activity_person ap LEFT JOIN person pe ON ap.person_id=person.id WHERE ap.activity_id=$1", id)
	if err == nil {
		defer rows.Close()
		result["persons"] = []string(nil)
		for rows.Next() {
			var person string
			err = rows.Scan(&person)
			if err == nil {
				result["persons"] = append(result["persons"].([]string), person)
			}
		}
		if len(result["persons"].([]string)) == 0 {
			delete(result, "persons")
		}
	}
*/
	return
}

func db_list_items() (result []map[string]interface{}) {
	rows, err := db.Query(`
		SELECT
			id, type, subtype, makemodel, diameter, length, prestable, alquilable, cost
		FROM item;
        `)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			item := db_fill_item(rows)
			result = append(result, item)
		}
	}
	return
}

func db_fill_item(rows *sql.Rows) (result map[string]interface{}) {
	var Type, subtype, makemodel, diameter string
	var id, cost, length int
	var alquilable, prestable bool
	err := rows.Scan(&id, &Type, &subtype, &makemodel, &diameter, &length, &prestable, &alquilable, &cost)
	if err == nil {
		result = make(map[string]interface{})
		result["id"] = id
		result["type"] = Type
		result["subtype"] = subtype
		result["makemodel"] = makemodel
		result["diameter"] = diameter
		result["length"] = length
		result["prestable"] = prestable
		result["alquilable"] = alquilable
		result["cost"] = cost
	}
	return
}

func db_get_accounts() (result []map[string]interface{}) {
	rows, err := db.Queryx(`
		SELECT a.id,a.parent_id,a.name,a.code,to_char(sum(s.value),'FM999990.00') AS balance,to_char(max(s.datetime),'DD-MM-YYYY') AS date
                FROM account a
                  LEFT JOIN split s ON a.id=s.account_id
                GROUP BY a.id
                ORDER BY a.id;
        `)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			accounts := make(map[string]interface{})
			err = rows.MapScan(accounts)
			if err == nil {
				result = append(result, accounts)
			}
		}
	}
	return
}

func db_get_money(account int, from string) (result []map[string]interface{}) {
	rows, err := db.Queryx(`
		SELECT
			account_id AS id, to_char(datetime,'DD-MM-YYYY') AS date, description, to_char(value,'FM999990.00') AS value, to_char(balance,'FM999990.00') AS balance
		FROM money
		WHERE account_id=$1 AND datetime >= $2
        `, account, from)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			money := make(map[string]interface{})
			err = rows.MapScan(money)
			if err == nil {
				result = append(result, money)
			}
		}
	}
	return
}

func db_money_add(t Transaction) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	var id int
	err = tx.QueryRow("INSERT INTO transaction (description) VALUES ($1) RETURNING id", t.Description).Scan(&id)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, e := range t.Entries {
		_, err = tx.Exec("INSERT INTO split (datetime, transaction_id, account_id, value) VALUES ($1, $2, $3, $4::numeric/100)",
			e.Date, id, e.Account, e.Value)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
