// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"fmt"
	"testing"
	"time"

	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
)

// =============================================================================
// Tests for mock event repository
// =============================================================================

func TestMockEventRepository_List(t *testing.T) {
	t.Run("List success", func(t *testing.T) {
		events := []*calendar.Event{
			{ID: "event1", Title: "Meeting 1"},
			{ID: "event2", Title: "Meeting 2"},
		}
		repo := &MockEventRepository{Events: events}

		timeMin := time.Now()
		timeMax := timeMin.Add(24 * time.Hour)
		result, err := repo.List(nil, "primary", timeMin, timeMax)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 events, got %d", len(result))
		}
	})

	t.Run("List error", func(t *testing.T) {
		repo := &MockEventRepository{ListErr: fmt.Errorf("list error")}

		timeMin := time.Now()
		timeMax := timeMin.Add(24 * time.Hour)
		_, err := repo.List(nil, "primary", timeMin, timeMax)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockEventRepository_Get(t *testing.T) {
	t.Run("Get success", func(t *testing.T) {
		event := &calendar.Event{ID: "event1", Title: "Test Event"}
		repo := &MockEventRepository{Event: event}

		result, err := repo.Get(nil, "primary", "event1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "event1" {
			t.Errorf("expected ID 'event1', got %s", result.ID)
		}
	})

	t.Run("Get error", func(t *testing.T) {
		repo := &MockEventRepository{GetErr: fmt.Errorf("not found")}

		_, err := repo.Get(nil, "primary", "nonexistent")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockEventRepository_Create(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		repo := &MockEventRepository{}
		event := &calendar.Event{Title: "New Event"}

		result, err := repo.Create(nil, "primary", event)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "mock-event-id" {
			t.Errorf("expected mock ID, got %s", result.ID)
		}
	})

	t.Run("Create with custom result", func(t *testing.T) {
		repo := &MockEventRepository{
			CreateResult: &calendar.Event{ID: "custom-id", Title: "Custom"},
		}
		event := &calendar.Event{Title: "New Event"}

		result, err := repo.Create(nil, "primary", event)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "custom-id" {
			t.Errorf("expected 'custom-id', got %s", result.ID)
		}
	})

	t.Run("Create error", func(t *testing.T) {
		repo := &MockEventRepository{CreateErr: fmt.Errorf("create error")}
		event := &calendar.Event{Title: "New Event"}

		_, err := repo.Create(nil, "primary", event)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockEventRepository_Update(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		repo := &MockEventRepository{}
		event := &calendar.Event{ID: "event1", Title: "Updated Event"}

		result, err := repo.Update(nil, "primary", event)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.Title != "Updated Event" {
			t.Errorf("expected 'Updated Event', got %s", result.Title)
		}
	})

	t.Run("Update error", func(t *testing.T) {
		repo := &MockEventRepository{UpdateErr: fmt.Errorf("update error")}
		event := &calendar.Event{ID: "event1", Title: "Updated"}

		_, err := repo.Update(nil, "primary", event)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockEventRepository_Delete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		repo := &MockEventRepository{}

		err := repo.Delete(nil, "primary", "event1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Delete error", func(t *testing.T) {
		repo := &MockEventRepository{DeleteErr: fmt.Errorf("delete error")}

		err := repo.Delete(nil, "primary", "event1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockEventRepository_Move(t *testing.T) {
	t.Run("Move success", func(t *testing.T) {
		event := &calendar.Event{ID: "event1", Title: "Moved Event"}
		repo := &MockEventRepository{Event: event}

		result, err := repo.Move(nil, "source-cal", "event1", "dest-cal")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "event1" {
			t.Errorf("expected ID 'event1', got %s", result.ID)
		}
	})

	t.Run("Move with custom result", func(t *testing.T) {
		repo := &MockEventRepository{
			MoveResult: &calendar.Event{ID: "moved-id"},
		}

		result, err := repo.Move(nil, "source-cal", "event1", "dest-cal")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "moved-id" {
			t.Errorf("expected 'moved-id', got %s", result.ID)
		}
	})

	t.Run("Move error", func(t *testing.T) {
		repo := &MockEventRepository{MoveErr: fmt.Errorf("move error")}

		_, err := repo.Move(nil, "source-cal", "event1", "dest-cal")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockEventRepository_QuickAdd(t *testing.T) {
	t.Run("QuickAdd success", func(t *testing.T) {
		repo := &MockEventRepository{}

		result, err := repo.QuickAdd(nil, "primary", "Meeting at 3pm")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "quick-add-id" {
			t.Errorf("expected 'quick-add-id', got %s", result.ID)
		}
		if result.Title != "Meeting at 3pm" {
			t.Errorf("expected title 'Meeting at 3pm', got %s", result.Title)
		}
	})

	t.Run("QuickAdd with custom result", func(t *testing.T) {
		repo := &MockEventRepository{
			QuickAddResult: &calendar.Event{ID: "custom-quick-id", Title: "Custom"},
		}

		result, err := repo.QuickAdd(nil, "primary", "Meeting")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "custom-quick-id" {
			t.Errorf("expected 'custom-quick-id', got %s", result.ID)
		}
	})

	t.Run("QuickAdd error", func(t *testing.T) {
		repo := &MockEventRepository{QuickAddErr: fmt.Errorf("quick add error")}

		_, err := repo.QuickAdd(nil, "primary", "Meeting")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockEventRepository_Instances(t *testing.T) {
	t.Run("Instances success", func(t *testing.T) {
		instances := []*calendar.Event{
			{ID: "instance1", Title: "Recurring Event"},
			{ID: "instance2", Title: "Recurring Event"},
		}
		repo := &MockEventRepository{Events: instances}

		timeMin := time.Now()
		timeMax := timeMin.Add(30 * 24 * time.Hour)
		result, err := repo.Instances(nil, "primary", "recurring-id", timeMin, timeMax)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 instances, got %d", len(result))
		}
	})

	t.Run("Instances error", func(t *testing.T) {
		repo := &MockEventRepository{InstancesErr: fmt.Errorf("instances error")}

		timeMin := time.Now()
		timeMax := timeMin.Add(30 * 24 * time.Hour)
		_, err := repo.Instances(nil, "primary", "recurring-id", timeMin, timeMax)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockEventRepository_RSVP(t *testing.T) {
	t.Run("RSVP success", func(t *testing.T) {
		repo := &MockEventRepository{}

		err := repo.RSVP(nil, "primary", "event1", "accepted")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("RSVP error", func(t *testing.T) {
		repo := &MockEventRepository{RSVPErr: fmt.Errorf("rsvp error")}

		err := repo.RSVP(nil, "primary", "event1", "accepted")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockFreeBusyRepository(t *testing.T) {
	t.Run("Query success", func(t *testing.T) {
		response := &calendar.FreeBusyResponse{
			Calendars: map[string][]*calendar.TimePeriod{
				"primary": {
					{Start: time.Now(), End: time.Now().Add(time.Hour)},
				},
			},
		}
		repo := &MockFreeBusyRepository{Response: response}

		request := &calendar.FreeBusyRequest{
			TimeMin:     time.Now(),
			TimeMax:     time.Now().Add(24 * time.Hour),
			CalendarIDs: []string{"primary"},
		}
		result, err := repo.Query(nil, request)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result.Calendars) != 1 {
			t.Errorf("expected 1 calendar, got %d", len(result.Calendars))
		}
	})

	t.Run("Query default response", func(t *testing.T) {
		repo := &MockFreeBusyRepository{}

		request := &calendar.FreeBusyRequest{
			TimeMin:     time.Now(),
			TimeMax:     time.Now().Add(24 * time.Hour),
			CalendarIDs: []string{"primary"},
		}
		result, err := repo.Query(nil, request)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.Calendars == nil {
			t.Error("expected non-nil calendars map")
		}
	})

	t.Run("Query error", func(t *testing.T) {
		repo := &MockFreeBusyRepository{QueryErr: fmt.Errorf("query error")}

		request := &calendar.FreeBusyRequest{}
		_, err := repo.Query(nil, request)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
