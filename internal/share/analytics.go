package share

import (
	"math"
	"strings"
	"time"
)

type Analytics struct {
	Downloads, Successes, Failures int
	ByReason                       map[string]int
	ByHour                         map[int]int
}

func Analyze(ds []Download) Analytics {
	a := Analytics{ByReason: map[string]int{}, ByHour: map[int]int{}}
	for _, d := range ds {
		a.Downloads++
		a.ByHour[d.At.Hour()]++
		if d.Success {
			a.Successes++
		} else {
			a.Failures++
			a.ByReason[d.Reason]++
		}
	}
	return a
}
func (a Analytics) SuccessRate() float64 {
	if a.Downloads == 0 {
		return 0
	}
	return float64(a.Successes) / float64(a.Downloads)
}
func (a Analytics) FailureRate() float64 {
	if a.Downloads == 0 {
		return 0
	}
	return float64(a.Failures) / float64(a.Downloads)
}
func (a Analytics) PeakHour() int {
	best, n := 0, 0
	for h, v := range a.ByHour {
		if v > n {
			best, n = h, v
		}
	}
	return best
}
func (a Analytics) Reason(name string) int { return a.ByReason[strings.ToLower(name)] }
func (a Analytics) Stable() bool           { return a.Failures == 0 || a.SuccessRate() > .5 }
func (a Analytics) Score() float64         { return math.Round(a.SuccessRate()*10000) / 100 }
func GroupDownloads(ds []Download) map[string][]Download {
	m := map[string][]Download{}
	for _, d := range ds {
		m[d.CardID] = append(m[d.CardID], d)
	}
	return m
}
func RecentDownloads(ds []Download, now time.Time, window time.Duration) []Download {
	out := []Download{}
	for _, d := range ds {
		if now.Sub(d.At) <= window {
			out = append(out, d)
		}
	}
	return out
}
func Successful(ds []Download) []Download {
	out := []Download{}
	for _, d := range ds {
		if d.Success {
			out = append(out, d)
		}
	}
	return out
}
func Failed(ds []Download) []Download {
	out := []Download{}
	for _, d := range ds {
		if !d.Success {
			out = append(out, d)
		}
	}
	return out
}
func UniqueIPs(ds []Download) int {
	m := map[string]bool{}
	for _, d := range ds {
		m[d.IP] = true
	}
	return len(m)
}
func CardDownloads(ds []Download, id string) int {
	n := 0
	for _, d := range ds {
		if d.CardID == id {
			n++
		}
	}
	return n
}
func HasFailure(ds []Download, reason string) bool {
	for _, d := range ds {
		if !d.Success && d.Reason == reason {
			return true
		}
	}
	return false
}
func MergeDownloads(a, b []Download) []Download { return append(append([]Download{}, a...), b...) }
