package main

import (
	"fmt"
	"strings"
	"strconv"
	"net/http"
	"unicode/utf8"
	"encoding/csv"
	"github.com/gorilla/mux"
)

func sociosHandler(ctx *Context) {
	log(ctx, LOG_DEBUG, "Page /socios")

	p := make(map[string]interface{})
	if ctx.admin {
		p["admin"] = true
	}
	if ctx.board {
		p["board"] = true
	}
	renderTemplate(ctx, "socios", p)
}

func ajaxSociosHandler(ctx *Context) {
	log(ctx, LOG_DEBUG, "Page /ajax/socios")

	ctx.r.ParseForm()

	var order []string
	for i := 0; i<100; i++ {
		switch ctx.r.FormValue(fmt.Sprintf("order-%d", i)) {
			case "name":
				order = append(order, "name")
			case "surname":
				order = append(order, "surname")
			case "gender":
				order = append(order, "gender")
			case "birth":
				order = append(order, "birth")
			case "cumple":
				if (ctx.board || ctx.admin) {
					order = append(order, "date_part('month',birth),date_part('day',birth)")
				}
			case "federation":
				if (ctx.board || ctx.admin) {
					order = append(order, "federation")
				}
			case "alta":
				order = append(order, "alta")
			case "baja":
				order = append(order, "baja")
			default:
				break
		}
	}

	fields := []string{ "id" }
	for i := 0; i<100; i++ {
		switch ctx.r.FormValue(fmt.Sprintf("field-%d", i)) {
			case "row":
				fields = append(fields, fmt.Sprintf(`row_number() OVER (ORDER BY %s) AS "Línea"`, strings.Join(order, ",")))
			case "name":
				fields = append(fields, `name AS "Nombre"`)
			case "surname":
				fields = append(fields, `surname AS "Apellidos"`)
			case "gender":
				fields = append(fields, `CASE WHEN gender='M' THEN 'Masculino' WHEN gender='F' THEN 'Femenino' ELSE '' END AS "Género"`)
			case "dni":
				if (ctx.board || ctx.admin) {
					fields = append(fields, `dni AS "DNI"`)
				}
			case "birth":
				if (ctx.board || ctx.admin) {
					fields = append(fields, `COALESCE(birth::TEXT,'') AS "Nacimiento"`)
				}
			case "city":
				if (ctx.board || ctx.admin) {
					fields = append(fields, `city AS "Ciudad"`)
				}
			case "federation":
				if (ctx.board || ctx.admin) {
					fields = append(fields, `COALESCE(federation,'') AS "Federación"`)
				}
			case "type":
				fields = append(fields, `CASE WHEN type=2 THEN 'Ex-socio' WHEN type=3 THEN 'Baja temporal' WHEN type=4 THEN 'Socio activo' ELSE '???' END AS "Tipo"`)
			case "alta":
				fields = append(fields, `CASE WHEN alta IS NULL THEN '' ELSE alta::TEXT END AS "Alta"`)
			case "baja":
				fields = append(fields, `CASE WHEN baja IS NULL OR baja='infinity' THEN '' ELSE baja::TEXT END AS "Baja"`)
			default:
				break
		}
	}

	var filter []string
	var filter_gender []string
	if ctx.r.FormValue("filter-male") != "" {
		filter_gender = append(filter_gender, "gender='M'")
	}
	if ctx.r.FormValue("filter-female") != "" {
		filter_gender = append(filter_gender, "gender='F'")
	}
	if len(filter_gender) > 0 {
		filter = append(filter, "(" + strings.Join(filter_gender, " OR ") + ")")
	}

	var filter_type []string
	if ctx.r.FormValue("filter-ex-socio") != "" {
		filter_type = append(filter_type, "type <= 2")
	}
	if ctx.r.FormValue("filter-baja-temporal") != "" {
		filter_type = append(filter_type, "type = 3")
	}
	if ctx.r.FormValue("filter-socio-activo") != "" {
		filter_type = append(filter_type, "type >= 4")
	}
	if len(filter_type) > 0 {
		filter = append(filter, "(" + strings.Join(filter_type, " OR ") + ")")
	}

	var filter_category []string
	if ctx.r.FormValue("filter-infantiles") != "" {
		filter_category = append(filter_category, "date_part('year',age(birth))<14")
	}
	if ctx.r.FormValue("filter-juveniles") != "" {
		filter_category = append(filter_category, "date_part('year',age(birth)) between 14 and 17")
	}
	if ctx.r.FormValue("filter-mayores") != "" {
		filter_category = append(filter_category, "date_part('year',age(birth))>17")
	}
	if len(filter_category) > 0 {
		filter = append(filter, "(" + strings.Join(filter_category, " OR ") + ")")
	}

	if len(filter)==0 {
		filter = []string{"true"}
	}

	sql := fmt.Sprintf("SELECT %s FROM vperson WHERE %s ORDER BY %s",
		strings.Join(fields, ","), strings.Join(filter, " AND "), strings.Join(order, ","))
//	fmt.Fprintln(ctx.w, sql, "<br>")
//	fmt.Fprintf(ctx.w, "fields=%v, order=%v, filter=%v\ndata=%v\n", fields, order, filter, ctx.r.Form)

	if rows, err := db.Query(sql); err == nil {
		defer rows.Close()

		var columns []string
		var num_rows int
		var data [][]string

		columns, err = rows.Columns()
		if err != nil {
			fmt.Fprintf(ctx.w, "error 1\n")
			return
		}
		num_columns := len(columns)
		for rows.Next() {
			num_rows++
			data_row1 := make([]string, num_columns)
			data_row2 := make([]interface{}, num_columns)
			for i:=0; i<num_columns; i++ {
				data_row2[i] = &data_row1[i]
			}
			err = rows.Scan(data_row2...)
			if err != nil {
				fmt.Fprintf(ctx.w, "error 2: %v\n", err.Error())
				return
			}
			data = append(data, data_row1)
		}
		switch ctx.r.FormValue("result-type") {
			case "org":
				socios_display_org(ctx, columns, data)
			case "csv":
				socios_display_csv(ctx, columns, data)
			case "html":
				socios_display_html(ctx, columns, data)
			default:
				socios_display_html(ctx, columns, data)
		}
	}
}

