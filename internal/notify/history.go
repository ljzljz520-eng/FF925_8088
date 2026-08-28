package notify

import (
	"sync"
	"time"
)

type Delivery struct {
	Message  Message
	SentAt   time.Time
	Attempts int
	Error    string
}
type History struct {
	mu    sync.RWMutex
	items []Delivery
}

func NewHistory() *History        { return &History{items: []Delivery{}} }
func (h *History) Add(d Delivery) { h.mu.Lock(); h.items = append(h.items, d); h.mu.Unlock() }
func (h *History) All() []Delivery {
	h.mu.RLock()
	defer h.mu.RUnlock()
	o := make([]Delivery, len(h.items))
	copy(o, h.items)
	return o
}
func (h *History) For(to string) []Delivery {
	out := []Delivery{}
	for _, d := range h.All() {
		if d.Message.To == to {
			out = append(out, d)
		}
	}
	return out
}
func (h *History) Failed() []Delivery {
	out := []Delivery{}
	for _, d := range h.All() {
		if d.Error != "" {
			out = append(out, d)
		}
	}
	return out
}
func (h *History) Count() int { return len(h.All()) }
