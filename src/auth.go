package main

import (
        "os"
        "fmt"
        "time"
        "strings"
        "io/ioutil"
        "net/url"
        "net/http"
        "net/smtp"
        "encoding/json"
        "crypto/sha256"
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

func authGoogle(w http.ResponseWriter, r *http.Request) {
	const authorization_endpoint = "https://accounts.google.com/o/oauth2/v2/auth"
	const token_endpoint = "https://www.googleapis.com/oauth2/v4/token"
	code := r.URL.Query().Get("code")
	if len(code)==0 {
		err := r.URL.Query().Get("error")
		if len(err) != 0 {
			fmt.Fprintf(w, "Google returned the error: %s\n", err)
			return
		}
		v := url.Values{}
		v.Set("client_id", Google_auth_client_id)
		v.Add("response_type", "code")
		v.Add("scope", "openid profile email")
		v.Add("redirect_uri", Google_auth_redirect_uri)
		http.Redirect(w, r, authorization_endpoint + "?" + v.Encode(), http.StatusFound)
//		fmt.Fprintln(w, "I would redirect to", authorization_endpoint + "?" + v.Encode())
		return
	}
	resp, err := http.PostForm(token_endpoint,
			url.Values{
				"code": {code},
				"client_id": {Google_auth_client_id},
				"client_secret": {Google_auth_client_secret},
				"redirect_uri": {Google_auth_redirect_uri},
				"grant_type": {"authorization_code"},
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
	id_token, ok := things["id_token"].(string)
	if !ok {
		fmt.Fprintln(w, "Google id_token is not a string")
		return
	}
//		fmt.Fprintln(w, "id_token:", id_token)
//		fmt.Fprintln(w, "response = " + string(contents))
	resp, err = http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + id_token)
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
	email, ok := things["email"].(string)

//	Ctx(r).session.Values["name"], ok = things["name"].(string)
//	Ctx(r).session.Values["picture"], ok = things["picture"].(string)
	Log(r, LOG_INFO, fmt.Sprintf("Usuario autenticado en la Intranet (via Google): %s", email))
//		fmt.Fprintln(w, "response2 = " + string(contents))
/* Sample response:
response2 = {
 "iss": "https://accounts.google.com",
 "at_hash": "xom_Ml4HXUuQswIwAkw32w",
 "aud": "434725510955-2lhlvbbdum01g8akgigk2v7123rpadid.apps.googleusercontent.com",
 "sub": "114151858104579138691",
 "email_verified": "true",
 "azp": "434725510955-2lhlvbbdum01g8akgigk2v7123rpadid.apps.googleusercontent.com",
 "email": "espeleo.katiuskas@gmail.com",
 "iat": "1460630724",
 "exp": "1460634324",
 "name": "Club D.E. de Espeleología Katiuskas",
 "given_name": "Club D.E. de Espeleología",
 "family_name": "Katiuskas",
 "alg": "RS256",
 "kid": "08ff58ef6a5f48d96fe609726351ba6df277e79b"
}
*/
	id, person_type := db_mail_2_id(email)
	if person_type==NoUser {
		fmt.Fprintln(w, "ERR: NoUser (?)")
	} else if person_type==NoSocio {
		p := make(map[string]interface{})
		p["email"] = email
		renderTemplate(w, r, "auth-wrongdata", p)
	} else {
		Ctx(r).session.Values["auth"] = "google"
		Ctx(r).session.Values["id"] = id
		Ctx(r).session.Values["type"] = person_type
		Ctx(r).session.Values["roles"] = db_get_roles(id)
		Ctx(r).Save(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func auth_get_hash(id int, timeout int64) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%d-%d-%s", id, timeout, Secret_auth_hash)))
	return fmt.Sprintf("%.16x", h.Sum(nil))
}

func authMail(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form.Get("email")
	phone := r.Form.Get("phone")
	if email=="" || phone=="" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	p := make(map[string]interface{})
	p["email"] = email
	p["phone"] = phone

	var id int
	var name, surname string
	err := db.QueryRow("SELECT a.id_person AS id, c.name, c.surname FROM person_email a INNER JOIN person_phone b ON a.id_person=b.id_person LEFT JOIN person c ON a.id_person=c.id WHERE email=$1 AND phone=$2", email, phone).Scan(&id, &name, &surname)
	if err != nil {
		Log(r, LOG_INFO, fmt.Sprintf("auth: email+phone no válidos: %s / %s", email, phone))
		renderTemplate(w, r, "auth-wrongdata", p)
	} else {
		Log(r, LOG_INFO, fmt.Sprintf("auth: enviando enlace a %s %s <%s>", name, surname, email))
		auth := smtp.PlainAuth("", SMTP_From, SMTP_Pass, SMTP_Server)
		to := []string{email}
		timeout := time.Now().Unix() + 2*60*60
		key := fmt.Sprintf("%d-%d-%s", id, timeout, auth_get_hash(id, timeout))
		msg := []byte(
			"From: Intranet de Katiuskas <intranet@katiuskas.es>\r\n" +
			"To: " + name + " " + surname + " <" + email + ">\r\n" +
			"Subject: Acceso a la Intranet de Katiuskas\r\n" +
			"\r\n" +
			"Hola, " + name + ".\r\n" +
			"\r\n" +
			"Para poder acceder a la Intranet de Katiuskas debes hacer clic en el siguiente enlace:\r\n" +
			"\r\n" +
			"https://intranet.katiuskas.es/auth/hash?code=" + key + "\r\n" +
			"\r\n" +
			"Un saludo,\r\n" +
			"\r\n" +
			"La Intranet de Katiuskas.\r\n")
		go smtp.SendMail(fmt.Sprintf("%s:%d", SMTP_Server, SMTP_Port), auth, SMTP_From, to, msg)
		p["name"] = name
		renderTemplate(w, r, "auth-sendmail", p)
	}
}

func authHash(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.Form.Get("code")
	s := strings.Split(code, "-")
	if len(s) != 3 {
		Log(r, LOG_INFO, fmt.Sprintf("auth: código erróneo: %s", code))
		renderTemplate(w, r, "auth-wronghash", make(map[string]interface{}))
		return
	}
	var id int
	var timeout int64
	fmt.Sscan(s[0], &id)
	fmt.Sscan(s[1], &timeout)
	hash := s[2]
	if hash != auth_get_hash(id, timeout) {
		Log(r, LOG_INFO, fmt.Sprintf("auth: código erróneo: %s", code))
		renderTemplate(w, r, "auth-wronghash", make(map[string]interface{}))
		return
	}
	if time.Now().Unix() > timeout {
		Log(r, LOG_INFO, fmt.Sprintf("auth: código caducado: %s", code))
		renderTemplate(w, r, "auth-timeout", make(map[string]interface{}))
		return
	}
	person_type := db_id_2_type(id)
	if person_type==NoUser || person_type==NoSocio {
		Log(r, LOG_ERR, fmt.Sprintf("Error identifying person_id %d from hash", id))
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	Log(r, LOG_INFO, fmt.Sprintf("Usuario autenticado en la Intranet (via hash): %d", id))
	Ctx(r).session.Values["auth"] = "hash"
	Ctx(r).session.Values["id"] = id
	Ctx(r).session.Values["type"] = person_type
	Ctx(r).session.Values["roles"] = db_get_roles(id)
	Ctx(r).Save(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}
