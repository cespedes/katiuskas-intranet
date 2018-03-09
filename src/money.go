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

	action := ctx.r.FormValue("action")
	if action == "show-money" {
		ajaxMoneyShow(ctx)
	} else if action == "add-entry" {
		ajaxMoneyAddEntry(ctx)
	}
}

type TransactionEntry struct {
	Date    time.Time
	Account int
	Value   int // 100*(real value)
}

type Transaction struct {
	Description string
	Entries []TransactionEntry
}

func ajaxMoneyAddEntry(ctx *Context) {
	var t Transaction
	log(ctx, LOG_DEBUG, "func ajaxMoneyAddEntry()")
	ctx.r.ParseForm()
	t.Description = ctx.r.FormValue("entry-description")
	for i:=1; ; i++ {
		date, err := time.Parse("2006-01-02", ctx.r.FormValue("entry" + strconv.Itoa(i) + "-date"))
		if err != nil {
			break
		}
		account, err := strconv.Atoi(ctx.r.FormValue("entry" + strconv.Itoa(i) + "-account"))
		if err != nil || account < 100 {
			break
		}
		value_, err := strconv.ParseFloat(ctx.r.FormValue("entry" + strconv.Itoa(i) + "-value"), 64)
		if err != nil {
			break
		}
		value := round(100.0*value_)
		if value==0 {
			break
		}
		t.Entries = append(t.Entries, TransactionEntry{Date: date, Account: account, Value: value})
	}
	log(ctx, LOG_DEBUG, fmt.Sprintf("ajaxMoneyAddEntry(): t=%v", t))
	err := db_money_add(t)
	if err != nil {
		log(ctx, LOG_ERR, "Error addding transaction: " + err.Error())
	}
}

func round(val float64) int {
	if val < 0 { return int(val-0.5) }
	return int(val+0.5)
}

func ajaxMoneyShow(ctx *Context) {
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
