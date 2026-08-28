package share

import (
	"sort"
	"time"
)

type Ledger struct{ events []ShareEvent }

func NewLedger() *Ledger { return &Ledger{events: []ShareEvent{}} }
func (l *Ledger) Add(card, kind, actor, detail string) ShareEvent {
	e := ShareEvent{ID: card + "-" + kind + "-" + time.Now().Format("150405.000000"), CardID: card, Kind: kind, Actor: actor, At: time.Now().UTC(), Detail: detail}
	l.events = append(l.events, e)
	return e
}
func (l *Ledger) Events(card string) []ShareEvent {
	out := []ShareEvent{}
	for _, e := range l.events {
		if e.CardID == card {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out
}
func (l *Ledger) Latest(card string) (ShareEvent, bool) {
	v := l.Events(card)
	if len(v) == 0 {
		return ShareEvent{}, false
	}
	return v[len(v)-1], true
}
func (l *Ledger) Count(kind string) int {
	n := 0
	for _, e := range l.events {
		if e.Kind == kind {
			n++
		}
	}
	return n
}
