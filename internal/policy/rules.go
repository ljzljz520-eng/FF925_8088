package policy

import (
	"errors"
	"strings"
	"tempcards/internal/share"
	"time"
)

type Rule struct {
	Name             string
	MaxSize          int64
	TTL              time.Duration
	RequireRecipient bool
}

func Default() Rule {
	return Rule{Name: "standard", MaxSize: 100 << 20, TTL: 24 * time.Hour, RequireRecipient: true}
}
func (r Rule) Validate(filename string, size int64) error {
	if strings.TrimSpace(filename) == "" {
		return errors.New("filename required")
	}
	if size <= 0 || size > r.MaxSize {
		return errors.New("size outside policy")
	}
	return nil
}
func (r Rule) CanDownload(c share.ShareCard, now time.Time) bool {
	if !share.ValidAt(c, now) {
		return false
	}
	return c.Size <= r.MaxSize
}
func (r Rule) Expiry(now time.Time) time.Time {
	if r.TTL <= 0 {
		return now
	}
	return now.Add(r.TTL)
}
