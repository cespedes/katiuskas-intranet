package main

import (
	"net/http"
	"html/template"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, p interface{}) {
	session := context.Get(r, "session").(*sessions.Session)
	session_saved := context.Get(r, "session_saved")
	if session_saved == nil {
		session.Save(r, w)
		context.Set(r, "session_saved", true)
	}
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var templates *template.Template

func templates_init() {
        funcMap := template.FuncMap{
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
        }
	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))
}
