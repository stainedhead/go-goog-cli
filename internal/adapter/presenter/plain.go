package presenter

import (
	"fmt"
	"strings"

	"github.com/stainedhead/go-goog-cli/internal/domain/account"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	domaincontacts "github.com/stainedhead/go-goog-cli/internal/domain/contacts"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
	domaintasks "github.com/stainedhead/go-goog-cli/internal/domain/tasks"
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

// RenderTaskList renders a single task list as key-value pairs.
func (p *PlainPresenter) RenderTaskList(taskList *domaintasks.TaskList) string {
	if taskList == nil {
		return ""
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("ID: %s", taskList.ID))
	lines = append(lines, fmt.Sprintf("Title: %s", taskList.Title))
	lines = append(lines, fmt.Sprintf("Updated: %s", taskList.Updated.Format("2006-01-02 15:04:05")))

	return strings.Join(lines, "\n")
}

// RenderTaskLists renders multiple task lists, one per line.
func (p *PlainPresenter) RenderTaskLists(taskLists []*domaintasks.TaskList) string {
	if len(taskLists) == 0 {
		return ""
	}

	var lines []string
	for _, tl := range taskLists {
		if tl == nil {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s",
			tl.ID,
			tl.Title,
			tl.Updated.Format("2006-01-02 15:04:05"),
		))
	}
	return strings.Join(lines, "\n")
}

// RenderTask renders a single task as key-value pairs.
func (p *PlainPresenter) RenderTask(task *domaintasks.Task) string {
	if task == nil {
		return ""
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("ID: %s", task.ID))
	lines = append(lines, fmt.Sprintf("Title: %s", task.Title))
	lines = append(lines, fmt.Sprintf("Status: %s", task.Status))
	if task.Notes != "" {
		lines = append(lines, fmt.Sprintf("Notes: %s", task.Notes))
	}
	if task.Due != nil {
		lines = append(lines, fmt.Sprintf("Due: %s", task.Due.Format("2006-01-02")))
	}
	if task.Completed != nil {
		lines = append(lines, fmt.Sprintf("Completed: %s", task.Completed.Format("2006-01-02 15:04:05")))
	}
	if task.Parent != nil {
		lines = append(lines, fmt.Sprintf("Parent: %s", *task.Parent))
	}
	lines = append(lines, fmt.Sprintf("Updated: %s", task.Updated.Format("2006-01-02 15:04:05")))

	return strings.Join(lines, "\n")
}

// RenderTasks renders multiple tasks, one per line.
func (p *PlainPresenter) RenderTasks(tasks []*domaintasks.Task) string {
	if len(tasks) == 0 {
		return ""
	}

	var lines []string
	for _, task := range tasks {
		if task == nil {
			continue
		}
		dueStr := ""
		if task.Due != nil {
			dueStr = task.Due.Format("2006-01-02")
		}
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s\t%s\t%s",
			task.ID,
			task.Title,
			task.Status,
			dueStr,
			task.Updated.Format("2006-01-02 15:04:05"),
		))
	}
	return strings.Join(lines, "\n")
}

