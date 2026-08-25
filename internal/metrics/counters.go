package metrics

import (
	"sync"
	"sync/atomic"
)

type Counters struct {
	created, downloads, rejected int64
	mu                           sync.RWMutex
	labels                       map[string]int64
}

func New() *Counters            { return &Counters{labels: map[string]int64{}} }
func (c *Counters) Created()    { atomic.AddInt64(&c.created, 1) }
func (c *Counters) Downloaded() { atomic.AddInt64(&c.downloads, 1) }
func (c *Counters) Rejected()   { atomic.AddInt64(&c.rejected, 1) }
func (c *Counters) Snapshot() map[string]int64 {
	return map[string]int64{"created": atomic.LoadInt64(&c.created), "downloads": atomic.LoadInt64(&c.downloads), "rejected": atomic.LoadInt64(&c.rejected)}
}
func (c *Counters) Tag(k string) { c.mu.Lock(); c.labels[k]++; c.mu.Unlock() }
func (c *Counters) Tags() map[string]int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	o := map[string]int64{}
	for k, v := range c.labels {
		o[k] = v
	}
	return o
}
