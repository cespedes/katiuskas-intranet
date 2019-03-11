package main

import (
	"fmt"
	"runtime"
	"net/http"
	"net/url"
)

const (
	LOG_EMERG int = iota    /* system is unusable */
	LOG_ALERT               /* action must be taken immediately */
	LOG_CRIT                /* critical conditions */
	LOG_ERR                 /* error conditions */
	LOG_WARNING             /* warning conditions */
	LOG_NOTICE              /* normal but significant condition */
	LOG_INFO                /* informational */
	LOG_DEBUG               /* debug-level messages */
)

func Log(r *http.Request, severity int, text string) {
	if severity <= LOG_ERR {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			text = fmt.Sprintf("(file=%v line=%v) %s", file, line, text)
		}
	}
	if severity <= LOG_NOTICE {
		http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendmessage?chat_id=%s&text=%s",
			config["telegram_bot_token"], config["telegram_log_chat_id"], url.QueryEscape(text)))
	}

	db.Exec("INSERT INTO log (severity, ipaddr, uid, text) VALUES ($1,$2,$3,$4)", severity, Ctx(r).ipaddr, Ctx(r).id, text)
}
