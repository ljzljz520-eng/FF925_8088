package web

import (
	"net/http/httptest"
	"tempcards/internal/share"
	"tempcards/internal/storage"
	"testing"
)

func TestHealth(t *testing.T) {
	s, _ := storage.Open(t.TempDir() + "/x.db")
	defer s.Close()
	r := httptest.NewRecorder()
	httptest.NewRequest("GET", "/health", nil)
	New(share.NewService(s)).health(r, nil)
	if r.Code != 204 {
		t.Fatal(r.Code)
	}
}
