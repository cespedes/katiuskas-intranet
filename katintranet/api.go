package katintranet

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type moneyAccount struct {
	ID       int      `json:"id"`
	ParentID *int     `db:"parent_id" json:"parent_id,omitempty"`
	Name     string   `json:"name"`
	Code     *string  `json:"code,omitempty"`
	Balance  *float64 `json:"balance,omitempty"`
}

type moneySplit struct {
	TransactionID int     `json:"transaction_id"`
	AccountID     int     `json:"account_id"`
	Value         float64 `json:"value"`
	Balance       float64 `json:"balance"` // only relevany after queries, to show balance of account after this
}

type moneyTransaction struct {
	ID          int
	Description string
	Datetime    time.Time
	Splits      []moneySplit
}

func (s *server) apiHandler() http.Handler {
	api := http.NewServeMux()

	api.HandleFunc("GET /user", s.apiGetUser)

	api.HandleFunc("GET /money/accounts", s.apiGetMoneyAccounts)
	api.HandleFunc("GET /money/accounts/{account}", s.apiGetMoneyAccountsAccount)
	api.HandleFunc("POST /money/accounts", s.apiPostMoneyAccounts)
	api.HandleFunc("DELETE /money/accounts/{account}", s.apiDeleteMoneyAccountsAccount)
	api.HandleFunc("GET /money/transactions/{transaction}", s.apiGetMoneyTransactionsTransaction)
	api.HandleFunc("POST /money/transactions", s.apiPostMoneyTransactions)
	api.HandleFunc("PUT /money/transactions/{transaction}", s.apiPutMoneyTransactions)
	api.HandleFunc("DELETE /money/transactions/{transaction}", s.apiDeleteMoneyTransactionsTransaction)

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Called API\n")
		fmt.Fprintf(w, "Path: %v\n", r.URL)
		fmt.Fprintf(w, "Headers: %v\n", r.Header)
	})

	// return api
	mux := http.NewServeMux()
	mux.Handle("/", middlewareAuth(api))
	return mux
}

func clientLogError(r *http.Request, format string, a ...any) {
	args := []any{r.RemoteAddr, r.Method, r.URL.Path}
	args = append(args, a...)
	log.Printf("%s: %s %s: "+format, args...)
}

func middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			clientLogError(r, "no 'Authorization' header (not authorized)")
			httpError(w, "Not authorized.", http.StatusUnauthorized)
			return
		}

		fmt.Println("middleware")
		fmt.Printf("Path: %v\n", r.URL)
		fmt.Printf("Headers: %v\n", r.Header)
		// Our middleware logic goes here...
		next.ServeHTTP(w, r)
	})
}

func (s *server) apiGetUser(w http.ResponseWriter, r *http.Request) {
	ctx := Ctx(r)
	user := s.DBgetUserinfo(ctx.id)
	s.outJSON(w, user)
}

func (s *server) apiGetMoneyAccounts(w http.ResponseWriter, r *http.Request) {
	var accounts []moneyAccount

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
	var err error
	var account moneyAccount
	decoder := json.NewDecoder(r.Body)

	if err = decoder.Decode(&account); err != nil {
		httpError(w, err)
		return
	}

	if account.ID == 0 {
		httpError(w, "must specify `id` for new account")
		return
	}
	if account.Name == "" {
		httpError(w, "must specify `name` for new account")
		return
	}

	_, err = s.db.Exec(`
		INSERT INTO "account" (id,parent_id,name,code)
		VALUES ($1,$2,$3,$4)
	`, account.ID, account.ParentID, account.Name, account.Code)
	if err != nil {
		httpError(w, err)
		return
	}
	s.outJSON(w, account)
}

func (s *server) apiDeleteMoneyAccountsAccount(w http.ResponseWriter, r *http.Request) {
	var err error
	var account moneyAccount

	id := r.PathValue("account")

	err = s.db.Get(&account, `
		DELETE FROM "account" WHERE id=$1
		RETURNING id,parent_id,name,code
	`, id)
	if err != nil {
		httpError(w, err)
		return
	}
	s.outJSON(w, account)
}

func (s *server) apiGetMoneyTransactionsTransaction(w http.ResponseWriter, r *http.Request) {
	httpError(w, "Unimplemented: GET /money/accounts/transactions/transaction")
}

func (s *server) apiPostMoneyTransactions(w http.ResponseWriter, r *http.Request) {
	httpError(w, "Unimplemented: POST /money/transactions")
}

func (s *server) apiPutMoneyTransactions(w http.ResponseWriter, r *http.Request) {
	httpError(w, "Unimplemented: PUT /money/transactions")
}

func (s *server) apiDeleteMoneyTransactionsTransaction(w http.ResponseWriter, r *http.Request) {
	httpError(w, "Unimplemented: DELETE /money/transactions/(transactions)")
}
