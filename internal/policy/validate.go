package policy

import (
	"errors"
	"strings"
	"time"
)

func ValidateOwner(owner string) error {
	if len(strings.TrimSpace(owner)) < 2 {
		return errors.New("owner too short")
	}
	return nil
}
func ValidateCodeLength(code string, min, max int) error {
	n := len(strings.TrimSpace(code))
	if n < min || n > max {
		return errors.New("code length invalid")
	}
	return nil
}
func ValidateExpiry(created, expiry time.Time) error {
	if !expiry.After(created) {
		return errors.New("expiry must be future")
	}
	if expiry.Sub(created) > 31*24*time.Hour {
		return errors.New("expiry too distant")
	}
	return nil
}
func IsSafeSize(size int64) bool        { return size > 0 && size <= 1<<30 }
func NormalizeFilename(v string) string { return strings.TrimSpace(strings.ReplaceAll(v, "..", "")) }
func PolicyReason(size int64, ttl time.Duration) string {
	if size <= 0 {
		return "empty"
	}
	if ttl <= 0 {
		return "no ttl"
	}
	return "accepted"
}
