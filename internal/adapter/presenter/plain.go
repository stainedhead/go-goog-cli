package presenter

import (
	"fmt"
	"strings"

	"github.com/stainedhead/go-goog-cli/internal/domain/account"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
)

// PlainPresenter formats output as plain text, suitable for piping.
type PlainPresenter struct{}

// NewPlainPresenter creates a new PlainPresenter.
func NewPlainPresenter() *PlainPresenter {
	return &PlainPresenter{}
}

// RenderMessage renders a single message as key-value pairs.
func (p *PlainPresenter) RenderMessage(msg *mail.Message) string {
	if msg == nil {
		return ""
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("ID: %s", msg.ID))
	lines = append(lines, fmt.Sprintf("ThreadID: %s", msg.ThreadID))
	lines = append(lines, fmt.Sprintf("From: %s", msg.From))
	lines = append(lines, fmt.Sprintf("To: %s", strings.Join(msg.To, ", ")))
	if len(msg.Cc) > 0 {
		lines = append(lines, fmt.Sprintf("Cc: %s", strings.Join(msg.Cc, ", ")))
	}
	if len(msg.Bcc) > 0 {
		lines = append(lines, fmt.Sprintf("Bcc: %s", strings.Join(msg.Bcc, ", ")))
	}
	lines = append(lines, fmt.Sprintf("Subject: %s", msg.Subject))
	lines = append(lines, fmt.Sprintf("Date: %s", msg.Date.Format("2006-01-02 15:04:05")))
	lines = append(lines, fmt.Sprintf("Labels: %s", strings.Join(msg.Labels, ", ")))
	lines = append(lines, fmt.Sprintf("Read: %v", msg.IsRead))
	lines = append(lines, fmt.Sprintf("Starred: %v", msg.IsStarred))
	if msg.Snippet != "" {
		lines = append(lines, fmt.Sprintf("Snippet: %s", msg.Snippet))
	}
	if msg.Body != "" {
		lines = append(lines, fmt.Sprintf("Body: %s", msg.Body))
	}

	return strings.Join(lines, "\n")
}

// RenderMessages renders multiple messages, one per line.
func (p *PlainPresenter) RenderMessages(msgs []*mail.Message) string {
	if len(msgs) == 0 {
		return ""
	}

	var lines []string
	for _, msg := range msgs {
		if msg == nil {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s\t%s",
			msg.ID,
			msg.From,
			msg.Subject,
			msg.Date.Format("2006-01-02"),
		))
	}
	return strings.Join(lines, "\n")
}

// RenderDraft renders a single draft as key-value pairs.
func (p *PlainPresenter) RenderDraft(draft *mail.Draft) string {
	if draft == nil {
		return ""
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("ID: %s", draft.ID))
	lines = append(lines, fmt.Sprintf("Created: %s", draft.Created.Format("2006-01-02 15:04:05")))
	lines = append(lines, fmt.Sprintf("Updated: %s", draft.Updated.Format("2006-01-02 15:04:05")))

	if draft.Message != nil {
		lines = append(lines, fmt.Sprintf("MessageID: %s", draft.Message.ID))
		lines = append(lines, fmt.Sprintf("To: %s", strings.Join(draft.Message.To, ", ")))
		lines = append(lines, fmt.Sprintf("Subject: %s", draft.Message.Subject))
	}

	return strings.Join(lines, "\n")
}

// RenderDrafts renders multiple drafts, one per line.
func (p *PlainPresenter) RenderDrafts(drafts []*mail.Draft) string {
	if len(drafts) == 0 {
		return ""
	}

	var lines []string
	for _, draft := range drafts {
		if draft == nil {
			continue
		}
		subject := ""
		if draft.Message != nil {
			subject = draft.Message.Subject
		}
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s",
			draft.ID,
			subject,
			draft.Updated.Format("2006-01-02"),
		))
	}
	return strings.Join(lines, "\n")
}

