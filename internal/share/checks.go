package share

import (
	"errors"
	"time"
)

func CheckCard(c ShareCard) error {
	if c.ID == "" {
		return errors.New("id required")
	}
	if e := ValidateCode(c.Code); e != nil {
		return e
	}
	if e := SafeFilename(c.Filename); e != nil {
		return e
	}
	if c.Size <= 0 {
		return errors.New("size required")
	}
	return nil
}
func CheckLifecycle(c ShareCard, now time.Time) error {
	if c.Deleted {
		return ErrDeleted
	}
	if !now.Before(c.ExpiresAt) {
		return ErrExpired
	}
	return nil
}
func CheckDownload(c ShareCard, now time.Time) error {
	if e := CheckCard(c); e != nil {
		return e
	}
	return CheckLifecycle(c, now)
}
func CheckGrant(g AccessGrant, now time.Time) error {
	if g.ID == "" || g.CardID == "" {
		return errors.New("grant identifiers required")
	}
	if !Active(g, now) {
		return errors.New("grant inactive")
	}
	return nil
}
func CheckEvent(e ShareEvent) error {
	if e.ID == "" || e.CardID == "" || e.Kind == "" {
		return errors.New("event incomplete")
	}
	return nil
}
func CheckTTL(ttl time.Duration) error { return ValidateTTL(ttl) }
func CheckOwner(owner string) error {
	if owner == "" {
		return errors.New("owner required")
	}
	return nil
}
func CheckCode(code string) error     { return ValidateCode(code) }
func CheckFilename(name string) error { return SafeFilename(name) }
func CheckPositive(n int64) error {
	if n <= 0 {
		return errors.New("value must be positive")
	}
	return nil
}
