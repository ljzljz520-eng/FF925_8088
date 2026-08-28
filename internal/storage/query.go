package storage

import (
	"encoding/json"
	"go.etcd.io/bbolt"
	"sort"
	"strings"
	"time"
)

func (s *Store) ListCards(owner string) ([]ShareCard, error) {
	out := []ShareCard{}
	e := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(cardsBucket).ForEach(func(_, v []byte) error {
			var c ShareCard
			if e := json.Unmarshal(v, &c); e != nil {
				return e
			}
			if owner == "" || c.Owner == owner {
				out = append(out, c)
			}
			return nil
		})
	})
	return out, e
}
func FilterCards(cards []ShareCard, owner string, includeDeleted bool) []ShareCard {
	out := []ShareCard{}
	for _, c := range cards {
		if owner != "" && c.Owner != owner {
			continue
		}
		if !includeDeleted && c.Deleted {
			continue
		}
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out
}
func Expiring(cards []ShareCard, now time.Time, window time.Duration) []ShareCard {
	out := []ShareCard{}
	for _, c := range cards {
		d := c.ExpiresAt.Sub(now)
		if d >= 0 && d <= window && !c.Deleted {
			out = append(out, c)
		}
	}
	return out
}
func Search(cards []ShareCard, term string) []ShareCard {
	term = strings.ToLower(strings.TrimSpace(term))
	out := []ShareCard{}
	for _, c := range cards {
		if strings.Contains(strings.ToLower(c.Filename), term) || strings.Contains(strings.ToLower(c.Owner), term) {
			out = append(out, c)
		}
	}
	return out
}
