package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"
)

func (s *server) ConfigInit() error {
	var defaultConfig = map[string]string{
		"http_host":                 "intranet.katiuskas.es",
		"http_listen_addr":          "localhost:8000",
		"cookie_secret":             "some random value",
		"secret_db_conn":            "host=localhost user=katiuskas dbname=katiuskas password=secretpwd",
		"telegram_bot_token":        "token as given by @BotFather on Telegram",
		"telegram_log_chat_id":      "chat ID in Telegram",
		"telegram_webhook_path":     "/some-random-path",
		"auth_hash_secret":          "some random value",
		"smtp_from":                 "intranet@katiuskas.es",
		"smtp_pass":                 "password by your ISP",
		"smtp_server":               "mail.your-isp.com",
		"smtp_port":                 "587",
		"google_auth_client_id":     "Client ID by console.developers.google.com",
		"google_auth_client_secret": "Client secret by console.developers.google.com",
		"google_auth_redirect_uri":  "https://intranet.katiuskas.es/auth/google",
	}

	configFile := flag.String("c", "config.ini", "config file")
	flag.Parse()

	file, err := os.Open(*configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	lineno := 0

	regDoubleQuote := regexp.MustCompile("^([^= \t]+)[ \t]*=[ \t]*\"([^\"]*)\"$")
	for {
		l, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		lineno++
		line := strings.TrimSpace(string(l))
		if len(line) == 0 || line[0] == '#' || line[0] == ';' {
			continue
		}
		if m := regDoubleQuote.FindAllStringSubmatch(line, 1); m != nil {
			defaultConfig[m[0][1]] = m[0][2]
		} else {
			log.Printf("Syntax error in %s:%d: unexpected \"%s\"", *configFile, lineno, line)
		}
	}

	s.config = defaultConfig

	return nil
}