// RenderThread renders a single thread as key-value pairs.
func (p *PlainPresenter) RenderThread(thread *mail.Thread) string {
	if thread == nil {
		return ""
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("ID: %s", thread.ID))
	lines = append(lines, fmt.Sprintf("Messages: %d", thread.MessageCount()))
	lines = append(lines, fmt.Sprintf("Labels: %s", strings.Join(thread.Labels, ", ")))
	if thread.Snippet != "" {
		lines = append(lines, fmt.Sprintf("Snippet: %s", thread.Snippet))
	}

	for i, msg := range thread.Messages {
		if msg == nil {
			continue
		}
		lines = append(lines, fmt.Sprintf("Message[%d]: %s from %s: %s",
			i, msg.ID, msg.From, msg.Subject))
	}

	return strings.Join(lines, "\n")
}

// RenderThreads renders multiple threads, one per line.
func (p *PlainPresenter) RenderThreads(threads []*mail.Thread) string {
	if len(threads) == 0 {
		return ""
	}

	var lines []string
	for _, thread := range threads {
		if thread == nil {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s\t%d\t%s",
			thread.ID,
			thread.MessageCount(),
			thread.Snippet,
		))
	}
	return strings.Join(lines, "\n")
}

// RenderLabel renders a single label as key-value pairs.
func (p *PlainPresenter) RenderLabel(label *mail.Label) string {
	if label == nil {
		return ""
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("ID: %s", label.ID))
	lines = append(lines, fmt.Sprintf("Name: %s", label.Name))
	lines = append(lines, fmt.Sprintf("Type: %s", label.Type))
	if label.MessageListVisibility != "" {
		lines = append(lines, fmt.Sprintf("MessageVisibility: %s", label.MessageListVisibility))
	}
	if label.LabelListVisibility != "" {
		lines = append(lines, fmt.Sprintf("LabelVisibility: %s", label.LabelListVisibility))
	}
	if label.Color != nil {
		lines = append(lines, fmt.Sprintf("Background: %s", label.Color.Background))
		lines = append(lines, fmt.Sprintf("TextColor: %s", label.Color.Text))
	}

	return strings.Join(lines, "\n")
}

// RenderLabels renders multiple labels, one per line.
func (p *PlainPresenter) RenderLabels(labels []*mail.Label) string {
	if len(labels) == 0 {
		return ""
	}

	var lines []string
	for _, label := range labels {
		if label == nil {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s",
			label.ID,
			label.Name,
			label.Type,
		))
	}
	return strings.Join(lines, "\n")
}

// RenderEvent renders a single event as key-value pairs.
func (p *PlainPresenter) RenderEvent(event *calendar.Event) string {
	if event == nil {
		return ""
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("ID: %s", event.ID))
	lines = append(lines, fmt.Sprintf("Title: %s", event.Title))
	if event.Description != "" {
		lines = append(lines, fmt.Sprintf("Description: %s", event.Description))
	}
	if event.Location != "" {
		lines = append(lines, fmt.Sprintf("Location: %s", event.Location))
	}

	if event.AllDay {
		lines = append(lines, fmt.Sprintf("Date: %s (All Day)", event.Start.Format("2006-01-02")))
	} else {
		lines = append(lines, fmt.Sprintf("Start: %s", event.Start.Format("2006-01-02 15:04")))
		lines = append(lines, fmt.Sprintf("End: %s", event.End.Format("2006-01-02 15:04")))
	}

	lines = append(lines, fmt.Sprintf("Status: %s", event.Status))
	if event.CalendarID != "" {
		lines = append(lines, fmt.Sprintf("Calendar: %s", event.CalendarID))
	}
	if event.HasConference() {
		lines = append(lines, fmt.Sprintf("Conference: %s", event.ConferenceData.URI))
	}
	if event.HTMLLink != "" {
		lines = append(lines, fmt.Sprintf("Link: %s", event.HTMLLink))
	}

	return strings.Join(lines, "\n")
}

// RenderEvents renders multiple events, one per line.
func (p *PlainPresenter) RenderEvents(events []*calendar.Event) string {
	if len(events) == 0 {
		return ""
	}

	var lines []string
	for _, event := range events {
		if event == nil {
			continue
		}
		timeStr := event.Start.Format("2006-01-02 15:04")
		if event.AllDay {
			timeStr = event.Start.Format("2006-01-02") + " (All Day)"
		}
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s\t%s",
			event.ID,
			event.Title,
			timeStr,
			event.Location,
		))
	}
	return strings.Join(lines, "\n")
}

