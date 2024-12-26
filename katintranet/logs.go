package katintranet

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
)

const (
	LOG_EMERG   int = iota /* system is unusable */
	LOG_ALERT              /* action must be taken immediately */
	LOG_CRIT               /* critical conditions */
	LOG_ERR                /* error conditions */
	LOG_WARNING            /* warning conditions */
	LOG_NOTICE             /* normal but significant condition */
	LOG_INFO               /* informational */
	LOG_DEBUG              /* debug-level messages */
)

var log_level = [...]string{
	"EMERG",
	"ALERT",
	"CRIT",
	"ERROR",
	"WARNING",
	"NOTICE",
	"INFO",
	"DEBUG",
}

func (s *server) Log(r *http.Request, severity int, text string) {
	var pref1, pref2 string

	if severity <= LOG_ERR {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			pref1 = fmt.Sprintf("(file=%v line=%v) ", file, line)
		}
	}
	pref2 = log_level[severity] + ": " + Ctx(r).ipaddr + ": "

	if severity <= LOG_NOTICE {
		http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendmessage?chat_id=%s&text=%s",
			s.config["telegram_bot_token"], s.config["telegram_log_chat_id"],
			url.QueryEscape(pref1+pref2+text)))
	}

	s.db.Exec("INSERT INTO log (severity, ipaddr, uid, text) VALUES ($1,$2,$3,$4)", severity, Ctx(r).ipaddr, Ctx(r).id, pref1+text)
}
