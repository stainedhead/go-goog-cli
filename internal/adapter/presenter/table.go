package presenter

import (
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/stainedhead/go-goog-cli/internal/domain/account"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
)

// TablePresenter formats output as ASCII tables.
type TablePresenter struct{}

// NewTablePresenter creates a new TablePresenter.
func NewTablePresenter() *TablePresenter {
	return &TablePresenter{}
}

// truncate shortens s to maxLen characters, appending "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// createTable creates a new tablewriter with standard settings.
func createTable(buf *strings.Builder, headers []string) *tablewriter.Table {
	table := tablewriter.NewTable(buf)
	table.Header(headers)
	return table
}

// RenderMessage renders a single message as a table.
func (p *TablePresenter) RenderMessage(msg *mail.Message) string {
	if msg == nil {
		return "No message found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"Field", "Value"})

	_ = table.Append([]string{"ID", msg.ID})
	_ = table.Append([]string{"Thread ID", msg.ThreadID})
	_ = table.Append([]string{"From", msg.From})
	_ = table.Append([]string{"To", strings.Join(msg.To, ", ")})
	if len(msg.Cc) > 0 {
		_ = table.Append([]string{"Cc", strings.Join(msg.Cc, ", ")})
	}
	_ = table.Append([]string{"Subject", msg.Subject})
	_ = table.Append([]string{"Date", msg.Date.Format("2006-01-02 15:04")})
	_ = table.Append([]string{"Labels", strings.Join(msg.Labels, ", ")})
	_ = table.Append([]string{"Read", fmt.Sprintf("%v", msg.IsRead)})
	_ = table.Append([]string{"Starred", fmt.Sprintf("%v", msg.IsStarred)})
	if msg.Snippet != "" {
		_ = table.Append([]string{"Snippet", truncate(msg.Snippet, 60)})
	}

	_ = table.Render()
	return buf.String()
}

// RenderMessages renders multiple messages as a table.
func (p *TablePresenter) RenderMessages(msgs []*mail.Message) string {
	if len(msgs) == 0 {
		return "No messages found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"ID", "From", "Subject", "Date", "Labels"})

	for _, msg := range msgs {
		if msg == nil {
			continue
		}
		_ = table.Append([]string{
			truncate(msg.ID, 12),
			truncate(msg.From, 25),
			truncate(msg.Subject, 40),
			msg.Date.Format("2006-01-02"),
			truncate(strings.Join(msg.Labels, ", "), 20),
		})
	}

	_ = table.Render()
	return buf.String()
}

// RenderDraft renders a single draft as a table.
func (p *TablePresenter) RenderDraft(draft *mail.Draft) string {
	if draft == nil {
		return "No draft found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"Field", "Value"})

	_ = table.Append([]string{"Draft ID", draft.ID})
	_ = table.Append([]string{"Created", draft.Created.Format("2006-01-02 15:04")})
	_ = table.Append([]string{"Updated", draft.Updated.Format("2006-01-02 15:04")})

	if draft.Message != nil {
		_ = table.Append([]string{"Message ID", draft.Message.ID})
		_ = table.Append([]string{"To", strings.Join(draft.Message.To, ", ")})
		_ = table.Append([]string{"Subject", draft.Message.Subject})
	}

	_ = table.Render()
	return buf.String()
}

// RenderDrafts renders multiple drafts as a table.
func (p *TablePresenter) RenderDrafts(drafts []*mail.Draft) string {
	if len(drafts) == 0 {
		return "No drafts found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"ID", "Subject", "To", "Updated"})

	for _, draft := range drafts {
		if draft == nil {
			continue
		}
		subject := ""
		to := ""
		if draft.Message != nil {
			subject = draft.Message.Subject
			to = strings.Join(draft.Message.To, ", ")
		}
		_ = table.Append([]string{
			truncate(draft.ID, 12),
			truncate(subject, 40),
			truncate(to, 25),
			draft.Updated.Format("2006-01-02"),
		})
	}

	_ = table.Render()
	return buf.String()
}

// RenderThread renders a single thread as a table.
func (p *TablePresenter) RenderThread(thread *mail.Thread) string {
	if thread == nil {
		return "No thread found"
	}

	var buf strings.Builder

	// Thread info table
	infoTable := createTable(&buf, []string{"Field", "Value"})
	_ = infoTable.Append([]string{"Thread ID", thread.ID})
	_ = infoTable.Append([]string{"Message Count", fmt.Sprintf("%d", thread.MessageCount())})
	_ = infoTable.Append([]string{"Labels", strings.Join(thread.Labels, ", ")})
	if thread.Snippet != "" {
		_ = infoTable.Append([]string{"Snippet", truncate(thread.Snippet, 60)})
	}
	_ = infoTable.Render()

	// Messages table if present
	if len(thread.Messages) > 0 {
		buf.WriteString("\nMessages:\n")
		msgTable := createTable(&buf, []string{"ID", "From", "Subject", "Date"})
		for _, msg := range thread.Messages {
			if msg == nil {
				continue
			}
			_ = msgTable.Append([]string{
				truncate(msg.ID, 12),
				truncate(msg.From, 25),
				truncate(msg.Subject, 40),
				msg.Date.Format("2006-01-02"),
			})
		}
		_ = msgTable.Render()
	}

	return buf.String()
}

// RenderThreads renders multiple threads as a table.
func (p *TablePresenter) RenderThreads(threads []*mail.Thread) string {
	if len(threads) == 0 {
		return "No threads found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"ID", "Messages", "Snippet", "Labels"})

	for _, thread := range threads {
		if thread == nil {
			continue
		}
		_ = table.Append([]string{
			truncate(thread.ID, 12),
			fmt.Sprintf("%d", thread.MessageCount()),
			truncate(thread.Snippet, 40),
			truncate(strings.Join(thread.Labels, ", "), 20),
		})
	}

	_ = table.Render()
	return buf.String()
}

// RenderLabel renders a single label as a table.
func (p *TablePresenter) RenderLabel(label *mail.Label) string {
	if label == nil {
		return "No label found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"Field", "Value"})

	_ = table.Append([]string{"ID", label.ID})
	_ = table.Append([]string{"Name", label.Name})
	_ = table.Append([]string{"Type", label.Type})
	if label.MessageListVisibility != "" {
		_ = table.Append([]string{"Message Visibility", label.MessageListVisibility})
	}
	if label.LabelListVisibility != "" {
		_ = table.Append([]string{"Label Visibility", label.LabelListVisibility})
	}
	if label.Color != nil {
		_ = table.Append([]string{"Background", label.Color.Background})
		_ = table.Append([]string{"Text Color", label.Color.Text})
	}

	_ = table.Render()
	return buf.String()
}

// RenderLabels renders multiple labels as a table.
func (p *TablePresenter) RenderLabels(labels []*mail.Label) string {
	if len(labels) == 0 {
		return "No labels found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"ID", "Name", "Type"})

	for _, label := range labels {
		if label == nil {
			continue
		}
		_ = table.Append([]string{
			label.ID,
			label.Name,
			label.Type,
		})
	}

	_ = table.Render()
	return buf.String()
}

// RenderEvent renders a single event as a table.
func (p *TablePresenter) RenderEvent(event *calendar.Event) string {
	if event == nil {
		return "No event found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"Field", "Value"})

	_ = table.Append([]string{"ID", event.ID})
	_ = table.Append([]string{"Title", event.Title})
	if event.Description != "" {
		_ = table.Append([]string{"Description", truncate(event.Description, 60)})
	}
	if event.Location != "" {
		_ = table.Append([]string{"Location", event.Location})
	}

	if event.AllDay {
		_ = table.Append([]string{"Date", event.Start.Format("2006-01-02") + " (All Day)"})
	} else {
		_ = table.Append([]string{"Start", event.Start.Format("2006-01-02 15:04")})
		_ = table.Append([]string{"End", event.End.Format("2006-01-02 15:04")})
	}

	_ = table.Append([]string{"Status", event.Status})
	if event.CalendarID != "" {
		_ = table.Append([]string{"Calendar", event.CalendarID})
	}
	if event.HasConference() {
		_ = table.Append([]string{"Conference", event.ConferenceData.URI})
	}

	_ = table.Render()
	return buf.String()
}

// RenderEvents renders multiple events as a table.
func (p *TablePresenter) RenderEvents(events []*calendar.Event) string {
	if len(events) == 0 {
		return "No events found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"ID", "Title", "Start", "End", "Location"})

	for _, event := range events {
		if event == nil {
			continue
		}
		startStr := event.Start.Format("2006-01-02 15:04")
		endStr := event.End.Format("2006-01-02 15:04")
		if event.AllDay {
			startStr = event.Start.Format("2006-01-02")
			endStr = "(All Day)"
		}
		_ = table.Append([]string{
			truncate(event.ID, 12),
			truncate(event.Title, 30),
			startStr,
			endStr,
			truncate(event.Location, 20),
		})
	}

	_ = table.Render()
	return buf.String()
}

// RenderCalendar renders a single calendar as a table.
func (p *TablePresenter) RenderCalendar(cal *calendar.Calendar) string {
	if cal == nil {
		return "No calendar found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"Field", "Value"})

	_ = table.Append([]string{"ID", cal.ID})
	_ = table.Append([]string{"Title", cal.Title})
	if cal.Description != "" {
		_ = table.Append([]string{"Description", truncate(cal.Description, 60)})
	}
	if cal.TimeZone != "" {
		_ = table.Append([]string{"Time Zone", cal.TimeZone})
	}
	_ = table.Append([]string{"Primary", fmt.Sprintf("%v", cal.Primary)})
	_ = table.Append([]string{"Access Role", cal.AccessRole})

	_ = table.Render()
	return buf.String()
}

// RenderCalendars renders multiple calendars as a table.
func (p *TablePresenter) RenderCalendars(cals []*calendar.Calendar) string {
	if len(cals) == 0 {
		return "No calendars found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"ID", "Title", "Primary", "Access Role", "Time Zone"})

	for _, cal := range cals {
		if cal == nil {
			continue
		}
		primary := ""
		if cal.Primary {
			primary = "Yes"
		}
		_ = table.Append([]string{
			truncate(cal.ID, 20),
			truncate(cal.Title, 25),
			primary,
			cal.AccessRole,
			cal.TimeZone,
		})
	}

	_ = table.Render()
	return buf.String()
}

// RenderAccount renders a single account as a table.
func (p *TablePresenter) RenderAccount(acct *account.Account) string {
	if acct == nil {
		return "No account found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"Field", "Value"})

	_ = table.Append([]string{"Alias", acct.Alias})
	_ = table.Append([]string{"Email", acct.Email})
	_ = table.Append([]string{"Default", fmt.Sprintf("%v", acct.IsDefault)})
	_ = table.Append([]string{"Scopes", fmt.Sprintf("%d", len(acct.Scopes))})
	_ = table.Append([]string{"Added", acct.Added.Format("2006-01-02")})
	if !acct.LastUsed.IsZero() {
		_ = table.Append([]string{"Last Used", acct.LastUsed.Format("2006-01-02")})
	}

	if len(acct.Scopes) > 0 {
		_ = table.Append([]string{"Scope List", strings.Join(acct.Scopes, "\n")})
	}

	_ = table.Render()
	return buf.String()
}

// RenderAccounts renders multiple accounts as a table.
func (p *TablePresenter) RenderAccounts(accts []*account.Account) string {
	if len(accts) == 0 {
		return "No accounts found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"Alias", "Email", "Default", "Scopes"})

	for _, acct := range accts {
		if acct == nil {
			continue
		}
		isDefault := ""
		if acct.IsDefault {
			isDefault = "Yes"
		}
		_ = table.Append([]string{
			acct.Alias,
			acct.Email,
			isDefault,
			fmt.Sprintf("%d", len(acct.Scopes)),
		})
	}

	_ = table.Render()
	return buf.String()
}

// RenderError renders an error as formatted text.
func (p *TablePresenter) RenderError(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("Error: %s", err.Error())
}

// RenderSuccess renders a success message as formatted text.
func (p *TablePresenter) RenderSuccess(msg string) string {
	return fmt.Sprintf("Success: %s", msg)
}
