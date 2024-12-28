package katintranet

import (
	"net/http"

	"github.com/cespedes/api"
)

type server struct {
	handler        *api.Server
	config         map[string]string
	db             *DB
	mux            *Mux
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

	err = s.SessionInit()
	if err != nil {
		panic("SessionInit(): " + err.Error())
	}

	s.routes()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
