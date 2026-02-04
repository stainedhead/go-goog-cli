package repository

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domaintasks "github.com/stainedhead/go-goog-cli/internal/domain/tasks"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

// TestAPITaskListToDomain tests conversion from API TaskList to domain TaskList.
func TestAPITaskListToDomain(t *testing.T) {
	tests := []struct {
		name          string
		apiTaskList   *tasks.TaskList
		expectedID    string
		expectedTitle string
	}{
		{
			name: "basic task list",
			apiTaskList: &tasks.TaskList{
				Id:       "list1",
				Title:    "My Tasks",
				Updated:  "2024-01-15T10:30:00Z",
				SelfLink: "https://example.com/tasklists/list1",
			},
			expectedID:    "list1",
			expectedTitle: "My Tasks",
		},
		{
			name: "task list with special characters",
			apiTaskList: &tasks.TaskList{
				Id:      "list2",
				Title:   "Work & Projects!",
				Updated: "2024-01-15T10:30:00Z",
			},
			expectedID:    "list2",
			expectedTitle: "Work & Projects!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := apiTaskListToDomain(tt.apiTaskList)

			if result.ID != tt.expectedID {
				t.Errorf("ID = %q, want %q", result.ID, tt.expectedID)
			}
			if result.Title != tt.expectedTitle {
				t.Errorf("Title = %q, want %q", result.Title, tt.expectedTitle)
			}
		})
	}
}

// TestDomainTaskListToAPI tests conversion from domain TaskList to API TaskList.
func TestDomainTaskListToAPI(t *testing.T) {
	updated := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	taskList := &domaintasks.TaskList{
		ID:       "list1",
		Title:    "My Tasks",
		Updated:  updated,
		SelfLink: "https://example.com/tasklists/list1",
	}

	result := domainTaskListToAPI(taskList)

	if result.Id != "list1" {
		t.Errorf("Id = %q, want %q", result.Id, "list1")
	}
	if result.Title != "My Tasks" {
		t.Errorf("Title = %q, want %q", result.Title, "My Tasks")
	}
}

// TestAPITaskToDomain tests conversion from API Task to domain Task.
func TestAPITaskToDomain(t *testing.T) {
	parentID := "parent1"
	completedStr := "2024-01-15T12:00:00Z"

	tests := []struct {
		name           string
		apiTask        *tasks.Task
		taskListID     string
		expectedTitle  string
		expectedStatus string
	}{
		{
			name: "basic task",
			apiTask: &tasks.Task{
				Id:      "task1",
				Title:   "Test Task",
				Status:  "needsAction",
				Updated: "2024-01-15T10:30:00Z",
			},
			taskListID:     "list1",
			expectedTitle:  "Test Task",
			expectedStatus: "needsAction",
		},
		{
			name: "completed task with parent",
			apiTask: &tasks.Task{
				Id:        "task2",
				Title:     "Subtask",
				Status:    "completed",
				Updated:   "2024-01-15T10:30:00Z",
				Completed: &completedStr,
				Parent:    parentID,
			},
			taskListID:     "list1",
			expectedTitle:  "Subtask",
			expectedStatus: "completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := apiTaskToDomain(tt.apiTask, tt.taskListID)

			if result.Title != tt.expectedTitle {
				t.Errorf("Title = %q, want %q", result.Title, tt.expectedTitle)
			}
			if result.Status != tt.expectedStatus {
				t.Errorf("Status = %q, want %q", result.Status, tt.expectedStatus)
			}
			if result.TaskListID != tt.taskListID {
				t.Errorf("TaskListID = %q, want %q", result.TaskListID, tt.taskListID)
			}
		})
	}
}

// TestDomainTaskToAPI tests conversion from domain Task to API Task.
func TestDomainTaskToAPI(t *testing.T) {
	parentID := "parent1"
	dueDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	completedDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	task := &domaintasks.Task{
		ID:         "task1",
		TaskListID: "list1",
		Title:      "Test Task",
		Notes:      "Some notes",
		Status:     "needsAction",
		Due:        &dueDate,
		Completed:  &completedDate,
		Parent:     &parentID,
	}

	result := domainTaskToAPI(task)

	if result.Id != "task1" {
		t.Errorf("Id = %q, want %q", result.Id, "task1")
	}
	if result.Title != "Test Task" {
		t.Errorf("Title = %q, want %q", result.Title, "Test Task")
	}
	if result.Notes != "Some notes" {
		t.Errorf("Notes = %q, want %q", result.Notes, "Some notes")
	}
	if result.Parent != parentID {
		t.Errorf("Parent = %q, want %q", result.Parent, parentID)
	}
}

