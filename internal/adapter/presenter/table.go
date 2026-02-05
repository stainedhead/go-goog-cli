package presenter

import (
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/stainedhead/go-goog-cli/internal/domain/account"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	domaincontacts "github.com/stainedhead/go-goog-cli/internal/domain/contacts"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
	domaintasks "github.com/stainedhead/go-goog-cli/internal/domain/tasks"
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

// formatFieldWithTypeAndPrimary formats a value with optional type and primary markers.
func formatFieldWithTypeAndPrimary(value, fieldType string, primary bool) string {
	result := value
	if fieldType != "" {
		result += " (" + fieldType + ")"
	}
	if primary {
		result += " [primary]"
	}
	return result
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

// RenderACLRule renders a single ACL rule as a table.
func (p *TablePresenter) RenderACLRule(rule *calendar.ACLRule) string {
	if rule == nil {
		return "No ACL rule found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"Field", "Value"})

	_ = table.Append([]string{"ID", rule.ID})
	_ = table.Append([]string{"Role", rule.Role})
	if rule.Scope != nil {
		_ = table.Append([]string{"Scope Type", rule.Scope.Type})
		if rule.Scope.Value != "" {
			_ = table.Append([]string{"Scope Value", rule.Scope.Value})
		}
	}

	_ = table.Render()
	return buf.String()
}

// RenderACLRules renders multiple ACL rules as a table.
func (p *TablePresenter) RenderACLRules(rules []*calendar.ACLRule) string {
	if len(rules) == 0 {
		return "No ACL rules found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"ID", "Role", "Scope Type", "Scope Value"})

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
		_ = table.Append([]string{
			truncate(rule.ID, 30),
			rule.Role,
			scopeType,
			truncate(scopeValue, 30),
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

// RenderTaskList renders a single task list as a table.
func (p *TablePresenter) RenderTaskList(taskList *domaintasks.TaskList) string {
	if taskList == nil {
		return "No task list found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"Field", "Value"})

	_ = table.Append([]string{"ID", taskList.ID})
	_ = table.Append([]string{"Title", taskList.Title})
	_ = table.Append([]string{"Updated", taskList.Updated.Format("2006-01-02 15:04")})

	_ = table.Render()
	return buf.String()
}

// RenderTaskLists renders multiple task lists as a table.
func (p *TablePresenter) RenderTaskLists(taskLists []*domaintasks.TaskList) string {
	if len(taskLists) == 0 {
		return "No task lists found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"ID", "Title", "Updated"})

	for _, tl := range taskLists {
		if tl == nil {
			continue
		}
		_ = table.Append([]string{
			truncate(tl.ID, 30),
			truncate(tl.Title, 40),
			tl.Updated.Format("2006-01-02 15:04"),
		})
	}

	_ = table.Render()
	return buf.String()
}

// RenderTask renders a single task as a table.
func (p *TablePresenter) RenderTask(task *domaintasks.Task) string {
	if task == nil {
		return "No task found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"Field", "Value"})

	_ = table.Append([]string{"ID", task.ID})
	_ = table.Append([]string{"Title", task.Title})
	_ = table.Append([]string{"Status", task.Status})
	if task.Notes != "" {
		_ = table.Append([]string{"Notes", truncate(task.Notes, 60)})
	}
	if task.Due != nil {
		_ = table.Append([]string{"Due", task.Due.Format("2006-01-02")})
	}
	if task.Completed != nil {
		_ = table.Append([]string{"Completed", task.Completed.Format("2006-01-02 15:04")})
	}
	if task.Parent != nil {
		_ = table.Append([]string{"Parent", *task.Parent})
	}
	_ = table.Append([]string{"Updated", task.Updated.Format("2006-01-02 15:04")})

	_ = table.Render()
	return buf.String()
}

// RenderTasks renders multiple tasks as a table.
func (p *TablePresenter) RenderTasks(tasks []*domaintasks.Task) string {
	if len(tasks) == 0 {
		return "No tasks found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"ID", "Title", "Status", "Due", "Updated"})

	for _, task := range tasks {
		if task == nil {
			continue
		}
		dueStr := ""
		if task.Due != nil {
			dueStr = task.Due.Format("2006-01-02")
		}
		_ = table.Append([]string{
			truncate(task.ID, 20),
			truncate(task.Title, 40),
			task.Status,
			dueStr,
			task.Updated.Format("2006-01-02 15:04"),
		})
	}

	_ = table.Render()
	return buf.String()
}

// RenderContact renders a single contact as a table.
func (p *TablePresenter) RenderContact(contact *domaincontacts.Contact) string {
	if contact == nil {
		return "No contact found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"Field", "Value"})

	_ = table.Append([]string{"ResourceName", contact.ResourceName})
	_ = table.Append([]string{"Name", contact.GetDisplayName()})

	if len(contact.EmailAddresses) > 0 {
		var emails []string
		for _, email := range contact.EmailAddresses {
			emails = append(emails, formatFieldWithTypeAndPrimary(email.Value, email.Type, email.Primary))
		}
		_ = table.Append([]string{"Emails", strings.Join(emails, "\n")})
	}

	if len(contact.PhoneNumbers) > 0 {
		var phones []string
		for _, phone := range contact.PhoneNumbers {
			phones = append(phones, formatFieldWithTypeAndPrimary(phone.Value, phone.Type, phone.Primary))
		}
		_ = table.Append([]string{"Phones", strings.Join(phones, "\n")})
	}

	if len(contact.Addresses) > 0 {
		var addresses []string
		for _, addr := range contact.Addresses {
			if addr.FormattedValue != "" {
				typeStr := ""
				if addr.Type != "" {
					typeStr = " (" + addr.Type + ")"
				}
				addresses = append(addresses, addr.FormattedValue+typeStr)
			}
		}
		if len(addresses) > 0 {
			_ = table.Append([]string{"Addresses", strings.Join(addresses, "\n")})
		}
	}

	if len(contact.Organizations) > 0 {
		var orgs []string
		for _, org := range contact.Organizations {
			orgStr := org.Name
			if org.Title != "" {
				orgStr += " - " + org.Title
			}
			if org.Department != "" {
				orgStr += " (" + org.Department + ")"
			}
			orgs = append(orgs, orgStr)
		}
		_ = table.Append([]string{"Organizations", strings.Join(orgs, "\n")})
	}

	if len(contact.Birthdays) > 0 {
		if contact.Birthdays[0].Date != nil {
			_ = table.Append([]string{"Birthday", contact.Birthdays[0].Date.FormatDate()})
		} else if contact.Birthdays[0].Text != "" {
			_ = table.Append([]string{"Birthday", contact.Birthdays[0].Text})
		}
	}

	if len(contact.Biographies) > 0 && contact.Biographies[0].Value != "" {
		_ = table.Append([]string{"Biography", truncate(contact.Biographies[0].Value, 60)})
	}

	if contact.ETag != "" {
		_ = table.Append([]string{"ETag", contact.ETag})
	}

	_ = table.Render()
	return buf.String()
}

// RenderContacts renders multiple contacts as a table.
func (p *TablePresenter) RenderContacts(contacts []*domaincontacts.Contact) string {
	if len(contacts) == 0 {
		return "No contacts found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"ResourceName", "Name", "Email", "Phone"})

	for _, c := range contacts {
		if c == nil {
			continue
		}

		name := c.GetDisplayName()

		email, _ := c.GetPrimaryEmail()

		phone := ""
		if len(c.PhoneNumbers) > 0 {
			for _, p := range c.PhoneNumbers {
				if p.Primary {
					phone = p.Value
					break
				}
			}
			if phone == "" {
				phone = c.PhoneNumbers[0].Value
			}
		}

		_ = table.Append([]string{
			truncate(c.ResourceName, 25),
			truncate(name, 30),
			truncate(email, 30),
			truncate(phone, 20),
		})
	}

	_ = table.Render()
	return buf.String()
}

// RenderContactGroup renders a single contact group as a table.
func (p *TablePresenter) RenderContactGroup(group *domaincontacts.ContactGroup) string {
	if group == nil {
		return "No contact group found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"Field", "Value"})

	_ = table.Append([]string{"ResourceName", group.ResourceName})
	_ = table.Append([]string{"Name", group.Name})
	if group.FormattedName != "" {
		_ = table.Append([]string{"Formatted Name", group.FormattedName})
	}
	_ = table.Append([]string{"Type", group.GroupType})
	_ = table.Append([]string{"Member Count", fmt.Sprintf("%d", group.MemberCount)})
	if group.ETag != "" {
		_ = table.Append([]string{"ETag", group.ETag})
	}
	if group.Metadata != nil && !group.Metadata.UpdateTime.IsZero() {
		_ = table.Append([]string{"Updated", group.Metadata.UpdateTime.Format("2006-01-02 15:04")})
	}

	_ = table.Render()
	return buf.String()
}

// RenderContactGroups renders multiple contact groups as a table.
func (p *TablePresenter) RenderContactGroups(groups []*domaincontacts.ContactGroup) string {
	if len(groups) == 0 {
		return "No contact groups found"
	}

	var buf strings.Builder
	table := createTable(&buf, []string{"ResourceName", "Name", "Type", "MemberCount"})

	for _, g := range groups {
		if g == nil {
			continue
		}
		_ = table.Append([]string{
			truncate(g.ResourceName, 30),
			truncate(g.Name, 30),
			g.GroupType,
			fmt.Sprintf("%d", g.MemberCount),
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
