package repository

import (
	"context"
	"fmt"
	"time"

	domaintasks "github.com/stainedhead/go-goog-cli/internal/domain/tasks"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

// GTasksRepository is the base repository that wraps the Google Tasks API service.
type GTasksRepository struct {
	service     *tasks.Service
	maxRetries  int
	baseBackoff time.Duration
}

// GTaskListRepository implements TaskListRepository using the Google Tasks API.
type GTaskListRepository struct {
	*GTasksRepository
}

// GTaskRepository implements TaskRepository using the Google Tasks API.
type GTaskRepository struct {
	*GTasksRepository
}

// Compile-time interface compliance checks.
var (
	_ domaintasks.TaskListRepository = (*GTaskListRepository)(nil)
	_ domaintasks.TaskRepository     = (*GTaskRepository)(nil)
)

// NewGTasksRepository creates a new GTasksRepository with the given OAuth2 token source.
func NewGTasksRepository(ctx context.Context, tokenSource oauth2.TokenSource) (*GTasksRepository, error) {
	httpClient := oauth2.NewClient(ctx, tokenSource)

	service, err := tasks.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("failed to create Tasks service: %w", err)
	}

	return &GTasksRepository{
		service:     service,
		maxRetries:  defaultMaxRetries,
		baseBackoff: defaultBaseBackoff,
	}, nil
}

// NewGTasksRepositoryWithService creates a GTasksRepository with a pre-configured service.
// This is useful for testing with mock servers.
func NewGTasksRepositoryWithService(service *tasks.Service) *GTasksRepository {
	return &GTasksRepository{
		service:     service,
		maxRetries:  defaultMaxRetries,
		baseBackoff: defaultBaseBackoff,
	}
}

// NewGTaskListRepository creates a new GTaskListRepository.
func NewGTaskListRepository(repo *GTasksRepository) *GTaskListRepository {
	return &GTaskListRepository{GTasksRepository: repo}
}

// NewGTaskRepository creates a new GTaskRepository.
func NewGTaskRepository(repo *GTasksRepository) *GTaskRepository {
	return &GTaskRepository{GTasksRepository: repo}
}

// =============================================================================
// TaskListRepository Implementation
// =============================================================================

// List retrieves all task lists.
func (r *GTaskListRepository) List(ctx context.Context) ([]*domaintasks.TaskList, error) {
	result, err := retryWithBackoff(ctx, r.maxRetries, r.baseBackoff, func() (*tasks.TaskLists, error) {
		call := r.service.Tasklists.List()
		return call.Do()
	})
	if err != nil {
		return nil, mapTasksError(err, "list task lists")
	}

	taskLists := make([]*domaintasks.TaskList, 0, len(result.Items))
	for _, item := range result.Items {
		taskLists = append(taskLists, apiTaskListToDomain(item))
	}

	return taskLists, nil
}

// Get retrieves a specific task list by ID.
func (r *GTaskListRepository) Get(ctx context.Context, taskListID string) (*domaintasks.TaskList, error) {
	result, err := retryWithBackoff(ctx, r.maxRetries, r.baseBackoff, func() (*tasks.TaskList, error) {
		return r.service.Tasklists.Get(taskListID).Do()
	})
	if err != nil {
		return nil, mapTasksError(err, "get task list")
	}

	return apiTaskListToDomain(result), nil
}

// Create creates a new task list.
func (r *GTaskListRepository) Create(ctx context.Context, taskList *domaintasks.TaskList) (*domaintasks.TaskList, error) {
	apiTaskList := domainTaskListToAPI(taskList)

	result, err := retryWithBackoff(ctx, r.maxRetries, r.baseBackoff, func() (*tasks.TaskList, error) {
		return r.service.Tasklists.Insert(apiTaskList).Do()
	})
	if err != nil {
		return nil, mapTasksError(err, "create task list")
	}

	return apiTaskListToDomain(result), nil
}

// Update updates an existing task list.
func (r *GTaskListRepository) Update(ctx context.Context, taskList *domaintasks.TaskList) (*domaintasks.TaskList, error) {
	apiTaskList := domainTaskListToAPI(taskList)

	result, err := retryWithBackoff(ctx, r.maxRetries, r.baseBackoff, func() (*tasks.TaskList, error) {
		return r.service.Tasklists.Update(taskList.ID, apiTaskList).Do()
	})
	if err != nil {
		return nil, mapTasksError(err, "update task list")
	}

	return apiTaskListToDomain(result), nil
}

