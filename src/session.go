package main

import (
	"fmt"
	"time"
	"net/http"
	"github.com/gorilla/sessions"
	"github.com/gorilla/context"
)

var _session_store = sessions.NewCookieStore([]byte("11UinL5BLSMVqivclTDo27qLVhIahkJM"))

func session_get(w http.ResponseWriter, r *http.Request) map[interface{}]interface{} {
	session, ok := context.GetOk(r, "session")

	if ok == false {
		sess, err := _session_store.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return nil
		}
		if sess.Values["start"] == nil {
			sess.Values["start"] = time.Now().Unix()
		}
		if sess.Values["count"] == nil {
			sess.Values["count"] = 1
		} else {
			sess.Values["count"] = sess.Values["count"].(int) + 1
		}
		context.Set(r, "session", sess)
		session = sess
	}

	return session.(*sessions.Session).Values
}

func session_save(w http.ResponseWriter, r *http.Request) {
	if context.Get(r, "session_saved") == nil {
		_, ok := context.GetOk(r, "session")
		if ok == false {
			session_get(w, r)
		}
		session := context.Get(r, "session").(*sessions.Session)
		if session != nil {
			err := session.Save(r, w)
			if err != nil {
				fmt.Println("session_save:", err)
			}
		}
	}
	context.Set(r, "session_saved", true)
}
