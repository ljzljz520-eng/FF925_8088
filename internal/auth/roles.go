package auth

import (
	"errors"
	"sort"
	"strings"
)

var roleRank = map[string]int{"viewer": 1, "editor": 2, "owner": 3, "admin": 4}

func Roles() []string {
	out := []string{}
	for r := range roleRank {
		out = append(out, r)
	}
	sort.Strings(out)
	return out
}
func ValidRole(r string) bool { _, ok := roleRank[strings.ToLower(r)]; return ok }
func AtLeast(actual, required string) bool {
	return roleRank[strings.ToLower(actual)] >= roleRank[strings.ToLower(required)]
}
func RequireRole(c Claims, r string) error {
	if !ValidRole(r) {
		return errors.New("unknown role")
	}
	if !c.Allows(r) {
		return errors.New("forbidden")
	}
	return nil
}
func NormalizeRole(r string) string { return strings.ToLower(strings.TrimSpace(r)) }
