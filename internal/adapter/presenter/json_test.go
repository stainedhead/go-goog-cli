package presenter

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stainedhead/go-goog-cli/internal/domain/account"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
)

func TestNewJSONPresenter(t *testing.T) {
	p := NewJSONPresenter()
	if p == nil {
		t.Error("NewJSONPresenter() returned nil")
	}
}

func TestJSONPresenter_RenderMessage(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders message as JSON", func(t *testing.T) {
		msg := mail.NewMessage("msg-123", "thread-456", "sender@example.com", "Test Subject", "Body text")
		msg.Date = time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
		msg.AddLabel("INBOX")

		result := p.RenderMessage(msg)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, "msg-123") {
			t.Error("Result should contain message ID")
		}
		if !strings.Contains(result, "Test Subject") {
			t.Error("Result should contain subject")
		}
		if !strings.Contains(result, "INBOX") {
			t.Error("Result should contain label")
		}
	})

	t.Run("renders nil message", func(t *testing.T) {
		result := p.RenderMessage(nil)
		if result != "null" {
			t.Errorf("Expected 'null', got %q", result)
		}
	})
}

func TestJSONPresenter_RenderMessages(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders multiple messages", func(t *testing.T) {
		msgs := []*mail.Message{
			mail.NewMessage("msg-1", "t-1", "a@example.com", "Subject 1", "Body 1"),
			mail.NewMessage("msg-2", "t-2", "b@example.com", "Subject 2", "Body 2"),
		}

		result := p.RenderMessages(msgs)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.HasPrefix(result, "[") {
			t.Error("Result should be a JSON array")
		}
		if !strings.Contains(result, "msg-1") || !strings.Contains(result, "msg-2") {
			t.Error("Result should contain both message IDs")
		}
	})

	t.Run("renders empty list", func(t *testing.T) {
		result := p.RenderMessages([]*mail.Message{})
		if result != "[]" {
			t.Errorf("Expected '[]', got %q", result)
		}
	})

	t.Run("renders nil list", func(t *testing.T) {
		result := p.RenderMessages(nil)
		if result != "[]" {
			t.Errorf("Expected '[]', got %q", result)
		}
	})
}

func TestJSONPresenter_RenderDraft(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders draft as JSON", func(t *testing.T) {
		msg := mail.NewMessage("msg-123", "thread-456", "sender@example.com", "Draft Subject", "Draft body")
		draft := mail.NewDraft("draft-789", msg)

		result := p.RenderDraft(draft)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, "draft-789") {
			t.Error("Result should contain draft ID")
		}
		if !strings.Contains(result, "Draft Subject") {
			t.Error("Result should contain message subject")
		}
	})

	t.Run("renders nil draft", func(t *testing.T) {
		result := p.RenderDraft(nil)
		if result != "null" {
			t.Errorf("Expected 'null', got %q", result)
		}
	})
}

func TestJSONPresenter_RenderDrafts(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders multiple drafts", func(t *testing.T) {
		drafts := []*mail.Draft{
			mail.NewDraft("d-1", mail.NewMessage("m-1", "t-1", "a@ex.com", "S1", "B1")),
			mail.NewDraft("d-2", mail.NewMessage("m-2", "t-2", "b@ex.com", "S2", "B2")),
		}

		result := p.RenderDrafts(drafts)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, "d-1") || !strings.Contains(result, "d-2") {
			t.Error("Result should contain both draft IDs")
		}
	})

	t.Run("renders nil list", func(t *testing.T) {
		result := p.RenderDrafts(nil)
		if result != "[]" {
			t.Errorf("Expected '[]', got %q", result)
		}
	})
}

func TestJSONPresenter_RenderThread(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders thread as JSON", func(t *testing.T) {
		thread := mail.NewThread("thread-123")
		thread.AddMessage(mail.NewMessage("msg-1", "thread-123", "a@ex.com", "Re: Test", "Reply"))
		thread.Snippet = "This is a snippet"

		result := p.RenderThread(thread)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, "thread-123") {
			t.Error("Result should contain thread ID")
		}
		if !strings.Contains(result, "This is a snippet") {
			t.Error("Result should contain snippet")
		}
	})

	t.Run("renders nil thread", func(t *testing.T) {
		result := p.RenderThread(nil)
		if result != "null" {
			t.Errorf("Expected 'null', got %q", result)
		}
	})
}

