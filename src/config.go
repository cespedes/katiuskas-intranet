package main

import (
	"os"
	"log"
	"flag"
	"bufio"
	"regexp"
	"strings"
	"unicode"
)

var config = map[string]string{
	"http_host":                 "intranet.katiuskas.es",
	"http_listen_addr":          "localhost:8081",
	"cookie_secret":             "11UinL5BLSMVqivclTDo27qLVhIahkJM",
	"secret_db_conn":            "host=localhost user=katiuskas dbname=katiuskas password=Ohqu8Get",
	"telegram_bot_token":        "204701695:AAFkgoxJPCUWpXTWDQco33I97y5BJIHmOKU",
	"telegram_log_chat_id":      "-147649668",
	"telegram_webhook_path":     "/tgbot.aif7eoca",
	"auth_hash_secret":          "ruucaish2yiesaep6ailotae7sooto5U",
	"smtp_from":                 "intranet@katiuskas.es",
	"smtp_pass":                 "ahch0Vieg",
	"smtp_server":               "ssl0.ovh.net",
	"smtp_port":                 "587",
	"google_auth_client_id":     "739018663335-rcrta00jqv86lonvl9hhgn7afvjhp4ic.apps.googleusercontent.com",
	"google_auth_client_secret": "uCP5xO1nz6msnQ7cWFrhUX02",
	"google_auth_redirect_uri":  "https://intranet.katiuskas.es/auth/google",
}

var (
	regDoubleQuote = regexp.MustCompile("^([^= \t]+)[ \t]*=[ \t]*\"([^\"]*)\"$")
)

func init() {
	config_file := flag.String("c", "config.ini", "config file")
	flag.Parse()

	file, err := os.Open(*config_file)
	if err != nil {
		log.Printf("Reading configuration: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	lineno := 0

	for {
		l, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		lineno++
		line := strings.TrimFunc(string(l), unicode.IsSpace)
		if len(line)==0 || line[0]=='#' || line[0]==';' {
			continue
		}
		if m := regDoubleQuote.FindAllStringSubmatch(line, 1); m != nil {
			config[m[0][1]] = m[0][2];
		} else {
			log.Printf("Syntax error in %s:%d: unexpected \"%s\"", *config_file, lineno, line)
		}
	}
}
