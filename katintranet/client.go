package katintranet

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cespedes/katiuskas-intranet/katintranet/token"
	"github.com/gorilla/sessions"
)

// Do we really need a global variable for this?  Ugly, ugly!
var _session_store *sessions.CookieStore

func (s *server) SessionInit() error {
	gob.Register(map[string]bool{})
	_session_store = sessions.NewCookieStore([]byte(s.config["cookie_secret"]))
	return nil
}

type Client struct {
	session       *sessions.Session
	session_saved bool
	ipaddr        string
	id            int
	roles         map[string]bool
}
type clientKey struct{}

func C(r *http.Request) *Client {
	c := r.Context().Value(clientKey{})
	return c.(*Client)
}

func HasRole(r *http.Request, roles ...string) bool {
	for _, role := range roles {
		if C(r).roles[role] {
			return true
		}
	}
	return false
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
	c := new(Client)
	/* ipaddr */
	c.ipaddr = getIPFromRequest(r)

	/* session */
	sess, err := _session_store.Get(r, "session")
	if err != nil {
		s.Log(r, LOG_WARNING, fmt.Sprintf("session_get: %q", err.Error()))
	}
	// sess.Save(r, w)
	c.session = sess

	if id, ok := c.session.Values["id"].(int); ok {
		c.id = id
	}
	if roles, ok := c.session.Values["roles"].(map[string]bool); ok {
		c.roles = roles
	}

	/* token */
	token, err := token.GetFromHeaders(r)
	if err != nil {
		log.Println(err.Error())
	}
	var id int
	err = s.db.Get(&id, `SELECT person_id FROM token WHERE token=$1`, token)
	if err == nil {
		c.id = id
	}

	return r.WithContext(context.WithValue(r.Context(), clientKey{}, c))
}

func (c *Client) Save(w http.ResponseWriter, r *http.Request) {
	if !c.session_saved {
		if c.session != nil {
			err := c.session.Save(r, w)
			if err != nil {
				log.Printf("session_save: %q", err.Error())
			}
		}
	}
	c.session_saved = true
}
