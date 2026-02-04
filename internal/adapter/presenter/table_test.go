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

func TestNewTablePresenter(t *testing.T) {
	p := NewTablePresenter()
	if p == nil {
		t.Error("NewTablePresenter() returned nil")
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{"short string unchanged", "hello", 10, "hello"},
		{"exact length unchanged", "hello", 5, "hello"},
		{"long string truncated", "hello world", 8, "hello..."},
		{"very short maxLen", "hello", 3, "hel"},
		{"empty string", "", 10, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncate(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

// isTableFormatted checks if the output looks like a table (has multiple lines and structured output)
func isTableFormatted(s string) bool {
	lines := strings.Split(s, "\n")
	return len(lines) > 2 // At least header, separator, and data
}

func TestTablePresenter_RenderMessage(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders message as table", func(t *testing.T) {
		msg := mail.NewMessage("msg-123456789012", "thread-456", "sender@example.com", "Test Subject", "Body")
		msg.Date = time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
		msg.AddLabel("INBOX")
		msg.IsRead = true

		result := p.RenderMessage(msg)

		if !strings.Contains(result, "msg-123456789012") {
			t.Error("Result should contain message ID")
		}
		if !strings.Contains(result, "sender@example.com") {
			t.Error("Result should contain From")
		}
		if !strings.Contains(result, "Test Subject") {
			t.Error("Result should contain Subject")
		}
		if !strings.Contains(result, "2024-01-15") {
			t.Error("Result should contain Date")
		}
		if !strings.Contains(result, "INBOX") {
			t.Error("Result should contain Labels")
		}
		// Check it's formatted as a table (has multiple lines)
		if !isTableFormatted(result) {
			t.Error("Result should be formatted as a table")
		}
	})

	t.Run("renders nil message", func(t *testing.T) {
		result := p.RenderMessage(nil)
		if result != "No message found" {
			t.Errorf("Expected 'No message found', got %q", result)
		}
	})
}

func TestTablePresenter_RenderMessages(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders multiple messages with truncation", func(t *testing.T) {
		msgs := []*mail.Message{
			mail.NewMessage("msg-123456789012345", "t-1", "sender@example.com", "This is a very long subject that should be truncated", "Body"),
			mail.NewMessage("msg-2", "t-2", "another@example.com", "Short", "Body"),
		}
		msgs[0].Date = time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
		msgs[1].Date = time.Date(2024, 1, 16, 11, 0, 0, 0, time.UTC)

		result := p.RenderMessages(msgs)

		// Check truncation is happening (... should appear)
		if !strings.Contains(result, "...") {
			t.Error("Long fields should be truncated with ...")
		}
		// Check table structure - headers should be in uppercase in new tablewriter
		upperResult := strings.ToUpper(result)
		if !strings.Contains(upperResult, "ID") || !strings.Contains(upperResult, "FROM") {
			t.Error("Table should have headers")
		}
	})

	t.Run("renders empty list", func(t *testing.T) {
		result := p.RenderMessages([]*mail.Message{})
		if result != "No messages found" {
			t.Errorf("Expected 'No messages found', got %q", result)
		}
	})

	t.Run("renders nil list", func(t *testing.T) {
		result := p.RenderMessages(nil)
		if result != "No messages found" {
			t.Errorf("Expected 'No messages found', got %q", result)
		}
	})
}

func TestTablePresenter_RenderDraft(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders draft as table", func(t *testing.T) {
		msg := mail.NewMessage("msg-123", "thread-456", "sender@example.com", "Draft Subject", "Body")
		draft := mail.NewDraft("draft-789", msg)

		result := p.RenderDraft(draft)

		if !strings.Contains(result, "draft-789") {
			t.Error("Result should contain draft ID")
		}
		if !strings.Contains(result, "Draft Subject") {
			t.Error("Result should contain subject")
		}
	})

	t.Run("renders nil draft", func(t *testing.T) {
		result := p.RenderDraft(nil)
		if result != "No draft found" {
			t.Errorf("Expected 'No draft found', got %q", result)
		}
	})
}

