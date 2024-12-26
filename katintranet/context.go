package katintranet

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

// Do we really need a global variable for this?  Ugly, ugly!
var _session_store *sessions.CookieStore

func (s *server) SessionInit() error {
	gob.Register(map[string]bool{})
	_session_store = sessions.NewCookieStore([]byte(s.config["cookie_secret"]))
	return nil
}

type Context struct {
	session       *sessions.Session
	session_saved bool
	ipaddr        string
	id            int
	roles         map[string]bool
}
type contextKey struct{}

func Ctx(r *http.Request) *Context {
	ctx := r.Context().Value(contextKey{})
	return ctx.(*Context)
}

func getIPFromRequest(r *http.Request) string {
	if tmp := r.Header["X-Forwarded-For"]; len(tmp) > 0 {
		return tmp[0]
	}
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx > -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

func (s *server) NewContext(r *http.Request) *http.Request {
	ctx := new(Context)
	/* ipaddr */
	ctx.ipaddr = getIPFromRequest(r)

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

	/* token */
	token, err := getTokenFromHeaders(r)
	if err != nil {
		log.Println(err.Error())
	}
	var id int
	err = s.db.Get(&id, `SELECT person_id FROM token WHERE token=$1`, token)
	if err == nil {
		ctx.id = id
	}

	return r.WithContext(context.WithValue(r.Context(), contextKey{}, ctx))
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
