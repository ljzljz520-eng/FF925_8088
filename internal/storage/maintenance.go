package storage

import (
	"encoding/json"
	"go.etcd.io/bbolt"
	"time"
)

func (s *Store) UpdateCard(id string, fn func(*ShareCard) error) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(cardsBucket)
		v := b.Get([]byte(id))
		if v == nil {
			return nil
		}
		var c ShareCard
		if e := json.Unmarshal(v, &c); e != nil {
			return e
		}
		if e := fn(&c); e != nil {
			return e
		}
		nv, e := json.Marshal(c)
		if e != nil {
			return e
		}
		return b.Put([]byte(id), nv)
	})
}
func (s *Store) DeleteCard(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(cardsBucket).Delete([]byte(id)) })
}
func (s *Store) TouchCard(id string, at time.Time) error {
	return s.UpdateCard(id, func(c *ShareCard) error { c.CreatedAt = at; return nil })
}
func (s *Store) MarkExpired(now time.Time) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(cardsBucket).ForEach(func(k, v []byte) error {
			var c ShareCard
			if e := json.Unmarshal(v, &c); e != nil {
				return e
			}
			if now.After(c.ExpiresAt) && !c.Deleted {
				c.Deleted = true
				nv, _ := json.Marshal(c)
				return tx.Bucket(cardsBucket).Put(k, nv)
			}
			return nil
		})
	})
}
func (s *Store) SnapshotCards() ([]ShareCard, error) { return s.ListCards("") }
