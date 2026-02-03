package mail

import "time"

// Draft represents an email draft.
type Draft struct {
	ID      string
	Message *Message
	Created time.Time
	Updated time.Time
}

// NewDraft creates a new Draft with the given ID and message.
func NewDraft(id string, message *Message) *Draft {
	now := time.Now()
	return &Draft{
		ID:      id,
		Message: message,
		Created: now,
		Updated: now,
	}
}

// UpdateMessage updates the draft's message and sets the Updated timestamp.
func (d *Draft) UpdateMessage(message *Message) {
	d.Message = message
	d.Updated = time.Now()
}

// Touch updates the Updated timestamp without changing the message.
func (d *Draft) Touch() {
	d.Updated = time.Now()
}
