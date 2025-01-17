package katintranet

import (
	"html/template"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, p map[string]interface{}) {
	p["id"] = C(r).id
	p["roles"] = C(r).roles
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var templates *template.Template

func init() {
	funcMap := template.FuncMap{
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))
}
