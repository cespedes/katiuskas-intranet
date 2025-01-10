package katintranet

import (
	"net/http"

	"github.com/cespedes/api"
)

type server struct {
	handler        *api.Server
	config         map[string]string
	db             *DB
	telegramBotAPI *TelegramBotAPI
}

func NewServer(args []string) *server {
	var err error

	s := new(server)

	// Configuration must be initialized before anything else:
	err = s.ConfigInit(args)
	if err != nil {
		panic("ConfigInit(): " + err.Error())
	}

	err = s.DBInit()
	if err != nil {
		panic("DBInit(): " + err.Error())
	}

	err = s.TelegramInit()
	if err != nil {
		panic("TelegramInit(): " + err.Error())
	}

	s.routes()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// s.mux.ServeHTTP(w, r)
	s.handler.ServeHTTP(w, r)
}

// requireRole is a permission function that reports
// if the user is root or belongs to any of the specified roles.
func requireRole(roles ...string) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		return hasRole(r, roles...)
	}
}

func hasRole(r *http.Request, roles ...string) bool {
	c := C(r)
	if c == nil {
		return false
	}

	for _, role := range append(roles, "root") {
		if c.roles[role] {
			return true
		}
	}
	return false
}
