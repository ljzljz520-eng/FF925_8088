package auth

import (
	"sync"
	"time"
)

type Session struct {
	mu     sync.RWMutex
	Tokens map[string]Token
}

func NewSession() *Session     { return &Session{Tokens: map[string]Token{}} }
func (s *Session) Put(t Token) { s.mu.Lock(); s.Tokens[t.Value] = t; s.mu.Unlock() }
func (s *Session) Get(v string) (Token, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.Tokens[v]
	return t, ok
}
func (s *Session) Valid(v string, now time.Time) bool { t, ok := s.Get(v); return ok && t.Valid(now) }
func (s *Session) Revoke(v string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.Tokens[v]
	if !ok {
		return false
	}
	t.Revoked = true
	s.Tokens[v] = t
	return true
}
func (s *Session) Purge(now time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for k, t := range s.Tokens {
		if !t.Valid(now) {
			delete(s.Tokens, k)
			n++
		}
	}
	return n
}
