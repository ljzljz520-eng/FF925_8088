package share

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var codePattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func NormalizeCode(v string) string { return strings.ToUpper(strings.TrimSpace(v)) }
func ValidateCode(v string) error {
	v = NormalizeCode(v)
	if len(v) < 6 {
		return errors.New("code too short")
	}
	if !codePattern.MatchString(v) {
		return errors.New("code contains unsupported characters")
	}
	return nil
}
func ValidateTTL(ttl time.Duration) error {
	if ttl < time.Minute {
		return errors.New("ttl too short")
	}
	if ttl > 30*24*time.Hour {
		return errors.New("ttl too long")
	}
	return nil
}
func SafeFilename(v string) error {
	if strings.TrimSpace(v) == "" || strings.ContainsAny(v, "/\\\x00") {
		return errors.New("unsafe filename")
	}
	return nil
}
func CallbackURL(v string) error {
	u, e := url.Parse(v)
	if e != nil || u.Scheme != "https" || u.Host == "" {
		return errors.New("callback must use https")
	}
	return nil
}
func IsDownloadable(c ShareCard, now time.Time) bool { return ValidAt(c, now) && c.Size > 0 }
