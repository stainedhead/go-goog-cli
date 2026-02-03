package calendar

import "testing"

func TestNewAttendee(t *testing.T) {
	email := "test@example.com"
	attendee := NewAttendee(email)

	if attendee.Email != email {
		t.Errorf("expected email %q, got %q", email, attendee.Email)
	}

	if attendee.ResponseStatus != ResponseNeedsAction {
		t.Errorf("expected response status %q, got %q", ResponseNeedsAction, attendee.ResponseStatus)
	}

	if attendee.Optional {
		t.Error("expected Optional to be false by default")
	}

	if attendee.Organizer {
		t.Error("expected Organizer to be false by default")
	}

	if attendee.Self {
		t.Error("expected Self to be false by default")
	}
}

func TestIsValidResponseStatus(t *testing.T) {
	tests := []struct {
		status string
		valid  bool
	}{
		{ResponseNeedsAction, true},
		{ResponseDeclined, true},
		{ResponseTentative, true},
		{ResponseAccepted, true},
		{"needsAction", true},
		{"declined", true},
		{"tentative", true},
		{"accepted", true},
		{"invalid", false},
		{"", false},
		{"ACCEPTED", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := IsValidResponseStatus(tt.status)
			if got != tt.valid {
				t.Errorf("IsValidResponseStatus(%q) = %v, want %v", tt.status, got, tt.valid)
			}
		})
	}
}

func TestAttendeeFields(t *testing.T) {
	attendee := &Attendee{
		Email:          "organizer@example.com",
		DisplayName:    "Event Organizer",
		ResponseStatus: ResponseAccepted,
		Optional:       false,
		Organizer:      true,
		Self:           true,
	}

	if attendee.Email != "organizer@example.com" {
		t.Errorf("unexpected Email: %s", attendee.Email)
	}

	if attendee.DisplayName != "Event Organizer" {
		t.Errorf("unexpected DisplayName: %s", attendee.DisplayName)
	}

	if attendee.ResponseStatus != ResponseAccepted {
		t.Errorf("unexpected ResponseStatus: %s", attendee.ResponseStatus)
	}

	if attendee.Optional {
		t.Error("expected Optional to be false")
	}

	if !attendee.Organizer {
		t.Error("expected Organizer to be true")
	}

	if !attendee.Self {
		t.Error("expected Self to be true")
	}
}
