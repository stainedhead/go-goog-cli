// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"fmt"
	"testing"

	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
)

// =============================================================================
// Tests for mock label repository
// =============================================================================

func TestMockLabelRepository_List(t *testing.T) {
	t.Run("List success", func(t *testing.T) {
		labels := []*mail.Label{
			{ID: "INBOX", Name: "INBOX", Type: "system"},
			{ID: "STARRED", Name: "STARRED", Type: "system"},
			{ID: "Label_1", Name: "Work", Type: "user"},
		}
		repo := &MockLabelRepository{Labels: labels}

		result, err := repo.List(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 labels, got %d", len(result))
		}
	})

	t.Run("List error", func(t *testing.T) {
		repo := &MockLabelRepository{ListErr: fmt.Errorf("list error")}

		_, err := repo.List(nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockLabelRepository_Get(t *testing.T) {
	t.Run("Get success", func(t *testing.T) {
		label := &mail.Label{ID: "Label_1", Name: "Work", Type: "user"}
		repo := &MockLabelRepository{Label: label}

		result, err := repo.Get(nil, "Label_1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "Label_1" {
			t.Errorf("expected ID 'Label_1', got %s", result.ID)
		}
		if result.Name != "Work" {
			t.Errorf("expected name 'Work', got %s", result.Name)
		}
	})

	t.Run("Get error", func(t *testing.T) {
		repo := &MockLabelRepository{GetErr: fmt.Errorf("not found")}

		_, err := repo.Get(nil, "nonexistent")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockLabelRepository_GetByName(t *testing.T) {
	t.Run("GetByName success", func(t *testing.T) {
		label := &mail.Label{ID: "Label_1", Name: "Work", Type: "user"}
		repo := &MockLabelRepository{Label: label}

		result, err := repo.GetByName(nil, "Work")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.Name != "Work" {
			t.Errorf("expected name 'Work', got %s", result.Name)
		}
	})

	t.Run("GetByName error", func(t *testing.T) {
		repo := &MockLabelRepository{GetByNameErr: fmt.Errorf("not found")}

		_, err := repo.GetByName(nil, "NonExistent")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockLabelRepository_Create(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		repo := &MockLabelRepository{}
		label := &mail.Label{Name: "New Label"}

		result, err := repo.Create(nil, label)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "mock-label-id" {
			t.Errorf("expected mock ID, got %s", result.ID)
		}
	})

	t.Run("Create with custom result", func(t *testing.T) {
		repo := &MockLabelRepository{
			CreateResult: &mail.Label{ID: "custom-id", Name: "Custom"},
		}
		label := &mail.Label{Name: "New Label"}

		result, err := repo.Create(nil, label)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "custom-id" {
			t.Errorf("expected 'custom-id', got %s", result.ID)
		}
	})

	t.Run("Create error", func(t *testing.T) {
		repo := &MockLabelRepository{CreateErr: fmt.Errorf("create error")}
		label := &mail.Label{Name: "New Label"}

		_, err := repo.Create(nil, label)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockLabelRepository_Update(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		repo := &MockLabelRepository{}
		label := &mail.Label{ID: "Label_1", Name: "Updated Name"}

		result, err := repo.Update(nil, label)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.Name != "Updated Name" {
			t.Errorf("expected 'Updated Name', got %s", result.Name)
		}
	})

	t.Run("Update with custom result", func(t *testing.T) {
		repo := &MockLabelRepository{
			UpdateResult: &mail.Label{ID: "Label_1", Name: "Custom Update"},
		}
		label := &mail.Label{ID: "Label_1", Name: "New Name"}

		result, err := repo.Update(nil, label)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.Name != "Custom Update" {
			t.Errorf("expected 'Custom Update', got %s", result.Name)
		}
	})

	t.Run("Update error", func(t *testing.T) {
		repo := &MockLabelRepository{UpdateErr: fmt.Errorf("update error")}
		label := &mail.Label{ID: "Label_1", Name: "Updated"}

		_, err := repo.Update(nil, label)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockLabelRepository_Delete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		repo := &MockLabelRepository{}

		err := repo.Delete(nil, "Label_1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Delete error", func(t *testing.T) {
		repo := &MockLabelRepository{DeleteErr: fmt.Errorf("delete error")}

		err := repo.Delete(nil, "Label_1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
