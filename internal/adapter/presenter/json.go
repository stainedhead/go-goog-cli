package presenter

import (
	"encoding/json"

	"github.com/stainedhead/go-goog-cli/internal/domain/account"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	domaincontacts "github.com/stainedhead/go-goog-cli/internal/domain/contacts"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
	domaintasks "github.com/stainedhead/go-goog-cli/internal/domain/tasks"
)

// JSONPresenter formats output as indented JSON.
type JSONPresenter struct{}

// NewJSONPresenter creates a new JSONPresenter.
func NewJSONPresenter() *JSONPresenter {
	return &JSONPresenter{}
}

// marshalJSON marshals v to indented JSON, returning an empty object on error.
func (p *JSONPresenter) marshalJSON(v interface{}) string {
	if v == nil {
		return "null"
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(data)
}

// RenderMessage renders a single message as JSON.
func (p *JSONPresenter) RenderMessage(msg *mail.Message) string {
	return p.marshalJSON(msg)
}

// RenderMessages renders multiple messages as JSON.
func (p *JSONPresenter) RenderMessages(msgs []*mail.Message) string {
	if msgs == nil {
		return "[]"
	}
	return p.marshalJSON(msgs)
}

// RenderDraft renders a single draft as JSON.
func (p *JSONPresenter) RenderDraft(draft *mail.Draft) string {
	return p.marshalJSON(draft)
}

// RenderDrafts renders multiple drafts as JSON.
func (p *JSONPresenter) RenderDrafts(drafts []*mail.Draft) string {
	if drafts == nil {
		return "[]"
	}
	return p.marshalJSON(drafts)
}

// RenderThread renders a single thread as JSON.
func (p *JSONPresenter) RenderThread(thread *mail.Thread) string {
	return p.marshalJSON(thread)
}

// RenderThreads renders multiple threads as JSON.
func (p *JSONPresenter) RenderThreads(threads []*mail.Thread) string {
	if threads == nil {
		return "[]"
	}
	return p.marshalJSON(threads)
}

// RenderLabel renders a single label as JSON.
func (p *JSONPresenter) RenderLabel(label *mail.Label) string {
	return p.marshalJSON(label)
}

// RenderLabels renders multiple labels as JSON.
func (p *JSONPresenter) RenderLabels(labels []*mail.Label) string {
	if labels == nil {
		return "[]"
	}
	return p.marshalJSON(labels)
}

// RenderEvent renders a single event as JSON.
func (p *JSONPresenter) RenderEvent(event *calendar.Event) string {
	return p.marshalJSON(event)
}

// RenderEvents renders multiple events as JSON.
func (p *JSONPresenter) RenderEvents(events []*calendar.Event) string {
	if events == nil {
		return "[]"
	}
	return p.marshalJSON(events)
}

// RenderCalendar renders a single calendar as JSON.
func (p *JSONPresenter) RenderCalendar(cal *calendar.Calendar) string {
	return p.marshalJSON(cal)
}

// RenderCalendars renders multiple calendars as JSON.
func (p *JSONPresenter) RenderCalendars(cals []*calendar.Calendar) string {
	if cals == nil {
		return "[]"
	}
	return p.marshalJSON(cals)
}

// RenderACLRule renders a single ACL rule as JSON.
func (p *JSONPresenter) RenderACLRule(rule *calendar.ACLRule) string {
	return p.marshalJSON(rule)
}

// RenderACLRules renders multiple ACL rules as JSON.
func (p *JSONPresenter) RenderACLRules(rules []*calendar.ACLRule) string {
	if rules == nil {
		return "[]"
	}
	return p.marshalJSON(rules)
}

// RenderAccount renders a single account as JSON.
func (p *JSONPresenter) RenderAccount(acct *account.Account) string {
	return p.marshalJSON(acct)
}

// RenderAccounts renders multiple accounts as JSON.
func (p *JSONPresenter) RenderAccounts(accts []*account.Account) string {
	if accts == nil {
		return "[]"
	}
	return p.marshalJSON(accts)
}

// RenderTaskList renders a single task list as JSON.
func (p *JSONPresenter) RenderTaskList(taskList *domaintasks.TaskList) string {
	return p.marshalJSON(taskList)
}

// RenderTaskLists renders multiple task lists as JSON.
func (p *JSONPresenter) RenderTaskLists(taskLists []*domaintasks.TaskList) string {
	if taskLists == nil {
		return "[]"
	}
	return p.marshalJSON(taskLists)
}

// RenderTask renders a single task as JSON.
func (p *JSONPresenter) RenderTask(task *domaintasks.Task) string {
	return p.marshalJSON(task)
}

// RenderTasks renders multiple tasks as JSON.
func (p *JSONPresenter) RenderTasks(tasks []*domaintasks.Task) string {
	if tasks == nil {
		return "[]"
	}
	return p.marshalJSON(tasks)
}

// errorResponse is the JSON structure for error output.
type errorResponse struct {
	Error string `json:"error"`
}

// successResponse is the JSON structure for success output.
type successResponse struct {
	Message string `json:"message"`
}

// RenderContact renders a single contact as JSON.
func (p *JSONPresenter) RenderContact(contact *domaincontacts.Contact) string {
	return p.marshalJSON(contact)
}

// RenderContacts renders multiple contacts as JSON.
func (p *JSONPresenter) RenderContacts(contacts []*domaincontacts.Contact) string {
	if contacts == nil {
		return "[]"
	}
	return p.marshalJSON(contacts)
}

// RenderContactGroup renders a single contact group as JSON.
func (p *JSONPresenter) RenderContactGroup(group *domaincontacts.ContactGroup) string {
	return p.marshalJSON(group)
}

// RenderContactGroups renders multiple contact groups as JSON.
func (p *JSONPresenter) RenderContactGroups(groups []*domaincontacts.ContactGroup) string {
	if groups == nil {
		return "[]"
	}
	return p.marshalJSON(groups)
}

// RenderError renders an error as JSON.
func (p *JSONPresenter) RenderError(err error) string {
	if err == nil {
		return p.marshalJSON(errorResponse{Error: ""})
	}
	return p.marshalJSON(errorResponse{Error: err.Error()})
}

// RenderSuccess renders a success message as JSON.
func (p *JSONPresenter) RenderSuccess(msg string) string {
	return p.marshalJSON(successResponse{Message: msg})
}
