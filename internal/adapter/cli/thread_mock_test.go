// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"fmt"
	"testing"

	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
)

// =============================================================================
// Tests for mock thread repository
// =============================================================================

func TestMockThreadRepository_List(t *testing.T) {
	t.Run("List success", func(t *testing.T) {
		threads := []*mail.Thread{
			{ID: "thread1", Snippet: "Thread 1"},
			{ID: "thread2", Snippet: "Thread 2"},
		}
		repo := &MockThreadRepository{Threads: threads}

		result, err := repo.List(nil, mail.ListOptions{MaxResults: 10})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result.Items) != 2 {
			t.Errorf("expected 2 threads, got %d", len(result.Items))
		}
	})

	t.Run("List with custom result", func(t *testing.T) {
		listResult := &mail.ListResult[*mail.Thread]{
			Items: []*mail.Thread{{ID: "custom"}},
			Total: 100,
		}
		repo := &MockThreadRepository{ListResult: listResult}

		result, err := repo.List(nil, mail.ListOptions{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.Total != 100 {
			t.Errorf("expected total 100, got %d", result.Total)
		}
	})

	t.Run("List error", func(t *testing.T) {
		repo := &MockThreadRepository{ListErr: fmt.Errorf("list error")}

		_, err := repo.List(nil, mail.ListOptions{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockThreadRepository_Get(t *testing.T) {
	t.Run("Get success", func(t *testing.T) {
		thread := &mail.Thread{
			ID:       "thread1",
			Snippet:  "Test Thread",
			Messages: []*mail.Message{{ID: "msg1"}},
		}
		repo := &MockThreadRepository{Thread: thread}

		result, err := repo.Get(nil, "thread1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "thread1" {
			t.Errorf("expected ID 'thread1', got %s", result.ID)
		}
		if len(result.Messages) != 1 {
			t.Errorf("expected 1 message, got %d", len(result.Messages))
		}
	})

	t.Run("Get error", func(t *testing.T) {
		repo := &MockThreadRepository{GetErr: fmt.Errorf("not found")}

		_, err := repo.Get(nil, "nonexistent")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockThreadRepository_Modify(t *testing.T) {
	t.Run("Modify success", func(t *testing.T) {
		thread := &mail.Thread{ID: "thread1", Labels: []string{"STARRED"}}
		repo := &MockThreadRepository{Thread: thread}

		req := mail.ModifyRequest{AddLabels: []string{"STARRED"}}
		result, err := repo.Modify(nil, "thread1", req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "thread1" {
			t.Errorf("expected ID 'thread1', got %s", result.ID)
		}
	})

	t.Run("Modify with custom result", func(t *testing.T) {
		repo := &MockThreadRepository{
			ModifyResult: &mail.Thread{ID: "modified", Labels: []string{"IMPORTANT"}},
		}

		req := mail.ModifyRequest{AddLabels: []string{"IMPORTANT"}}
		result, err := repo.Modify(nil, "thread1", req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "modified" {
			t.Errorf("expected ID 'modified', got %s", result.ID)
		}
	})

	t.Run("Modify error", func(t *testing.T) {
		repo := &MockThreadRepository{ModifyErr: fmt.Errorf("modify error")}

		req := mail.ModifyRequest{}
		_, err := repo.Modify(nil, "thread1", req)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockThreadRepository_Trash(t *testing.T) {
	t.Run("Trash success", func(t *testing.T) {
		repo := &MockThreadRepository{}

		err := repo.Trash(nil, "thread1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Trash error", func(t *testing.T) {
		repo := &MockThreadRepository{TrashErr: fmt.Errorf("trash error")}

		err := repo.Trash(nil, "thread1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockThreadRepository_Untrash(t *testing.T) {
	t.Run("Untrash success", func(t *testing.T) {
		repo := &MockThreadRepository{}

		err := repo.Untrash(nil, "thread1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Untrash error", func(t *testing.T) {
		repo := &MockThreadRepository{UntrashErr: fmt.Errorf("untrash error")}

		err := repo.Untrash(nil, "thread1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockThreadRepository_Delete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		repo := &MockThreadRepository{}

		err := repo.Delete(nil, "thread1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Delete error", func(t *testing.T) {
		repo := &MockThreadRepository{DeleteErr: fmt.Errorf("delete error")}

		err := repo.Delete(nil, "thread1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
