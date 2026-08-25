package storage

import (
	"encoding/json"
)

func EncodeCard(c ShareCard) ([]byte, error) { return json.Marshal(c) }
func DecodeCard(b []byte) (ShareCard, error) {
	var c ShareCard
	e := json.Unmarshal(b, &c)
	return c, e
}
func Export(cards []ShareCard) ([]byte, error) { return json.Marshal(cards) }
func Import(b []byte) ([]ShareCard, error) {
	var c []ShareCard
	e := json.Unmarshal(b, &c)
	return c, e
}
func CardType(c ShareCard) string {
	if c.ContentType != "" {
		return c.ContentType
	}
	return "application/octet-stream"
}
func IsPersistable(c ShareCard) bool {
	return c.ID != "" && c.Code != "" && c.ExpiresAt.After(c.CreatedAt)
}
func Normalize(c ShareCard) ShareCard {
	if c.ContentType == "" {
		c.ContentType = CardType(c)
	}
	return c
}
func (s *Store) SaveCards(cards []ShareCard) error {
	for _, c := range cards {
		if e := s.SaveCard(Normalize(c)); e != nil {
			return e
		}
	}
	return nil
}
func (s *Store) LoadCards() ([]ShareCard, error) { return s.ListCards("") }
