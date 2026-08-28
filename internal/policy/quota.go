package policy

import (
	"errors"
	"sync"
	"time"
)

type Quota struct {
	mu    sync.Mutex
	Used  int64
	Limit int64
	Reset time.Time
}

func NewQuota(limit int64) *Quota {
	return &Quota{Limit: limit, Reset: time.Now().UTC().Add(24 * time.Hour)}
}
func (q *Quota) Consume(n int64) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if n <= 0 {
		return errors.New("invalid amount")
	}
	if q.Used+n > q.Limit {
		return errors.New("quota exceeded")
	}
	q.Used += n
	return nil
}
func (q *Quota) Remaining() int64 { q.mu.Lock(); defer q.mu.Unlock(); return q.Limit - q.Used }
func (q *Quota) ResetIfNeeded(now time.Time) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if now.After(q.Reset) {
		q.Used = 0
		q.Reset = now.Add(24 * time.Hour)
	}
}
func (q *Quota) Percentage() float64 {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.Limit == 0 {
		return 1
	}
	return float64(q.Used) / float64(q.Limit)
}
