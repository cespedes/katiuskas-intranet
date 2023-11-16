package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/sessions"
)

func init() {
	// to be able to store "roles" (which is a map[string]bool) in session)
	gob.Register(map[string]bool{})
}

func Ctx(r *http.Request) *Context {
	ctx := r.Context().Value(0)
	return ctx.(*Context)
}

type Context struct {
	session       *sessions.Session
	session_saved bool
	ipaddr        string
	id            int
	roles         map[string]bool
}

var _session_store = sessions.NewCookieStore([]byte(config("cookie_secret")))

func (s *server) NewContext(r *http.Request) *http.Request {
	ctx := new(Context)
	/* ipaddr */
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx > -1 {
		ctx.ipaddr = r.RemoteAddr[:idx]
	}
	if tmp := r.Header["X-Cespedes-Remote-Addr"]; len(tmp) > 0 {
		ctx.ipaddr = tmp[0]
	}

	/* session */
	sess, err := _session_store.Get(r, "session")
	if err != nil {
		s.Log(r, LOG_WARNING, fmt.Sprintf("session_get: %q", err.Error()))
	}
	if sess.Values["start"] == nil {
		sess.Values["start"] = time.Now().Unix()
	}
	if count, ok := sess.Values["count"].(int); ok {
		sess.Values["count"] = count + 1
	} else {
		sess.Values["count"] = 1
	}
	// sess.Save(r, w)
	ctx.session = sess

	if id, ok := ctx.session.Values["id"].(int); ok {
		ctx.id = id
	}
	if roles, ok := ctx.session.Values["roles"].(map[string]bool); ok {
		ctx.roles = roles
	}
	return r.WithContext(context.WithValue(r.Context(), 0, ctx))
}

func (ctx *Context) Save(w http.ResponseWriter, r *http.Request) {
	if !ctx.session_saved {
		if ctx.session != nil {
			err := ctx.session.Save(r, w)
			if err != nil {
				log.Printf("session_save: %q", err.Error())
			}
		}
	}
	ctx.session_saved = true
}
