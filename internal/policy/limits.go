package policy

import (
	"errors"
	"time"
)

type Limits struct {
	Concurrent, Daily int
	Window            time.Duration
}

func StandardLimits() Limits { return Limits{Concurrent: 3, Daily: 100, Window: 24 * time.Hour} }
func (l Limits) CheckConcurrent(n int) error {
	if n < 0 {
		return errors.New("negative count")
	}
	if n >= l.Concurrent {
		return errors.New("concurrency limit")
	}
	return nil
}
func (l Limits) CheckDaily(n int) error {
	if n > l.Daily {
		return errors.New("daily limit")
	}
	return nil
}
func (l Limits) ResetAt(now time.Time) time.Time { return now.Add(l.Window) }
