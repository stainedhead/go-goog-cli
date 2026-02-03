// Package presenter provides output formatting for CLI commands.
package presenter

import (
	"github.com/stainedhead/go-goog-cli/internal/domain/account"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
)

// Format constants for presenter output types.
const (
	FormatJSON  = "json"
	FormatTable = "table"
	FormatPlain = "plain"
)

// Presenter defines the interface for rendering domain entities as formatted output.
type Presenter interface {
	// Mail entities
	RenderMessage(msg *mail.Message) string
	RenderMessages(msgs []*mail.Message) string
	RenderDraft(draft *mail.Draft) string
	RenderDrafts(drafts []*mail.Draft) string
	RenderThread(thread *mail.Thread) string
	RenderThreads(threads []*mail.Thread) string
	RenderLabel(label *mail.Label) string
	RenderLabels(labels []*mail.Label) string

	// Calendar entities
	RenderEvent(event *calendar.Event) string
	RenderEvents(events []*calendar.Event) string
	RenderCalendar(cal *calendar.Calendar) string
	RenderCalendars(cals []*calendar.Calendar) string

	// Account
	RenderAccount(acct *account.Account) string
	RenderAccounts(accts []*account.Account) string

	// Generic
	RenderError(err error) string
	RenderSuccess(msg string) string
}

// New creates a new Presenter based on the specified format.
// Supported formats: "json", "table", "plain".
// Returns a TablePresenter as the default if the format is not recognized.
func New(format string) Presenter {
	switch format {
	case FormatJSON:
		return NewJSONPresenter()
	case FormatTable:
		return NewTablePresenter()
	case FormatPlain:
		return NewPlainPresenter()
	default:
		return NewTablePresenter()
	}
}
