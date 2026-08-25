package storage

import (
	"encoding/json"
	"os"
	"time"
)

type Backup struct {
	Created time.Time
	Cards   []ShareCard
}

func (s *Store) MakeBackup() (Backup, error) {
	c, e := s.LoadCards()
	return Backup{Created: time.Now().UTC(), Cards: c}, e
}
func (b Backup) Bytes() ([]byte, error)   { return json.Marshal(b) }
func LoadBackup(v []byte) (Backup, error) { var b Backup; e := json.Unmarshal(v, &b); return b, e }
func WriteBackup(path string, b Backup) error {
	v, e := b.Bytes()
	if e != nil {
		return e
	}
	return os.WriteFile(path, v, 0600)
}
func ReadBackup(path string) (Backup, error) {
	v, e := os.ReadFile(path)
	if e != nil {
		return Backup{}, e
	}
	return LoadBackup(v)
}
func (s *Store) RestoreBackup(b Backup) error         { return s.SaveCards(b.Cards) }
func BackupAge(b Backup, now time.Time) time.Duration { return now.Sub(b.Created) }
func ValidBackup(b Backup) bool                       { return !b.Created.IsZero() }
