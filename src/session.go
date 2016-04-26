package main

import (
	"time"
	"net/http"
	"github.com/gorilla/sessions"
	"github.com/gorilla/context"
)

var store = sessions.NewCookieStore([]byte("11UinL5BLSMVqivclTDo27qLVhIahkJM"))

func session_init(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "session")

	if err != nil {
		http.Error(w, err.Error(), 500)
		return err
	}

	if session.Values["start"] == nil {
		session.Values["start"] = time.Now().Unix()
	}
	if session.Values["count"] == nil {
		session.Values["count"] = 1
	} else {
		session.Values["count"] = session.Values["count"].(int) + 1
	}

	context.Set(r, "session", session)
	session.Save(r, w)
	return nil
}
