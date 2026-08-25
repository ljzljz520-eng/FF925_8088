package auth

import (
	"testing"
	"time"
)

func TestTokenValid(t *testing.T) {
	if !Issue("u", "s", time.Hour).Valid(time.Now()) {
		t.Fatal()
	}
}
