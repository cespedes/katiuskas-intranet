package katintranet

import (
	"net/http"

	"github.com/cespedes/api"
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

func (s *server) routes() {
	s.handler = api.NewServer()
	s.handler.AddMiddleware(s.MyHandler) // soon-to-be obsoleted (I think)
	s.handler.Set("server", s)

	/* Main page */
	s.handler.Handle("/", s.rootHandler)

	/* Static files */
	s.handler.Handle("/css/", StaticDir("/css/", "css"))
	s.handler.Handle("/img/", StaticDir("/img/", "img"))
	s.handler.Handle("/js/", StaticDir("/js/", "js"))

	/* Auth */
	s.handler.Handle("/auth/google", s.authGoogle)
	s.handler.Handle("/auth/mail", s.authMail)
	s.handler.Handle("/auth/hash", s.authHash)

	/* Telegram: */
	s.handler.Handle(s.config["telegram_webhook_path"], s.telegramBotHandler)

	/* Users */
	s.handler.Handle("/files/", StaticDir("/files/", "files"), requireRole("user"))
	s.handler.Handle("/public/", StaticDir("/public/", "../katiuskas/public"), requireRole("user"))

	/* Other pages: */
	s.handler.Handle("/my", s.myHandler, requireRole("user"))
	s.handler.Handle("/info", s.infoHandler, requireRole("user"))
	s.handler.Handle("/socios", s.sociosHandler, requireRole("user"))
	s.handler.Handle("/ajax/socios", s.ajaxSociosHandler, requireRole("user"))

	s.handler.Handle("/socio/{id}", s.viewSocioHandler, requireRole("board"))

	s.handler.Handle("/socio/new", s.socioNewHandler, requireRole("admin"))
	s.handler.Handle("/admin", s.adminHandler, requireRole("admin"))
	s.handler.Handle("/ajax/admin", s.ajaxAdminHandler, requireRole("admin"))
	s.handler.Handle("/actividades", s.activitiesHandler, requireRole("admin"))
	s.handler.Handle("/actividad/{id}", s.activityHandler, requireRole("admin"))
	s.handler.Handle("/ajax/activity", s.ajaxActivityHandler, requireRole("admin"))
	s.handler.Handle("/items", s.itemsHandler, requireRole("admin"))

	s.handler.Handle("/money", s.moneyHandler, requireRole("money"))
	s.handler.Handle("/money/{id}", s.moneyHandler, requireRole("money"))
	s.handler.Handle("/money/summary", s.moneySummaryHandler, requireRole("money"))
	s.handler.Handle("/ajax/money", s.ajaxMoneyHandler, requireRole("money"))

	s.handler.Handle("/repo/", StaticDir("/repo/", "../katiuskas"), requireRole("repo"))

	/* API */
	apiMux := api.NewServer()
	s.handler.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))
	apiMux.Handle("GET /user", s.apiGetUser)
	return
}
