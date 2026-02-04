package tasks

import (
	"errors"
	"time"
)

// Status constants
const (
	StatusNeedsAction = "needsAction"
	StatusCompleted   = "completed"
)

// Task represents a Google Task with business logic
type Task struct {
	ID         string
	TaskListID string
	Title      string
	Notes      string
	Status     string
	Due        *time.Time
	Completed  *time.Time
	Parent     *string
	Position   string
	Updated    time.Time
	Links      []TaskLink
	Hidden     bool
	Deleted    bool
}

// TaskLink represents a link associated with a task
type TaskLink struct {
	Type        string
	Description string
	Link        string
}

// NewTask creates a new task with the given title and list ID
func NewTask(title, listID string) (*Task, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if listID == "" {
		return nil, errors.New("list ID cannot be empty")
	}

	return &Task{
		TaskListID: listID,
		Title:      title,
		Status:     StatusNeedsAction,
		Updated:    time.Now(),
	}, nil
}

// Complete marks the task as completed
func (t *Task) Complete() error {
	if t.Status != StatusCompleted {
		t.Status = StatusCompleted
		now := time.Now()
		t.Completed = &now
		t.Updated = now
	}
	return nil
}

// Reopen marks a completed task as needing action
func (t *Task) Reopen() error {
	if t.Status == StatusCompleted {
		t.Status = StatusNeedsAction
		t.Completed = nil
		t.Updated = time.Now()
	}
	return nil
}

// IsCompleted returns true if the task is completed
func (t *Task) IsCompleted() bool {
	return t.Status == StatusCompleted
}

// IsSubtask returns true if the task has a parent
func (t *Task) IsSubtask() bool {
	return t.Parent != nil
}

// CanHaveSubtasks returns true if the task can have subtasks
// Only top-level tasks can have subtasks (not subtasks themselves)
func (t *Task) CanHaveSubtasks() bool {
	return !t.IsSubtask()
}

// SetParent sets the parent task ID for this task
func (t *Task) SetParent(parentID string) error {
	if parentID == "" {
		return errors.New("parent ID cannot be empty")
	}
	t.Parent = &parentID
	t.Updated = time.Now()
	return nil
}

// IsOverdue returns true if the task has a due date in the past and is not completed
func (t *Task) IsOverdue() bool {
	if t.Due == nil || t.IsCompleted() {
		return false
	}
	return t.Due.Before(time.Now())
}
