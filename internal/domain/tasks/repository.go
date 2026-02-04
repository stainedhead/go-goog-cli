package tasks

import (
	"context"
	"errors"
	"time"
)

// Domain errors
var (
	ErrTaskNotFound     = errors.New("task not found")
	ErrTaskListNotFound = errors.New("task list not found")
	ErrInvalidParent    = errors.New("invalid parent task")
	ErrInvalidStatus    = errors.New("invalid task status")
)

// ListOptions contains options for listing tasks
type ListOptions struct {
	MaxResults    int64
	PageToken     string
	ShowCompleted bool
	ShowHidden    bool
	ShowDeleted   bool
	DueMin        *time.Time
	DueMax        *time.Time
	UpdatedMin    *time.Time
	CompletedMin  *time.Time
	CompletedMax  *time.Time
}

// ListResult contains a paginated list of items
type ListResult[T any] struct {
	Items         []T
	NextPageToken string
}

// TaskListRepository defines operations for managing task lists
type TaskListRepository interface {
	List(ctx context.Context) ([]*TaskList, error)
	Get(ctx context.Context, taskListID string) (*TaskList, error)
	Create(ctx context.Context, taskList *TaskList) (*TaskList, error)
	Update(ctx context.Context, taskList *TaskList) (*TaskList, error)
	Delete(ctx context.Context, taskListID string) error
}

// TaskRepository defines operations for managing tasks
type TaskRepository interface {
	List(ctx context.Context, taskListID string, opts ListOptions) (*ListResult[*Task], error)
	Get(ctx context.Context, taskListID, taskID string) (*Task, error)
	Create(ctx context.Context, taskListID string, task *Task) (*Task, error)
	Update(ctx context.Context, taskListID string, task *Task) (*Task, error)
	Delete(ctx context.Context, taskListID, taskID string) error
	Move(ctx context.Context, taskListID, taskID, parent, previous string) (*Task, error)
	Clear(ctx context.Context, taskListID string) error
}
