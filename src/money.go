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

func moneyHandler(w http.ResponseWriter, r *http.Request) {
	Log(r, LOG_DEBUG, "Page /money")

	if !(Ctx(r).roles["admin"] || Ctx(r).roles["money"]) {
		http.Redirect(w, r, "/", http.StatusFound)
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
	renderTemplate(w, r, "money", p)
}

func moneySummaryHandler(w http.ResponseWriter, r *http.Request) {
	Log(r, LOG_DEBUG, "Page /money/summary")

	if !(Ctx(r).roles["admin"] || Ctx(r).roles["money"]) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	p := make(map[string]interface{})
	today := time.Now()

	p["today"] = today.Format("2006-01-02")
	p["last_365d"] = today.Add(-365 * 24 * time.Hour).Format("2006-01-02")
	p["this_year"] = today.Year()
	p["last_year"] = today.Year() - 1
	p["second_to_last_year"] = today.Year() - 2
	p["accounts"] = db_get_accounts()
	renderTemplate(w, r, "money-summary", p)
}

func ajaxMoneyHandler(w http.ResponseWriter, r *http.Request) {
	Log(r, LOG_DEBUG, "Page /ajax/money")

	action := r.FormValue("action")
	if action == "show-money" {
		ajaxMoneyShow(w, r)
	} else if action == "add-entry" {
		ajaxMoneyAddEntry(w, r)
	} else if action == "show-money-summary" {
		ajaxMoneySummaryShow(w, r)
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

func ajaxMoneyAddEntry(w http.ResponseWriter, r *http.Request) {
	var t Transaction
	Log(r, LOG_DEBUG, "func ajaxMoneyAddEntry()")
	r.ParseForm()
	t.Description = r.FormValue("entry-description")
	for i:=1; ; i++ {
		date, err := time.Parse("2006-01-02", r.FormValue("entry" + strconv.Itoa(i) + "-date"))
		if err != nil {
			break
		}
		account, err := strconv.Atoi(r.FormValue("entry" + strconv.Itoa(i) + "-account"))
		if err != nil || account < 100 {
			break
		}
		value_, err := strconv.ParseFloat(r.FormValue("entry" + strconv.Itoa(i) + "-value"), 64)
		if err != nil {
			break
		}
		value := round(100.0*value_)
		if value==0 {
			break
		}
		t.Entries = append(t.Entries, TransactionEntry{Date: date, Account: account, Value: value})
	}
	Log(r, LOG_DEBUG, fmt.Sprintf("ajaxMoneyAddEntry(): t=%v", t))
	err := db_money_add(t)
	if err != nil {
		Log(r, LOG_ERR, "Error addding transaction: " + err.Error())
	}
}

func round(val float64) int {
	if val < 0 { return int(val-0.5) }
	return int(val+0.5)
}

func ajaxMoneyShow(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	account, _ := strconv.Atoi(r.FormValue("account"))
	from := r.FormValue("from")
	lines := db_get_money(account, from)
	if len(lines)==0 {
		fmt.Fprint(w, "No lines to display\n")
		return
	}
	fmt.Fprint(w, `
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

	fmt.Fprint(w, "<table class=\"money\">\n")
	for _, line := range lines {
		fmt.Fprint(w, "<tr>")
		fmt.Fprintf(w, "<td>%v</td>", line["date"])
		fmt.Fprintf(w, "<td>%v</td>", line["description"])
		fmt.Fprintf(w, "<td>%v</td>", line["value"])
		fmt.Fprintf(w, "<td>%v</td>", line["balance"])
		fmt.Fprint(w, "</tr>\n")
	}
	fmt.Fprint(w, "</table>\n")
}

func ajaxMoneySummaryShow(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	from := r.FormValue("from")
	lines := db_get_money_summary(from)
	if len(lines)==0 {
		fmt.Fprint(w, "No lines to display\n")
		return
	}
	fmt.Fprint(w, `
<style>
table.money tr td:first-child {
  width: 6ex;
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

	fmt.Fprint(w, "<table class=\"money\">\n")
	for _, line := range lines {
		fmt.Fprint(w, "<tr>")
		fmt.Fprintf(w, "<td>%v</td>", line["id"])
		fmt.Fprintf(w, "<td>%v</td>", line["account"])
		fmt.Fprintf(w, "<td>%v</td>", line["value"])
		fmt.Fprintf(w, "<td>%v</td>", line["balance"])
		fmt.Fprint(w, "</tr>\n")
	}
	fmt.Fprint(w, "</table>\n")
}
