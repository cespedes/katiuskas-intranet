package main

import (
	"net/http"
)

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	Log(r, LOG_DEBUG, "Page /items")

	p := make(map[string]interface{})

	p["items"] = db_list_items()
	renderTemplate(w, r, "items", p)
}
