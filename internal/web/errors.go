package web

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct{ Code, Message string }

func WriteError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorBody{Code: code, Message: message})
}
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func MethodAllowed(w http.ResponseWriter, allowed string) {
	w.Header().Set("Allow", allowed)
	WriteError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
}
func RequestID(r *http.Request) string {
	v := r.Header.Get("X-Request-ID")
	if v == "" {
		return "anonymous"
	}
	return v
}
func ClientIP(r *http.Request) string {
	if v := r.Header.Get("X-Forwarded-For"); v != "" {
		return v
	}
	return r.RemoteAddr
}
