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

func authGoogle(ctx *Context) {
	const client_id = "739018663335-rcrta00jqv86lonvl9hhgn7afvjhp4ic.apps.googleusercontent.com"
	const client_secret = "uCP5xO1nz6msnQ7cWFrhUX02"
	const redirect_uri = "https://intranet.katiuskas.es/auth/google"
	const authorization_endpoint = "https://accounts.google.com/o/oauth2/v2/auth"
	const token_endpoint = "https://www.googleapis.com/oauth2/v4/token"
	code := ctx.r.URL.Query().Get("code")
	if len(code)==0 {
		err := ctx.r.URL.Query().Get("error")
		if len(err) != 0 {
			fmt.Fprintf(ctx.w, "Google returned the error: %s\n", err)
			return
		}
		v := url.Values{}
		v.Set("client_id", client_id)
		v.Add("response_type", "code")
		v.Add("scope", "openid profile email")
		v.Add("redirect_uri", redirect_uri)
		http.Redirect(ctx.w, ctx.r, authorization_endpoint + "?" + v.Encode(), http.StatusFound)
//		fmt.Fprintln(w, "I would redirect to", authorization_endpoint + "?" + v.Encode())
		return
	}
	resp, err := http.PostForm(token_endpoint,
			url.Values{
				"code": {code},
				"client_id": {client_id},
				"client_secret": {client_secret},
				"redirect_uri": {redirect_uri},
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
		fmt.Fprintln(ctx.w, "Google id_token is not a string")
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

//	ctx.session.Values["name"], ok = things["name"].(string)
//	ctx.session.Values["picture"], ok = things["picture"].(string)
	ctx.session.Values["auth"] = "google"
	ctx.session.Values["email"] = email
	Log(ctx, LOG_INFO, fmt.Sprintf("Usuario autenticado en la Intranet (via Google): %s", email))
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
	id, person_type, board := db_mail_2_id(email)
	ctx.session.Values["id"] = id
	ctx.session.Values["type"] = person_type
	ctx.session.Values["board"] = board
	ctx.session.Values["roles"] = db_get_roles(id)
	ctx.Save()
	http.Redirect(ctx.w, ctx.r, "/", http.StatusFound)
}

func authFacebook(ctx *Context) {
	const client_id = "1692390947679031"
	const client_secret = "e06952f7f1208c7fd4d6d93d145be3e5"
	const redirect_uri = "https://intranet.katiuskas.es/auth/facebook"
	const authorization_endpoint = "https://www.facebook.com/dialog/oauth"
	const token_endpoint = "https://graph.facebook.com/v2.3/oauth/access_token"
	code := ctx.r.URL.Query().Get("code")
	if len(code)==0 {
		err := ctx.r.URL.Query().Get("error")
		if len(err) != 0 {
			fmt.Fprintf(ctx.w, "Facebook returned the error: %s\n", err)
			return
		}
		v := url.Values{}
		v.Set("client_id", client_id)
		v.Add("response_type", "code")
		v.Add("scope", "email")
		v.Add("redirect_uri", redirect_uri)
		http.Redirect(ctx.w, ctx.r, authorization_endpoint + "?" + v.Encode(), http.StatusFound)
//		fmt.Fprintln(ctx.w, "I would redirect to", authorization_endpoint + "?" + v.Encode())
		return
	}
	resp, err := http.PostForm(token_endpoint,
			url.Values{
				"client_id": {client_id},
				"redirect_uri": {redirect_uri},
				"client_secret": {client_secret},
				"code": {code},
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
	access_token, ok := things["access_token"].(string)
	if !ok {
		fmt.Fprintln(ctx.w, "Facebook access_token is not a string")
		return
	}
	resp, err = http.Get("https://graph.facebook.com/me?fields=name,email&access_token=" + access_token)

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
	err = json.Unmarshal(contents, &things)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	email, ok := things["email"].(string)

//	ctx.session.Values["name"], ok = things["name"].(string)
	ctx.session.Values["auth"] = "facebook"
	ctx.session.Values["email"] = email
	Log(ctx, LOG_INFO, fmt.Sprintf("Usuario autenticado en la Intranet (via Facebook): %s", email))
	id, person_type, board := db_mail_2_id(email)
	ctx.session.Values["id"] = id
	ctx.session.Values["type"] = person_type
	ctx.session.Values["board"] = board
	ctx.session.Values["roles"] = db_get_roles(id)
	ctx.Save()
	http.Redirect(ctx.w, ctx.r, "/", http.StatusFound)
}

const auth_hash_secret = "ruucaish2yiesaep6ailotae7sooto5U"

func auth_get_hash(id int, timeout int64) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%d-%d-%s", id, timeout, auth_hash_secret)))
	return fmt.Sprintf("%.16x", h.Sum(nil))
}

func authMail(ctx *Context) {
	ctx.r.ParseForm()
	email := ctx.r.Form.Get("email")
	phone := ctx.r.Form.Get("phone")
	if email=="" || phone=="" {
		http.Redirect(ctx.w, ctx.r, "/", http.StatusFound)
		return
	}

	p := make(map[string]interface{})
	p["email"] = email
	p["phone"] = phone

	var id int
	var name, surname string
	err := db.QueryRow("SELECT a.id_person AS id, c.name, c.surname FROM person_email a INNER JOIN person_phone b ON a.id_person=b.id_person LEFT JOIN person c ON a.id_person=c.id WHERE email=$1 AND phone=$2", email, phone).Scan(&id, &name, &surname)
	if err != nil {
		Log(ctx, LOG_INFO, fmt.Sprintf("auth: email+phone no válidos: %s / %s", email, phone))
		renderTemplate(ctx, "auth-wrongdata", p)
	} else {
		Log(ctx, LOG_INFO, fmt.Sprintf("auth: enviando enlace a %s %s <%s>", name, surname, email))
		auth := smtp.PlainAuth("", "intranet@katiuskas.es", "ahch0Vieg", "ssl0.ovh.net")
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
		go smtp.SendMail("ssl0.ovh.net:587", auth, "intranet@katiuskas.es", to, msg)
		p["name"] = name
		renderTemplate(ctx, "auth-sendmail", p)
	}
}

func authHash(ctx *Context) {
	ctx.r.ParseForm()
	code := ctx.r.Form.Get("code")
	s := strings.Split(code, "-")
	if len(s) != 3 {
		Log(ctx, LOG_INFO, fmt.Sprintf("auth: código erróneo: %s", code))
		renderTemplate(ctx, "auth-wronghash", make(map[string]interface{}))
		return
	}
	var id int
	var timeout int64
	fmt.Sscan(s[0], &id)
	fmt.Sscan(s[1], &timeout)
	hash := s[2]
	if hash != auth_get_hash(id, timeout) {
		Log(ctx, LOG_INFO, fmt.Sprintf("auth: código erróneo: %s", code))
		renderTemplate(ctx, "auth-wronghash", make(map[string]interface{}))
		return
	}
	if time.Now().Unix() > timeout {
		Log(ctx, LOG_INFO, fmt.Sprintf("auth: código caducado: %s", code))
		renderTemplate(ctx, "auth-timeout", make(map[string]interface{}))
		return
	}
	Log(ctx, LOG_INFO, fmt.Sprintf("Usuario autenticado en la Intranet (via hash): %d", id))
	person_type, board := db_id_2_type(id)
	ctx.session.Values["auth"] = "hash"
	ctx.session.Values["id"] = id
	ctx.session.Values["type"] = person_type
	ctx.session.Values["board"] = board
	ctx.session.Values["roles"] = db_get_roles(id)
	ctx.Save()
	http.Redirect(ctx.w, ctx.r, "/", http.StatusFound)
}
