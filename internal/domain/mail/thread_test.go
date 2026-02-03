package mail

import "testing"

func TestNewThread(t *testing.T) {
	thread := NewThread("thread-123")

	if thread.ID != "thread-123" {
		t.Errorf("expected ID 'thread-123', got '%s'", thread.ID)
	}
	if len(thread.Messages) != 0 {
		t.Errorf("expected empty Messages slice, got %d messages", len(thread.Messages))
	}
	if len(thread.Labels) != 0 {
		t.Errorf("expected empty Labels slice, got %v", thread.Labels)
	}
}

func TestThread_AddMessage(t *testing.T) {
	thread := NewThread("thread-123")
	msg1 := NewMessage("msg-1", "thread-123", "from@example.com", "Subject 1", "Body 1")
	msg2 := NewMessage("msg-2", "thread-123", "from@example.com", "Subject 2", "Body 2")

	thread.AddMessage(msg1)
	thread.AddMessage(msg2)

	if thread.MessageCount() != 2 {
		t.Errorf("expected 2 messages, got %d", thread.MessageCount())
	}
	if thread.Messages[0] != msg1 {
		t.Error("expected first message to be msg1")
	}
	if thread.Messages[1] != msg2 {
		t.Error("expected second message to be msg2")
	}
}

func TestThread_MessageCount(t *testing.T) {
	thread := NewThread("thread-123")

	if thread.MessageCount() != 0 {
		t.Errorf("expected 0 messages in new thread, got %d", thread.MessageCount())
	}

	thread.AddMessage(NewMessage("msg-1", "thread-123", "from@example.com", "Subject", "Body"))
	if thread.MessageCount() != 1 {
		t.Errorf("expected 1 message, got %d", thread.MessageCount())
	}
}

func TestThread_Labels(t *testing.T) {
	thread := NewThread("thread-123")

	// Test AddLabel
	thread.AddLabel("INBOX")
	thread.AddLabel("IMPORTANT")

	if len(thread.Labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(thread.Labels))
	}

	// Test HasLabel
	if !thread.HasLabel("INBOX") {
		t.Error("expected thread to have label 'INBOX'")
	}
	if thread.HasLabel("SPAM") {
		t.Error("expected thread not to have label 'SPAM'")
	}

	// Test RemoveLabel
	thread.RemoveLabel("INBOX")
	if len(thread.Labels) != 1 {
		t.Errorf("expected 1 label after removal, got %d", len(thread.Labels))
	}
	if thread.HasLabel("INBOX") {
		t.Error("expected thread not to have label 'INBOX' after removal")
	}

	// Test removing non-existent label (should not panic)
	thread.RemoveLabel("NONEXISTENT")
	if len(thread.Labels) != 1 {
		t.Errorf("expected 1 label after removing non-existent, got %d", len(thread.Labels))
	}
}

func TestThread_LatestMessage(t *testing.T) {
	thread := NewThread("thread-123")

	// Test empty thread
	if thread.LatestMessage() != nil {
		t.Error("expected nil for LatestMessage on empty thread")
	}

	msg1 := NewMessage("msg-1", "thread-123", "from@example.com", "Subject 1", "Body 1")
	msg2 := NewMessage("msg-2", "thread-123", "from@example.com", "Subject 2", "Body 2")

	thread.AddMessage(msg1)
	if thread.LatestMessage() != msg1 {
		t.Error("expected LatestMessage to be msg1")
	}

	thread.AddMessage(msg2)
	if thread.LatestMessage() != msg2 {
		t.Error("expected LatestMessage to be msg2 after adding second message")
	}
}

func TestThread_FirstMessage(t *testing.T) {
	thread := NewThread("thread-123")

	// Test empty thread
	if thread.FirstMessage() != nil {
		t.Error("expected nil for FirstMessage on empty thread")
	}

	msg1 := NewMessage("msg-1", "thread-123", "from@example.com", "Subject 1", "Body 1")
	msg2 := NewMessage("msg-2", "thread-123", "from@example.com", "Subject 2", "Body 2")

	thread.AddMessage(msg1)
	if thread.FirstMessage() != msg1 {
		t.Error("expected FirstMessage to be msg1")
	}

	thread.AddMessage(msg2)
	if thread.FirstMessage() != msg1 {
		t.Error("expected FirstMessage to still be msg1 after adding second message")
	}
}
