package storage

import (
	"errors"
	"go.etcd.io/bbolt"
	"time"
)

func (s *Store) Healthy() bool { return s != nil && s.db != nil }
func (s *Store) EnsureHealthy() error {
	if !s.Healthy() {
		return errors.New("store closed")
	}
	return nil
}
func (s *Store) Ping() error {
	if e := s.EnsureHealthy(); e != nil {
		return e
	}
	return s.db.View(func(_ *bbolt.Tx) error { return nil })
}
func (s *Store) OpenedAt() time.Time { return time.Now().UTC() }
