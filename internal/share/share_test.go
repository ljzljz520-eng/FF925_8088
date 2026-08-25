package share

import (
	"tempcards/internal/storage"
	"testing"
	"time"
)

func TestDeletedCodeRejected(t *testing.T) {
	s, _ := storage.Open(t.TempDir() + "/x.db")
	defer s.Close()
	svc := NewService(s)
	c, _ := svc.Create("a", "f.txt", "ABC123", 2, time.Hour)
	_ = svc.Delete(c.ID)
	if _, e := svc.Download(c.Code, "ip"); e != ErrDeleted {
		t.Fatalf("expected deleted error: %v", e)
	}
}
func TestWorkflowCreateShare(t *testing.T) {
	s, _ := storage.Open(t.TempDir() + "/x.db")
	defer s.Close()
	c, e := NewService(s).Create("a", "f", "ABC123", 2, time.Hour)
	if e != nil || c.ID == "" {
		t.Fatal(e)
	}
}
func TestWorkflowDownloadShare(t *testing.T) {
	s, _ := storage.Open(t.TempDir() + "/x.db")
	defer s.Close()
	svc := NewService(s)
	c, _ := svc.Create("a", "f", "ABC123", 2, time.Hour)
	if _, e := svc.Download(c.Code, "ip"); e != nil {
		t.Fatal(e)
	}
}
func TestWorkflowDeleteShare(t *testing.T) {
	s, _ := storage.Open(t.TempDir() + "/x.db")
	defer s.Close()
	svc := NewService(s)
	c, _ := svc.Create("a", "f", "ABC123", 2, time.Hour)
	if e := svc.Delete(c.ID); e != nil {
		t.Fatal(e)
	}
}
func TestWorkflowGrantAccess(t *testing.T) {
	s, _ := storage.Open(t.TempDir() + "/x.db")
	defer s.Close()
	svc := NewService(s)
	c, _ := svc.Create("a", "f", "ABC123", 2, time.Hour)
	if _, e := svc.Grant(c.ID, "x@example.com", "read", time.Hour); e != nil {
		t.Fatal(e)
	}
}
