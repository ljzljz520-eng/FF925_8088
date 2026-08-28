package metrics

import (
	"testing"
	"time"
)

func TestMetricsCounters(t *testing.T) {
	c := New()
	c.Created()
	c.Downloaded()
	if c.Snapshot()["created"] != 1 {
		t.Fatal()
	}
	tm := NewTimer()
	tm.Observe(time.Second)
	if tm.Count() != 1 {
		t.Fatal()
	}
}
