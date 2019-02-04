package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

func StaticDir(prefix, dir string) (http.Handler) {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(prefix) > len(path) {
			http.NotFound(w, r)
		}
		if path[len(path)-1:] == "/" {
			path = path + "index.html" /* serve index.html instead of directory index */
		}
		http.ServeFile(w, r, dir + "/" + path[len(prefix):])
	})
}

/*
func (r * MyRoute) Roles(roles ...string) *MyRoute {
	return (*mux.Router)(r).MatcherFunc(func(r *http.Request, rm *RouteMatch) bool {
	})
}

func roleMatcher(roles ...string) mux.MatcherFunc {
	return func(r *http.Request, rm *mux.RouteMatch) bool {
		return true
	}
}
*/

func router() *mux.Router {
	r := mux.NewRouter()

	/* Main page */
	r.HandleFunc("/", rootHandler)
//	r.NewRoute().Roles("admin", "money").Path("/").HandlerFunc(rootHandler)

	/* Auth */
	r.Path("/auth/google").  HandlerFunc(authGoogle)
	r.Path("/auth/facebook").HandlerFunc(authFacebook)
	r.Path("/auth/mail").    HandlerFunc(authMail)
	r.Path("/auth/hash").    HandlerFunc(authHash)

	/* Static files: */
	r.PathPrefix("/html/").Handler(StaticDir("/html/", "html"))
	r.PathPrefix("/css/").Handler(StaticDir("/css/", "css"))
	r.PathPrefix("/img/").Handler(StaticDir("/img/", "img"))
	r.PathPrefix("/js/").Handler(StaticDir("/js/", "js"))
	r.PathPrefix("/files/").Handler(StaticDir("/files/", "files"))

	/* Letsencrypt */
	r.PathPrefix("/.well-known/acme-challenge/").Handler(StaticDir("/.well-known/acme-challenge/", "/var/www/html/.well-known/acme-challenge"))

	/* Other pages: */
	r.HandleFunc("/my", myHandler)
	r.HandleFunc("/info", infoHandler)
	r.HandleFunc("/socios", sociosHandler)
	r.HandleFunc("/socio/new", socioNewHandler)
	r.HandleFunc("/socio/id={id:[0-9]+}", viewSocioHandler)
	r.HandleFunc("/actividades", activitiesHandler)
	r.HandleFunc("/actividad/id={id:[0-9]+}", activityHandler)
	r.HandleFunc("/items", itemsHandler)
	r.HandleFunc("/money", moneyHandler)
	r.HandleFunc("/money/summary", moneySummaryHandler)
	r.HandleFunc("/tgbot", tgbotHandler)

	/* AJAX */
	r.HandleFunc("/ajax/admin",    ajaxAdminHandler)
	r.HandleFunc("/ajax/socios",   ajaxSociosHandler)
	r.HandleFunc("/ajax/activity", ajaxActivityHandler)
	r.HandleFunc("/ajax/money",    ajaxMoneyHandler)

	return r
}
