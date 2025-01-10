package katintranet

import (
	"crypto/rand"
	"errors"
	"net/http"
	"strings"
	"time"
)

const (
	tokenLength      = 32
	tokenCookieName  = "TOKEN"
	sessionRenewWhen = 6 * 24 * time.Hour
	sessionRenewTime = 7 * 24 * time.Hour
)

type Token struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Session   bool      `json:"session"    db:"session"`
	Token     string    `json:"token"      db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	// ReadOnly  bool      `json:"read_only"  db:"read_only"`
}

// NewToken creates a random token
func NewToken() *Token {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789" // 62 chars

	var token Token

	b := make([]byte, 128)
	for {
		_, err := rand.Read(b)
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(b); i++ {
			d := b[i] % 64
			if d < 62 {
				token.Token = token.Token + string(chars[d])
				if len(token.Token) == tokenLength {
					break
				}
			}
		}
		if len(token.Token) == tokenLength {
			break
		}
	}

	token.CreatedAt = time.Now()
	token.Session = true
	token.ExpiresAt = token.CreatedAt.Add(sessionRenewTime)

	return &token
}

// AddToken stores a token in the database for user personID
func AddToken(token *Token, db *DB, personID int) error {
	var id int
	var err error

	err = db.Get(&id, `
                INSERT INTO token(person_id, session, token, expires_at)
		VALUES ($1, $2, $3, $4)
                RETURNING id
	`, personID, token.Session, token.Token, token.ExpiresAt)
	if err != nil {
		return err
	}

	return nil
}

// GetTokenFromHeaders returns the token received from the headers, if any.
// A token can be set:
// - Using an "Authorization: Bearer" header
// - Using a "TOKEN" cookie
func GetTokenFromHeaders(r *http.Request) (string, error) {
	var token string
	token = r.Header.Get("Authorization")

	if token != "" {
		token = strings.TrimSpace(token)
		words := strings.Fields(strings.TrimSpace(token))
		if len(words) != 2 || words[0] != "Bearer" {
			return "", errors.New("syntax error in `Authorization` header")
		}
		token = words[1]
	}

	if token == "" {
		cookie, err := r.Cookie(tokenCookieName)
		if err != nil {
			token = cookie.Value
		}
	}
	if token == "" {
		// No authorization header:
		return "", nil
	}
	return token, nil
}