// RenderCalendar renders a single calendar as key-value pairs.
func (p *PlainPresenter) RenderCalendar(cal *calendar.Calendar) string {
	if cal == nil {
		return ""
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("ID: %s", cal.ID))
	lines = append(lines, fmt.Sprintf("Title: %s", cal.Title))
	if cal.Description != "" {
		lines = append(lines, fmt.Sprintf("Description: %s", cal.Description))
	}
	if cal.TimeZone != "" {
		lines = append(lines, fmt.Sprintf("TimeZone: %s", cal.TimeZone))
	}
	lines = append(lines, fmt.Sprintf("Primary: %v", cal.Primary))
	lines = append(lines, fmt.Sprintf("AccessRole: %s", cal.AccessRole))

	return strings.Join(lines, "\n")
}

// RenderCalendars renders multiple calendars, one per line.
func (p *PlainPresenter) RenderCalendars(cals []*calendar.Calendar) string {
	if len(cals) == 0 {
		return ""
	}

	var lines []string
	for _, cal := range cals {
		if cal == nil {
			continue
		}
		primary := ""
		if cal.Primary {
			primary = "*"
		}
		lines = append(lines, fmt.Sprintf("%s%s\t%s\t%s",
			primary,
			cal.ID,
			cal.Title,
			cal.AccessRole,
		))
	}
	return strings.Join(lines, "\n")
}

// RenderACLRule renders a single ACL rule as key-value pairs.
func (p *PlainPresenter) RenderACLRule(rule *calendar.ACLRule) string {
	if rule == nil {
		return ""
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("ID: %s", rule.ID))
	lines = append(lines, fmt.Sprintf("Role: %s", rule.Role))
	if rule.Scope != nil {
		lines = append(lines, fmt.Sprintf("ScopeType: %s", rule.Scope.Type))
		if rule.Scope.Value != "" {
			lines = append(lines, fmt.Sprintf("ScopeValue: %s", rule.Scope.Value))
		}
	}

	return strings.Join(lines, "\n")
}

// RenderACLRules renders multiple ACL rules, one per line.
func (p *PlainPresenter) RenderACLRules(rules []*calendar.ACLRule) string {
	if len(rules) == 0 {
		return ""
	}

	var lines []string
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		scopeType := ""
		scopeValue := ""
		if rule.Scope != nil {
			scopeType = rule.Scope.Type
			scopeValue = rule.Scope.Value
		}
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s\t%s",
			rule.ID,
			rule.Role,
			scopeType,
			scopeValue,
		))
	}
	return strings.Join(lines, "\n")
}

// RenderAccount renders a single account as key-value pairs.
func (p *PlainPresenter) RenderAccount(acct *account.Account) string {
	if acct == nil {
		return ""
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("Alias: %s", acct.Alias))
	lines = append(lines, fmt.Sprintf("Email: %s", acct.Email))
	lines = append(lines, fmt.Sprintf("Default: %v", acct.IsDefault))
	lines = append(lines, fmt.Sprintf("Added: %s", acct.Added.Format("2006-01-02")))
	if !acct.LastUsed.IsZero() {
		lines = append(lines, fmt.Sprintf("LastUsed: %s", acct.LastUsed.Format("2006-01-02")))
	}
	lines = append(lines, fmt.Sprintf("Scopes: %d", len(acct.Scopes)))
	for _, scope := range acct.Scopes {
		lines = append(lines, fmt.Sprintf("  - %s", scope))
	}

	return strings.Join(lines, "\n")
}

// RenderAccounts renders multiple accounts, one per line.
func (p *PlainPresenter) RenderAccounts(accts []*account.Account) string {
	if len(accts) == 0 {
		return ""
	}

	var lines []string
	for _, acct := range accts {
		if acct == nil {
			continue
		}
		defaultMark := ""
		if acct.IsDefault {
			defaultMark = "*"
		}
		lines = append(lines, fmt.Sprintf("%s%s\t%s\t%d",
			defaultMark,
			acct.Alias,
			acct.Email,
			len(acct.Scopes),
		))
	}
	return strings.Join(lines, "\n")
}

// RenderError renders an error as plain text.
func (p *PlainPresenter) RenderError(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("error: %s", err.Error())
}

// RenderSuccess renders a success message as plain text.
func (p *PlainPresenter) RenderSuccess(msg string) string {
	return msg
}
