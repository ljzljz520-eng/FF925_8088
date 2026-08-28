package auth

import (
	"errors"
	"strings"
	"time"
)

type Claims struct {
	Subject, Role, CardID string
	IssuedAt, ExpiresAt   time.Time
}

func NewClaims(subject, role, card string, ttl time.Duration) Claims {
	n := time.Now().UTC()
	return Claims{Subject: subject, Role: role, CardID: card, IssuedAt: n, ExpiresAt: n.Add(ttl)}
}
func (c Claims) Expired(now time.Time) bool { return !now.Before(c.ExpiresAt) }
func (c Claims) Allows(role string) bool    { return c.Role == role || c.Role == "admin" }
func (c Claims) Validate(now time.Time) error {
	if strings.TrimSpace(c.Subject) == "" {
		return errors.New("subject missing")
	}
	if c.Expired(now) {
		return errors.New("claims expired")
	}
	return nil
}
