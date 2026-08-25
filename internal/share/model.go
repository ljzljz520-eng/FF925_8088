package share

import (
	"tempcards/internal/storage"
	"time"
)

type ShareCard = storage.ShareCard
type Download = storage.Download
type AccessGrant = storage.AccessGrant
type ShareEvent struct {
	ID, CardID, Kind, Actor string
	At                      time.Time
	Detail                  string
}

func ValidAt(c ShareCard, now time.Time) bool {
	return !c.Deleted && now.Before(c.ExpiresAt) && c.Code != ""
}
func Remaining(c ShareCard, now time.Time) time.Duration {
	if now.After(c.ExpiresAt) {
		return 0
	}
	return c.ExpiresAt.Sub(now)
}
func Active(g AccessGrant, now time.Time) bool { return !g.Revoked && now.Before(g.ExpiresAt) }
func NewShareCard(id, code, owner, filename string, size int64, ttl time.Duration) ShareCard {
	now := time.Now().UTC()
	return ShareCard{ID: id, Code: code, Owner: owner, Filename: filename, Size: size, CreatedAt: now, ExpiresAt: now.Add(ttl)}
}
