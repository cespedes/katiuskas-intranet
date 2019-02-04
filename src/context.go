package main

import (
	"fmt"
	"time"
	"strings"
	"net/http"
	"encoding/gob"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

func init() {
	// to be able to store "roles" (which is a map[string]bool) in session)
	gob.Register(map[string]bool{})
}

func Ctx(r *http.Request) *Context {
	c := context.Get(r, 0)
	if c == nil {
		ctx := new(Context)
		context.Set(r, 0, ctx)
		ctx.Get(r)
	}
	return context.Get(r, 0).(*Context)
}

type Context struct {
	session       *sessions.Session
	session_saved bool
	ipaddr        string
	id            int
	email         string
	person_type   int
	board         bool
	roles         map[string]bool
}

var _session_store = sessions.NewCookieStore([]byte("11UinL5BLSMVqivclTDo27qLVhIahkJM"))

func (ctx *Context) Get(r *http.Request) {
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
		Log(r, LOG_ERR, fmt.Sprintf("session_get: %q", err.Error()))
		return
	}
	if sess.Values["start"] == nil {
		sess.Values["start"] = time.Now().Unix()
	}
	if count, ok := sess.Values["count"].(int); ok==true {
		sess.Values["count"] = count + 1
	} else {
		sess.Values["count"] = 1
	}
	ctx.session = sess


	if email, ok := ctx.session.Values["email"].(string); ok {
		ctx.email = email
	}
	if id, ok := ctx.session.Values["id"].(int); ok {
		ctx.id = id
	}
	if person_type, ok := ctx.session.Values["type"].(int); ok {
		ctx.person_type = person_type
	}
	if board, ok := ctx.session.Values["board"].(bool); ok {
		ctx.board = board
	}
	if roles, ok := ctx.session.Values["roles"].(map[string]bool); ok {
		ctx.roles = roles
	}
	if ctx.id==0 && ctx.email!="" {
		id, person_type, board := db_mail_2_id(ctx.email)
		roles := db_get_roles(id)
		ctx.session.Values["id"] = id
		ctx.session.Values["type"] = person_type
		ctx.session.Values["board"] = board
		ctx.session.Values["roles"] = roles
		ctx.id = id
		ctx.person_type = person_type
		ctx.roles = roles
	}
}

func (ctx * Context) Save(w http.ResponseWriter, r *http.Request) {
	if ctx.session_saved == false {
		if ctx.session != nil {
			err := ctx.session.Save(r, w)
			if err != nil {
				Log(r, LOG_ERR, fmt.Sprintf("session_save: %q", err.Error()))
			}
		}
	}
	ctx.session_saved = true
}
