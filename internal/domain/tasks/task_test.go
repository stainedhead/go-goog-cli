package tasks

import (
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		listID  string
		wantErr bool
	}{
		{
			name:    "valid task",
			title:   "Test Task",
			listID:  "list1",
			wantErr: false,
		},
		{
			name:    "empty title",
			title:   "",
			listID:  "list1",
			wantErr: true,
		},
		{
			name:    "empty list ID",
			title:   "Test Task",
			listID:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewTask(tt.title, tt.listID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if task.Title != tt.title {
					t.Errorf("NewTask() title = %v, want %v", task.Title, tt.title)
				}
				if task.TaskListID != tt.listID {
					t.Errorf("NewTask() listID = %v, want %v", task.TaskListID, tt.listID)
				}
				if task.Status != StatusNeedsAction {
					t.Errorf("NewTask() status = %v, want %v", task.Status, StatusNeedsAction)
				}
			}
		})
	}
}

func TestTask_Complete(t *testing.T) {
	task := &Task{
		ID:         "task1",
		TaskListID: "list1",
		Title:      "Test Task",
		Status:     StatusNeedsAction,
	}

	err := task.Complete()
	if err != nil {
		t.Errorf("Complete() error = %v", err)
	}

	if task.Status != StatusCompleted {
		t.Errorf("Complete() status = %v, want %v", task.Status, StatusCompleted)
	}

	if task.Completed == nil {
		t.Error("Complete() completed time not set")
	}
}

func TestTask_CompleteAlreadyCompleted(t *testing.T) {
	completed := time.Now()
	task := &Task{
		ID:         "task1",
		TaskListID: "list1",
		Title:      "Test Task",
		Status:     StatusCompleted,
		Completed:  &completed,
	}

	err := task.Complete()
	if err != nil {
		t.Errorf("Complete() error = %v", err)
	}

	// Should remain completed with original time
	if task.Status != StatusCompleted {
		t.Errorf("Complete() status = %v, want %v", task.Status, StatusCompleted)
	}
}

func TestTask_Reopen(t *testing.T) {
	completed := time.Now()
	task := &Task{
		ID:         "task1",
		TaskListID: "list1",
		Title:      "Test Task",
		Status:     StatusCompleted,
		Completed:  &completed,
	}

	err := task.Reopen()
	if err != nil {
		t.Errorf("Reopen() error = %v", err)
	}

	if task.Status != StatusNeedsAction {
		t.Errorf("Reopen() status = %v, want %v", task.Status, StatusNeedsAction)
	}

	if task.Completed != nil {
		t.Error("Reopen() completed time not cleared")
	}
}

func TestTask_ReopenNotCompleted(t *testing.T) {
	task := &Task{
		ID:         "task1",
		TaskListID: "list1",
		Title:      "Test Task",
		Status:     StatusNeedsAction,
	}

	err := task.Reopen()
	if err != nil {
		t.Errorf("Reopen() error = %v", err)
	}

	// Should remain needsAction
	if task.Status != StatusNeedsAction {
		t.Errorf("Reopen() status = %v, want %v", task.Status, StatusNeedsAction)
	}
}

func TestTask_IsCompleted(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{
			name:   "completed task",
			status: StatusCompleted,
			want:   true,
		},
		{
			name:   "needs action task",
			status: StatusNeedsAction,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{Status: tt.status}
			if got := task.IsCompleted(); got != tt.want {
				t.Errorf("IsCompleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_IsSubtask(t *testing.T) {
	tests := []struct {
		name   string
		parent *string
		want   bool
	}{
		{
			name:   "has parent",
			parent: strPtr("parent1"),
			want:   true,
		},
		{
			name:   "no parent",
			parent: nil,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{Parent: tt.parent}
			if got := task.IsSubtask(); got != tt.want {
				t.Errorf("IsSubtask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_CanHaveSubtasks(t *testing.T) {
	tests := []struct {
		name   string
		parent *string
		want   bool
	}{
		{
			name:   "top-level task",
			parent: nil,
			want:   true,
		},
		{
			name:   "subtask cannot have children",
			parent: strPtr("parent1"),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{Parent: tt.parent}
			if got := task.CanHaveSubtasks(); got != tt.want {
				t.Errorf("CanHaveSubtasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_SetParent(t *testing.T) {
	tests := []struct {
		name     string
		parentID string
		wantErr  bool
	}{
		{
			name:     "valid parent",
			parentID: "parent1",
			wantErr:  false,
		},
		{
			name:     "empty parent ID",
			parentID: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{
				ID:         "task1",
				TaskListID: "list1",
				Title:      "Test Task",
			}

			err := task.SetParent(tt.parentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetParent() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if task.Parent == nil || *task.Parent != tt.parentID {
					t.Errorf("SetParent() parent = %v, want %v", task.Parent, tt.parentID)
				}
			}
		})
	}
}

func TestTask_IsOverdue(t *testing.T) {
	now := time.Now()
	pastDue := now.Add(-24 * time.Hour)
	futureDue := now.Add(24 * time.Hour)

	tests := []struct {
		name   string
		due    *time.Time
		status string
		want   bool
	}{
		{
			name:   "past due and not completed",
			due:    &pastDue,
			status: StatusNeedsAction,
			want:   true,
		},
		{
			name:   "future due",
			due:    &futureDue,
			status: StatusNeedsAction,
			want:   false,
		},
		{
			name:   "no due date",
			due:    nil,
			status: StatusNeedsAction,
			want:   false,
		},
		{
			name:   "past due but completed",
			due:    &pastDue,
			status: StatusCompleted,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{
				Due:    tt.due,
				Status: tt.status,
			}
			if got := task.IsOverdue(); got != tt.want {
				t.Errorf("IsOverdue() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function
func strPtr(s string) *string {
	return &s
}