func TestTablePresenter_RenderDrafts(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders multiple drafts", func(t *testing.T) {
		drafts := []*mail.Draft{
			mail.NewDraft("d-1", mail.NewMessage("m-1", "t-1", "a@ex.com", "Subject One", "Body")),
			mail.NewDraft("d-2", mail.NewMessage("m-2", "t-2", "b@ex.com", "Subject Two", "Body")),
		}

		result := p.RenderDrafts(drafts)

		if !strings.Contains(result, "d-1") && !strings.Contains(result, "d-2") {
			t.Error("Result should contain draft IDs")
		}
		if !strings.Contains(result, "Subject One") {
			t.Error("Result should contain subjects")
		}
	})

	t.Run("renders empty list", func(t *testing.T) {
		result := p.RenderDrafts([]*mail.Draft{})
		if result != "No drafts found" {
			t.Errorf("Expected 'No drafts found', got %q", result)
		}
	})
}

func TestTablePresenter_RenderThread(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders thread as table", func(t *testing.T) {
		thread := mail.NewThread("thread-123")
		thread.AddMessage(mail.NewMessage("msg-1", "thread-123", "a@ex.com", "Original", "Body"))
		thread.AddMessage(mail.NewMessage("msg-2", "thread-123", "b@ex.com", "Re: Original", "Reply"))
		thread.Snippet = "This is the thread snippet"

		result := p.RenderThread(thread)

		if !strings.Contains(result, "thread-123") {
			t.Error("Result should contain thread ID")
		}
		if !strings.Contains(result, "2") {
			t.Error("Result should show message count")
		}
		if !strings.Contains(result, "Messages:") {
			t.Error("Result should contain messages section")
		}
	})

	t.Run("renders nil thread", func(t *testing.T) {
		result := p.RenderThread(nil)
		if result != "No thread found" {
			t.Errorf("Expected 'No thread found', got %q", result)
		}
	})
}

func TestTablePresenter_RenderThreads(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders multiple threads", func(t *testing.T) {
		threads := []*mail.Thread{
			mail.NewThread("t-1"),
			mail.NewThread("t-2"),
		}
		threads[0].Snippet = "First thread snippet"
		threads[1].Snippet = "Second thread snippet"

		result := p.RenderThreads(threads)

		if !strings.Contains(result, "t-1") || !strings.Contains(result, "t-2") {
			t.Error("Result should contain thread IDs")
		}
	})

	t.Run("renders empty list", func(t *testing.T) {
		result := p.RenderThreads([]*mail.Thread{})
		if result != "No threads found" {
			t.Errorf("Expected 'No threads found', got %q", result)
		}
	})
}

func TestTablePresenter_RenderLabel(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders label as table", func(t *testing.T) {
		label := mail.NewLabel("label-123", "Important")
		label.SetColor("#ff0000", "#ffffff")

		result := p.RenderLabel(label)

		if !strings.Contains(result, "label-123") {
			t.Error("Result should contain label ID")
		}
		if !strings.Contains(result, "Important") {
			t.Error("Result should contain label name")
		}
		if !strings.Contains(result, "user") {
			t.Error("Result should contain label type")
		}
		if !strings.Contains(result, "#ff0000") {
			t.Error("Result should contain background color")
		}
	})

	t.Run("renders nil label", func(t *testing.T) {
		result := p.RenderLabel(nil)
		if result != "No label found" {
			t.Errorf("Expected 'No label found', got %q", result)
		}
	})
}

func TestTablePresenter_RenderLabels(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders multiple labels", func(t *testing.T) {
		labels := []*mail.Label{
			mail.NewLabel("l-1", "Work"),
			mail.NewSystemLabel("l-2", "INBOX"),
		}

		result := p.RenderLabels(labels)

		if !strings.Contains(result, "Work") {
			t.Error("Result should contain Work label")
		}
		if !strings.Contains(result, "INBOX") {
			t.Error("Result should contain INBOX label")
		}
		if !strings.Contains(result, "user") {
			t.Error("Result should contain user type")
		}
		if !strings.Contains(result, "system") {
			t.Error("Result should contain system type")
		}
	})

	t.Run("renders empty list", func(t *testing.T) {
		result := p.RenderLabels([]*mail.Label{})
		if result != "No labels found" {
			t.Errorf("Expected 'No labels found', got %q", result)
		}
	})
}

