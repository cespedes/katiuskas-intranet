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

	/* Letsencrypt */
	r.PathPrefix("/.well-known/acme-challenge/").Handler(StaticDir("/.well-known/acme-challenge/", "/var/www/html/.well-known/acme-challenge"))

	/* Auth */
	r.Path("/auth/google").  HandlerFunc(authGoogle)
	r.Path("/auth/mail").    HandlerFunc(authMail)
	r.Path("/auth/hash").    HandlerFunc(authHash)

	/* Static files: */
	r.PathPrefix("/css/").Handler(StaticDir("/css/", "css"))
	r.PathPrefix("/img/").Handler(StaticDir("/img/", "img"))
	r.PathPrefix("/js/").Handler(StaticDir("/js/", "js"))

	/* Telegram: */
	r.Path("/tgbot.aif7eoca").          HandlerFunc(tgbotHandler)

	users := r.MatcherFunc(roleMatcher("user")).Subrouter()
	users.PathPrefix("/files/").Handler(StaticDir("/files/", "files"))

	/* Other pages: */
	users.Path("/my").         HandlerFunc(myHandler)
	users.Path("/info").       HandlerFunc(infoHandler)
	users.Path("/socios").     HandlerFunc(sociosHandler)
	users.Path("/ajax/socios").HandlerFunc(ajaxSociosHandler)

	board := r.MatcherFunc(roleMatcher("board")).Subrouter()
	board.Path("/socio/id={id:[0-9]+}").    HandlerFunc(viewSocioHandler)

	admin := r.MatcherFunc(roleMatcher("admin")).Subrouter()
	admin.Path("/socio/new").               HandlerFunc(socioNewHandler)
	admin.Path("/admin").                   HandlerFunc(adminHandler)
	admin.Path("/ajax/admin").              HandlerFunc(ajaxAdminHandler)
	admin.Path("/actividades").             HandlerFunc(activitiesHandler)
	admin.Path("/actividad/id={id:[0-9]+}").HandlerFunc(activityHandler)
	admin.Path("/ajax/activity").           HandlerFunc(ajaxActivityHandler)
	admin.Path("/items").                   HandlerFunc(itemsHandler)

	money := r.MatcherFunc(roleMatcher("money")).Subrouter()
	money.Path("/money").                   HandlerFunc(moneyHandler)
	money.Path("/money/summary").           HandlerFunc(moneySummaryHandler)
	money.Path("/ajax/money").              HandlerFunc(ajaxMoneyHandler)

	return r
}
