package notify

import (
	"errors"
	"strings"
	"sync"
	"tempcards/internal/share"
)

type Message struct{ To, Subject, Body string }
type Sender interface{ Send(Message) error }
type Dispatcher struct {
	mu     sync.Mutex
	sender Sender
	sent   []Message
}

func New(s Sender) *Dispatcher { return &Dispatcher{sender: s} }
func (d *Dispatcher) SendCard(card share.ShareCard, to string) error {
	if !strings.Contains(to, "@") {
		return errors.New("invalid recipient")
	}
	m := Message{To: to, Subject: "Temporary file share", Body: card.Filename + " is ready"}
	d.mu.Lock()
	d.sent = append(d.sent, m)
	d.mu.Unlock()
	if d.sender != nil {
		return d.sender.Send(m)
	}
	return nil
}
func (d *Dispatcher) Messages() []Message {
	d.mu.Lock()
	defer d.mu.Unlock()
	o := make([]Message, len(d.sent))
	copy(o, d.sent)
	return o
}

type MemorySender struct{ MessagesSeen []Message }

func (m *MemorySender) Send(v Message) error { m.MessagesSeen = append(m.MessagesSeen, v); return nil }