func TestTablePresenter_RenderEvent(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders event as table", func(t *testing.T) {
		start := time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC)
		end := time.Date(2024, 1, 15, 15, 0, 0, 0, time.UTC)
		event := calendar.NewEvent("Team Meeting", start, end)
		event.ID = "event-123"
		event.Location = "Conference Room A"
		event.Description = "Weekly sync"

		result := p.RenderEvent(event)

		if !strings.Contains(result, "event-123") {
			t.Error("Result should contain event ID")
		}
		if !strings.Contains(result, "Team Meeting") {
			t.Error("Result should contain event title")
		}
		if !strings.Contains(result, "Conference Room A") {
			t.Error("Result should contain location")
		}
		if !strings.Contains(result, "2024-01-15 14:00") {
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

	t.Run("renders nil event", func(t *testing.T) {
		result := p.RenderEvent(nil)
		if result != "No event found" {
			t.Errorf("Expected 'No event found', got %q", result)
		}
	})
}

func TestTablePresenter_RenderEvents(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders multiple events with truncation", func(t *testing.T) {
		now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
		events := []*calendar.Event{
			calendar.NewEvent("This is a very long event title that needs truncation", now, now.Add(time.Hour)),
			calendar.NewEvent("Short", now.Add(2*time.Hour), now.Add(3*time.Hour)),
		}
		events[0].ID = "event-123456789012345"
		events[1].ID = "e-2"
		events[0].Location = "A very long location name that should be truncated"

		result := p.RenderEvents(events)

		// Check headers (case-insensitive since tablewriter may uppercase)
		upperResult := strings.ToUpper(result)
		if !strings.Contains(upperResult, "ID") || !strings.Contains(upperResult, "TITLE") {
			t.Error("Table should have headers")
		}
		// Check truncation
		if !strings.Contains(result, "...") {
			t.Error("Long fields should be truncated")
		}
	})

	t.Run("renders empty list", func(t *testing.T) {
		result := p.RenderEvents([]*calendar.Event{})
		if result != "No events found" {
			t.Errorf("Expected 'No events found', got %q", result)
		}
	})
}

func TestTablePresenter_RenderCalendar(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders calendar as table", func(t *testing.T) {
		cal := calendar.NewCalendar("My Calendar")
		cal.ID = "cal-123"
		cal.TimeZone = "America/New_York"
		cal.Primary = true
		cal.Description = "My primary calendar"

		result := p.RenderCalendar(cal)

		if !strings.Contains(result, "cal-123") {
			t.Error("Result should contain calendar ID")
		}
		if !strings.Contains(result, "My Calendar") {
			t.Error("Result should contain title")
		}
		if !strings.Contains(result, "America/New_York") {
			t.Error("Result should contain time zone")
		}
		if !strings.Contains(result, "true") {
			t.Error("Result should show primary status")
		}
	})

	t.Run("renders nil calendar", func(t *testing.T) {
		result := p.RenderCalendar(nil)
		if result != "No calendar found" {
			t.Errorf("Expected 'No calendar found', got %q", result)
		}
	})
}

func TestTablePresenter_RenderCalendars(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders multiple calendars", func(t *testing.T) {
		cals := []*calendar.Calendar{
			calendar.NewCalendar("Work"),
			calendar.NewCalendar("Personal"),
		}
		cals[0].ID = "c-1"
		cals[0].Primary = true
		cals[1].ID = "c-2"

		result := p.RenderCalendars(cals)

		if !strings.Contains(result, "Work") || !strings.Contains(result, "Personal") {
			t.Error("Result should contain both calendar titles")
		}
		if !strings.Contains(result, "Yes") {
			t.Error("Result should show 'Yes' for primary calendar")
		}
	})

	t.Run("renders empty list", func(t *testing.T) {
		result := p.RenderCalendars([]*calendar.Calendar{})
		if result != "No calendars found" {
			t.Errorf("Expected 'No calendars found', got %q", result)
		}
	})
}

