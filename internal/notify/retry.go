package notify

import (
	"errors"
	"time"
)

type RetryPolicy struct {
	Attempts int
	Delay    time.Duration
}

func DefaultRetry() RetryPolicy { return RetryPolicy{Attempts: 3, Delay: time.Millisecond} }
func (r RetryPolicy) Send(s Sender, m Message) error {
	if r.Attempts < 1 {
		return errors.New("attempts required")
	}
	var e error
	for i := 0; i < r.Attempts; i++ {
		e = s.Send(m)
		if e == nil {
			return nil
		}
		if r.Delay > 0 {
			time.Sleep(r.Delay)
		}
	}
	return e
}
func (r RetryPolicy) Backoff(n int) time.Duration {
	if n < 0 {
		return 0
	}
	return r.Delay * time.Duration(n+1)
}
func (r RetryPolicy) Valid() bool { return r.Attempts > 0 && r.Attempts < 20 }
