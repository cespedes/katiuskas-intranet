package katintranet

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"strings"
	"time"
)

/*
 * Documentation: https://developers.google.com/identity/protocols/OpenIDConnect
 * Google Developers console: https://console.developers.google.com/
 *
 * Accessing https://accounts.google.com/.well-known/openid-configuration
 * yields:
 * "authorization_endpoint": "https://accounts.google.com/o/oauth2/v2/auth"
 * "token_endpoint": "https://www.googleapis.com/oauth2/v4/token"
 */

/* To get "authorization_endpoint" and "token_endpoint" we can use somthing like this:
   resp, err := http.Get("https://accounts.google.com/.well-known/openid-configuration")
   if err != nil {
           fmt.Printf("%s", err)
           os.Exit(1)
   }
   defer resp.Body.Close()
   contents, err := ioutil.ReadAll(resp.Body)
   if err != nil {
           fmt.Printf("%s", err)
           os.Exit(1)
   }
   var things = make(map[string]interface{})
   err = json.Unmarshal(contents, &things)
   if err != nil {
           fmt.Printf("%s", err)
           os.Exit(1)
   }
   fmt.Println("authorization_endpoint: ", things["authorization_endpoint"].(string))
   fmt.Println("token_endpoint: ", things["token_endpoint"].(string))
*/

func (s *server) authGoogle(w http.ResponseWriter, r *http.Request) {
	const authorization_endpoint = "https://accounts.google.com/o/oauth2/v2/auth"
	const token_endpoint = "https://www.googleapis.com/oauth2/v4/token"
	code := r.URL.Query().Get("code")
	if len(code) == 0 {
		err := r.URL.Query().Get("error")
		if len(err) != 0 {
			fmt.Fprintf(w, "Google returned the error: %s\n", err)
			return
		}
		v := url.Values{}
		v.Set("client_id", s.config["google_auth_client_id"])
		v.Add("response_type", "code")
		v.Add("scope", "openid profile email")
		v.Add("redirect_uri", s.config["google_auth_redirect_uri"])
		http.Redirect(w, r, authorization_endpoint+"?"+v.Encode(), http.StatusFound)
		//		fmt.Fprintln(w, "I would redirect to", authorization_endpoint + "?" + v.Encode())
		return
	}
	resp, err := http.PostForm(token_endpoint,
		url.Values{
			"code":          {code},
			"client_id":     {s.config["google_auth_client_id"]},
			"client_secret": {s.config["google_auth_client_secret"]},
			"redirect_uri":  {s.config["google_auth_redirect_uri"]},
			"grant_type":    {"authorization_code"},
		})
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	var things = make(map[string]interface{})
	err = json.Unmarshal(contents, &things)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	idToken, ok := things["id_token"].(string)
	if !ok {
		fmt.Fprintln(w, "Google id_token is not a string")
		return
	}
	//		fmt.Fprintln(w, "id_token:", id_token)
	//		fmt.Fprintln(w, "response = " + string(contents))
	resp, err = http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + idToken)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	//		fmt.Fprintln(w, "response2 = " + string(contents))
	err = json.Unmarshal(contents, &things)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	email, _ := things["email"].(string)

	//	Ctx(r).session.Values["name"], ok = things["name"].(string)
	//	Ctx(r).session.Values["picture"], ok = things["picture"].(string)
	s.Log(r, LOG_NOTICE, fmt.Sprintf("Usuario autenticado en la Intranet (via Google): %s", email))

	id, personType := s.DBmail2id(email)
	if personType == NoUser {
		fmt.Fprintln(w, "ERR: NoUser (?)")
	} else if personType == NoSocio {
		p := make(map[string]interface{})
		p["email"] = email
		renderTemplate(w, r, "auth-wrongdata", p)
	} else {
		sess, _ := _session_store.Get(r, "session")
		sess.Values["auth"] = "google"
		sess.Values["id"] = id
		sess.Values["type"] = personType
		sess.Values["roles"] = s.DBgetRoles(id)
		sess.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func authGetHash(secret string, id int, timeout int64) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%d-%d-%s", id, timeout, secret)))
	return fmt.Sprintf("%.16x", h.Sum(nil))
}

