package main

/*
import (
	"fmt"
	"time"
	"strconv"
	"github.com/gorilla/mux"
)
*/

func itemsHandler(ctx *Context) {
	log(ctx, LOG_DEBUG, "Page /items")

	p := make(map[string]interface{})

	p["items"] = db_list_items()
	renderTemplate(ctx, "items", p)
}
