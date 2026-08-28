package policy

import (
	"sort"
	"strings"
)

func NormalizeLabels(in []string) []string {
	m := map[string]bool{}
	for _, v := range in {
		v = strings.ToLower(strings.TrimSpace(v))
		if v != "" {
			m[v] = true
		}
	}
	out := []string{}
	for v := range m {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}
func HasLabel(labels []string, want string) bool {
	want = strings.ToLower(strings.TrimSpace(want))
	for _, v := range labels {
		if strings.ToLower(v) == want {
			return true
		}
	}
	return false
}
func AddLabel(labels []string, label string) []string { return NormalizeLabels(append(labels, label)) }
func RemoveLabel(labels []string, label string) []string {
	out := []string{}
	for _, v := range labels {
		if !strings.EqualFold(v, label) {
			out = append(out, v)
		}
	}
	return out
}
func LabelCount(labels []string) map[string]int {
	m := map[string]int{}
	for _, v := range NormalizeLabels(labels) {
		m[v]++
	}
	return m
}
