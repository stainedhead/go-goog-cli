package mail

// Thread represents an email conversation thread.
type Thread struct {
	ID       string
	Messages []*Message
	Snippet  string
	Labels   []string
}

// NewThread creates a new Thread with the given ID.
func NewThread(id string) *Thread {
	return &Thread{
		ID:       id,
		Messages: []*Message{},
		Labels:   []string{},
	}
}

// AddMessage adds a message to the thread.
func (t *Thread) AddMessage(msg *Message) {
	t.Messages = append(t.Messages, msg)
}

// MessageCount returns the number of messages in the thread.
func (t *Thread) MessageCount() int {
	return len(t.Messages)
}

// AddLabel adds a label to the thread.
func (t *Thread) AddLabel(label string) {
	t.Labels = append(t.Labels, label)
}

// RemoveLabel removes a label from the thread.
func (t *Thread) RemoveLabel(label string) {
	for i, l := range t.Labels {
		if l == label {
			t.Labels = append(t.Labels[:i], t.Labels[i+1:]...)
			return
		}
	}
}

// HasLabel checks if the thread has a specific label.
func (t *Thread) HasLabel(label string) bool {
	for _, l := range t.Labels {
		if l == label {
			return true
		}
	}
	return false
}

// LatestMessage returns the most recently added message in the thread, or nil if empty.
func (t *Thread) LatestMessage() *Message {
	if len(t.Messages) == 0 {
		return nil
	}
	return t.Messages[len(t.Messages)-1]
}

// FirstMessage returns the first message in the thread, or nil if empty.
func (t *Thread) FirstMessage() *Message {
	if len(t.Messages) == 0 {
		return nil
	}
	return t.Messages[0]
}
