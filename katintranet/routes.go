package katintranet

import (
	"net/http"

	"github.com/cespedes/api"
	"github.com/gorilla/mux"
)

type Mux = mux.Router

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
	s.handler = api.NewServer()
	s.handler.Set("server", s)
	s.mux = mux.NewRouter()

	/* Main page */
	s.mux.HandleFunc("/", s.rootHandler)

	/* API */
	s.mux.PathPrefix("/api/v1").Handler(http.StripPrefix("/api/v1", s.apiHandler()))

	/* Lets' Encrypt */
	s.mux.PathPrefix("/.well-known/acme-challenge/").Handler(StaticDir("/.well-known/acme-challenge/", "/var/www/html/.well-known/acme-challenge"))

	/* Static files (no authentication, no context): */
	s.mux.PathPrefix("/css/").Handler(StaticDir("/css/", "css"))
	s.mux.PathPrefix("/img/").Handler(StaticDir("/img/", "img"))
	s.mux.PathPrefix("/js/").Handler(StaticDir("/js/", "js"))

	/* Auth */
	s.mux.Path("/auth/google").HandlerFunc(s.authGoogle)
	s.mux.Path("/auth/mail").HandlerFunc(s.authMail)
	s.mux.Path("/auth/hash").HandlerFunc(s.authHash)

	/* Telegram: */
	s.mux.Path(s.config["telegram_webhook_path"]).HandlerFunc(s.telegramBotHandler)

	users := s.mux.MatcherFunc(roleMatcher("user")).Subrouter()
	users.PathPrefix("/files/").Handler(StaticDir("/files/", "files"))
	users.PathPrefix("/public/").Handler(StaticDir("/public/", "../katiuskas/public"))

	/* Other pages: */
	users.Path("/my").HandlerFunc(s.myHandler)
	users.Path("/info").HandlerFunc(s.infoHandler)
	users.Path("/socios").HandlerFunc(s.sociosHandler)
	users.Path("/ajax/socios").HandlerFunc(s.ajaxSociosHandler)

	board := s.mux.MatcherFunc(roleMatcher("board")).Subrouter()
	board.Path("/socio/id={id:[0-9]+}").HandlerFunc(s.viewSocioHandler)

	admin := s.mux.MatcherFunc(roleMatcher("admin")).Subrouter()
	admin.Path("/socio/new").HandlerFunc(s.socioNewHandler)
	admin.Path("/admin").HandlerFunc(s.adminHandler)
	admin.Path("/ajax/admin").HandlerFunc(s.ajaxAdminHandler)
	admin.Path("/actividades").HandlerFunc(s.activitiesHandler)
	admin.Path("/actividad/id={id:[0-9]+}").HandlerFunc(s.activityHandler)
	admin.Path("/ajax/activity").HandlerFunc(s.ajaxActivityHandler)
	admin.Path("/items").HandlerFunc(s.itemsHandler)

	money := s.mux.MatcherFunc(roleMatcher("money")).Subrouter()
	money.Path("/money").HandlerFunc(s.moneyHandler)
	money.Path("/money/id={id:[0-9]+}").HandlerFunc(s.moneyHandler)
	money.Path("/money/summary").HandlerFunc(s.moneySummaryHandler)
	money.Path("/ajax/money").HandlerFunc(s.ajaxMoneyHandler)

	repo := s.mux.MatcherFunc(roleMatcher("repo")).Subrouter()
	// repo.Path("/repo").HandlerFunc(s.repoHandler)
	//repo.PathPrefix("/repo/").Handler(StaticDir("/repo/", "../katiuskas"))
	repo.PathPrefix("/repo/").Handler(http.StripPrefix("/repo/", http.FileServer(http.Dir("../katiuskas"))))

}
