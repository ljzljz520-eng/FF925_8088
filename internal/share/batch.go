package share

import (
	"errors"
	"time"
)

type BatchResult struct {
	Processed, Deleted, Skipped int
	Errors                      []error
}

func DeleteExpired(cards []ShareCard, now time.Time) (BatchResult, []ShareCard) {
	r := BatchResult{}
	out := append([]ShareCard{}, cards...)
	for i := range out {
		r.Processed++
		if out[i].Deleted {
			r.Skipped++
			continue
		}
		if !now.Before(out[i].ExpiresAt) {
			out[i].Deleted = true
			r.Deleted++
		}
	}
	return r, out
}
func ValidateBatch(cards []ShareCard) error {
	seen := map[string]bool{}
	for _, c := range cards {
		if c.ID == "" || seen[c.ID] {
			return errors.New("duplicate or empty id")
		}
		seen[c.ID] = true
		if c.Code == "" {
			return errors.New("empty code")
		}
	}
	return nil
}
func Merge(a, b []ShareCard) []ShareCard {
	m := map[string]ShareCard{}
	for _, c := range a {
		m[c.ID] = c
	}
	for _, c := range b {
		m[c.ID] = c
	}
	out := make([]ShareCard, 0, len(m))
	for _, c := range m {
		out = append(out, c)
	}
	return SortCreated(out)
}
func Partition(cards []ShareCard, now time.Time) ([]ShareCard, []ShareCard) {
	active := []ShareCard{}
	inactive := []ShareCard{}
	for _, c := range cards {
		if ValidAt(c, now) {
			active = append(active, c)
		} else {
			inactive = append(inactive, c)
		}
	}
	return active, inactive
}
func MapOwners(cards []ShareCard) map[string][]ShareCard {
	m := map[string][]ShareCard{}
	for _, c := range cards {
		m[c.Owner] = append(m[c.Owner], c)
	}
	return m
}
func MapCodes(cards []ShareCard) map[string]ShareCard {
	m := map[string]ShareCard{}
	for _, c := range cards {
		m[c.Code] = c
	}
	return m
}
func TotalSize(cards []ShareCard) int64 {
	var n int64
	for _, c := range cards {
		n += c.Size
	}
	return n
}
func AverageSize(cards []ShareCard) int64 {
	if len(cards) == 0 {
		return 0
	}
	return TotalSize(cards) / int64(len(cards))
}
func TouchAll(cards []ShareCard, at time.Time) []ShareCard {
	out := append([]ShareCard{}, cards...)
	for i := range out {
		out[i].CreatedAt = at
	}
	return out
}
func ExtendAll(cards []ShareCard, ttl time.Duration) []ShareCard {
	out := append([]ShareCard{}, cards...)
	for i := range out {
		out[i] = Extend(out[i], ttl)
	}
	return out
}
