package mail

import (
	"context"
	"errors"
)

// Domain errors for mail operations.
var (
	ErrMessageNotFound = errors.New("message not found")
	ErrDraftNotFound   = errors.New("draft not found")
	ErrThreadNotFound  = errors.New("thread not found")
	ErrLabelNotFound   = errors.New("label not found")
	ErrFilterNotFound  = errors.New("filter not found")
)

// ListOptions contains common options for list operations.
type ListOptions struct {
	MaxResults int
	PageToken  string
	Query      string
	LabelIDs   []string
}

// ListResult contains the result of a list operation with pagination.
type ListResult[T any] struct {
	Items         []T
	NextPageToken string
	Total         int
}

// ModifyRequest contains labels to add and remove from a message or thread.
type ModifyRequest struct {
	AddLabels    []string
	RemoveLabels []string
}

// VacationSettings represents auto-reply vacation settings.
type VacationSettings struct {
	EnableAutoReply    bool
	StartTime          int64
	EndTime            int64
	ResponseSubject    string
	ResponseBodyPlain  string
	ResponseBodyHTML   string
	RestrictToContacts bool
	RestrictToDomain   bool
}

// MessageRepository defines operations for managing email messages.
type MessageRepository interface {
	// List retrieves a list of messages matching the given options.
	List(ctx context.Context, opts ListOptions) (*ListResult[*Message], error)

	// Get retrieves a single message by ID.
	Get(ctx context.Context, id string) (*Message, error)

	// Send sends a new message.
	Send(ctx context.Context, msg *Message) (*Message, error)

	// Reply sends a reply to an existing message.
	Reply(ctx context.Context, messageID string, reply *Message) (*Message, error)

	// Forward forwards an existing message.
	Forward(ctx context.Context, messageID string, forward *Message) (*Message, error)

	// Trash moves a message to trash.
	Trash(ctx context.Context, id string) error

	// Untrash removes a message from trash.
	Untrash(ctx context.Context, id string) error

	// Delete permanently deletes a message.
	Delete(ctx context.Context, id string) error

	// Archive archives a message (removes INBOX label).
	Archive(ctx context.Context, id string) error

	// Modify modifies the labels on a message.
	Modify(ctx context.Context, id string, req ModifyRequest) (*Message, error)

	// Search searches for messages matching the query.
	Search(ctx context.Context, query string, opts ListOptions) (*ListResult[*Message], error)
}

// DraftRepository defines operations for managing email drafts.
type DraftRepository interface {
	// List retrieves a list of drafts.
	List(ctx context.Context, opts ListOptions) (*ListResult[*Draft], error)

	// Get retrieves a single draft by ID.
	Get(ctx context.Context, id string) (*Draft, error)

	// Create creates a new draft.
	Create(ctx context.Context, draft *Draft) (*Draft, error)

	// Update updates an existing draft.
	Update(ctx context.Context, draft *Draft) (*Draft, error)

	// Send sends a draft.
	Send(ctx context.Context, id string) (*Message, error)

	// Delete deletes a draft.
	Delete(ctx context.Context, id string) error
}

// ThreadRepository defines operations for managing email threads.
type ThreadRepository interface {
	// List retrieves a list of threads.
	List(ctx context.Context, opts ListOptions) (*ListResult[*Thread], error)

	// Get retrieves a single thread by ID.
	Get(ctx context.Context, id string) (*Thread, error)

	// Modify modifies the labels on a thread.
	Modify(ctx context.Context, id string, req ModifyRequest) (*Thread, error)

	// Trash moves a thread to trash.
	Trash(ctx context.Context, id string) error

	// Untrash removes a thread from trash.
	Untrash(ctx context.Context, id string) error

	// Delete permanently deletes a thread.
	Delete(ctx context.Context, id string) error
}

// LabelRepository defines operations for managing email labels.
type LabelRepository interface {
	// List retrieves all labels.
	List(ctx context.Context) ([]*Label, error)

	// Get retrieves a single label by ID.
	Get(ctx context.Context, id string) (*Label, error)

	// Create creates a new label.
	Create(ctx context.Context, label *Label) (*Label, error)

	// Update updates an existing label.
	Update(ctx context.Context, label *Label) (*Label, error)

	// Delete deletes a label.
	Delete(ctx context.Context, id string) error
}

// SettingsRepository defines operations for managing email settings.
type SettingsRepository interface {
	// GetVacation retrieves the vacation auto-reply settings.
	GetVacation(ctx context.Context) (*VacationSettings, error)

	// SetVacation updates the vacation auto-reply settings.
	SetVacation(ctx context.Context, settings *VacationSettings) (*VacationSettings, error)

	// ListFilters retrieves all filters.
	ListFilters(ctx context.Context) ([]*Filter, error)

	// CreateFilter creates a new filter.
	CreateFilter(ctx context.Context, filter *Filter) (*Filter, error)

	// DeleteFilter deletes a filter.
	DeleteFilter(ctx context.Context, id string) error
}