// Delete deletes a task list.
func (r *GTaskListRepository) Delete(ctx context.Context, taskListID string) error {
	_, err := retryWithBackoff(ctx, r.maxRetries, r.baseBackoff, func() (struct{}, error) {
		return struct{}{}, r.service.Tasklists.Delete(taskListID).Do()
	})
	if err != nil {
		return mapTasksError(err, "delete task list")
	}

	return nil
}

// =============================================================================
// TaskRepository Implementation
// =============================================================================

// List retrieves tasks from a specific task list.
func (r *GTaskRepository) List(ctx context.Context, taskListID string, opts domaintasks.ListOptions) (*domaintasks.ListResult[*domaintasks.Task], error) {
	call := r.service.Tasks.List(taskListID)

	if opts.MaxResults > 0 {
		call = call.MaxResults(opts.MaxResults)
	}
	if opts.PageToken != "" {
		call = call.PageToken(opts.PageToken)
	}
	if opts.ShowCompleted {
		call = call.ShowCompleted(true)
	}
	if opts.ShowHidden {
		call = call.ShowHidden(true)
	}
	if opts.ShowDeleted {
		call = call.ShowDeleted(true)
	}
	if opts.DueMin != nil {
		call = call.DueMin(opts.DueMin.Format(time.RFC3339))
	}
	if opts.DueMax != nil {
		call = call.DueMax(opts.DueMax.Format(time.RFC3339))
	}
	if opts.UpdatedMin != nil {
		call = call.UpdatedMin(opts.UpdatedMin.Format(time.RFC3339))
	}
	if opts.CompletedMin != nil {
		call = call.CompletedMin(opts.CompletedMin.Format(time.RFC3339))
	}
	if opts.CompletedMax != nil {
		call = call.CompletedMax(opts.CompletedMax.Format(time.RFC3339))
	}

	result, err := retryWithBackoff(ctx, r.maxRetries, r.baseBackoff, func() (*tasks.Tasks, error) {
		return call.Do()
	})
	if err != nil {
		return nil, mapTasksError(err, "list tasks")
	}

	domainTasks := make([]*domaintasks.Task, 0, len(result.Items))
	for _, item := range result.Items {
		domainTasks = append(domainTasks, apiTaskToDomain(item, taskListID))
	}

	return &domaintasks.ListResult[*domaintasks.Task]{
		Items:         domainTasks,
		NextPageToken: result.NextPageToken,
	}, nil
}

// Get retrieves a specific task.
func (r *GTaskRepository) Get(ctx context.Context, taskListID, taskID string) (*domaintasks.Task, error) {
	result, err := retryWithBackoff(ctx, r.maxRetries, r.baseBackoff, func() (*tasks.Task, error) {
		return r.service.Tasks.Get(taskListID, taskID).Do()
	})
	if err != nil {
		return nil, mapTasksError(err, "get task")
	}

	return apiTaskToDomain(result, taskListID), nil
}

// Create creates a new task.
func (r *GTaskRepository) Create(ctx context.Context, taskListID string, task *domaintasks.Task) (*domaintasks.Task, error) {
	apiTask := domainTaskToAPI(task)

	result, err := retryWithBackoff(ctx, r.maxRetries, r.baseBackoff, func() (*tasks.Task, error) {
		call := r.service.Tasks.Insert(taskListID, apiTask)
		if task.Parent != nil {
			call = call.Parent(*task.Parent)
		}
		return call.Do()
	})
	if err != nil {
		return nil, mapTasksError(err, "create task")
	}

	return apiTaskToDomain(result, taskListID), nil
}

// Update updates an existing task.
func (r *GTaskRepository) Update(ctx context.Context, taskListID string, task *domaintasks.Task) (*domaintasks.Task, error) {
	apiTask := domainTaskToAPI(task)

	result, err := retryWithBackoff(ctx, r.maxRetries, r.baseBackoff, func() (*tasks.Task, error) {
		return r.service.Tasks.Update(taskListID, task.ID, apiTask).Do()
	})
	if err != nil {
		return nil, mapTasksError(err, "update task")
	}

	return apiTaskToDomain(result, taskListID), nil
}

// Delete deletes a task.
func (r *GTaskRepository) Delete(ctx context.Context, taskListID, taskID string) error {
	_, err := retryWithBackoff(ctx, r.maxRetries, r.baseBackoff, func() (struct{}, error) {
		return struct{}{}, r.service.Tasks.Delete(taskListID, taskID).Do()
	})
	if err != nil {
		return mapTasksError(err, "delete task")
	}

	return nil
}

