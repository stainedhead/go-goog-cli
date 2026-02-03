package calendar

import (
	"testing"
	"time"
)

func TestNewEvent(t *testing.T) {
	title := "Test Meeting"
	start := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)

	event := NewEvent(title, start, end)

	if event.Title != title {
		t.Errorf("expected title %q, got %q", title, event.Title)
	}

	if !event.Start.Equal(start) {
		t.Errorf("expected start %v, got %v", start, event.Start)
	}

	if !event.End.Equal(end) {
		t.Errorf("expected end %v, got %v", end, event.End)
	}

	if event.AllDay {
		t.Error("expected AllDay to be false")
	}

	if event.Status != StatusConfirmed {
		t.Errorf("expected status %q, got %q", StatusConfirmed, event.Status)
	}

	if event.Visibility != VisibilityPrivate {
		t.Errorf("expected visibility %q, got %q", VisibilityPrivate, event.Visibility)
	}

	if event.Attendees == nil {
		t.Error("expected Attendees to be initialized")
	}

	if event.Reminders == nil {
		t.Error("expected Reminders to be initialized")
	}
}

func TestNewAllDayEvent(t *testing.T) {
	title := "All Day Event"
	date := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)

	event := NewAllDayEvent(title, date)

	if event.Title != title {
		t.Errorf("expected title %q, got %q", title, event.Title)
	}

	if !event.AllDay {
		t.Error("expected AllDay to be true")
	}

	// Start should be at midnight
	expectedStart := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	if !event.Start.Equal(expectedStart) {
		t.Errorf("expected start %v, got %v", expectedStart, event.Start)
	}

	// End should be next day at midnight
	expectedEnd := time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC)
	if !event.End.Equal(expectedEnd) {
		t.Errorf("expected end %v, got %v", expectedEnd, event.End)
	}
}

func TestNewReminder(t *testing.T) {
	method := ReminderMethodPopup
	minutes := 15

	reminder := NewReminder(method, minutes)

	if reminder.Method != method {
		t.Errorf("expected method %q, got %q", method, reminder.Method)
	}

	if reminder.Minutes != minutes {
		t.Errorf("expected minutes %d, got %d", minutes, reminder.Minutes)
	}
}

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		status string
		valid  bool
	}{
		{StatusConfirmed, true},
		{StatusTentative, true},
		{StatusCancelled, true},
		{"confirmed", true},
		{"tentative", true},
		{"cancelled", true},
		{"invalid", false},
		{"", false},
		{"CONFIRMED", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := IsValidStatus(tt.status)
			if got != tt.valid {
				t.Errorf("IsValidStatus(%q) = %v, want %v", tt.status, got, tt.valid)
			}
		})
	}
}

func TestIsValidVisibility(t *testing.T) {
	tests := []struct {
		visibility string
		valid      bool
	}{
		{VisibilityPublic, true},
		{VisibilityPrivate, true},
		{"public", true},
		{"private", true},
		{"invalid", false},
		{"", false},
		{"PUBLIC", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.visibility, func(t *testing.T) {
			got := IsValidVisibility(tt.visibility)
			if got != tt.valid {
				t.Errorf("IsValidVisibility(%q) = %v, want %v", tt.visibility, got, tt.valid)
			}
		})
	}
}

func TestIsValidReminderMethod(t *testing.T) {
	tests := []struct {
		method string
		valid  bool
	}{
		{ReminderMethodEmail, true},
		{ReminderMethodPopup, true},
		{"email", true},
		{"popup", true},
		{"sms", false},
		{"", false},
		{"EMAIL", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			got := IsValidReminderMethod(tt.method)
			if got != tt.valid {
				t.Errorf("IsValidReminderMethod(%q) = %v, want %v", tt.method, got, tt.valid)
			}
		})
	}
}

func TestEventDuration(t *testing.T) {
	start := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 11, 30, 0, 0, time.UTC)
	event := NewEvent("Test", start, end)

	expected := 90 * time.Minute
	if event.Duration() != expected {
		t.Errorf("expected duration %v, got %v", expected, event.Duration())
	}
}

func TestEventIsRecurring(t *testing.T) {
	event := NewEvent("Test", time.Now(), time.Now().Add(time.Hour))

	if event.IsRecurring() {
		t.Error("expected IsRecurring to be false for event without recurrence")
	}

	event.Recurrence = []string{"RRULE:FREQ=WEEKLY;BYDAY=MO"}

	if !event.IsRecurring() {
		t.Error("expected IsRecurring to be true for event with recurrence")
	}
}

