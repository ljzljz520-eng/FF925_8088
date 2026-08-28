package share

import (
	"tempcards/internal/storage"
	"testing"
	"time"
)

func TestServiceAuthorize(t *testing.T) {
	s, _ := storage.Open(t.TempDir() + "/x.db")
	defer s.Close()
	svc := NewService(s)
	c, _ := svc.Create("o", "f", "ABC123", 1, time.Hour)
	if _, e := svc.Authorize(c.Code); e != nil {
		t.Fatal(e)
	}
}
