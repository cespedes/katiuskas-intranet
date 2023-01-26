package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func StaticDir(prefix, dir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(prefix) > len(path) {
			http.NotFound(w, r)
		}
		if path[len(path)-1:] == "/" {
			path = path + "index.html" /* serve index.html instead of directory index */
		}
		http.ServeFile(w, r, dir+"/"+path[len(prefix):])
	})
}

func roleMatcher(role string) mux.MatcherFunc {
	return func(r *http.Request, rm *mux.RouteMatch) bool {
		return Ctx(r).roles[role]
	}
}

func (s *server) routes() {
	s.r = mux.NewRouter()

	/* Main page */
	s.r.HandleFunc("/", s.rootHandler)

	/* Lets' Encrypt */
	s.r.PathPrefix("/.well-known/acme-challenge/").Handler(StaticDir("/.well-known/acme-challenge/", "/var/www/html/.well-known/acme-challenge"))

	/* Static files (no authentication, no context): */
	s.r.PathPrefix("/css/").Handler(StaticDir("/css/", "css"))
	s.r.PathPrefix("/img/").Handler(StaticDir("/img/", "img"))
	s.r.PathPrefix("/js/").Handler(StaticDir("/js/", "js"))

	/* Auth */
	s.r.Path("/auth/google").HandlerFunc(s.authGoogle)
	s.r.Path("/auth/mail").HandlerFunc(s.authMail)
	s.r.Path("/auth/hash").HandlerFunc(s.authHash)

	/* Telegram: */
	s.r.Path(config("telegram_webhook_path")).HandlerFunc(s.telegramBotHandler)

	users := s.r.MatcherFunc(roleMatcher("user")).Subrouter()
	users.PathPrefix("/files/").Handler(StaticDir("/files/", "files"))
	users.PathPrefix("/public/").Handler(StaticDir("/public/", "../katiuskas/public"))

	/* Other pages: */
	users.Path("/my").HandlerFunc(s.myHandler)
	users.Path("/info").HandlerFunc(s.infoHandler)
	users.Path("/socios").HandlerFunc(s.sociosHandler)
	users.Path("/ajax/socios").HandlerFunc(s.ajaxSociosHandler)

	board := s.r.MatcherFunc(roleMatcher("board")).Subrouter()
	board.Path("/socio/id={id:[0-9]+}").HandlerFunc(s.viewSocioHandler)

	admin := s.r.MatcherFunc(roleMatcher("admin")).Subrouter()
	admin.Path("/socio/new").HandlerFunc(s.socioNewHandler)
	admin.Path("/admin").HandlerFunc(s.adminHandler)
	admin.Path("/ajax/admin").HandlerFunc(s.ajaxAdminHandler)
	admin.Path("/actividades").HandlerFunc(s.activitiesHandler)
	admin.Path("/actividad/id={id:[0-9]+}").HandlerFunc(s.activityHandler)
	admin.Path("/ajax/activity").HandlerFunc(s.ajaxActivityHandler)
	admin.Path("/items").HandlerFunc(s.itemsHandler)

	money := s.r.MatcherFunc(roleMatcher("money")).Subrouter()
	money.Path("/money").HandlerFunc(s.moneyHandler)
	money.Path("/money/id={id:[0-9]+}").HandlerFunc(s.moneyHandler)
	money.Path("/money/summary").HandlerFunc(s.moneySummaryHandler)
	money.Path("/ajax/money").HandlerFunc(s.ajaxMoneyHandler)

	repo := s.r.MatcherFunc(roleMatcher("repo")).Subrouter()
	// repo.Path("/repo").HandlerFunc(s.repoHandler)
	//repo.PathPrefix("/repo/").Handler(StaticDir("/repo/", "../katiuskas"))
	repo.PathPrefix("/repo/").Handler(http.StripPrefix("/repo/", http.FileServer(http.Dir("../katiuskas"))))

}