func TestEventHasConference(t *testing.T) {
	event := NewEvent("Test", time.Now(), time.Now().Add(time.Hour))

	if event.HasConference() {
		t.Error("expected HasConference to be false for event without conference data")
	}

	event.ConferenceData = &ConferenceData{}
	if event.HasConference() {
		t.Error("expected HasConference to be false for event with empty conference URI")
	}

	event.ConferenceData = &ConferenceData{
		Type: "hangoutsMeet",
		URI:  "https://meet.google.com/abc-defg-hij",
	}

	if !event.HasConference() {
		t.Error("expected HasConference to be true for event with conference data")
	}
}

func TestEventAddAttendee(t *testing.T) {
	event := NewEvent("Test", time.Now(), time.Now().Add(time.Hour))
	attendee := NewAttendee("test@example.com")

	event.AddAttendee(attendee)

	if len(event.Attendees) != 1 {
		t.Errorf("expected 1 attendee, got %d", len(event.Attendees))
	}

	if event.Attendees[0].Email != "test@example.com" {
		t.Errorf("expected attendee email %q, got %q", "test@example.com", event.Attendees[0].Email)
	}
}

func TestEventAddAttendeeNilSlice(t *testing.T) {
	event := &Event{Title: "Test"}
	attendee := NewAttendee("test@example.com")

	event.AddAttendee(attendee)

	if len(event.Attendees) != 1 {
		t.Errorf("expected 1 attendee, got %d", len(event.Attendees))
	}
}

func TestEventAddReminder(t *testing.T) {
	event := NewEvent("Test", time.Now(), time.Now().Add(time.Hour))
	reminder := NewReminder(ReminderMethodPopup, 10)

	event.AddReminder(reminder)

	if len(event.Reminders) != 1 {
		t.Errorf("expected 1 reminder, got %d", len(event.Reminders))
	}

	if event.Reminders[0].Minutes != 10 {
		t.Errorf("expected reminder minutes %d, got %d", 10, event.Reminders[0].Minutes)
	}
}

func TestEventAddReminderNilSlice(t *testing.T) {
	event := &Event{Title: "Test"}
	reminder := NewReminder(ReminderMethodEmail, 30)

	event.AddReminder(reminder)

	if len(event.Reminders) != 1 {
		t.Errorf("expected 1 reminder, got %d", len(event.Reminders))
	}
}

func TestConferenceData(t *testing.T) {
	conf := &ConferenceData{
		Type: "hangoutsMeet",
		URI:  "https://meet.google.com/abc-defg-hij",
	}

	if conf.Type != "hangoutsMeet" {
		t.Errorf("unexpected Type: %s", conf.Type)
	}

	if conf.URI != "https://meet.google.com/abc-defg-hij" {
		t.Errorf("unexpected URI: %s", conf.URI)
	}
}

func TestEventFields(t *testing.T) {
	now := time.Now()
	event := &Event{
		ID:          "event123",
		CalendarID:  "primary",
		Title:       "Team Meeting",
		Description: "Weekly sync",
		Location:    "Conference Room A",
		Start:       now,
		End:         now.Add(time.Hour),
		AllDay:      false,
		Recurrence:  []string{"RRULE:FREQ=WEEKLY"},
		Status:      StatusConfirmed,
		Visibility:  VisibilityPublic,
		ColorID:     "1",
		Created:     now,
		Updated:     now,
		HTMLLink:    "https://calendar.google.com/event/123",
	}

	if event.ID != "event123" {
		t.Errorf("unexpected ID: %s", event.ID)
	}

	if event.CalendarID != "primary" {
		t.Errorf("unexpected CalendarID: %s", event.CalendarID)
	}

	if event.Location != "Conference Room A" {
		t.Errorf("unexpected Location: %s", event.Location)
	}

	if event.Description != "Weekly sync" {
		t.Errorf("unexpected Description: %s", event.Description)
	}

	if event.ColorID != "1" {
		t.Errorf("unexpected ColorID: %s", event.ColorID)
	}

	if event.HTMLLink != "https://calendar.google.com/event/123" {
		t.Errorf("unexpected HTMLLink: %s", event.HTMLLink)
	}
}
