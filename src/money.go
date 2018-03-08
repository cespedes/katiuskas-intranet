package main

import (
	"fmt"
	"time"
	"strconv"
	"net/http"
)
/*
import (
	"fmt"
	"time"
	"strconv"
	"github.com/gorilla/mux"
)
*/

func moneyHandler(ctx *Context) {
	log(ctx, LOG_DEBUG, "Page /money")

	if !(ctx.roles["admin"] || ctx.roles["money"]) {
		http.Redirect(ctx.w, ctx.r, "/", http.StatusFound)
		return
	}

	p := make(map[string]interface{})
	today := time.Now()

	p["today"] = today.Format("2006-01-02")
	p["last_30d"] = today.Add(-30 * 24 * time.Hour).Format("2006-01-02")
	p["last_365d"] = today.Add(-365 * 24 * time.Hour).Format("2006-01-02")
	p["last_year"] = today.Year() - 1
	p["year"] = today.Year()
	p["accounts"] = db_get_accounts()
	renderTemplate(ctx, "money", p)
}

func ajaxMoneyHandler(ctx *Context) {
	log(ctx, LOG_DEBUG, "Page /ajax/money")

	ctx.r.ParseForm()

	account, _ := strconv.Atoi(ctx.r.FormValue("account"))
	from := ctx.r.FormValue("from")
	lines := db_get_money(account, from)
	if len(lines)==0 {
		fmt.Fprint(ctx.w, "No lines to display\n")
		return
	}
	fmt.Fprint(ctx.w, `
<style>
table.money tr td:first-child {
  width: 12ex;
}
table.money tr td:nth-child(3) {
  width: 12ex;
  text-align: right;
}
table.money tr td:nth-child(4) {
  width: 12ex;
  text-align: right;
}
</style>`)

	fmt.Fprint(ctx.w, "<table class=\"money\">\n")
	for _, line := range lines {
		fmt.Fprint(ctx.w, "<tr>")
		fmt.Fprintf(ctx.w, "<td>%v</td>", line["date"])
		fmt.Fprintf(ctx.w, "<td>%v</td>", line["description"])
		fmt.Fprintf(ctx.w, "<td>%v</td>", line["value"])
		fmt.Fprintf(ctx.w, "<td>%v</td>", line["balance"])
		fmt.Fprint(ctx.w, "</tr>\n")
	}
	fmt.Fprint(ctx.w, "</table>\n")
}
