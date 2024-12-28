package katintranet

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// outJSON writes the JSON-encoded object v to the http.ResponseWriter w.
func (s *server) outJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	err := e.Encode(v)
	if err != nil {
		httpError(w, err)
		return
	}
}

// httpError sends a HTTP error as a response
func httpError(w http.ResponseWriter, err any, codes ...int) {
	code := http.StatusInternalServerError

	if err == sql.ErrNoRows {
		code = http.StatusNotFound
	}

	if er, ok := err.(error); ok {
		var e ErrHttpStatus
		if errors.As(er, &e) {
			code = e.Status
		}
	}

	if len(codes) > 0 {
		code = codes[0]
	}

	httpMessage(w, code, "error", fmt.Sprint(err))
}

func httpMessage(w http.ResponseWriter, code int, label string, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{%q: %q}\n", label, msg)
}

type ErrHttpStatus struct {
	Status int
	Err    error
}

func (e ErrHttpStatus) Error() string {
	return e.Err.Error()
}

func (e ErrHttpStatus) Unwrap() error {
	return e.Err
}

func httpInfo(w http.ResponseWriter, msg any, codes ...int) {
	code := http.StatusOK

	if len(codes) > 0 {
		code = codes[0]
	}

	httpMessage(w, code, "info", fmt.Sprint(msg))
}
