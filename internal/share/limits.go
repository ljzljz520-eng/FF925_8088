package share

import (
	"errors"
	"sort"
	"time"
)

type DownloadWindow struct {
	Started time.Time
	Count   int
	Limit   int
}

func (w DownloadWindow) Allowed(now time.Time) bool {
	return w.Count < w.Limit && now.Sub(w.Started) < time.Hour
}
func (w *DownloadWindow) Record() { w.Count++ }
func (w DownloadWindow) Remaining() int {
	if w.Limit <= w.Count {
		return 0
	}
	return w.Limit - w.Count
}
func ValidateSize(size, max int64) error {
	if size <= 0 {
		return errors.New("empty file")
	}
	if max > 0 && size > max {
		return errors.New("file too large")
	}
	return nil
}
func SortByExpiry(cards []ShareCard) []ShareCard {
	out := append([]ShareCard{}, cards...)
	sort.Slice(out, func(i, j int) bool { return out[i].ExpiresAt.Before(out[j].ExpiresAt) })
	return out
}
func Earliest(cards []ShareCard) (ShareCard, bool) {
	v := SortByExpiry(cards)
	if len(v) == 0 {
		return ShareCard{}, false
	}
	return v[0], true
}
func Extend(c ShareCard, ttl time.Duration) ShareCard {
	if ttl > 0 {
		c.ExpiresAt = c.ExpiresAt.Add(ttl)
	}
	return c
}