// TestGTaskListRepository_List tests listing task lists with a mock server.
func TestGTaskListRepository_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/v1/users/@me/lists" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		response := tasks.TaskLists{
			Items: []*tasks.TaskList{
				{
					Id:      "list1",
					Title:   "My Tasks",
					Updated: "2024-01-15T10:30:00Z",
				},
				{
					Id:      "list2",
					Title:   "Work Tasks",
					Updated: "2024-01-14T15:45:00Z",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	service, err := tasks.NewService(context.Background(),
		option.WithEndpoint(server.URL),
		option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	repo := NewGTaskListRepository(NewGTasksRepositoryWithService(service))
	result, err := repo.List(context.Background())

	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(result) != 2 {
		t.Errorf("List() returned %d items, want 2", len(result))
	}

	if result[0].Title != "My Tasks" {
		t.Errorf("result[0].Title = %q, want %q", result[0].Title, "My Tasks")
	}
}

// TestGTaskListRepository_Get tests getting a specific task list.
func TestGTaskListRepository_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/v1/users/@me/lists/list1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		response := tasks.TaskList{
			Id:      "list1",
			Title:   "My Tasks",
			Updated: "2024-01-15T10:30:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	service, err := tasks.NewService(context.Background(),
		option.WithEndpoint(server.URL),
		option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	repo := NewGTaskListRepository(NewGTasksRepositoryWithService(service))
	result, err := repo.Get(context.Background(), "list1")

	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if result.ID != "list1" {
		t.Errorf("result.ID = %q, want %q", result.ID, "list1")
	}
	if result.Title != "My Tasks" {
		t.Errorf("result.Title = %q, want %q", result.Title, "My Tasks")
	}
}

// TestGTaskRepository_List tests listing tasks with filters.
func TestGTaskRepository_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/v1/lists/list1/tasks" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Check query parameters
		if r.URL.Query().Get("showCompleted") != "true" {
			t.Error("expected showCompleted=true")
		}

		response := tasks.Tasks{
			Items: []*tasks.Task{
				{
					Id:      "task1",
					Title:   "Task 1",
					Status:  "needsAction",
					Updated: "2024-01-15T10:30:00Z",
				},
				{
					Id:      "task2",
					Title:   "Task 2",
					Status:  "completed",
					Updated: "2024-01-14T15:45:00Z",
				},
			},
			NextPageToken: "nextpage",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	service, err := tasks.NewService(context.Background(),
		option.WithEndpoint(server.URL),
		option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	repo := NewGTaskRepository(NewGTasksRepositoryWithService(service))
	result, err := repo.List(context.Background(), "list1", domaintasks.ListOptions{
		ShowCompleted: true,
	})

	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(result.Items) != 2 {
		t.Errorf("List() returned %d items, want 2", len(result.Items))
	}

	if result.NextPageToken != "nextpage" {
		t.Errorf("NextPageToken = %q, want %q", result.NextPageToken, "nextpage")
	}

	if result.Items[0].Title != "Task 1" {
		t.Errorf("result.Items[0].Title = %q, want %q", result.Items[0].Title, "Task 1")
	}
}

// TestGTaskRepository_Create tests creating a new task.
func TestGTaskRepository_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/v1/lists/list1/tasks" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		response := tasks.Task{
			Id:      "task1",
			Title:   "New Task",
			Status:  "needsAction",
			Updated: "2024-01-15T10:30:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	service, err := tasks.NewService(context.Background(),
		option.WithEndpoint(server.URL),
		option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	repo := NewGTaskRepository(NewGTasksRepositoryWithService(service))
	task := &domaintasks.Task{
		Title:  "New Task",
		Status: "needsAction",
	}

	result, err := repo.Create(context.Background(), "list1", task)

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if result.ID != "task1" {
		t.Errorf("result.ID = %q, want %q", result.ID, "task1")
	}
	if result.Title != "New Task" {
		t.Errorf("result.Title = %q, want %q", result.Title, "New Task")
	}
}

// TestGTaskRepository_Delete tests deleting a task.
func TestGTaskRepository_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/v1/lists/list1/tasks/task1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	service, err := tasks.NewService(context.Background(),
		option.WithEndpoint(server.URL),
		option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	repo := NewGTaskRepository(NewGTasksRepositoryWithService(service))
	err = repo.Delete(context.Background(), "list1", "task1")

	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

// TestGTaskRepository_Clear tests clearing completed tasks.
func TestGTaskRepository_Clear(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/v1/lists/list1/clear" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	service, err := tasks.NewService(context.Background(),
		option.WithEndpoint(server.URL),
		option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	repo := NewGTaskRepository(NewGTasksRepositoryWithService(service))
	err = repo.Clear(context.Background(), "list1")

	if err != nil {
		t.Fatalf("Clear() error = %v", err)
	}
}
