package share

import "time"

type Retention struct {
	Grace, MaxAge time.Duration
	PurgeDeleted  bool
}

func DefaultRetention() Retention {
	return Retention{Grace: time.Hour, MaxAge: 7 * 24 * time.Hour, PurgeDeleted: true}
}
func (r Retention) ShouldPurge(c ShareCard, now time.Time) bool {
	if c.Deleted && r.PurgeDeleted && now.Sub(c.CreatedAt) > r.Grace {
		return true
	}
	return now.Sub(c.CreatedAt) > r.MaxAge
}
func (r Retention) ShouldNotify(c ShareCard, now time.Time) bool {
	return !c.Deleted && c.ExpiresAt.After(now) && c.ExpiresAt.Sub(now) < r.Grace
}
func (r Retention) NextSweep(now time.Time) time.Time { return now.Add(r.Grace) }
func (r Retention) Age(c ShareCard, now time.Time) time.Duration {
	if now.Before(c.CreatedAt) {
		return 0
	}
	return now.Sub(c.CreatedAt)
}
func (r Retention) Eligible(cards []ShareCard, now time.Time) []ShareCard {
	out := []ShareCard{}
	for _, c := range cards {
		if r.ShouldPurge(c, now) {
			out = append(out, c)
		}
	}
	return out
}
func (r Retention) Active(cards []ShareCard, now time.Time) []ShareCard {
	out := []ShareCard{}
	for _, c := range cards {
		if ValidAt(c, now) {
			out = append(out, c)
		}
	}
	return out
}
