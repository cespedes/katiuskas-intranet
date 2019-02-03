package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

type Router struct {
	r *mux.Router
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.r.ServeHTTP(w, req)
}

func NewRouter() *Router {
	r := new(Router)
	r.r = mux.NewRouter()
	return r
}

func (r Router) HandleFunc(path string, f func(ctx *Context)) {
	r.r.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		ctx := new(Context)
		ctx.w = w
		ctx.r = req
		ctx.Get()
		f(ctx)
	})
}

func (r Router) StaticDir(prefix, dir string) {
	r.r.PathPrefix(prefix).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func router() *Router {
	r := NewRouter()

	/* Main page */
	r.HandleFunc("/", rootHandler)

	/* Auth */
	r.HandleFunc("/auth/google", authGoogle)
	r.HandleFunc("/auth/facebook", authFacebook)
	r.HandleFunc("/auth/mail", authMail)
	r.HandleFunc("/auth/hash", authHash)

	/* Static files: */
	r.StaticDir("/html/", "html")
	r.StaticDir("/css/", "css")
	r.StaticDir("/img/", "img")
	r.StaticDir("/js/", "js")
	r.StaticDir("/files/", "files")

	/* Letsencrypt */
	r.StaticDir("/.well-known/acme-challenge/", "/var/www/html/.well-known/acme-challenge")

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
	r.HandleFunc("/ajax/admin", ajaxAdminHandler)
	r.HandleFunc("/ajax/socios", ajaxSociosHandler)
	r.HandleFunc("/ajax/activity", ajaxActivityHandler)
	r.HandleFunc("/ajax/money", ajaxMoneyHandler)

	return r
}
