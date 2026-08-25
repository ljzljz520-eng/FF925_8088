package audit

import (
	"sort"
	"strings"
	"time"
)

func Filter(entries []Entry, action string) []Entry {
	out := []Entry{}
	for _, e := range entries {
		if action == "" || e.Action == action {
			out = append(out, e)
		}
	}
	return out
}
func Sort(entries []Entry) []Entry {
	out := append([]Entry{}, entries...)
	sort.Slice(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out
}
func Search(entries []Entry, text string) []Entry {
	text = strings.ToLower(text)
	out := []Entry{}
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.Action), text) || strings.Contains(strings.ToLower(e.Target), text) {
			out = append(out, e)
		}
	}
	return out
}
func Between(entries []Entry, start, end time.Time) []Entry {
	out := []Entry{}
	for _, e := range entries {
		if !e.At.Before(start) && e.At.Before(end) {
			out = append(out, e)
		}
	}
	return out
}
func Actions(entries []Entry) map[string]int {
	m := map[string]int{}
	for _, e := range entries {
		m[e.Action]++
	}
	return m
}
func Actors(entries []Entry) map[string]int {
	m := map[string]int{}
	for _, e := range entries {
		m[e.Actor]++
	}
	return m
}
func Latest(entries []Entry) (Entry, bool) {
	if len(entries) == 0 {
		return Entry{}, false
	}
	v := Sort(entries)
	return v[len(v)-1], true
}
func Has(entries []Entry, action, target string) bool {
	for _, e := range entries {
		if e.Action == action && e.Target == target {
			return true
		}
	}
	return false
}
