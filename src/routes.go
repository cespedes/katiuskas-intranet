package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

func router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)

	/* Static files: */
	r.PathPrefix("/html/{file}").Handler(http.StripPrefix("/html/", http.FileServer(http.Dir("html"))))
	r.PathPrefix("/css/{file}").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	r.PathPrefix("/img/{file}").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	r.PathPrefix("/js/{file}").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

//	/* Wiki */
//	r.HandleFunc("/wiki/{title}", wikiViewHandler)
//	r.HandleFunc("/wiki-edit/{title}", wikiEditHandler)
//	r.HandleFunc("/wiki-save/{title}", wikiSaveHandler)
//
//	/* Blog */
//	r.HandleFunc("/blog/{blog_id}", blogHandler)
//	r.HandleFunc("/blog/{blog_id}/{blog_name}", blogHandler)

	/* Auth */
	r.HandleFunc("/auth/google", authGoogle)
	r.HandleFunc("/auth/facebook", authFacebook)

//	/* AJAX */
//	r.HandleFunc("/ajax/markdown", ajaxMarkdownHandler)
//
//	/* Info */
//	r.HandleFunc("/info", infoHandler)
//
//	/* Sessions */
//	r.HandleFunc("/session", sessionHandler)

	/* Letsencrypt */
	r.PathPrefix("/.well-known/acme-challenge/").Handler(http.StripPrefix("/.well-known/acme-challenge/", http.FileServer(http.Dir("/var/www/html/.well-known/acme-challenge"))))

	return r
}
