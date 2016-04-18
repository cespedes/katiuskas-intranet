package main

import (
	"fmt"
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

func sessionHandler(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "session").(*sessions.Session)

	if session.Values["in-session"] == nil {
		session.Values["in-session"] = 1
	} else {
		session.Values["in-session"] = session.Values["in-session"].(int) + 1
	}
	session.Save(r, w)
	fmt.Fprintln(w, "Current session:", session.Values)

	fmt.Fprintf(w, "Current session: %d requests during %s\n", session.Values["count"].(int),
		time.Now().Sub(time.Unix(session.Values["start"].(int64), 0)))

//	// Set some session values.
//	session.Values["foo"] = "bar"
//	session.Values[42] = 43
//	// Save it before we write to the response/return from the handler.
//	session.Save(r, w)
}
