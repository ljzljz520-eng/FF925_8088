package web

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func Query(r *http.Request, key string) string { return strings.TrimSpace(r.URL.Query().Get(key)) }
func IntQuery(r *http.Request, key string, def int) (int, error) {
	v := Query(r, key)
	if v == "" {
		return def, nil
	}
	n, e := strconv.Atoi(v)
	return n, e
}
func BoolQuery(r *http.Request, key string, def bool) (bool, error) {
	v := Query(r, key)
	if v == "" {
		return def, nil
	}
	switch strings.ToLower(v) {
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	default:
		return def, errors.New("invalid boolean")
	}
}
func RequireQuery(r *http.Request, key string) (string, error) {
	v := Query(r, key)
	if v == "" {
		return "", errors.New("missing query parameter")
	}
	return v, nil
}
func IsJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}
func IsMethod(r *http.Request, m string) bool   { return strings.EqualFold(r.Method, m) }
func Header(r *http.Request, key string) string { return strings.TrimSpace(r.Header.Get(key)) }
func AcceptsJSON(r *http.Request) bool {
	return Header(r, "Accept") == "" || strings.Contains(Header(r, "Accept"), "json")
}
func CacheControl(w http.ResponseWriter) { w.Header().Set("Cache-Control", "no-store") }
func NoSniff(w http.ResponseWriter)      { w.Header().Set("X-Content-Type-Options", "nosniff") }
