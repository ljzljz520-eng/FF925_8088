package share

import (
	"errors"
	"fmt"
	"strings"
	"tempcards/internal/storage"
	"time"
)

var ErrInvalidCode = errors.New("invalid card code")
var ErrExpired = errors.New("card expired")
var ErrDeleted = errors.New("card deleted")

type Service struct {
	store *storage.Store
	clock func() time.Time
}

func NewService(s *storage.Store) *Service {
	return &Service{store: s, clock: func() time.Time { return time.Now().UTC() }}
}
func (s *Service) Create(owner, filename, code string, size int64, ttl time.Duration) (ShareCard, error) {
	if strings.TrimSpace(owner) == "" || strings.TrimSpace(filename) == "" {
		return ShareCard{}, errors.New("owner and filename required")
	}
	if len(code) < 6 {
		return ShareCard{}, errors.New("code too short")
	}
	c := NewShareCard(fmt.Sprintf("card-%d", s.clock().UnixNano()), code, owner, filename, size, ttl)
	return c, s.store.SaveCard(c)
}
func (s *Service) Find(code string) (*ShareCard, error) {
	c, err := s.store.GetCardByCode(code)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}
	return c, nil
}
func (s *Service) Authorize(code string) (ShareCard, error) {
	c, err := s.Find(code)
	if err != nil {
		return ShareCard{}, err
	}
	if c == nil {
		return ShareCard{}, ErrInvalidCode
	}
	if c.Deleted {
		return ShareCard{}, ErrDeleted
	}
	if !ValidAt(*c, s.clock()) {
		return ShareCard{}, ErrExpired
	}
	return *c, nil
}
func (s *Service) Download(code, ip string) (Download, error) {
	c, err := s.Find(code)
	if err != nil {
		return Download{}, err
	}
	if c == nil {
		return Download{Code: code, IP: ip, At: s.clock(), Success: false, Reason: "invalid"}, ErrInvalidCode
	}
	if c.Deleted {
		return Download{CardID: c.ID, Code: code, IP: ip, At: s.clock(), Success: false, Reason: "deleted"}, ErrDeleted
	}
	if !ValidAt(*c, s.clock()) {
		return Download{CardID: c.ID, Code: code, IP: ip, At: s.clock(), Success: false, Reason: "expired"}, ErrExpired
	}
	c.DownloadCount++
	if err := s.store.SaveCard(*c); err != nil {
		return Download{}, err
	}
	d := Download{ID: fmt.Sprintf("dl-%d", s.clock().UnixNano()), CardID: c.ID, Code: code, IP: ip, At: s.clock(), Success: true}
	return d, s.store.SaveDownload(d)
}
func (s *Service) Delete(id string) error {
	c, err := s.store.GetCard(id)
	if err != nil {
		return err
	}
	if c == nil {
		return ErrInvalidCode
	}
	c.Deleted = true
	return s.store.SaveCard(*c)
}
func (s *Service) Grant(cardID, recipient, scope string, ttl time.Duration) (AccessGrant, error) {
	if recipient == "" {
		return AccessGrant{}, errors.New("recipient required")
	}
	g := AccessGrant{ID: fmt.Sprintf("grant-%d", s.clock().UnixNano()), CardID: cardID, Recipient: recipient, Scope: scope, IssuedAt: s.clock(), ExpiresAt: s.clock().Add(ttl)}
	return g, s.store.SaveGrant(g)
}
func (s *Service) RevokeGrant(id string) error {
	g, e := s.store.GetGrant(id)
	if e != nil {
		return e
	}
	if g == nil {
		return errors.New("grant missing")
	}
	g.Revoked = true
	return s.store.SaveGrant(*g)
}
