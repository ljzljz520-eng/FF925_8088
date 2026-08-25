package share

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type CardReport struct {
	Total, Active, Deleted, Expired int
	Bytes                           int64
	Owners                          map[string]int
}

func BuildReport(cards []ShareCard, now time.Time) CardReport {
	r := CardReport{Owners: map[string]int{}}
	for _, c := range cards {
		r.Total++
		r.Bytes += c.Size
		r.Owners[c.Owner]++
		if c.Deleted {
			r.Deleted++
		} else if !now.Before(c.ExpiresAt) {
			r.Expired++
		} else {
			r.Active++
		}
	}
	return r
}
func (r CardReport) Utilization() float64 {
	if r.Total == 0 {
		return 0
	}
	return float64(r.Active) / float64(r.Total)
}
func (r CardReport) OwnerCount() int { return len(r.Owners) }
func (r CardReport) TopOwner() string {
	best := ""
	n := 0
	for o, v := range r.Owners {
		if v > n {
			best = o
			n = v
		}
	}
	return best
}
func (r CardReport) Summary() string {
	return fmt.Sprintf("total=%d active=%d deleted=%d expired=%d bytes=%d", r.Total, r.Active, r.Deleted, r.Expired, r.Bytes)
}
func FilterActive(cards []ShareCard, now time.Time) []ShareCard {
	out := []ShareCard{}
	for _, c := range cards {
		if ValidAt(c, now) {
			out = append(out, c)
		}
	}
	return out
}
func FilterDeleted(cards []ShareCard) []ShareCard {
	out := []ShareCard{}
	for _, c := range cards {
		if c.Deleted {
			out = append(out, c)
		}
	}
	return out
}
func FilterOwner(cards []ShareCard, owner string) []ShareCard {
	out := []ShareCard{}
	for _, c := range cards {
		if c.Owner == owner {
			out = append(out, c)
		}
	}
	return out
}
func SortCreated(cards []ShareCard) []ShareCard {
	out := append([]ShareCard{}, cards...)
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out
}
func SortSize(cards []ShareCard) []ShareCard {
	out := append([]ShareCard{}, cards...)
	sort.Slice(out, func(i, j int) bool { return out[i].Size > out[j].Size })
	return out
}
func Codes(cards []ShareCard) []string {
	out := []string{}
	for _, c := range cards {
		out = append(out, c.Code)
	}
	return out
}
func JoinCodes(cards []ShareCard) string { return strings.Join(Codes(cards), ",") }
func HasCode(cards []ShareCard, code string) bool {
	for _, c := range cards {
		if c.Code == code {
			return true
		}
	}
	return false
}
func FindOwner(cards []ShareCard, code string) string {
	for _, c := range cards {
		if c.Code == code {
			return c.Owner
		}
	}
	return ""
}
func CountByExpiry(cards []ShareCard, now time.Time) map[string]int {
	m := map[string]int{"day": 0, "week": 0, "later": 0}
	for _, c := range cards {
		d := c.ExpiresAt.Sub(now)
		if d < 24*time.Hour {
			m["day"]++
		} else if d < 7*24*time.Hour {
			m["week"]++
		} else {
			m["later"]++
		}
	}
	return m
}
func CanShare(c ShareCard, now time.Time) error {
	if c.Deleted {
		return ErrDeleted
	}
	if !now.Before(c.ExpiresAt) {
		return ErrExpired
	}
	if c.Size <= 0 {
		return fmt.Errorf("invalid size")
	}
	return nil
}
func Clone(c ShareCard) ShareCard                      { return c }
func WithContentType(c ShareCard, ct string) ShareCard { c.ContentType = ct; return c }
func WithOwner(c ShareCard, o string) ShareCard        { c.Owner = o; return c }
func WithExpiry(c ShareCard, t time.Time) ShareCard    { c.ExpiresAt = t; return c }
func IsOwned(c ShareCard, o string) bool               { return c.Owner == o && !c.Deleted }
func IsExpired(c ShareCard, now time.Time) bool        { return !now.Before(c.ExpiresAt) }
func IsDeleted(c ShareCard) bool                       { return c.Deleted }
func CardAge(c ShareCard, now time.Time) time.Duration { return now.Sub(c.CreatedAt) }
