package tasks

import (
	"testing"
)

func TestNewTaskList(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{
			name:    "valid task list",
			title:   "My Tasks",
			wantErr: false,
		},
		{
			name:    "empty title",
			title:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskList, err := NewTaskList(tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTaskList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if taskList.Title != tt.title {
					t.Errorf("NewTaskList() title = %v, want %v", taskList.Title, tt.title)
				}
			}
		})
	}
}

func TestTaskList_IsDefault(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want bool
	}{
		{
			name: "default list",
			id:   "@default",
			want: true,
		},
		{
			name: "custom list",
			id:   "custom-list-123",
			want: false,
		},
		{
			name: "empty ID",
			id:   "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskList := &TaskList{ID: tt.id}
			if got := taskList.IsDefault(); got != tt.want {
				t.Errorf("IsDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
