package mail

import (
	"testing"
	"time"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage("msg-123", "thread-456", "sender@example.com", "Test Subject", "Test body")

	if msg.ID != "msg-123" {
		t.Errorf("expected ID 'msg-123', got '%s'", msg.ID)
	}
	if msg.ThreadID != "thread-456" {
		t.Errorf("expected ThreadID 'thread-456', got '%s'", msg.ThreadID)
	}
	if msg.From != "sender@example.com" {
		t.Errorf("expected From 'sender@example.com', got '%s'", msg.From)
	}
	if msg.Subject != "Test Subject" {
		t.Errorf("expected Subject 'Test Subject', got '%s'", msg.Subject)
	}
	if msg.Body != "Test body" {
		t.Errorf("expected Body 'Test body', got '%s'", msg.Body)
	}
	if len(msg.To) != 0 {
		t.Errorf("expected empty To slice, got %v", msg.To)
	}
	if len(msg.Cc) != 0 {
		t.Errorf("expected empty Cc slice, got %v", msg.Cc)
	}
	if len(msg.Bcc) != 0 {
		t.Errorf("expected empty Bcc slice, got %v", msg.Bcc)
	}
	if len(msg.Labels) != 0 {
		t.Errorf("expected empty Labels slice, got %v", msg.Labels)
	}
	if msg.Date.IsZero() {
		t.Error("expected Date to be set")
	}
}

func TestMessage_AddRecipient(t *testing.T) {
	msg := NewMessage("1", "1", "from@example.com", "Subject", "Body")
	msg.AddRecipient("to1@example.com")
	msg.AddRecipient("to2@example.com")

	if len(msg.To) != 2 {
		t.Errorf("expected 2 recipients, got %d", len(msg.To))
	}
	if msg.To[0] != "to1@example.com" {
		t.Errorf("expected first recipient 'to1@example.com', got '%s'", msg.To[0])
	}
	if msg.To[1] != "to2@example.com" {
		t.Errorf("expected second recipient 'to2@example.com', got '%s'", msg.To[1])
	}
}

func TestMessage_AddCc(t *testing.T) {
	msg := NewMessage("1", "1", "from@example.com", "Subject", "Body")
	msg.AddCc("cc@example.com")

	if len(msg.Cc) != 1 {
		t.Errorf("expected 1 Cc recipient, got %d", len(msg.Cc))
	}
	if msg.Cc[0] != "cc@example.com" {
		t.Errorf("expected Cc 'cc@example.com', got '%s'", msg.Cc[0])
	}
}

func TestMessage_AddBcc(t *testing.T) {
	msg := NewMessage("1", "1", "from@example.com", "Subject", "Body")
	msg.AddBcc("bcc@example.com")

	if len(msg.Bcc) != 1 {
		t.Errorf("expected 1 Bcc recipient, got %d", len(msg.Bcc))
	}
	if msg.Bcc[0] != "bcc@example.com" {
		t.Errorf("expected Bcc 'bcc@example.com', got '%s'", msg.Bcc[0])
	}
}

func TestMessage_Labels(t *testing.T) {
	msg := NewMessage("1", "1", "from@example.com", "Subject", "Body")

	// Test AddLabel
	msg.AddLabel("INBOX")
	msg.AddLabel("UNREAD")

	if len(msg.Labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(msg.Labels))
	}

	// Test HasLabel
	if !msg.HasLabel("INBOX") {
		t.Error("expected message to have label 'INBOX'")
	}
	if msg.HasLabel("SPAM") {
		t.Error("expected message not to have label 'SPAM'")
	}

	// Test RemoveLabel
	msg.RemoveLabel("UNREAD")
	if len(msg.Labels) != 1 {
		t.Errorf("expected 1 label after removal, got %d", len(msg.Labels))
	}
	if msg.HasLabel("UNREAD") {
		t.Error("expected message not to have label 'UNREAD' after removal")
	}

	// Test removing non-existent label (should not panic)
	msg.RemoveLabel("NONEXISTENT")
	if len(msg.Labels) != 1 {
		t.Errorf("expected 1 label after removing non-existent, got %d", len(msg.Labels))
	}
}

func TestMessage_ReadStatus(t *testing.T) {
	msg := NewMessage("1", "1", "from@example.com", "Subject", "Body")

	if msg.IsRead {
		t.Error("expected new message to be unread")
	}

	msg.MarkAsRead()
	if !msg.IsRead {
		t.Error("expected message to be read after MarkAsRead")
	}

	msg.MarkAsUnread()
	if msg.IsRead {
		t.Error("expected message to be unread after MarkAsUnread")
	}
}

func TestMessage_StarStatus(t *testing.T) {
	msg := NewMessage("1", "1", "from@example.com", "Subject", "Body")

	if msg.IsStarred {
		t.Error("expected new message to be unstarred")
	}

	msg.Star()
	if !msg.IsStarred {
		t.Error("expected message to be starred after Star")
	}

	msg.Unstar()
	if msg.IsStarred {
		t.Error("expected message to be unstarred after Unstar")
	}
}

func TestMessage_DateIsSet(t *testing.T) {
	before := time.Now()
	msg := NewMessage("1", "1", "from@example.com", "Subject", "Body")
	after := time.Now()

	if msg.Date.Before(before) || msg.Date.After(after) {
		t.Error("expected Date to be set to current time")
	}
}
