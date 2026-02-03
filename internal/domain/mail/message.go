// Package mail provides domain entities for email operations.
package mail

import "time"

// Message represents an email message.
type Message struct {
	ID        string
	ThreadID  string
	From      string
	To        []string
	Cc        []string
	Bcc       []string
	Subject   string
	Body      string
	BodyHTML  string
	Labels    []string
	Date      time.Time
	IsRead    bool
	IsStarred bool
	Snippet   string
}

// NewMessage creates a new Message with the given parameters.
func NewMessage(id, threadID, from, subject, body string) *Message {
	return &Message{
		ID:       id,
		ThreadID: threadID,
		From:     from,
		Subject:  subject,
		Body:     body,
		To:       []string{},
		Cc:       []string{},
		Bcc:      []string{},
		Labels:   []string{},
		Date:     time.Now(),
	}
}

// AddRecipient adds a recipient to the To field.
func (m *Message) AddRecipient(email string) {
	m.To = append(m.To, email)
}

// AddCc adds a recipient to the Cc field.
func (m *Message) AddCc(email string) {
	m.Cc = append(m.Cc, email)
}

// AddBcc adds a recipient to the Bcc field.
func (m *Message) AddBcc(email string) {
	m.Bcc = append(m.Bcc, email)
}

// AddLabel adds a label to the message.
func (m *Message) AddLabel(label string) {
	m.Labels = append(m.Labels, label)
}

// RemoveLabel removes a label from the message.
func (m *Message) RemoveLabel(label string) {
	for i, l := range m.Labels {
		if l == label {
			m.Labels = append(m.Labels[:i], m.Labels[i+1:]...)
			return
		}
	}
}

// HasLabel checks if the message has a specific label.
func (m *Message) HasLabel(label string) bool {
	for _, l := range m.Labels {
		if l == label {
			return true
		}
	}
	return false
}

// MarkAsRead marks the message as read.
func (m *Message) MarkAsRead() {
	m.IsRead = true
}

// MarkAsUnread marks the message as unread.
func (m *Message) MarkAsUnread() {
	m.IsRead = false
}

// Star marks the message as starred.
func (m *Message) Star() {
	m.IsStarred = true
}

// Unstar removes the star from the message.
func (m *Message) Unstar() {
	m.IsStarred = false
}