func socios_display_html(ctx *Context, columns []string, data [][]string) {
	fmt.Fprintf(ctx.w, "<table>\n")
	fmt.Fprintf(ctx.w, "  <tr>\n")
	for _, x := range(columns[1:]) {
		fmt.Fprintf(ctx.w, "    <th>%s</th>\n", x)
	}
	fmt.Fprintf(ctx.w, "  </tr>\n")
	for _, x := range(data) {
		fmt.Fprintf(ctx.w, "  <tr>\n")
		for _, y := range(x[1:]) {
			if ctx.board || ctx.admin {
				fmt.Fprintf(ctx.w, "    <td><a href=\"/socio/id=%s\">%s</a></td>\n", x[0], y)
			} else {
				fmt.Fprintf(ctx.w, "    <td>%s</td>\n", y)
			}
		}
		fmt.Fprintf(ctx.w, "  </tr>\n")
	}
	fmt.Fprintf(ctx.w, "</table>\n")
}

func socios_display_org(ctx *Context, columns []string, data [][]string) {
	widths := make([]int, len(columns)-1)
	for i, x := range(columns[1:]) {
		widths[i] = utf8.RuneCountInString(x)
	}
	for _, x := range(data) {
		for i, y := range(x[1:]) {
			if utf8.RuneCountInString(y) > widths[i] {
				widths[i] = utf8.RuneCountInString(y)
			}
		}
	}
	fmt.Fprint(ctx.w, "<pre>\n")
	line := fmt.Sprint("|", strings.Repeat("-", widths[0]+2))
	for i, _ := range(columns[2:]) {
		line += "+" + strings.Repeat("-", widths[i+1]+2)
	}
	line += "|"
	fmt.Fprint(ctx.w, line, "\n|")
	for i, x := range(columns[1:]) {
		fmt.Fprintf(ctx.w, " %-*s |", widths[i], x)
	}
	fmt.Fprint(ctx.w, "\n", line, "\n")
	for _, x := range(data) {
		fmt.Fprintf(ctx.w, "|")
		for i, y := range(x[1:]) {
			fmt.Fprintf(ctx.w, " %-*s |", widths[i], y)
		}
		fmt.Fprintf(ctx.w, "\n")
	}
	fmt.Fprintln(ctx.w, line)
	fmt.Fprintf(ctx.w, "</pre>\n")
}

func socios_display_csv(ctx *Context, columns []string, data [][]string) {
	fmt.Fprintln(ctx.w, "<pre>")
	w := csv.NewWriter(ctx.w)
	for _, x := range(data) {
		w.Write(x[1:])
	}
	w.Flush()
	fmt.Fprintln(ctx.w, "</pre>")
}

func viewSocioHandler(ctx *Context) {
	if !(ctx.board || ctx.admin) {
		http.Redirect(ctx.w, ctx.r, "/", http.StatusFound)
		return
	}

	vars := mux.Vars(ctx.r)
	id, _ := strconv.Atoi(vars["id"])

	log(ctx, LOG_DEBUG, fmt.Sprintf("Page /socio/id=%d", id))

	p := make(map[string]interface{})
	p["userinfo"] = db_get_userinfo(id)
	p["altas_bajas"] = db_list_altas_bajas(id)
	if ctx.admin {
		p["admin"] = true
	}

	renderTemplate(ctx, "socio", p)
}

func socioNewHandler(ctx *Context) {
	var id int

	log(ctx, LOG_DEBUG, "Page /socio/new")
	err := db.QueryRow("INSERT INTO person DEFAULT VALUES RETURNING id").Scan(&id)
	if err != nil {
		http.Redirect(ctx.w, ctx.r, "/", http.StatusFound)
		return
	}
	http.Redirect(ctx.w, ctx.r, "/socio/id=" + strconv.Itoa(id), http.StatusFound)
}
