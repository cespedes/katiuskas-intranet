package main

import (
	"crypto/rand"
	"errors"
	"net/http"
	"strings"
)

const (
	// sessionRenewWhen = 6 * 24 * time.Hour
	// sessionRenewTime = 7 * 24 * time.Hour
	tokenLength = 32
)

// getRandomToken returns a random string of length tokenLength, composed by letters and numbers
func getRandomToken() string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789" // 62 chars
	result := ""

	b := make([]byte, 128)
	for {
		_, err := rand.Read(b)
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(b); i++ {
			d := b[i] % 64
			if d < 62 {
				result = result + string(chars[d])
				if len(result) == tokenLength {
					return result
				}
			}
		}
	}
}

func getTokenFromHeaders(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		// No authorization header:
		return "", nil
	}
	token = strings.TrimSpace(token)
	words := strings.Fields(strings.TrimSpace(token))
	if len(words) != 2 || words[0] != "Bearer" {
		return "", errors.New("syntax error in `Authorization` header")
	}
	token = words[1]
	return token, nil
}
