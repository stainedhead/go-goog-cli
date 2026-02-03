package mail

import (
	"errors"
	"testing"
)

func TestDomainErrors(t *testing.T) {
	// Test that domain errors are defined and distinct
	domainErrors := []error{
		ErrMessageNotFound,
		ErrDraftNotFound,
		ErrThreadNotFound,
		ErrLabelNotFound,
		ErrFilterNotFound,
	}

	for i, err1 := range domainErrors {
		if err1 == nil {
			t.Errorf("domain error at index %d is nil", i)
			continue
		}
		for j, err2 := range domainErrors {
			if i != j && errors.Is(err1, err2) {
				t.Errorf("domain errors at index %d and %d should not be equal", i, j)
			}
		}
	}
}

func TestErrMessageNotFound(t *testing.T) {
	if ErrMessageNotFound.Error() != "message not found" {
		t.Errorf("expected 'message not found', got '%s'", ErrMessageNotFound.Error())
	}
}

func TestErrDraftNotFound(t *testing.T) {
	if ErrDraftNotFound.Error() != "draft not found" {
		t.Errorf("expected 'draft not found', got '%s'", ErrDraftNotFound.Error())
	}
}

func TestErrThreadNotFound(t *testing.T) {
	if ErrThreadNotFound.Error() != "thread not found" {
		t.Errorf("expected 'thread not found', got '%s'", ErrThreadNotFound.Error())
	}
}

func TestErrLabelNotFound(t *testing.T) {
	if ErrLabelNotFound.Error() != "label not found" {
		t.Errorf("expected 'label not found', got '%s'", ErrLabelNotFound.Error())
	}
}

func TestErrFilterNotFound(t *testing.T) {
	if ErrFilterNotFound.Error() != "filter not found" {
		t.Errorf("expected 'filter not found', got '%s'", ErrFilterNotFound.Error())
	}
}

func TestListOptions(t *testing.T) {
	opts := ListOptions{
		MaxResults: 10,
		PageToken:  "next-page-token",
		Query:      "is:unread",
		LabelIDs:   []string{"INBOX", "IMPORTANT"},
	}

	if opts.MaxResults != 10 {
		t.Errorf("expected MaxResults 10, got %d", opts.MaxResults)
	}
	if opts.PageToken != "next-page-token" {
		t.Errorf("expected PageToken 'next-page-token', got '%s'", opts.PageToken)
	}
	if opts.Query != "is:unread" {
		t.Errorf("expected Query 'is:unread', got '%s'", opts.Query)
	}
	if len(opts.LabelIDs) != 2 {
		t.Errorf("expected 2 LabelIDs, got %d", len(opts.LabelIDs))
	}
}

func TestListResult(t *testing.T) {
	messages := []*Message{
		NewMessage("1", "t1", "from@example.com", "Subject 1", "Body 1"),
		NewMessage("2", "t2", "from@example.com", "Subject 2", "Body 2"),
	}

	result := ListResult[*Message]{
		Items:         messages,
		NextPageToken: "page-2",
		Total:         100,
	}

	if len(result.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(result.Items))
	}
	if result.NextPageToken != "page-2" {
		t.Errorf("expected NextPageToken 'page-2', got '%s'", result.NextPageToken)
	}
	if result.Total != 100 {
		t.Errorf("expected Total 100, got %d", result.Total)
	}
}

func TestModifyRequest(t *testing.T) {
	req := ModifyRequest{
		AddLabels:    []string{"IMPORTANT", "STARRED"},
		RemoveLabels: []string{"INBOX"},
	}

	if len(req.AddLabels) != 2 {
		t.Errorf("expected 2 AddLabels, got %d", len(req.AddLabels))
	}
	if len(req.RemoveLabels) != 1 {
		t.Errorf("expected 1 RemoveLabels, got %d", len(req.RemoveLabels))
	}
}

func TestVacationSettings(t *testing.T) {
	settings := VacationSettings{
		EnableAutoReply:    true,
		StartTime:          1609459200000, // 2021-01-01
		EndTime:            1609545600000, // 2021-01-02
		ResponseSubject:    "Out of Office",
		ResponseBodyPlain:  "I am currently out of the office.",
		ResponseBodyHTML:   "<p>I am currently out of the office.</p>",
		RestrictToContacts: true,
		RestrictToDomain:   false,
	}

	if !settings.EnableAutoReply {
		t.Error("expected EnableAutoReply to be true")
	}
	if settings.StartTime != 1609459200000 {
		t.Errorf("expected StartTime 1609459200000, got %d", settings.StartTime)
	}
	if settings.EndTime != 1609545600000 {
		t.Errorf("expected EndTime 1609545600000, got %d", settings.EndTime)
	}
	if settings.ResponseSubject != "Out of Office" {
		t.Errorf("expected ResponseSubject 'Out of Office', got '%s'", settings.ResponseSubject)
	}
	if settings.ResponseBodyPlain != "I am currently out of the office." {
		t.Errorf("expected ResponseBodyPlain to match, got '%s'", settings.ResponseBodyPlain)
	}
	if settings.ResponseBodyHTML != "<p>I am currently out of the office.</p>" {
		t.Errorf("expected ResponseBodyHTML to match, got '%s'", settings.ResponseBodyHTML)
	}
	if !settings.RestrictToContacts {
		t.Error("expected RestrictToContacts to be true")
	}
	if settings.RestrictToDomain {
		t.Error("expected RestrictToDomain to be false")
	}
}

// Compile-time interface implementation checks would go here
// if we had concrete implementations. The interfaces are tested
// implicitly by their method signatures.
