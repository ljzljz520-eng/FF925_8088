package audit

import (
	"fmt"
	"sync"
	"time"
)

type Entry struct {
	ID, Action, Actor, Target string
	At                        time.Time
	Metadata                  map[string]string
}
type Log struct {
	mu      sync.RWMutex
	entries []Entry
}

func New() *Log { return &Log{entries: make([]Entry, 0)} }
func (l *Log) Append(action, actor, target string, meta map[string]string) Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := Entry{ID: fmt.Sprintf("evt-%d", time.Now().UnixNano()), Action: action, Actor: actor, Target: target, At: time.Now().UTC(), Metadata: meta}
	l.entries = append(l.entries, e)
	return e
}
func (l *Log) List() []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}
func (l *Log) ForActor(actor string) []Entry {
	all := l.List()
	out := make([]Entry, 0)
	for _, e := range all {
		if e.Actor == actor {
			out = append(out, e)
		}
	}
	return out
}
func (l *Log) Since(t time.Time) []Entry {
	all := l.List()
	out := make([]Entry, 0)
	for _, e := range all {
		if e.At.After(t) {
			out = append(out, e)
		}
	}
	return out
}
