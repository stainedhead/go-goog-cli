// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"fmt"
	"testing"

	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
)

// =============================================================================
// Additional tests for mock draft repository
// =============================================================================

func TestMockDraftRepository_List(t *testing.T) {
	t.Run("List success", func(t *testing.T) {
		drafts := []*mail.Draft{
			{ID: "draft1", Message: &mail.Message{Subject: "Draft 1"}},
			{ID: "draft2", Message: &mail.Message{Subject: "Draft 2"}},
		}
		repo := &MockDraftRepository{Drafts: drafts}

		result, err := repo.List(nil, mail.ListOptions{MaxResults: 10})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result.Items) != 2 {
			t.Errorf("expected 2 drafts, got %d", len(result.Items))
		}
	})

	t.Run("List with custom result", func(t *testing.T) {
		listResult := &mail.ListResult[*mail.Draft]{
			Items: []*mail.Draft{{ID: "custom"}},
			Total: 50,
		}
		repo := &MockDraftRepository{ListResult: listResult}

		result, err := repo.List(nil, mail.ListOptions{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.Total != 50 {
			t.Errorf("expected total 50, got %d", result.Total)
		}
	})

	t.Run("List error", func(t *testing.T) {
		repo := &MockDraftRepository{ListErr: fmt.Errorf("list error")}

		_, err := repo.List(nil, mail.ListOptions{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockDraftRepository_Get(t *testing.T) {
	t.Run("Get success", func(t *testing.T) {
		draft := &mail.Draft{
			ID: "draft1",
			Message: &mail.Message{
				Subject: "Test Draft",
				To:      []string{"user@example.com"},
				Body:    "Draft body content",
			},
		}
		repo := &MockDraftRepository{Draft: draft}

		result, err := repo.Get(nil, "draft1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "draft1" {
			t.Errorf("expected ID 'draft1', got %s", result.ID)
		}
		if result.Message.Subject != "Test Draft" {
			t.Errorf("expected subject 'Test Draft', got %s", result.Message.Subject)
		}
	})

	t.Run("Get error", func(t *testing.T) {
		repo := &MockDraftRepository{GetErr: fmt.Errorf("not found")}

		_, err := repo.Get(nil, "nonexistent")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockDraftRepository_Create(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		repo := &MockDraftRepository{}
		draft := &mail.Draft{
			Message: &mail.Message{
				To:      []string{"user@example.com"},
				Subject: "New Draft",
			},
		}

		result, err := repo.Create(nil, draft)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "mock-draft-id" {
			t.Errorf("expected mock ID, got %s", result.ID)
		}
	})

	t.Run("Create with custom result", func(t *testing.T) {
		repo := &MockDraftRepository{
			CreateResult: &mail.Draft{ID: "custom-draft-id"},
		}
		draft := &mail.Draft{
			Message: &mail.Message{Subject: "New Draft"},
		}

		result, err := repo.Create(nil, draft)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "custom-draft-id" {
			t.Errorf("expected 'custom-draft-id', got %s", result.ID)
		}
	})

	t.Run("Create error", func(t *testing.T) {
		repo := &MockDraftRepository{CreateErr: fmt.Errorf("create error")}
		draft := &mail.Draft{
			Message: &mail.Message{Subject: "New Draft"},
		}

		_, err := repo.Create(nil, draft)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockDraftRepository_Update(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		repo := &MockDraftRepository{}
		draft := &mail.Draft{
			ID: "draft1",
			Message: &mail.Message{
				Subject: "Updated Subject",
			},
		}

		result, err := repo.Update(nil, draft)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.Message.Subject != "Updated Subject" {
			t.Errorf("expected 'Updated Subject', got %s", result.Message.Subject)
		}
	})

	t.Run("Update with custom result", func(t *testing.T) {
		repo := &MockDraftRepository{
			UpdateResult: &mail.Draft{
				ID:      "draft1",
				Message: &mail.Message{Subject: "Custom Updated"},
			},
		}
		draft := &mail.Draft{ID: "draft1"}

		result, err := repo.Update(nil, draft)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.Message.Subject != "Custom Updated" {
			t.Errorf("expected 'Custom Updated', got %s", result.Message.Subject)
		}
	})

	t.Run("Update error", func(t *testing.T) {
		repo := &MockDraftRepository{UpdateErr: fmt.Errorf("update error")}
		draft := &mail.Draft{ID: "draft1"}

		_, err := repo.Update(nil, draft)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockDraftRepository_Send(t *testing.T) {
	t.Run("Send success with default result", func(t *testing.T) {
		repo := &MockDraftRepository{}

		result, err := repo.Send(nil, "draft1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "sent-msg-id" {
			t.Errorf("expected 'sent-msg-id', got %s", result.ID)
		}
	})

	t.Run("Send with custom result", func(t *testing.T) {
		repo := &MockDraftRepository{
			SendResult: &mail.Message{
				ID:       "custom-sent-id",
				ThreadID: "thread-123",
			},
		}

		result, err := repo.Send(nil, "draft1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "custom-sent-id" {
			t.Errorf("expected 'custom-sent-id', got %s", result.ID)
		}
		if result.ThreadID != "thread-123" {
			t.Errorf("expected thread 'thread-123', got %s", result.ThreadID)
		}
	})

	t.Run("Send error", func(t *testing.T) {
		repo := &MockDraftRepository{SendErr: fmt.Errorf("send error")}

		_, err := repo.Send(nil, "draft1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockDraftRepository_Delete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		repo := &MockDraftRepository{}

		err := repo.Delete(nil, "draft1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Delete error", func(t *testing.T) {
		repo := &MockDraftRepository{DeleteErr: fmt.Errorf("delete error")}

		err := repo.Delete(nil, "draft1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
