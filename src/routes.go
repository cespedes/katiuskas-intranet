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

func roleMatcher(role string) mux.MatcherFunc {
	return func(r *http.Request, rm *mux.RouteMatch) bool {
		return Ctx(r).roles[role]
	}
}

func router() *mux.Router {
	r := mux.NewRouter()

	/* Main page */
	r.HandleFunc("/", rootHandler)
//	r.NewRoute().Roles("admin", "money").Path("/").HandlerFunc(rootHandler)

	/* Auth */
	r.Path("/auth/google").  HandlerFunc(authGoogle)
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
	r.Path("/my").                      HandlerFunc(myHandler)
	r.Path("/info").                    HandlerFunc(infoHandler)
	r.Path("/socios").                  HandlerFunc(sociosHandler)
	r.Path("/socio/new").               HandlerFunc(socioNewHandler)
	r.Path("/socio/id={id:[0-9]+}").    HandlerFunc(viewSocioHandler)
	r.Path("/actividades").             HandlerFunc(activitiesHandler)
	r.Path("/actividad/id={id:[0-9]+}").HandlerFunc(activityHandler)
	r.Path("/items").                   HandlerFunc(itemsHandler)
	r.Path("/money").                   HandlerFunc(moneyHandler)
	r.Path("/money/summary").           HandlerFunc(moneySummaryHandler)
	r.Path("/tgbot").                   HandlerFunc(tgbotHandler)
	r.Path("/tgbot.aif7eoca").          HandlerFunc(tgbotHandler)
	r.Path("/admin").MatcherFunc(roleMatcher("admin")).HandlerFunc(adminHandler)

	/* AJAX */
	r.HandleFunc("/ajax/admin",    ajaxAdminHandler)
	r.HandleFunc("/ajax/socios",   ajaxSociosHandler)
	r.HandleFunc("/ajax/activity", ajaxActivityHandler)
	r.HandleFunc("/ajax/money",    ajaxMoneyHandler)

	return r
}
