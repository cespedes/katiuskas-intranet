package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type server struct {
	db *sqlx.DB
	r  *mux.Router
}

func NewServer() *server {
	var err error
	s := new(server)
	s.db, err = sqlx.Open("postgres", config("secret_db_conn"))
	if err != nil {
		panic(err)
	}
	s.routes()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}