func TestJSONPresenter_RenderThreads(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders multiple threads", func(t *testing.T) {
		threads := []*mail.Thread{
			mail.NewThread("t-1"),
			mail.NewThread("t-2"),
		}

		result := p.RenderThreads(threads)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, "t-1") || !strings.Contains(result, "t-2") {
			t.Error("Result should contain both thread IDs")
		}
	})

	t.Run("renders nil list", func(t *testing.T) {
		result := p.RenderThreads(nil)
		if result != "[]" {
			t.Errorf("Expected '[]', got %q", result)
		}
	})
}

func TestJSONPresenter_RenderLabel(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders label as JSON", func(t *testing.T) {
		label := mail.NewLabel("label-123", "Important")
		label.SetColor("#ff0000", "#ffffff")

		result := p.RenderLabel(label)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, "label-123") {
			t.Error("Result should contain label ID")
		}
		if !strings.Contains(result, "Important") {
			t.Error("Result should contain label name")
		}
	})

	t.Run("renders nil label", func(t *testing.T) {
		result := p.RenderLabel(nil)
		if result != "null" {
			t.Errorf("Expected 'null', got %q", result)
		}
	})
}

func TestJSONPresenter_RenderLabels(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders multiple labels", func(t *testing.T) {
		labels := []*mail.Label{
			mail.NewLabel("l-1", "Work"),
			mail.NewSystemLabel("l-2", "INBOX"),
		}

		result := p.RenderLabels(labels)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, "Work") || !strings.Contains(result, "INBOX") {
			t.Error("Result should contain both label names")
		}
	})

	t.Run("renders nil list", func(t *testing.T) {
		result := p.RenderLabels(nil)
		if result != "[]" {
			t.Errorf("Expected '[]', got %q", result)
		}
	})
}

func TestJSONPresenter_RenderEvent(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders event as JSON", func(t *testing.T) {
		start := time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC)
		end := time.Date(2024, 1, 15, 15, 0, 0, 0, time.UTC)
		event := calendar.NewEvent("Team Meeting", start, end)
		event.ID = "event-123"
		event.Location = "Conference Room A"

		result := p.RenderEvent(event)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, "event-123") {
			t.Error("Result should contain event ID")
		}
		if !strings.Contains(result, "Team Meeting") {
			t.Error("Result should contain event title")
		}
		if !strings.Contains(result, "Conference Room A") {
			t.Error("Result should contain location")
		}
	})

	t.Run("renders nil event", func(t *testing.T) {
		result := p.RenderEvent(nil)
		if result != "null" {
			t.Errorf("Expected 'null', got %q", result)
		}
	})
}

func TestJSONPresenter_RenderEvents(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders multiple events", func(t *testing.T) {
		now := time.Now()
		events := []*calendar.Event{
			calendar.NewEvent("Event 1", now, now.Add(time.Hour)),
			calendar.NewEvent("Event 2", now.Add(2*time.Hour), now.Add(3*time.Hour)),
		}
		events[0].ID = "e-1"
		events[1].ID = "e-2"

		result := p.RenderEvents(events)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, "e-1") || !strings.Contains(result, "e-2") {
			t.Error("Result should contain both event IDs")
		}
	})

	t.Run("renders nil list", func(t *testing.T) {
		result := p.RenderEvents(nil)
		if result != "[]" {
			t.Errorf("Expected '[]', got %q", result)
		}
	})
}

func TestJSONPresenter_RenderCalendar(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders calendar as JSON", func(t *testing.T) {
		cal := calendar.NewCalendar("My Calendar")
		cal.ID = "cal-123"
		cal.TimeZone = "America/New_York"
		cal.Primary = true

		result := p.RenderCalendar(cal)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, "cal-123") {
			t.Error("Result should contain calendar ID")
		}
		if !strings.Contains(result, "My Calendar") {
			t.Error("Result should contain calendar title")
		}
		if !strings.Contains(result, "America/New_York") {
			t.Error("Result should contain time zone")
		}
	})

	t.Run("renders nil calendar", func(t *testing.T) {
		result := p.RenderCalendar(nil)
		if result != "null" {
			t.Errorf("Expected 'null', got %q", result)
		}
	})
}

