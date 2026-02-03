package mail

import "testing"

func TestNewFilter(t *testing.T) {
	filter := NewFilter("filter-123")

	if filter.ID != "filter-123" {
		t.Errorf("expected ID 'filter-123', got '%s'", filter.ID)
	}
	if filter.Criteria == nil {
		t.Error("expected Criteria to be initialized")
	}
	if filter.Action == nil {
		t.Error("expected Action to be initialized")
	}
	if filter.Action.AddLabels == nil {
		t.Error("expected AddLabels to be initialized")
	}
	if filter.Action.RemoveLabels == nil {
		t.Error("expected RemoveLabels to be initialized")
	}
}

func TestNewFilterWithCriteria(t *testing.T) {
	criteria := &FilterCriteria{
		From:    "sender@example.com",
		Subject: "Test",
	}
	filter := NewFilterWithCriteria("filter-123", criteria)

	if filter.ID != "filter-123" {
		t.Errorf("expected ID 'filter-123', got '%s'", filter.ID)
	}
	if filter.Criteria != criteria {
		t.Error("expected Criteria to be the provided criteria")
	}
	if filter.Criteria.From != "sender@example.com" {
		t.Errorf("expected From 'sender@example.com', got '%s'", filter.Criteria.From)
	}
}

func TestFilter_SetCriteria(t *testing.T) {
	filter := NewFilter("filter-123")
	criteria := &FilterCriteria{
		From: "new@example.com",
	}

	filter.SetCriteria(criteria)

	if filter.Criteria != criteria {
		t.Error("expected Criteria to be updated")
	}
}

func TestFilter_SetAction(t *testing.T) {
	filter := NewFilter("filter-123")
	action := &FilterAction{
		Archive:  true,
		MarkRead: true,
	}

	filter.SetAction(action)

	if filter.Action != action {
		t.Error("expected Action to be updated")
	}
}

func TestFilter_HasCriteria(t *testing.T) {
	filter := NewFilter("filter-123")

	if filter.HasCriteria() {
		t.Error("expected HasCriteria to return false for empty criteria")
	}

	filter.Criteria.From = "sender@example.com"
	if !filter.HasCriteria() {
		t.Error("expected HasCriteria to return true when From is set")
	}

	filter.Criteria = &FilterCriteria{To: "recipient@example.com"}
	if !filter.HasCriteria() {
		t.Error("expected HasCriteria to return true when To is set")
	}

	filter.Criteria = &FilterCriteria{Subject: "Test Subject"}
	if !filter.HasCriteria() {
		t.Error("expected HasCriteria to return true when Subject is set")
	}

	filter.Criteria = &FilterCriteria{Query: "is:important"}
	if !filter.HasCriteria() {
		t.Error("expected HasCriteria to return true when Query is set")
	}

	filter.Criteria = nil
	if filter.HasCriteria() {
		t.Error("expected HasCriteria to return false when Criteria is nil")
	}
}

func TestFilter_HasAction(t *testing.T) {
	filter := NewFilter("filter-123")

	if filter.HasAction() {
		t.Error("expected HasAction to return false for empty action")
	}

	filter.Action.AddLabels = []string{"IMPORTANT"}
	if !filter.HasAction() {
		t.Error("expected HasAction to return true when AddLabels is set")
	}

	filter.Action = &FilterAction{RemoveLabels: []string{"INBOX"}}
	if !filter.HasAction() {
		t.Error("expected HasAction to return true when RemoveLabels is set")
	}

	filter.Action = &FilterAction{Forward: "forward@example.com"}
	if !filter.HasAction() {
		t.Error("expected HasAction to return true when Forward is set")
	}

	filter.Action = &FilterAction{Archive: true}
	if !filter.HasAction() {
		t.Error("expected HasAction to return true when Archive is true")
	}

	filter.Action = &FilterAction{MarkRead: true}
	if !filter.HasAction() {
		t.Error("expected HasAction to return true when MarkRead is true")
	}

	filter.Action = &FilterAction{Star: true}
	if !filter.HasAction() {
		t.Error("expected HasAction to return true when Star is true")
	}

	filter.Action = &FilterAction{Trash: true}
	if !filter.HasAction() {
		t.Error("expected HasAction to return true when Trash is true")
	}

	filter.Action = nil
	if filter.HasAction() {
		t.Error("expected HasAction to return false when Action is nil")
	}
}

