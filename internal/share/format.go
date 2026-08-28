package share

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func MarshalCard(c ShareCard) ([]byte, error) { return json.Marshal(c) }
func UnmarshalCard(v []byte) (ShareCard, error) {
	var c ShareCard
	e := json.Unmarshal(v, &c)
	return c, e
}
func CardString(c ShareCard) string { return fmt.Sprintf("%s:%s:%s", c.ID, c.Code, c.Filename) }
func CardLabels(c ShareCard) []string {
	l := []string{c.Owner}
	if c.Deleted {
		l = append(l, "deleted")
	}
	if IsExpired(c, time.Now()) {
		l = append(l, "expired")
	}
	return l
}
func NormalizeOwner(v string) string { return strings.ToLower(strings.TrimSpace(v)) }
func NormalizeType(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "" {
		return "application/octet-stream"
	}
	return v
}
func NormalizeCard(c ShareCard) ShareCard {
	c.Owner = NormalizeOwner(c.Owner)
	c.ContentType = NormalizeType(c.ContentType)
	c.Code = NormalizeCode(c.Code)
	return c
}
func CardEqual(a, b ShareCard) bool {
	return a.ID == b.ID && a.Code == b.Code && a.Owner == b.Owner && a.Filename == b.Filename && a.Size == b.Size && a.Deleted == b.Deleted
}
func CardChanged(a, b ShareCard) bool { return !CardEqual(a, b) }
func CardVersion(c ShareCard) string  { return c.ID + "-" + c.ExpiresAt.UTC().Format(time.RFC3339Nano) }
func RedactCode(c ShareCard) string {
	if len(c.Code) < 4 {
		return "****"
	}
	return c.Code[:2] + "****" + c.Code[len(c.Code)-2:]
}
func RedactedString(c ShareCard) string {
	return fmt.Sprintf("%s:%s:%s", c.ID, RedactCode(c), c.Filename)
}
func IsLarge(c ShareCard) bool  { return c.Size > 10<<20 }
func IsSmall(c ShareCard) bool  { return c.Size > 0 && c.Size <= 1<<20 }
func IsMedium(c ShareCard) bool { return c.Size > 1<<20 && c.Size <= 10<<20 }
func SizeClass(c ShareCard) string {
	if IsSmall(c) {
		return "small"
	}
	if IsMedium(c) {
		return "medium"
	}
	if IsLarge(c) {
		return "large"
	}
	return "empty"
}
func ExpiryClass(c ShareCard, now time.Time) string {
	d := c.ExpiresAt.Sub(now)
	if d <= 0 {
		return "expired"
	}
	if d < time.Hour {
		return "urgent"
	}
	if d < 24*time.Hour {
		return "soon"
	}
	return "later"
}
func IsCodeMatch(c ShareCard, code string) bool {
	return strings.EqualFold(c.Code, strings.TrimSpace(code))
}
func FindCode(cards []ShareCard, code string) (ShareCard, bool) {
	for _, c := range cards {
		if IsCodeMatch(c, code) {
			return c, true
		}
	}
	return ShareCard{}, false
}
func CountOwner(cards []ShareCard, owner string) int {
	n := 0
	for _, c := range cards {
		if strings.EqualFold(c.Owner, owner) {
			n++
		}
	}
	return n
}
func CodesForOwner(cards []ShareCard, owner string) []string {
	out := []string{}
	for _, c := range cards {
		if strings.EqualFold(c.Owner, owner) {
			out = append(out, c.Code)
		}
	}
	return out
}
func ExpiredCount(cards []ShareCard, now time.Time) int {
	n := 0
	for _, c := range cards {
		if IsExpired(c, now) {
			n++
		}
	}
	return n
}
func ActiveCount(cards []ShareCard, now time.Time) int { return len(FilterActive(cards, now)) }
func DeletedCount(cards []ShareCard) int               { return len(FilterDeleted(cards)) }
