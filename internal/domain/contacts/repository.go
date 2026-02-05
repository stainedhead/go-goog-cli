package contacts

import (
	"context"
	"errors"
)

// Domain errors
var (
	ErrContactNotFound      = errors.New("contact not found")
	ErrContactGroupNotFound = errors.New("contact group not found")
	ErrInvalidContact       = errors.New("invalid contact data")
	ErrInvalidGroup         = errors.New("invalid contact group data")
	ErrCannotModifySystem   = errors.New("cannot modify system contact group")
)

// ListOptions contains options for listing contacts
type ListOptions struct {
	MaxResults int64
	PageToken  string
	SortOrder  string
}

// SearchOptions contains options for searching contacts
type SearchOptions struct {
	Query      string
	MaxResults int64
	PageToken  string
}

// ListResult contains a paginated list of items
type ListResult[T any] struct {
	Items         []T
	NextPageToken string
	TotalSize     int
}

// ContactRepository defines operations for managing contacts
type ContactRepository interface {
	List(ctx context.Context, opts ListOptions) (*ListResult[*Contact], error)
	Get(ctx context.Context, resourceName string) (*Contact, error)
	Create(ctx context.Context, contact *Contact) (*Contact, error)
	Update(ctx context.Context, contact *Contact, updateMask []string) (*Contact, error)
	Delete(ctx context.Context, resourceName string) error
	Search(ctx context.Context, opts SearchOptions) (*ListResult[*Contact], error)
	BatchGet(ctx context.Context, resourceNames []string) ([]*Contact, error)
}

// ContactGroupRepository defines operations for managing contact groups
type ContactGroupRepository interface {
	List(ctx context.Context) ([]*ContactGroup, error)
	Get(ctx context.Context, resourceName string) (*ContactGroup, error)
	Create(ctx context.Context, group *ContactGroup) (*ContactGroup, error)
	Update(ctx context.Context, group *ContactGroup) (*ContactGroup, error)
	Delete(ctx context.Context, resourceName string) error
	ListMembers(ctx context.Context, resourceName string, opts ListOptions) (*ListResult[*Contact], error)
	AddMembers(ctx context.Context, groupResourceName string, contactResourceNames []string) error
	RemoveMembers(ctx context.Context, groupResourceName string, contactResourceNames []string) error
}
