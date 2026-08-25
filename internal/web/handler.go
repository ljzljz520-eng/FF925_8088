package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"tempcards/internal/share"
)

type API struct{ Service *share.Service }

func New(s *share.Service) *API { return &API{Service: s} }
func (a *API) Routes() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/health", a.health)
	m.HandleFunc("/v1/download", a.download)
	m.HandleFunc("/v1/cards", a.cards)
	return m
}
func (a *API) health(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }
func (a *API) cards(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var in struct {
		Owner, Filename, Code string
		Size                  int64
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		http.Error(w, "bad json", 400)
		return
	}
	c, e := a.Service.Create(in.Owner, in.Filename, in.Code, in.Size, 24*60*60*1e9)
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(c)
}
func (a *API) download(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	d, e := a.Service.Download(code, r.RemoteAddr)
	if e != nil {
		switch {
		case errors.Is(e, share.ErrInvalidCode):
			WriteError(w, http.StatusNotFound, "invalid_code", "卡密无效")
		case errors.Is(e, share.ErrDeleted):
			WriteError(w, http.StatusNotFound, "card_deleted", "该卡密已被删除")
		case errors.Is(e, share.ErrExpired):
			WriteError(w, http.StatusGone, "card_expired", "该卡密已过期")
		default:
			WriteError(w, http.StatusNotFound, "download_failed", e.Error())
		}
		return
	}
	WriteJSON(w, http.StatusOK, d)
}