func TestTablePresenter_RenderAccount(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders account as table", func(t *testing.T) {
		acct := account.NewAccount("work", "user@company.com")
		acct.AddScope("gmail.readonly")
		acct.AddScope("calendar.readonly")
		acct.IsDefault = true

		result := p.RenderAccount(acct)

		if !strings.Contains(result, "work") {
			t.Error("Result should contain alias")
		}
		if !strings.Contains(result, "user@company.com") {
			t.Error("Result should contain email")
		}
		if !strings.Contains(result, "2") {
			t.Error("Result should show scope count")
		}
		if !strings.Contains(result, "gmail.readonly") {
			t.Error("Result should list scopes")
		}
	})

	t.Run("renders nil account", func(t *testing.T) {
		result := p.RenderAccount(nil)
		if result != "No account found" {
			t.Errorf("Expected 'No account found', got %q", result)
		}
	})
}

func TestTablePresenter_RenderAccounts(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders multiple accounts", func(t *testing.T) {
		accts := []*account.Account{
			account.NewAccount("work", "work@company.com"),
			account.NewAccount("personal", "me@gmail.com"),
		}
		accts[0].IsDefault = true
		accts[0].AddScope("gmail.readonly")
		accts[1].AddScope("calendar.readonly")
		accts[1].AddScope("gmail.send")

		result := p.RenderAccounts(accts)

		// Check headers (case-insensitive)
		upperResult := strings.ToUpper(result)
		if !strings.Contains(upperResult, "ALIAS") || !strings.Contains(upperResult, "EMAIL") {
			t.Error("Table should have headers")
		}
		if !strings.Contains(result, "work@company.com") || !strings.Contains(result, "me@gmail.com") {
			t.Error("Result should contain both emails")
		}
		if !strings.Contains(result, "Yes") {
			t.Error("Result should show 'Yes' for default account")
		}
	})

	t.Run("renders empty list", func(t *testing.T) {
		result := p.RenderAccounts([]*account.Account{})
		if result != "No accounts found" {
			t.Errorf("Expected 'No accounts found', got %q", result)
		}
	})
}

func TestTablePresenter_RenderError(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders error", func(t *testing.T) {
		err := errors.New("something went wrong")
		result := p.RenderError(err)

		if result != "Error: something went wrong" {
			t.Errorf("Expected 'Error: something went wrong', got %q", result)
		}
	})

	t.Run("renders nil error", func(t *testing.T) {
		result := p.RenderError(nil)
		if result != "" {
			t.Errorf("Expected empty string, got %q", result)
		}
	})
}

func TestTablePresenter_RenderSuccess(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders success message", func(t *testing.T) {
		result := p.RenderSuccess("Operation completed")

		if result != "Success: Operation completed" {
			t.Errorf("Expected 'Success: Operation completed', got %q", result)
		}
	})
}

func TestTablePresenter_SkipsNilItems(t *testing.T) {
	p := NewTablePresenter()

	t.Run("skips nil messages in list", func(t *testing.T) {
		msgs := []*mail.Message{
			mail.NewMessage("msg-1", "t-1", "a@ex.com", "Subject", "Body"),
			nil,
			mail.NewMessage("msg-2", "t-2", "b@ex.com", "Subject 2", "Body"),
		}

		result := p.RenderMessages(msgs)

		if !strings.Contains(result, "msg-1") || !strings.Contains(result, "msg-2") {
			t.Error("Result should contain non-nil message IDs")
		}
	})
}

// =============================================================================
// ACL Rule Tests (0% coverage improvement)
// =============================================================================

