package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

type Token struct {
	Value, Subject  string
	Issued, Expires time.Time
	Revoked         bool
}

func Issue(subject, secret string, ttl time.Duration) Token {
	h := sha256.Sum256([]byte(subject + secret))
	n := time.Now().UTC()
	return Token{Value: hex.EncodeToString(h[:]), Subject: subject, Issued: n, Expires: n.Add(ttl)}
}
func (t Token) Valid(now time.Time) bool { return t.Value != "" && !t.Revoked && now.Before(t.Expires) }
func Parse(raw string) (string, error) {
	p := strings.Split(raw, ".")
	if len(p) != 2 || p[0] == "" || p[1] == "" {
		return "", errors.New("malformed token")
	}
	return p[0], nil
}
func Revoke(t *Token) error {
	if t == nil {
		return errors.New("nil token")
	}
	t.Revoked = true
	return nil
}
