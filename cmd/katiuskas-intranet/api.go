package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (s *server) apiHandler() http.Handler {
	api := http.NewServeMux()

	api.HandleFunc("GET /user", s.apiGetUser)

	api.HandleFunc("GET /money/accounts", s.apiGetMoneyAccounts)
	api.HandleFunc("GET /money/accounts/{account}", s.apiGetMoneyAccountsAccount)
	api.HandleFunc("POST /money/accounts", s.apiPostMoneyAccounts)
	api.HandleFunc("DELETE /money/accounts/{account}", s.apiDeleteMoneyAccountsAccount)
	api.HandleFunc("GET /money/transactions/{transaction}", s.apiGetMoneyTransactionsTransaction)
	api.HandleFunc("POST /money/transactions", s.apiPostMoneyTransactions)
	api.HandleFunc("DELETE /money/transactions/{transaction}", s.apiDeleteMoneyTransactionsTransaction)

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Called API\n")
		fmt.Fprintf(w, "Path: %v\n", r.URL)
		fmt.Fprintf(w, "Headers: %v\n", r.Header)
	})

	return api
}

func (s *server) apiGetUser(w http.ResponseWriter, r *http.Request) {
	ctx := Ctx(r)
	user := s.DBgetUserinfo(ctx.id)
	s.outJSON(w, user)
}

func (s *server) apiGetMoneyAccounts(w http.ResponseWriter, r *http.Request) {
	var accounts []struct {
		ID       int      `json:"id"`
		ParentID *int     `db:"parent_id" json:"parent_id,omitempty"`
		Name     string   `json:"name"`
		Code     *string  `json:"code,omitempty"`
		Balance  *float64 `json:"balance,omitempty"`
	}

	sql := `
		SELECT
			a.id,
			a.parent_id,
			a.name,
			a.code,
			sum((s.value*100)::int)::real/100 AS balance
                FROM account a
                  LEFT JOIN split s ON a.id=s.account_id
                GROUP BY a.id
                ORDER BY a.id;
        `
	err := s.db.Select(&accounts, sql)
	if err != nil {
		httpError(w, err)
		return
	}
	s.outJSON(w, accounts)

}

func (s *server) apiGetMoneyAccountsAccount(w http.ResponseWriter, r *http.Request) {
	account := r.PathValue("account")

	var transactions []struct {
		ID          int     `json:"id"`
		Date        string  `json:"date"`
		Description string  `json:"description"`
		Value       float64 `json:"value"`
		Balance     float64 `json:"balance"`
	}

	sql := `
		SELECT
			transaction_id AS id,
			to_char(datetime,'YYYY-MM-DD') AS date,
			description,
			to_char(value,'FM999990.00') AS value,
			to_char(balance,'FM999990.00') AS balance
		FROM money
		WHERE account_id=$1
	`

	err := s.db.Select(&transactions, sql, account)
	if err != nil {
		httpError(w, err)
		return
	}
	s.outJSON(w, transactions)
}

func (s *server) apiPostMoneyAccounts(w http.ResponseWriter, r *http.Request) {
	httpError(w, "Unimplemented: POST /money/accounts")
}

func (s *server) apiDeleteMoneyAccountsAccount(w http.ResponseWriter, r *http.Request) {
	account := r.PathValue("account")
	httpError(w, fmt.Sprintf("Unimplemented: DELETE /money/accounts/%s", account))
}

func (s *server) apiGetMoneyTransactionsTransaction(w http.ResponseWriter, r *http.Request) {
	httpError(w, "Unimplemented: GET /money/accounts/transactions/transaction")
}
func (s *server) apiPostMoneyTransactions(w http.ResponseWriter, r *http.Request) {
	httpError(w, "Unimplemented: POST /money/transactions")
}
func (s *server) apiDeleteMoneyTransactionsTransaction(w http.ResponseWriter, r *http.Request) {
	httpError(w, "Unimplemented: DELETE /money/transactions/(transactions)")
}

// outJSON writes the JSON-encoded object v to the http.ResponseWriter w.
func (s *server) outJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	err := e.Encode(v)
	if err != nil {
		httpError(w, err)
		return
	}
}

// httpError sends a HTTP error as a response
func httpError(w http.ResponseWriter, err any, codes ...int) {
	code := http.StatusInternalServerError

	if err == sql.ErrNoRows {
		code = http.StatusNotFound
	}

	if er, ok := err.(error); ok {
		var e ErrHttpStatus
		if errors.As(er, &e) {
			code = e.Status
		}
	}

	if len(codes) > 0 {
		code = codes[0]
	}

	httpMessage(w, code, "error", fmt.Sprint(err))
}

func httpMessage(w http.ResponseWriter, code int, label string, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{%q: %q}\n", label, msg)
}

type ErrHttpStatus struct {
	Status int
	Err    error
}

func (e ErrHttpStatus) Error() string {
	return e.Err.Error()
}

func (e ErrHttpStatus) Unwrap() error {
	return e.Err
}

func httpInfo(w http.ResponseWriter, msg any, codes ...int) {
	code := http.StatusOK

	if len(codes) > 0 {
		code = codes[0]
	}

	httpMessage(w, code, "info", fmt.Sprint(msg))
}
