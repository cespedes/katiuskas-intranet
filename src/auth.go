package main

import (
        "fmt"
        "os"
        "io/ioutil"
        "net/http"
        "net/url"
        "encoding/json"
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
	const client_id = "739018663335-rcrta00jqv86lonvl9hhgn7afvjhp4ic.apps.googleusercontent.com"
	const client_secret = "uCP5xO1nz6msnQ7cWFrhUX02"
	const redirect_uri = "https://intranet.katiuskas.es/auth/google"
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
		v.Set("client_id", client_id)
		v.Add("response_type", "code")
		v.Add("scope", "openid profile email")
		v.Add("redirect_uri", redirect_uri)
		http.Redirect(w, r, authorization_endpoint + "?" + v.Encode(), http.StatusFound)
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

	session := session_get(w, r)
//	session.Values["name"], ok = things["name"].(string)
//	session.Values["picture"], ok = things["picture"].(string)
	session["auth"] = "google"
	session["email"] = email
	log(fmt.Sprintf("Usuario autenticado en la Intranet (via Google): %s", email))
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
	session["id"] = id
	session["type"] = person_type
	session_save(w, r)
	if err != nil {
		fmt.Println("auth: session.Save:", err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func authFacebook(w http.ResponseWriter, r *http.Request) {
	const client_id = "1692390947679031"
	const client_secret = "e06952f7f1208c7fd4d6d93d145be3e5"
	const redirect_uri = "https://intranet.katiuskas.es/auth/facebook"
	const authorization_endpoint = "https://www.facebook.com/dialog/oauth"
	const token_endpoint = "https://graph.facebook.com/v2.3/oauth/access_token"
	code := r.URL.Query().Get("code")
	if len(code)==0 {
		err := r.URL.Query().Get("error")
		if len(err) != 0 {
			fmt.Fprintf(w, "Facebook returned the error: %s\n", err)
			return
		}
		v := url.Values{}
		v.Set("client_id", client_id)
		v.Add("response_type", "code")
		v.Add("scope", "email")
		v.Add("redirect_uri", redirect_uri)
		http.Redirect(w, r, authorization_endpoint + "?" + v.Encode(), http.StatusFound)
//		fmt.Fprintln(w, "I would redirect to", authorization_endpoint + "?" + v.Encode())
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
		fmt.Fprintln(w, "Facebook access_token is not a string")
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

	session := session_get(w, r)
//	session.Values["name"], ok = things["name"].(string)
	session["auth"] = "facebook"
	session["email"] = email
	log(fmt.Sprintf("Usuario autenticado en la Intranet (via Facebook): %s", email))
	id, person_type := db_mail_2_id(email)
	session["id"] = id
	session["type"] = person_type
	session_save(w, r)
	if err != nil {
		fmt.Println("auth: session.Save:", err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
