package main

import (
	"net/http"
)

func infoHandler(w http.ResponseWriter, r *http.Request) {
	p := make(map[string]interface{})

	p["board"] = db_list_board()

	renderTemplate(w, r, "info", p)
}
