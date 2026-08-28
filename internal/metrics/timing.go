package metrics

import (
	"sync"
	"time"
)

type Timer struct {
	mu     sync.Mutex
	values []time.Duration
}

func NewTimer() *Timer                   { return &Timer{values: []time.Duration{}} }
func (t *Timer) Observe(d time.Duration) { t.mu.Lock(); t.values = append(t.values, d); t.mu.Unlock() }
func (t *Timer) Count() int              { t.mu.Lock(); defer t.mu.Unlock(); return len(t.values) }
func (t *Timer) Average() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.values) == 0 {
		return 0
	}
	var n time.Duration
	for _, v := range t.values {
		n += v
	}
	return n / time.Duration(len(t.values))
}
