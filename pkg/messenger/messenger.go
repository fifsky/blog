package messenger

import "context"

// Action represents a button action in the message.
type Action struct {
	Title string // Button title
	URL   string // Action URL (for URL-based actions)
}

// Message represents a notification message to be sent.
type Message struct {
	Title   string   // Message title
	Content string   // Message content (markdown supported)
	Time    string   // Optional time display
	Token   string   // Token to pass in callback value
	Actions []Action // Action buttons
}

// Sender is the interface for sending notification messages.
type Sender interface {
	Send(ctx context.Context, msg Message) error
}

// MultiSender sends messages to multiple channels.
type MultiSender struct {
	senders []Sender
}

// NewMultiSender creates a new MultiSender with the given senders.
func NewMultiSender(senders ...Sender) *MultiSender {
	return &MultiSender{senders: senders}
}

// Send sends the message to all configured channels.
func (m *MultiSender) Send(ctx context.Context, msg Message) error {
	var lastErr error
	for _, sender := range m.senders {
		if err := sender.Send(ctx, msg); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