func TestJSONPresenter_RenderCalendars(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders multiple calendars", func(t *testing.T) {
		cals := []*calendar.Calendar{
			calendar.NewCalendar("Work"),
			calendar.NewCalendar("Personal"),
		}
		cals[0].ID = "c-1"
		cals[1].ID = "c-2"

		result := p.RenderCalendars(cals)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, "c-1") || !strings.Contains(result, "c-2") {
			t.Error("Result should contain both calendar IDs")
		}
	})

	t.Run("renders nil list", func(t *testing.T) {
		result := p.RenderCalendars(nil)
		if result != "[]" {
			t.Errorf("Expected '[]', got %q", result)
		}
	})
}

func TestJSONPresenter_RenderAccount(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders account as JSON", func(t *testing.T) {
		acct := account.NewAccount("work", "user@company.com")
		acct.AddScope("gmail.readonly")
		acct.IsDefault = true

		result := p.RenderAccount(acct)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, "work") {
			t.Error("Result should contain alias")
		}
		if !strings.Contains(result, "user@company.com") {
			t.Error("Result should contain email")
		}
		if !strings.Contains(result, "gmail.readonly") {
			t.Error("Result should contain scope")
		}
	})

	t.Run("renders nil account", func(t *testing.T) {
		result := p.RenderAccount(nil)
		if result != "null" {
			t.Errorf("Expected 'null', got %q", result)
		}
	})
}

func TestJSONPresenter_RenderAccounts(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders multiple accounts", func(t *testing.T) {
		accts := []*account.Account{
			account.NewAccount("work", "work@company.com"),
			account.NewAccount("personal", "me@gmail.com"),
		}

		result := p.RenderAccounts(accts)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, "work@company.com") || !strings.Contains(result, "me@gmail.com") {
			t.Error("Result should contain both account emails")
		}
	})

	t.Run("renders nil list", func(t *testing.T) {
		result := p.RenderAccounts(nil)
		if result != "[]" {
			t.Errorf("Expected '[]', got %q", result)
		}
	})
}

func TestJSONPresenter_RenderError(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders error as JSON", func(t *testing.T) {
		err := errors.New("something went wrong")
		result := p.RenderError(err)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, "error") {
			t.Error("Result should contain error key")
		}
		if !strings.Contains(result, "something went wrong") {
			t.Error("Result should contain error message")
		}
	})

	t.Run("renders nil error", func(t *testing.T) {
		result := p.RenderError(nil)

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, `"error": ""`) {
			t.Errorf("Expected empty error field, got %q", result)
		}
	})
}

func TestJSONPresenter_RenderSuccess(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders success message as JSON", func(t *testing.T) {
		result := p.RenderSuccess("Operation completed successfully")

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, "message") {
			t.Error("Result should contain message key")
		}
		if !strings.Contains(result, "Operation completed successfully") {
			t.Error("Result should contain success message")
		}
	})

	t.Run("renders empty success message", func(t *testing.T) {
		result := p.RenderSuccess("")

		if !json.Valid([]byte(result)) {
			t.Error("Result is not valid JSON")
		}
		if !strings.Contains(result, `"message": ""`) {
			t.Errorf("Expected empty message field, got %q", result)
		}
	})
}

func TestJSONPresenter_IndentFormat(t *testing.T) {
	p := NewJSONPresenter()

	msg := mail.NewMessage("id", "tid", "from@ex.com", "Subject", "Body")
	result := p.RenderMessage(msg)

	// Check that it's indented with 2 spaces
	if !strings.Contains(result, "  ") {
		t.Error("JSON should be indented with 2 spaces")
	}
	// Check that it contains newlines (pretty-printed)
	if !strings.Contains(result, "\n") {
		t.Error("JSON should be pretty-printed with newlines")
	}
}
