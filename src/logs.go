package main

import (
	"fmt"
	"runtime"
	"net/http"
	"net/url"
)

func log(text string) {
	const bot_token = "204701695:AAFkgoxJPCUWpXTWDQco33I97y5BJIHmOKU" /* KatiuskasBot */
	const chat_id = "-147649668"                                      /* Intranet de Katiuskas */

	http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendmessage?chat_id=%s&text=%s",
		bot_token, chat_id, url.QueryEscape(text)))
}

func log_error(text string) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		text = fmt.Sprintf("(file=%v line=%v) %s", file, line, text)
	}

	log(text)
}
