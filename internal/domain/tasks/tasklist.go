package tasks

import (
	"errors"
	"time"
)

// TaskList represents a Google Tasks list
type TaskList struct {
	ID       string
	Title    string
	Updated  time.Time
	SelfLink string
}

// NewTaskList creates a new task list with the given title
func NewTaskList(title string) (*TaskList, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	return &TaskList{
		Title:   title,
		Updated: time.Now(),
	}, nil
}

// IsDefault returns true if this is the default task list
func (tl *TaskList) IsDefault() bool {
	return tl.ID == "@default"
}
