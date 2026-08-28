package storage

import (
	"encoding/json"
	"errors"
	"go.etcd.io/bbolt"
	"os"
	"sync"
	"time"
)

type ShareCard struct {
	ID, Code, Owner, Filename, ContentType string
	Size                                   int64
	CreatedAt, ExpiresAt                   time.Time
	Deleted                                bool
	DownloadCount                          int
}
type Download struct {
	ID, CardID, Code, IP string
	At                   time.Time
	Success              bool
	Reason               string
}
type AccessGrant struct {
	ID, CardID, Recipient, Scope string
	IssuedAt, ExpiresAt          time.Time
	Revoked                      bool
}
type Store struct {
	db *bbolt.DB
	mu sync.RWMutex
}

var cardsBucket = []byte("cards")
var downloadsBucket = []byte("downloads")
var grantsBucket = []byte("grants")

func Open(path string) (*Store, error) {
	db, e := bbolt.Open(path, 0600, nil)
	if e != nil {
		return nil, e
	}
	s := &Store{db: db}
	e = db.Update(func(tx *bbolt.Tx) error {
		for _, b := range [][]byte{cardsBucket, downloadsBucket, grantsBucket} {
			if _, x := tx.CreateBucketIfNotExists(b); x != nil {
				return x
			}
		}
		return nil
	})
	if e != nil {
		db.Close()
		return nil, e
	}
	return s, nil
}
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	e := s.db.Close()
	s.db = nil
	return e
}
func encode(v interface{}) ([]byte, error) { return json.Marshal(v) }
func (s *Store) SaveCard(c ShareCard) error {
	b, e := encode(c)
	if e != nil {
		return e
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(cardsBucket).Put([]byte(c.ID), b) })
}
func (s *Store) GetCard(id string) (*ShareCard, error) {
	var c ShareCard
	e := s.db.View(func(tx *bbolt.Tx) error {
		v := tx.Bucket(cardsBucket).Get([]byte(id))
		if v == nil {
			return nil
		}
		return json.Unmarshal(v, &c)
	})
	if e != nil {
		return nil, e
	}
	if c.ID == "" {
		return nil, nil
	}
	return &c, nil
}
func (s *Store) GetCardByCode(code string) (*ShareCard, error) {
	var out *ShareCard
	e := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(cardsBucket).ForEach(func(_, v []byte) error {
			var c ShareCard
			if e := json.Unmarshal(v, &c); e != nil {
				return e
			}
			if c.Code == code {
				cp := c
				out = &cp
			}
			return nil
		})
	})
	return out, e
}
func (s *Store) SaveDownload(d Download) error {
	b, e := encode(d)
	if e != nil {
		return e
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(downloadsBucket).Put([]byte(d.ID), b) })
}
func (s *Store) SaveGrant(g AccessGrant) error {
	b, e := encode(g)
	if e != nil {
		return e
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(grantsBucket).Put([]byte(g.ID), b) })
}
func (s *Store) GetGrant(id string) (*AccessGrant, error) {
	var g AccessGrant
	e := s.db.View(func(tx *bbolt.Tx) error {
		v := tx.Bucket(grantsBucket).Get([]byte(id))
		if v == nil {
			return errors.New("grant not found")
		}
		return json.Unmarshal(v, &g)
	})
	if e != nil {
		return nil, e
	}
	return &g, nil
}
func (s *Store) Count(bucket string) (int, error) {
	n := 0
	e := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return os.ErrNotExist
		}
		return b.ForEach(func(_, _ []byte) error { n++; return nil })
	})
	return n, e
}