func TestFilter_IsValid(t *testing.T) {
	filter := NewFilter("filter-123")

	if filter.IsValid() {
		t.Error("expected IsValid to return false for filter without criteria and actions")
	}

	filter.Criteria.From = "sender@example.com"
	if filter.IsValid() {
		t.Error("expected IsValid to return false for filter with criteria but no actions")
	}

	filter.Action.Archive = true
	if !filter.IsValid() {
		t.Error("expected IsValid to return true for filter with both criteria and actions")
	}

	filter.Criteria = &FilterCriteria{}
	if filter.IsValid() {
		t.Error("expected IsValid to return false for filter with actions but no criteria")
	}
}

func TestFilter_AddLabelToAction(t *testing.T) {
	filter := NewFilter("filter-123")

	filter.AddLabelToAction("IMPORTANT")
	filter.AddLabelToAction("STARRED")

	if len(filter.Action.AddLabels) != 2 {
		t.Errorf("expected 2 labels in AddLabels, got %d", len(filter.Action.AddLabels))
	}
	if filter.Action.AddLabels[0] != "IMPORTANT" {
		t.Errorf("expected first label 'IMPORTANT', got '%s'", filter.Action.AddLabels[0])
	}
	if filter.Action.AddLabels[1] != "STARRED" {
		t.Errorf("expected second label 'STARRED', got '%s'", filter.Action.AddLabels[1])
	}
}

func TestFilter_AddLabelToAction_NilAction(t *testing.T) {
	filter := &Filter{ID: "filter-123"}

	filter.AddLabelToAction("IMPORTANT")

	if filter.Action == nil {
		t.Error("expected Action to be initialized")
	}
	if len(filter.Action.AddLabels) != 1 {
		t.Errorf("expected 1 label in AddLabels, got %d", len(filter.Action.AddLabels))
	}
}

func TestFilter_AddRemoveLabelToAction(t *testing.T) {
	filter := NewFilter("filter-123")

	filter.AddRemoveLabelToAction("INBOX")
	filter.AddRemoveLabelToAction("UNREAD")

	if len(filter.Action.RemoveLabels) != 2 {
		t.Errorf("expected 2 labels in RemoveLabels, got %d", len(filter.Action.RemoveLabels))
	}
	if filter.Action.RemoveLabels[0] != "INBOX" {
		t.Errorf("expected first label 'INBOX', got '%s'", filter.Action.RemoveLabels[0])
	}
}

func TestFilter_AddRemoveLabelToAction_NilAction(t *testing.T) {
	filter := &Filter{ID: "filter-123"}

	filter.AddRemoveLabelToAction("INBOX")

	if filter.Action == nil {
		t.Error("expected Action to be initialized")
	}
	if len(filter.Action.RemoveLabels) != 1 {
		t.Errorf("expected 1 label in RemoveLabels, got %d", len(filter.Action.RemoveLabels))
	}
}

func TestFilterCriteria(t *testing.T) {
	criteria := &FilterCriteria{
		From:    "sender@example.com",
		To:      "recipient@example.com",
		Subject: "Test Subject",
		Query:   "is:important",
	}

	if criteria.From != "sender@example.com" {
		t.Errorf("expected From 'sender@example.com', got '%s'", criteria.From)
	}
	if criteria.To != "recipient@example.com" {
		t.Errorf("expected To 'recipient@example.com', got '%s'", criteria.To)
	}
	if criteria.Subject != "Test Subject" {
		t.Errorf("expected Subject 'Test Subject', got '%s'", criteria.Subject)
	}
	if criteria.Query != "is:important" {
		t.Errorf("expected Query 'is:important', got '%s'", criteria.Query)
	}
}

func TestFilterAction(t *testing.T) {
	action := &FilterAction{
		AddLabels:    []string{"IMPORTANT", "STARRED"},
		RemoveLabels: []string{"INBOX"},
		Forward:      "forward@example.com",
		Archive:      true,
		MarkRead:     true,
		Star:         true,
		Trash:        false,
	}

	if len(action.AddLabels) != 2 {
		t.Errorf("expected 2 AddLabels, got %d", len(action.AddLabels))
	}
	if len(action.RemoveLabels) != 1 {
		t.Errorf("expected 1 RemoveLabels, got %d", len(action.RemoveLabels))
	}
	if action.Forward != "forward@example.com" {
		t.Errorf("expected Forward 'forward@example.com', got '%s'", action.Forward)
	}
	if !action.Archive {
		t.Error("expected Archive to be true")
	}
	if !action.MarkRead {
		t.Error("expected MarkRead to be true")
	}
	if !action.Star {
		t.Error("expected Star to be true")
	}
	if action.Trash {
		t.Error("expected Trash to be false")
	}
}
