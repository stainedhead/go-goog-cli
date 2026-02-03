package presenter

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stainedhead/go-goog-cli/internal/domain/account"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
)

func TestNewPlainPresenter(t *testing.T) {
	p := NewPlainPresenter()
	if p == nil {
		t.Error("NewPlainPresenter() returned nil")
	}
}

func TestPlainPresenter_RenderMessage(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders message as key-value pairs", func(t *testing.T) {
		msg := mail.NewMessage("msg-123", "thread-456", "sender@example.com", "Test Subject", "Body text")
		msg.Date = time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
		msg.AddRecipient("recipient@example.com")
		msg.AddLabel("INBOX")
		msg.IsRead = true
		msg.Snippet = "This is a snippet"

		result := p.RenderMessage(msg)

		if !strings.Contains(result, "ID: msg-123") {
			t.Error("Result should contain ID")
		}
		if !strings.Contains(result, "From: sender@example.com") {
			t.Error("Result should contain From")
		}
		if !strings.Contains(result, "Subject: Test Subject") {
			t.Error("Result should contain Subject")
		}
		if !strings.Contains(result, "2024-01-15") {
			t.Error("Result should contain Date")
		}
		if !strings.Contains(result, "Read: true") {
			t.Error("Result should contain Read status")
		}
		// Check it's line-separated
		lines := strings.Split(result, "\n")
		if len(lines) < 5 {
			t.Error("Result should have multiple lines")
		}
	})

	t.Run("renders nil message as empty string", func(t *testing.T) {
		result := p.RenderMessage(nil)
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_RenderMessages(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders multiple messages with tabs", func(t *testing.T) {
		msgs := []*mail.Message{
			mail.NewMessage("msg-1", "t-1", "a@example.com", "Subject 1", "Body"),
			mail.NewMessage("msg-2", "t-2", "b@example.com", "Subject 2", "Body"),
		}
		msgs[0].Date = time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
		msgs[1].Date = time.Date(2024, 1, 16, 11, 0, 0, 0, time.UTC)

		result := p.RenderMessages(msgs)

		// Check tab-separated format
		if !strings.Contains(result, "\t") {
			t.Error("Result should be tab-separated")
		}
		// Check both messages present
		lines := strings.Split(result, "\n")
		if len(lines) != 2 {
			t.Errorf("Expected 2 lines, got %d", len(lines))
		}
		if !strings.Contains(result, "msg-1") || !strings.Contains(result, "msg-2") {
			t.Error("Result should contain both message IDs")
		}
	})

	t.Run("renders empty list as empty string", func(t *testing.T) {
		result := p.RenderMessages([]*mail.Message{})
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})

	t.Run("renders nil list as empty string", func(t *testing.T) {
		result := p.RenderMessages(nil)
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_RenderDraft(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders draft as key-value pairs", func(t *testing.T) {
		msg := mail.NewMessage("msg-123", "thread-456", "sender@example.com", "Draft Subject", "Body")
		msg.AddRecipient("recipient@example.com")
		draft := mail.NewDraft("draft-789", msg)

		result := p.RenderDraft(draft)

		if !strings.Contains(result, "ID: draft-789") {
			t.Error("Result should contain draft ID")
		}
		if !strings.Contains(result, "Subject: Draft Subject") {
			t.Error("Result should contain subject")
		}
		if !strings.Contains(result, "To: recipient@example.com") {
			t.Error("Result should contain recipient")
		}
	})

	t.Run("renders nil draft as empty string", func(t *testing.T) {
		result := p.RenderDraft(nil)
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_RenderDrafts(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders multiple drafts", func(t *testing.T) {
		drafts := []*mail.Draft{
			mail.NewDraft("d-1", mail.NewMessage("m-1", "t-1", "a@ex.com", "Subject 1", "Body")),
			mail.NewDraft("d-2", mail.NewMessage("m-2", "t-2", "b@ex.com", "Subject 2", "Body")),
		}

		result := p.RenderDrafts(drafts)

		lines := strings.Split(result, "\n")
		if len(lines) != 2 {
			t.Errorf("Expected 2 lines, got %d", len(lines))
		}
	})

	t.Run("renders empty list as empty string", func(t *testing.T) {
		result := p.RenderDrafts([]*mail.Draft{})
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_RenderThread(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders thread as key-value pairs", func(t *testing.T) {
		thread := mail.NewThread("thread-123")
		thread.AddMessage(mail.NewMessage("msg-1", "thread-123", "a@ex.com", "Original", "Body"))
		thread.AddMessage(mail.NewMessage("msg-2", "thread-123", "b@ex.com", "Re: Original", "Reply"))
		thread.Snippet = "Thread snippet"
		thread.AddLabel("INBOX")

		result := p.RenderThread(thread)

		if !strings.Contains(result, "ID: thread-123") {
			t.Error("Result should contain thread ID")
		}
		if !strings.Contains(result, "Messages: 2") {
			t.Error("Result should show message count")
		}
		if !strings.Contains(result, "Message[0]:") {
			t.Error("Result should list messages")
		}
	})

	t.Run("renders nil thread as empty string", func(t *testing.T) {
		result := p.RenderThread(nil)
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_RenderThreads(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders multiple threads", func(t *testing.T) {
		threads := []*mail.Thread{
			mail.NewThread("t-1"),
			mail.NewThread("t-2"),
		}
		threads[0].Snippet = "First snippet"
		threads[1].Snippet = "Second snippet"

		result := p.RenderThreads(threads)

		if !strings.Contains(result, "t-1") || !strings.Contains(result, "t-2") {
			t.Error("Result should contain both thread IDs")
		}
	})

	t.Run("renders empty list as empty string", func(t *testing.T) {
		result := p.RenderThreads([]*mail.Thread{})
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_RenderLabel(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders label as key-value pairs", func(t *testing.T) {
		label := mail.NewLabel("label-123", "Important")
		label.SetColor("#ff0000", "#ffffff")
		label.MessageListVisibility = "show"

		result := p.RenderLabel(label)

		if !strings.Contains(result, "ID: label-123") {
			t.Error("Result should contain label ID")
		}
		if !strings.Contains(result, "Name: Important") {
			t.Error("Result should contain label name")
		}
		if !strings.Contains(result, "Type: user") {
			t.Error("Result should contain label type")
		}
		if !strings.Contains(result, "Background: #ff0000") {
			t.Error("Result should contain background color")
		}
	})

	t.Run("renders nil label as empty string", func(t *testing.T) {
		result := p.RenderLabel(nil)
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_RenderLabels(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders multiple labels", func(t *testing.T) {
		labels := []*mail.Label{
			mail.NewLabel("l-1", "Work"),
			mail.NewSystemLabel("l-2", "INBOX"),
		}

		result := p.RenderLabels(labels)

		lines := strings.Split(result, "\n")
		if len(lines) != 2 {
			t.Errorf("Expected 2 lines, got %d", len(lines))
		}
		if !strings.Contains(result, "Work") || !strings.Contains(result, "INBOX") {
			t.Error("Result should contain both label names")
		}
	})

	t.Run("renders empty list as empty string", func(t *testing.T) {
		result := p.RenderLabels([]*mail.Label{})
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_RenderEvent(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders event as key-value pairs", func(t *testing.T) {
		start := time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC)
		end := time.Date(2024, 1, 15, 15, 0, 0, 0, time.UTC)
		event := calendar.NewEvent("Team Meeting", start, end)
		event.ID = "event-123"
		event.Location = "Conference Room A"
		event.Description = "Weekly sync"
		event.CalendarID = "cal-456"

		result := p.RenderEvent(event)

		if !strings.Contains(result, "ID: event-123") {
			t.Error("Result should contain event ID")
		}
		if !strings.Contains(result, "Title: Team Meeting") {
			t.Error("Result should contain title")
		}
		if !strings.Contains(result, "Location: Conference Room A") {
			t.Error("Result should contain location")
		}
		if !strings.Contains(result, "Start: 2024-01-15 14:00") {
			t.Error("Result should contain start time")
		}
	})

	t.Run("renders all-day event", func(t *testing.T) {
		event := calendar.NewAllDayEvent("Holiday", time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC))
		event.ID = "holiday-123"

		result := p.RenderEvent(event)

		if !strings.Contains(result, "All Day") {
			t.Error("Result should indicate all-day event")
		}
	})

	t.Run("renders nil event as empty string", func(t *testing.T) {
		result := p.RenderEvent(nil)
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_RenderEvents(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders multiple events", func(t *testing.T) {
		now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
		events := []*calendar.Event{
			calendar.NewEvent("Event 1", now, now.Add(time.Hour)),
			calendar.NewEvent("Event 2", now.Add(2*time.Hour), now.Add(3*time.Hour)),
		}
		events[0].ID = "e-1"
		events[1].ID = "e-2"

		result := p.RenderEvents(events)

		lines := strings.Split(result, "\n")
		if len(lines) != 2 {
			t.Errorf("Expected 2 lines, got %d", len(lines))
		}
	})

	t.Run("renders empty list as empty string", func(t *testing.T) {
		result := p.RenderEvents([]*calendar.Event{})
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_RenderCalendar(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders calendar as key-value pairs", func(t *testing.T) {
		cal := calendar.NewCalendar("My Calendar")
		cal.ID = "cal-123"
		cal.TimeZone = "America/New_York"
		cal.Primary = true
		cal.Description = "Personal calendar"

		result := p.RenderCalendar(cal)

		if !strings.Contains(result, "ID: cal-123") {
			t.Error("Result should contain calendar ID")
		}
		if !strings.Contains(result, "Title: My Calendar") {
			t.Error("Result should contain title")
		}
		if !strings.Contains(result, "TimeZone: America/New_York") {
			t.Error("Result should contain time zone")
		}
		if !strings.Contains(result, "Primary: true") {
			t.Error("Result should show primary status")
		}
	})

	t.Run("renders nil calendar as empty string", func(t *testing.T) {
		result := p.RenderCalendar(nil)
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_RenderCalendars(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders multiple calendars with primary marker", func(t *testing.T) {
		cals := []*calendar.Calendar{
			calendar.NewCalendar("Work"),
			calendar.NewCalendar("Personal"),
		}
		cals[0].ID = "c-1"
		cals[0].Primary = true
		cals[1].ID = "c-2"

		result := p.RenderCalendars(cals)

		// Primary calendar should have asterisk
		if !strings.Contains(result, "*c-1") {
			t.Error("Primary calendar should be marked with *")
		}
		if strings.Contains(result, "*c-2") {
			t.Error("Non-primary calendar should not have *")
		}
	})

	t.Run("renders empty list as empty string", func(t *testing.T) {
		result := p.RenderCalendars([]*calendar.Calendar{})
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_RenderAccount(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders account as key-value pairs", func(t *testing.T) {
		acct := account.NewAccount("work", "user@company.com")
		acct.AddScope("gmail.readonly")
		acct.AddScope("calendar.readonly")
		acct.IsDefault = true

		result := p.RenderAccount(acct)

		if !strings.Contains(result, "Alias: work") {
			t.Error("Result should contain alias")
		}
		if !strings.Contains(result, "Email: user@company.com") {
			t.Error("Result should contain email")
		}
		if !strings.Contains(result, "Default: true") {
			t.Error("Result should show default status")
		}
		if !strings.Contains(result, "Scopes: 2") {
			t.Error("Result should show scope count")
		}
		if !strings.Contains(result, "- gmail.readonly") {
			t.Error("Result should list scopes")
		}
	})

	t.Run("renders nil account as empty string", func(t *testing.T) {
		result := p.RenderAccount(nil)
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_RenderAccounts(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders multiple accounts with default marker", func(t *testing.T) {
		accts := []*account.Account{
			account.NewAccount("work", "work@company.com"),
			account.NewAccount("personal", "me@gmail.com"),
		}
		accts[0].IsDefault = true
		accts[0].AddScope("gmail.readonly")
		accts[1].AddScope("calendar.readonly")
		accts[1].AddScope("gmail.send")

		result := p.RenderAccounts(accts)

		// Default account should have asterisk
		if !strings.Contains(result, "*work") {
			t.Error("Default account should be marked with *")
		}
		// Check scope counts
		if !strings.Contains(result, "\t1") {
			t.Error("Result should show scope count of 1")
		}
		if !strings.Contains(result, "\t2") {
			t.Error("Result should show scope count of 2")
		}
	})

	t.Run("renders empty list as empty string", func(t *testing.T) {
		result := p.RenderAccounts([]*account.Account{})
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_RenderError(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders error", func(t *testing.T) {
		err := errors.New("something went wrong")
		result := p.RenderError(err)

		if result != "error: something went wrong" {
			t.Errorf("Expected 'error: something went wrong', got %q", result)
		}
	})

	t.Run("renders nil error as empty string", func(t *testing.T) {
		result := p.RenderError(nil)
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_RenderSuccess(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders success message directly", func(t *testing.T) {
		result := p.RenderSuccess("Operation completed")

		if result != "Operation completed" {
			t.Errorf("Expected 'Operation completed', got %q", result)
		}
	})

	t.Run("renders empty success message", func(t *testing.T) {
		result := p.RenderSuccess("")
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestPlainPresenter_SkipsNilItems(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("skips nil messages in list", func(t *testing.T) {
		msgs := []*mail.Message{
			mail.NewMessage("msg-1", "t-1", "a@ex.com", "Subject", "Body"),
			nil,
			mail.NewMessage("msg-2", "t-2", "b@ex.com", "Subject 2", "Body"),
		}

		result := p.RenderMessages(msgs)

		lines := strings.Split(result, "\n")
		if len(lines) != 2 {
			t.Errorf("Expected 2 lines (skipping nil), got %d", len(lines))
		}
	})

	t.Run("skips nil events in list", func(t *testing.T) {
		now := time.Now()
		events := []*calendar.Event{
			calendar.NewEvent("Event 1", now, now.Add(time.Hour)),
			nil,
		}
		events[0].ID = "e-1"

		result := p.RenderEvents(events)

		lines := strings.Split(result, "\n")
		if len(lines) != 1 {
			t.Errorf("Expected 1 line (skipping nil), got %d", len(lines))
		}
	})
}

func TestPlainPresenter_PipeableFriendly(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("messages output is grep-friendly", func(t *testing.T) {
		msgs := []*mail.Message{
			mail.NewMessage("msg-1", "t-1", "sender@ex.com", "Important meeting", "Body"),
		}
		msgs[0].Date = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

		result := p.RenderMessages(msgs)

		// Should be single line per item, tab-separated
		if strings.Count(result, "\n") != 0 {
			t.Error("Single item should be on one line")
		}
		if strings.Count(result, "\t") != 3 {
			t.Errorf("Expected 3 tabs (4 fields), got %d", strings.Count(result, "\t"))
		}
	})
}
