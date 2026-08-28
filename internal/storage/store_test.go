package storage

import (
	"testing"
	"time"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	p := t.TempDir() + "/x.db"
	s, e := Open(p)
	if e != nil {
		t.Fatal(e)
	}
	c := ShareCard{ID: "1", Code: "ABC123", Owner: "o", CreatedAt: time.Now(), ExpiresAt: time.Now().Add(time.Hour)}
	if e = s.SaveCard(c); e != nil {
		t.Fatal(e)
	}
	s.Close()
	s, e = Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	got, e := s.GetCard("1")
	if e != nil || got == nil {
		t.Fatal(e)
	}
}
