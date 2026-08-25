package policy

import (
	"testing"
	"time"
)

func TestPolicy(t *testing.T) {
	r := Default()
	if e := r.Validate("f", 2); e != nil {
		t.Fatal(e)
	}
	if r.Expiry(time.Now()).Before(time.Now()) {
		t.Fatal()
	}
}
