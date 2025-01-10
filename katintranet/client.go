package katintranet

import (
	"log"
	"net/http"
	"strings"

	"github.com/cespedes/api"
)

type Client struct {
	ipaddr string
	id     int
	token  string
	roles  map[string]bool
}

// C returns the client associated with the current request.
func C(r *http.Request) *Client {
	c, _ := api.Get(r, "client").(*Client)
	return c
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
	var err error

	c := new(Client)
	/* ipaddr */
	c.ipaddr = getIPFromRequest(r)

	/* token */
	token, err := GetTokenFromHeaders(r)
	if err != nil {
		log.Println(err.Error())
	}
	var id int
	err = s.db.Get(&id, `SELECT person_id FROM token WHERE token=$1`, token)
	if err == nil {
		c.id = id
		c.roles = s.DBgetRoles(id)
	}

	// Client
	r = api.Set(r, "client", c)

	return r
}