// RenderContact renders a single contact.
func (p *PlainPresenter) RenderContact(contact *domaincontacts.Contact) string {
	if contact == nil {
		return ""
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("ResourceName: %s", contact.ResourceName))
	lines = append(lines, fmt.Sprintf("Name: %s", contact.GetDisplayName()))

	if len(contact.EmailAddresses) > 0 {
		lines = append(lines, "Emails:")
		for _, email := range contact.EmailAddresses {
			typeStr := ""
			if email.Type != "" {
				typeStr = " (" + email.Type + ")"
			}
			primaryStr := ""
			if email.Primary {
				primaryStr = " [primary]"
			}
			lines = append(lines, fmt.Sprintf("  - %s%s%s", email.Value, typeStr, primaryStr))
		}
	}

	if len(contact.PhoneNumbers) > 0 {
		lines = append(lines, "Phones:")
		for _, phone := range contact.PhoneNumbers {
			typeStr := ""
			if phone.Type != "" {
				typeStr = " (" + phone.Type + ")"
			}
			primaryStr := ""
			if phone.Primary {
				primaryStr = " [primary]"
			}
			lines = append(lines, fmt.Sprintf("  - %s%s%s", phone.Value, typeStr, primaryStr))
		}
	}

	if len(contact.Addresses) > 0 {
		lines = append(lines, "Addresses:")
		for _, addr := range contact.Addresses {
			if addr.FormattedValue != "" {
				typeStr := ""
				if addr.Type != "" {
					typeStr = " (" + addr.Type + ")"
				}
				lines = append(lines, fmt.Sprintf("  - %s%s", addr.FormattedValue, typeStr))
			}
		}
	}

	if len(contact.Organizations) > 0 {
		lines = append(lines, "Organizations:")
		for _, org := range contact.Organizations {
			orgStr := org.Name
			if org.Title != "" {
				orgStr += " - " + org.Title
			}
			if org.Department != "" {
				orgStr += " (" + org.Department + ")"
			}
			lines = append(lines, fmt.Sprintf("  - %s", orgStr))
		}
	}

	if len(contact.Birthdays) > 0 {
		if contact.Birthdays[0].Date != nil {
			lines = append(lines, fmt.Sprintf("Birthday: %s", contact.Birthdays[0].Date.FormatDate()))
		} else if contact.Birthdays[0].Text != "" {
			lines = append(lines, fmt.Sprintf("Birthday: %s", contact.Birthdays[0].Text))
		}
	}

	if len(contact.Biographies) > 0 && contact.Biographies[0].Value != "" {
		lines = append(lines, fmt.Sprintf("Biography: %s", contact.Biographies[0].Value))
	}

	if len(contact.Nicknames) > 0 {
		var nicknames []string
		for _, nn := range contact.Nicknames {
			nicknames = append(nicknames, nn.Value)
		}
		lines = append(lines, fmt.Sprintf("Nicknames: %s", strings.Join(nicknames, ", ")))
	}

	if len(contact.URLs) > 0 {
		lines = append(lines, "URLs:")
		for _, url := range contact.URLs {
			typeStr := ""
			if url.Type != "" {
				typeStr = " (" + url.Type + ")"
			}
			lines = append(lines, fmt.Sprintf("  - %s%s", url.Value, typeStr))
		}
	}

	return strings.Join(lines, "\n")
}

// RenderContacts renders multiple contacts.
func (p *PlainPresenter) RenderContacts(contacts []*domaincontacts.Contact) string {
	if len(contacts) == 0 {
		return ""
	}

	var lines []string
	for _, c := range contacts {
		if c == nil {
			continue
		}

		name := c.GetDisplayName()

		email := ""
		if len(c.EmailAddresses) > 0 {
			primaryEmail, err := c.GetPrimaryEmail()
			if err == nil {
				email = primaryEmail
			}
		}

		phone := ""
		if len(c.PhoneNumbers) > 0 {
			for _, p := range c.PhoneNumbers {
				if p.Primary {
					phone = p.Value
					break
				}
			}
			if phone == "" && len(c.PhoneNumbers) > 0 {
				phone = c.PhoneNumbers[0].Value
			}
		}

		lines = append(lines, fmt.Sprintf("%s\t%s\t%s\t%s",
			c.ResourceName,
			name,
			email,
			phone,
		))
	}
	return strings.Join(lines, "\n")
}

// RenderContactGroup renders a single contact group.
func (p *PlainPresenter) RenderContactGroup(group *domaincontacts.ContactGroup) string {
	if group == nil {
		return ""
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("ResourceName: %s", group.ResourceName))
	lines = append(lines, fmt.Sprintf("Name: %s", group.Name))
	if group.FormattedName != "" {
		lines = append(lines, fmt.Sprintf("FormattedName: %s", group.FormattedName))
	}
	lines = append(lines, fmt.Sprintf("Type: %s", group.GroupType))
	lines = append(lines, fmt.Sprintf("MemberCount: %d", group.MemberCount))
	if group.Metadata != nil && !group.Metadata.UpdateTime.IsZero() {
		lines = append(lines, fmt.Sprintf("Updated: %s", group.Metadata.UpdateTime.Format("2006-01-02 15:04:05")))
	}

	return strings.Join(lines, "\n")
}

// RenderContactGroups renders multiple contact groups.
func (p *PlainPresenter) RenderContactGroups(groups []*domaincontacts.ContactGroup) string {
	if len(groups) == 0 {
		return ""
	}

	var lines []string
	for _, g := range groups {
		if g == nil {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s\t%d",
			g.ResourceName,
			g.Name,
			g.GroupType,
			g.MemberCount,
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
