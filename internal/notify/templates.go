package notify

import (
	"fmt"
	"strings"
	"tempcards/internal/share"
)

func Subject(kind string) string {
	switch strings.ToLower(kind) {
	case "created":
		return "Your share card is ready"
	case "expired":
		return "Your share card expired"
	case "deleted":
		return "Your share card was deleted"
	default:
		return "Share card update"
	}
}
func Body(card share.ShareCard, kind string) string {
	return fmt.Sprintf("Card %s for %s: %s", card.Code, kind, card.Filename)
}
func Digest(cards []share.ShareCard) string {
	parts := make([]string, 0, len(cards))
	for _, c := range cards {
		parts = append(parts, c.Filename)
	}
	return strings.Join(parts, ", ")
}