// Move moves a task to a different position.
func (r *GTaskRepository) Move(ctx context.Context, taskListID, taskID, parent, previous string) (*domaintasks.Task, error) {
	result, err := retryWithBackoff(ctx, r.maxRetries, r.baseBackoff, func() (*tasks.Task, error) {
		call := r.service.Tasks.Move(taskListID, taskID)
		if parent != "" {
			call = call.Parent(parent)
		}
		if previous != "" {
			call = call.Previous(previous)
		}
		return call.Do()
	})
	if err != nil {
		return nil, mapTasksError(err, "move task")
	}

	return apiTaskToDomain(result, taskListID), nil
}

// Clear clears all completed tasks from a task list.
func (r *GTaskRepository) Clear(ctx context.Context, taskListID string) error {
	_, err := retryWithBackoff(ctx, r.maxRetries, r.baseBackoff, func() (struct{}, error) {
		return struct{}{}, r.service.Tasks.Clear(taskListID).Do()
	})
	if err != nil {
		return mapTasksError(err, "clear completed tasks")
	}

	return nil
}

// =============================================================================
// Transformation Functions
// =============================================================================

// apiTaskListToDomain converts an API TaskList to domain TaskList.
func apiTaskListToDomain(api *tasks.TaskList) *domaintasks.TaskList {
	updated, _ := time.Parse(time.RFC3339, api.Updated)

	return &domaintasks.TaskList{
		ID:       api.Id,
		Title:    api.Title,
		Updated:  updated,
		SelfLink: api.SelfLink,
	}
}

// domainTaskListToAPI converts a domain TaskList to API TaskList.
func domainTaskListToAPI(domain *domaintasks.TaskList) *tasks.TaskList {
	return &tasks.TaskList{
		Id:       domain.ID,
		Title:    domain.Title,
		Updated:  domain.Updated.Format(time.RFC3339),
		SelfLink: domain.SelfLink,
	}
}

// apiTaskToDomain converts an API Task to domain Task.
func apiTaskToDomain(api *tasks.Task, taskListID string) *domaintasks.Task {
	task := &domaintasks.Task{
		ID:         api.Id,
		TaskListID: taskListID,
		Title:      api.Title,
		Notes:      api.Notes,
		Status:     api.Status,
		Position:   api.Position,
		Hidden:     api.Hidden,
		Deleted:    api.Deleted,
	}

	if api.Updated != "" {
		if updated, err := time.Parse(time.RFC3339, api.Updated); err == nil {
			task.Updated = updated
		}
	}

	if api.Due != "" {
		if due, err := time.Parse(time.RFC3339, api.Due); err == nil {
			task.Due = &due
		}
	}

	if api.Completed != nil && *api.Completed != "" {
		if completed, err := time.Parse(time.RFC3339, *api.Completed); err == nil {
			task.Completed = &completed
		}
	}

	if api.Parent != "" {
		task.Parent = &api.Parent
	}

	if api.Links != nil {
		task.Links = make([]domaintasks.TaskLink, 0, len(api.Links))
		for _, link := range api.Links {
			task.Links = append(task.Links, domaintasks.TaskLink{
				Type:        link.Type,
				Description: link.Description,
				Link:        link.Link,
			})
		}
	}

	return task
}

// domainTaskToAPI converts a domain Task to API Task.
func domainTaskToAPI(domain *domaintasks.Task) *tasks.Task {
	apiTask := &tasks.Task{
		Id:       domain.ID,
		Title:    domain.Title,
		Notes:    domain.Notes,
		Status:   domain.Status,
		Position: domain.Position,
		Hidden:   domain.Hidden,
		Deleted:  domain.Deleted,
	}

	if !domain.Updated.IsZero() {
		apiTask.Updated = domain.Updated.Format(time.RFC3339)
	}

	if domain.Due != nil {
		apiTask.Due = domain.Due.Format(time.RFC3339)
	}

	if domain.Completed != nil {
		completedStr := domain.Completed.Format(time.RFC3339)
		apiTask.Completed = &completedStr
	}

	if domain.Parent != nil {
		apiTask.Parent = *domain.Parent
	}

	if domain.Links != nil {
		apiTask.Links = make([]*tasks.TaskLinks, 0, len(domain.Links))
		for _, link := range domain.Links {
			apiTask.Links = append(apiTask.Links, &tasks.TaskLinks{
				Type:        link.Type,
				Description: link.Description,
				Link:        link.Link,
			})
		}
	}

	return apiTask
}

// mapTasksError maps Google Tasks API errors to domain errors.
func mapTasksError(err error, operation string) error {
	if err == nil {
		return nil
	}

	// Use the existing mapAPIError function from errors.go
	return mapAPIError(err, operation)
}
