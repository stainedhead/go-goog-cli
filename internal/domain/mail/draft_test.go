package mail

import (
	"testing"
	"time"
)

func TestNewDraft(t *testing.T) {
	msg := NewMessage("msg-1", "thread-1", "from@example.com", "Subject", "Body")
	draft := NewDraft("draft-123", msg)

	if draft.ID != "draft-123" {
		t.Errorf("expected ID 'draft-123', got '%s'", draft.ID)
	}
	if draft.Message != msg {
		t.Error("expected Message to be the provided message")
	}
	if draft.Created.IsZero() {
		t.Error("expected Created to be set")
	}
	if draft.Updated.IsZero() {
		t.Error("expected Updated to be set")
	}
	if !draft.Created.Equal(draft.Updated) {
		t.Error("expected Created and Updated to be equal for new draft")
	}
}

func TestDraft_UpdateMessage(t *testing.T) {
	msg1 := NewMessage("msg-1", "thread-1", "from@example.com", "Subject 1", "Body 1")
	draft := NewDraft("draft-123", msg1)

	originalUpdated := draft.Updated
	time.Sleep(time.Millisecond) // Ensure time difference

	msg2 := NewMessage("msg-2", "thread-1", "from@example.com", "Subject 2", "Body 2")
	draft.UpdateMessage(msg2)

	if draft.Message != msg2 {
		t.Error("expected Message to be updated to new message")
	}
	if !draft.Updated.After(originalUpdated) {
		t.Error("expected Updated timestamp to be updated")
	}
	if draft.Created.After(originalUpdated) {
		t.Error("expected Created timestamp to remain unchanged")
	}
}

func TestDraft_Touch(t *testing.T) {
	msg := NewMessage("msg-1", "thread-1", "from@example.com", "Subject", "Body")
	draft := NewDraft("draft-123", msg)

	originalUpdated := draft.Updated
	originalMessage := draft.Message
	time.Sleep(time.Millisecond) // Ensure time difference

	draft.Touch()

	if draft.Message != originalMessage {
		t.Error("expected Message to remain unchanged after Touch")
	}
	if !draft.Updated.After(originalUpdated) {
		t.Error("expected Updated timestamp to be updated after Touch")
	}
}

func TestNewDraft_WithNilMessage(t *testing.T) {
	draft := NewDraft("draft-123", nil)

	if draft.ID != "draft-123" {
		t.Errorf("expected ID 'draft-123', got '%s'", draft.ID)
	}
	if draft.Message != nil {
		t.Error("expected Message to be nil")
	}
}
