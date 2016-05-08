package main

import (
	"net/http"
	"html/template"
)

func renderTemplate(ctx *Context, tmpl string, p interface{}) {
	ctx.Save()
	err := templates.ExecuteTemplate(ctx.w, tmpl+".html", p)
	if err != nil {
		http.Error(ctx.w, err.Error(), http.StatusInternalServerError)
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