func TestTablePresenter_RenderACLRule(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders ACL rule as table", func(t *testing.T) {
		rule := &calendar.ACLRule{
			ID:   "user:alice@example.com",
			Role: "writer",
			Scope: &calendar.ACLScope{
				Type:  "user",
				Value: "alice@example.com",
			},
		}

		result := p.RenderACLRule(rule)

		if !strings.Contains(result, "user:alice@example.com") {
			t.Error("Result should contain rule ID")
		}
		if !strings.Contains(result, "writer") {
			t.Error("Result should contain role")
		}
		if !strings.Contains(result, "user") {
			t.Error("Result should contain scope type")
		}
		if !strings.Contains(result, "alice@example.com") {
			t.Error("Result should contain scope value")
		}
		// Should be formatted as table
		if !isTableFormatted(result) {
			t.Error("Result should be formatted as a table")
		}
	})

	t.Run("renders nil ACL rule", func(t *testing.T) {
		result := p.RenderACLRule(nil)
		if result != "No ACL rule found" {
			t.Errorf("Expected 'No ACL rule found', got %q", result)
		}
	})

	t.Run("renders ACL rule without scope", func(t *testing.T) {
		rule := &calendar.ACLRule{
			ID:   "default",
			Role: "freeBusyReader",
		}

		result := p.RenderACLRule(rule)

		if !strings.Contains(result, "default") {
			t.Error("Result should contain ID")
		}
		if !strings.Contains(result, "freeBusyReader") {
			t.Error("Result should contain Role")
		}
	})

	t.Run("renders ACL rule with scope but no value", func(t *testing.T) {
		rule := &calendar.ACLRule{
			ID:   "default",
			Role: "freeBusyReader",
			Scope: &calendar.ACLScope{
				Type: "default",
			},
		}

		result := p.RenderACLRule(rule)

		if !strings.Contains(result, "Scope Type") {
			t.Error("Result should contain Scope Type field")
		}
	})
}

func TestTablePresenter_RenderACLRules(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders multiple ACL rules as table", func(t *testing.T) {
		rules := []*calendar.ACLRule{
			{
				ID:   "user:alice@example.com",
				Role: "writer",
				Scope: &calendar.ACLScope{
					Type:  "user",
					Value: "alice@example.com",
				},
			},
			{
				ID:   "domain:example.com",
				Role: "reader",
				Scope: &calendar.ACLScope{
					Type:  "domain",
					Value: "example.com",
				},
			},
		}

		result := p.RenderACLRules(rules)

		// Check headers (case-insensitive)
		upperResult := strings.ToUpper(result)
		if !strings.Contains(upperResult, "ID") || !strings.Contains(upperResult, "ROLE") {
			t.Error("Table should have headers")
		}
		if !strings.Contains(result, "alice@example.com") {
			t.Error("Result should contain first rule")
		}
		if !strings.Contains(result, "example.com") {
			t.Error("Result should contain second rule")
		}
	})

	t.Run("renders empty list", func(t *testing.T) {
		result := p.RenderACLRules([]*calendar.ACLRule{})
		if result != "No ACL rules found" {
			t.Errorf("Expected 'No ACL rules found', got %q", result)
		}
	})

	t.Run("renders nil list", func(t *testing.T) {
		result := p.RenderACLRules(nil)
		if result != "No ACL rules found" {
			t.Errorf("Expected 'No ACL rules found', got %q", result)
		}
	})

	t.Run("skips nil rules in list", func(t *testing.T) {
		rules := []*calendar.ACLRule{
			{
				ID:   "user:alice@example.com",
				Role: "writer",
				Scope: &calendar.ACLScope{
					Type:  "user",
					Value: "alice@example.com",
				},
			},
			nil,
			{
				ID:   "user:bob@example.com",
				Role: "reader",
				Scope: &calendar.ACLScope{
					Type:  "user",
					Value: "bob@example.com",
				},
			},
		}

		result := p.RenderACLRules(rules)

		if !strings.Contains(result, "alice@example.com") {
			t.Error("Result should contain first rule")
		}
		if !strings.Contains(result, "bob@example.com") {
			t.Error("Result should contain third rule")
		}
	})

	t.Run("handles rules with nil scope", func(t *testing.T) {
		rules := []*calendar.ACLRule{
			{
				ID:   "rule1",
				Role: "reader",
			},
		}

		result := p.RenderACLRules(rules)

		if !strings.Contains(result, "rule1") {
			t.Error("Result should contain rule ID")
		}
		if !strings.Contains(result, "reader") {
			t.Error("Result should contain role")
		}
	})

	t.Run("truncates long values", func(t *testing.T) {
		rules := []*calendar.ACLRule{
			{
				ID:   "user:very.long.email.address.that.exceeds.the.truncation.limit@example.com",
				Role: "writer",
				Scope: &calendar.ACLScope{
					Type:  "user",
					Value: "very.long.email.address.that.exceeds.the.truncation.limit@example.com",
				},
			},
		}

		result := p.RenderACLRules(rules)

		// Should still contain some part of the long values
		if !strings.Contains(result, "very") {
			t.Error("Result should contain truncated value")
		}
	})
}
