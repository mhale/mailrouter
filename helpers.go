package main

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// Utility function for parsing URLs like /route/:id/:action
func ParsePath(path string) (base string, id string, action string) {
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		base = parts[0]
	}
	if len(parts) > 1 {
		id = parts[1]
	}
	if len(parts) > 2 {
		action = parts[2]
	}
	return base, id, action
}

// Wrapper function for setting flash messages via cookies.
func SetCookie(w http.ResponseWriter, name string, value string) {
	valueEnc := base64.StdEncoding.EncodeToString([]byte(value))
	http.SetCookie(w, &http.Cookie{Name: name, Value: valueEnc, Path: "/"})
}

// Wrapper function for getting flash messages via cookies.
func GetCookie(w http.ResponseWriter, req *http.Request, name string) string {
	value := ""
	cookie, _ := req.Cookie(name)
	if cookie != nil {
		valueDec, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err == nil {
			value = string(valueDec)
		}
		// The cookie has been read, so clear it by setting MaxAge < 0
		http.SetCookie(w, &http.Cookie{Name: name, MaxAge: -1, Path: "/"})
	}
	return value
}
