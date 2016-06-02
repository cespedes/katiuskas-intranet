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

func (r Router) StaticDir(prefix, dir string) {
	r.r.PathPrefix(prefix + "{file}").Handler(http.StripPrefix(prefix, http.FileServer(http.Dir(dir))))
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

func router() *Router {
	r := NewRouter()

	/* Main page */
	r.HandleFunc("/", rootHandler)

	/* Auth */
	r.HandleFunc("/auth/google", authGoogle)
	r.HandleFunc("/auth/facebook", authFacebook)

	/* Static files: */
	r.StaticDir("/html/", "html")
	r.StaticDir("/css/", "css")
	r.StaticDir("/img/", "img")
	r.StaticDir("/js/", "js")
	r.StaticDir("/files/", "files")

	/* Letsencrypt */
	r.StaticDir("/.well-known/acme-challenge/", "/var/www/html/.well-known/acme-challenge")

	/* Other: */
	r.HandleFunc("/info", infoHandler)
	r.HandleFunc("/query", queryHandler)
	r.HandleFunc("/query/person={id:[0-9]+}", queryPersonHandler)

	/* AJAX */
	r.HandleFunc("/ajax/admin", ajaxAdminHandler)
	r.HandleFunc("/ajax/query", ajaxQueryHandler)

//	/* Wiki */
//	r.HandleFunc("/wiki/{title}", wikiViewHandler)
//	r.HandleFunc("/wiki-edit/{title}", wikiEditHandler)
//	r.HandleFunc("/wiki-save/{title}", wikiSaveHandler)
//
//	/* Blog */
//	r.HandleFunc("/blog/{blog_id}", blogHandler)
//	r.HandleFunc("/blog/{blog_id}/{blog_name}", blogHandler)

//	/* AJAX */
//	r.HandleFunc("/ajax/markdown", ajaxMarkdownHandler)
//
//	/* Info */
//	r.HandleFunc("/info", infoHandler)
//
//	/* Sessions */
//	r.HandleFunc("/session", sessionHandler)

	return r
}
