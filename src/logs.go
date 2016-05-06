package main

import (
	"fmt"
	"strings"
	"runtime"
	"net/http"
	"net/url"
	"github.com/gorilla/context"
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

func log(w http.ResponseWriter, r *http.Request, severity int, text string) {
	if severity <= LOG_ERR {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			text = fmt.Sprintf("(file=%v line=%v) %s", file, line, text)
		}
	}
	if severity <= LOG_NOTICE {
		const bot_token = "204701695:AAFkgoxJPCUWpXTWDQco33I97y5BJIHmOKU" /* KatiuskasBot */
		const chat_id = "-147649668"                                      /* Intranet de Katiuskas */

		http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendmessage?chat_id=%s&text=%s",
			bot_token, chat_id, url.QueryEscape(text)))
	}
	var ipaddr string
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx > -1 {
		ipaddr = r.RemoteAddr[:idx]
	}
	if tmp := r.Header["X-Cespedes-Remote-Addr"]; len(tmp) > 0 {
		ipaddr = tmp[0]
	}
        id, _ := context.GetOk(r, "id")
	uid, _ := id.(int)

	db.Exec("INSERT INTO log (severity, ipaddr, uid, text) VALUES ($1,$2,$3,$4)", severity, ipaddr, uid, text)
}
