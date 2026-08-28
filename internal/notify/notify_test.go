package notify

import (
	"tempcards/internal/share"
	"testing"
	"time"
)

func TestNotify(t *testing.T) {
	m := &MemorySender{}
	d := New(m)
	if e := d.SendCard(share.ShareCard{Filename: "f"}, "a@b"); e != nil {
		t.Fatal(e)
	}
	_ = time.Now()
}
