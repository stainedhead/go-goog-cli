package tasks

import (
	"errors"
	"testing"
)

func TestDomainErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "task not found",
			err:  ErrTaskNotFound,
			want: "task not found",
		},
		{
			name: "task list not found",
			err:  ErrTaskListNotFound,
			want: "task list not found",
		},
		{
			name: "invalid parent",
			err:  ErrInvalidParent,
			want: "invalid parent task",
		},
		{
			name: "invalid status",
			err:  ErrInvalidStatus,
			want: "invalid task status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.want {
				t.Errorf("error message = %v, want %v", tt.err.Error(), tt.want)
			}
		})
	}
}

func TestDomainErrorsAreUnique(t *testing.T) {
	if errors.Is(ErrTaskNotFound, ErrTaskListNotFound) {
		t.Error("ErrTaskNotFound should not be ErrTaskListNotFound")
	}
	if errors.Is(ErrTaskNotFound, ErrInvalidParent) {
		t.Error("ErrTaskNotFound should not be ErrInvalidParent")
	}
	if errors.Is(ErrTaskNotFound, ErrInvalidStatus) {
		t.Error("ErrTaskNotFound should not be ErrInvalidStatus")
	}
}

func TestListOptions(t *testing.T) {
	opts := ListOptions{
		MaxResults:    100,
		PageToken:     "token123",
		ShowCompleted: true,
		ShowHidden:    false,
		ShowDeleted:   false,
	}

	if opts.MaxResults != 100 {
		t.Errorf("MaxResults = %v, want 100", opts.MaxResults)
	}
	if opts.PageToken != "token123" {
		t.Errorf("PageToken = %v, want token123", opts.PageToken)
	}
	if !opts.ShowCompleted {
		t.Error("ShowCompleted should be true")
	}
	if opts.ShowHidden {
		t.Error("ShowHidden should be false")
	}
	if opts.ShowDeleted {
		t.Error("ShowDeleted should be false")
	}
}

func TestListResult(t *testing.T) {
	tasks := []*Task{
		{ID: "task1", Title: "Task 1"},
		{ID: "task2", Title: "Task 2"},
	}

	result := ListResult[*Task]{
		Items:         tasks,
		NextPageToken: "nextPage",
	}

	if len(result.Items) != 2 {
		t.Errorf("Items length = %v, want 2", len(result.Items))
	}
	if result.NextPageToken != "nextPage" {
		t.Errorf("NextPageToken = %v, want nextPage", result.NextPageToken)
	}
}
