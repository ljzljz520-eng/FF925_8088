package metrics

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Report struct {
	Generated time.Time
	Values    map[string]int64
	Tags      map[string]int64
}

func (c *Counters) Report() Report {
	return Report{Generated: time.Now().UTC(), Values: c.Snapshot(), Tags: c.Tags()}
}
func (r Report) Value(k string) int64 { return r.Values[k] }
func (r Report) Format() string {
	keys := make([]string, 0, len(r.Values))
	for k := range r.Values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	p := make([]string, 0, len(keys))
	for _, k := range keys {
		p = append(p, fmt.Sprintf("%s=%d", k, r.Values[k]))
	}
	return strings.Join(p, " ")
}
func (r Report) HasActivity() bool {
	for _, v := range r.Values {
		if v > 0 {
			return true
		}
	}
	return false
}
