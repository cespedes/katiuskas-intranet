package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

func router() *mux.Router {
	r := mux.NewRouter()

	/* Main page */
	r.HandleFunc("/", rootHandler)

	/* Auth */
	r.HandleFunc("/auth/google", authGoogle)
	r.HandleFunc("/auth/facebook", authFacebook)

	/* Static files: */
	r.PathPrefix("/html/{file}").Handler(http.StripPrefix("/html/", http.FileServer(http.Dir("html"))))
	r.PathPrefix("/css/{file}").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	r.PathPrefix("/img/{file}").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	r.PathPrefix("/js/{file}").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	/* Letsencrypt */
	r.PathPrefix("/.well-known/acme-challenge/").Handler(http.StripPrefix("/.well-known/acme-challenge/", http.FileServer(http.Dir("/var/www/html/.well-known/acme-challenge"))))

	/* Other: */
	r.HandleFunc("/info", infoHandler)
	r.HandleFunc("/admin", adminHandler)
	r.HandleFunc("/admin/person={id:[0-9]+}", adminPersonHandler)

	/* AJAX */
	r.HandleFunc("/ajax/admin", ajaxAdminHandler)

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