func (s *server) authMail(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form.Get("email")
	phone := r.Form.Get("phone")
	if email == "" || phone == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	p := make(map[string]interface{})
	p["email"] = email
	p["phone"] = phone

	var id int
	var name, surname string
	err := s.db.QueryRow("SELECT a.person_id AS id, c.name, c.surname FROM person_email a INNER JOIN person_phone b ON a.person_id=b.person_id LEFT JOIN person c ON a.person_id=c.id WHERE email=$1 AND phone=$2", email, phone).Scan(&id, &name, &surname)
	if err != nil {
		s.Log(r, LOG_INFO, fmt.Sprintf("auth: email+phone no válidos: %s / %s", email, phone))
		renderTemplate(w, r, "auth-wrongdata", p)
	} else {
		s.Log(r, LOG_INFO, fmt.Sprintf("auth: enviando enlace a %s %s <%s>", name, surname, email))
		auth := smtp.PlainAuth("", s.config["smtp_from"], s.config["smtp_pass"], s.config["smtp_server"])
		to := []string{email}
		timeout := time.Now().Unix() + 2*60*60
		key := fmt.Sprintf("%d-%d-%s", id, timeout, authGetHash(s.config["auth_hash_secret"], id, timeout))
		msg := []byte(
			"From: Intranet de Katiuskas <" + s.config["smtp_from"] + ">\r\n" +
				"To: " + name + " " + surname + " <" + email + ">\r\n" +
				"Subject: Acceso a la Intranet de Katiuskas\r\n" +
				"\r\n" +
				"Hola, " + name + ".\r\n" +
				"\r\n" +
				"Para poder acceder a la Intranet de Katiuskas debes hacer clic en el siguiente enlace:\r\n" +
				"\r\n" +
				"https://" + s.config["http_host"] + "/auth/hash?code=" + key + "\r\n" +
				"\r\n" +
				"Un saludo,\r\n" +
				"\r\n" +
				"La Intranet de Katiuskas.\r\n")
		go smtp.SendMail(fmt.Sprintf("%s:%s", s.config["smtp_server"], s.config["smtp_port"]), auth, s.config["smtp_from"], to, msg)
		p["name"] = name
		renderTemplate(w, r, "auth-sendmail", p)
	}
}

func (s *server) authHash(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.Form.Get("code")
	str := strings.Split(code, "-")
	if len(str) != 3 {
		s.Log(r, LOG_INFO, fmt.Sprintf("auth: código erróneo: %s", code))
		renderTemplate(w, r, "auth-wronghash", make(map[string]interface{}))
		return
	}
	var id int
	var timeout int64
	fmt.Sscan(str[0], &id)
	fmt.Sscan(str[1], &timeout)
	hash := str[2]
	if hash != authGetHash(s.config["auth_hash_secret"], id, timeout) {
		s.Log(r, LOG_INFO, fmt.Sprintf("auth: código erróneo: %s", code))
		renderTemplate(w, r, "auth-wronghash", make(map[string]interface{}))
		return
	}
	if time.Now().Unix() > timeout {
		s.Log(r, LOG_INFO, fmt.Sprintf("auth: código caducado: %s", code))
		renderTemplate(w, r, "auth-timeout", make(map[string]interface{}))
		return
	}
	personType := s.DBidToType(id)
	if personType == NoUser || personType == NoSocio {
		s.Log(r, LOG_ERR, fmt.Sprintf("Error identifying person_id %d from hash", id))
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	s.Log(r, LOG_NOTICE, fmt.Sprintf("Usuario autenticado en la Intranet (via hash): %d", id))
	sess, _ := _session_store.Get(r, "session")
	sess.Values["auth"] = "hash"
	sess.Values["id"] = id
	sess.Values["type"] = personType
	sess.Values["roles"] = s.DBgetRoles(id)
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
