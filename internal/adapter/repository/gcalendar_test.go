package repository

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	gcal "google.golang.org/api/calendar/v3"
)

// TestParseEventDateTime tests the parseEventDateTime function.
func TestParseEventDateTime(t *testing.T) {
	tests := []struct {
		name       string
		input      *gcal.EventDateTime
		wantTime   time.Time
		wantAllDay bool
		wantErr    bool
	}{
		{
			name:       "nil input",
			input:      nil,
			wantTime:   time.Time{},
			wantAllDay: false,
			wantErr:    false,
		},
		{
			name: "datetime with RFC3339",
			input: &gcal.EventDateTime{
				DateTime: "2025-06-15T10:30:00-07:00",
				TimeZone: "America/Los_Angeles",
			},
			wantTime:   time.Date(2025, 6, 15, 10, 30, 0, 0, time.FixedZone("", -7*3600)),
			wantAllDay: false,
			wantErr:    false,
		},
		{
			name: "all-day event with date only",
			input: &gcal.EventDateTime{
				Date: "2025-06-15",
			},
			wantTime:   time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
			wantAllDay: true,
			wantErr:    false,
		},
		{
			name: "invalid datetime format",
			input: &gcal.EventDateTime{
				DateTime: "invalid-datetime",
			},
			wantTime:   time.Time{},
			wantAllDay: false,
			wantErr:    true,
		},
		{
			name: "invalid date format",
			input: &gcal.EventDateTime{
				Date: "invalid-date",
			},
			wantTime:   time.Time{},
			wantAllDay: true,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, gotAllDay, err := parseEventDateTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseEventDateTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !gotTime.Equal(tt.wantTime) {
				t.Errorf("parseEventDateTime() gotTime = %v, want %v", gotTime, tt.wantTime)
			}
			if gotAllDay != tt.wantAllDay {
				t.Errorf("parseEventDateTime() gotAllDay = %v, want %v", gotAllDay, tt.wantAllDay)
			}
		})
	}
}

// TestGcalEventToDomain tests the gcalEventToDomain conversion function.
func TestGcalEventToDomain(t *testing.T) {
	tests := []struct {
		name  string
		input *gcal.Event
		want  *calendar.Event
	}{
		{
			name:  "nil event",
			input: nil,
			want:  nil,
		},
		{
			name: "basic event",
			input: &gcal.Event{
				Id:          "event123",
				Summary:     "Test Meeting",
				Description: "A test meeting",
				Location:    "Conference Room",
				Start: &gcal.EventDateTime{
					DateTime: "2025-06-15T10:00:00Z",
				},
				End: &gcal.EventDateTime{
					DateTime: "2025-06-15T11:00:00Z",
				},
				Status:     "confirmed",
				Visibility: "private",
				HtmlLink:   "https://calendar.google.com/event?id=event123",
			},
			want: &calendar.Event{
				ID:          "event123",
				Title:       "Test Meeting",
				Description: "A test meeting",
				Location:    "Conference Room",
				Start:       time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC),
				End:         time.Date(2025, 6, 15, 11, 0, 0, 0, time.UTC),
				AllDay:      false,
				Status:      "confirmed",
				Visibility:  "private",
				HTMLLink:    "https://calendar.google.com/event?id=event123",
			},
		},
		{
			name: "all-day event",
			input: &gcal.Event{
				Id:      "allday123",
				Summary: "Holiday",
				Start: &gcal.EventDateTime{
					Date: "2025-12-25",
				},
				End: &gcal.EventDateTime{
					Date: "2025-12-26",
				},
			},
			want: &calendar.Event{
				ID:     "allday123",
				Title:  "Holiday",
				Start:  time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC),
				End:    time.Date(2025, 12, 26, 0, 0, 0, 0, time.UTC),
				AllDay: true,
			},
		},
		{
			name: "event with attendees",
			input: &gcal.Event{
				Id:      "event456",
				Summary: "Team Meeting",
				Start: &gcal.EventDateTime{
					DateTime: "2025-06-15T14:00:00Z",
				},
				End: &gcal.EventDateTime{
					DateTime: "2025-06-15T15:00:00Z",
				},
				Attendees: []*gcal.EventAttendee{
					{
						Email:          "alice@example.com",
						DisplayName:    "Alice",
						ResponseStatus: "accepted",
						Optional:       false,
						Organizer:      true,
					},
					{
						Email:          "bob@example.com",
						DisplayName:    "Bob",
						ResponseStatus: "needsAction",
						Optional:       true,
						Self:           true,
					},
				},
				Organizer: &gcal.EventOrganizer{
					Email:       "alice@example.com",
					DisplayName: "Alice",
				},
			},
			want: &calendar.Event{
				ID:    "event456",
				Title: "Team Meeting",
				Start: time.Date(2025, 6, 15, 14, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 6, 15, 15, 0, 0, 0, time.UTC),
				Attendees: []*calendar.Attendee{
					{
						Email:          "alice@example.com",
						DisplayName:    "Alice",
						ResponseStatus: "accepted",
						Optional:       false,
						Organizer:      true,
					},
					{
						Email:          "bob@example.com",
						DisplayName:    "Bob",
						ResponseStatus: "needsAction",
						Optional:       true,
						Self:           true,
					},
				},
				Organizer: &calendar.Attendee{
					Email:       "alice@example.com",
					DisplayName: "Alice",
				},
			},
		},
		{
			name: "event with recurrence",
			input: &gcal.Event{
				Id:      "recurring123",
				Summary: "Weekly Standup",
				Start: &gcal.EventDateTime{
					DateTime: "2025-06-15T09:00:00Z",
				},
				End: &gcal.EventDateTime{
					DateTime: "2025-06-15T09:30:00Z",
				},
				Recurrence: []string{
					"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR",
				},
			},
			want: &calendar.Event{
				ID:    "recurring123",
				Title: "Weekly Standup",
				Start: time.Date(2025, 6, 15, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 6, 15, 9, 30, 0, 0, time.UTC),
				Recurrence: []string{
					"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR",
				},
			},
		},
		{
			name: "event with reminders",
			input: &gcal.Event{
				Id:      "reminder123",
				Summary: "Important Meeting",
				Start: &gcal.EventDateTime{
					DateTime: "2025-06-15T15:00:00Z",
				},
				End: &gcal.EventDateTime{
					DateTime: "2025-06-15T16:00:00Z",
				},
				Reminders: &gcal.EventReminders{
					UseDefault: false,
					Overrides: []*gcal.EventReminder{
						{Method: "email", Minutes: 60},
						{Method: "popup", Minutes: 10},
					},
				},
			},
			want: &calendar.Event{
				ID:    "reminder123",
				Title: "Important Meeting",
				Start: time.Date(2025, 6, 15, 15, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 6, 15, 16, 0, 0, 0, time.UTC),
				Reminders: []*calendar.Reminder{
					{Method: "email", Minutes: 60},
					{Method: "popup", Minutes: 10},
				},
			},
		},
		{
			name: "event with conference data",
			input: &gcal.Event{
				Id:      "conf123",
				Summary: "Video Call",
				Start: &gcal.EventDateTime{
					DateTime: "2025-06-15T11:00:00Z",
				},
				End: &gcal.EventDateTime{
					DateTime: "2025-06-15T12:00:00Z",
				},
				ConferenceData: &gcal.ConferenceData{
					ConferenceSolution: &gcal.ConferenceSolution{
						Name: "Google Meet",
						Key:  &gcal.ConferenceSolutionKey{Type: "hangoutsMeet"},
					},
					EntryPoints: []*gcal.EntryPoint{
						{
							EntryPointType: "video",
							Uri:            "https://meet.google.com/abc-defg-hij",
						},
					},
				},
			},
			want: &calendar.Event{
				ID:    "conf123",
				Title: "Video Call",
				Start: time.Date(2025, 6, 15, 11, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC),
				ConferenceData: &calendar.ConferenceData{
					Type: "hangoutsMeet",
					URI:  "https://meet.google.com/abc-defg-hij",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gcalEventToDomain(tt.input)

			if tt.want == nil {
				if got != nil {
					t.Errorf("gcalEventToDomain() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Error("gcalEventToDomain() = nil, want non-nil")
				return
			}

			// Compare basic fields
			if got.ID != tt.want.ID {
				t.Errorf("ID = %v, want %v", got.ID, tt.want.ID)
			}
			if got.Title != tt.want.Title {
				t.Errorf("Title = %v, want %v", got.Title, tt.want.Title)
			}
			if got.Description != tt.want.Description {
				t.Errorf("Description = %v, want %v", got.Description, tt.want.Description)
			}
			if got.Location != tt.want.Location {
				t.Errorf("Location = %v, want %v", got.Location, tt.want.Location)
			}
			if !got.Start.Equal(tt.want.Start) {
				t.Errorf("Start = %v, want %v", got.Start, tt.want.Start)
			}
			if !got.End.Equal(tt.want.End) {
				t.Errorf("End = %v, want %v", got.End, tt.want.End)
			}
			if got.AllDay != tt.want.AllDay {
				t.Errorf("AllDay = %v, want %v", got.AllDay, tt.want.AllDay)
			}
			if got.Status != tt.want.Status {
				t.Errorf("Status = %v, want %v", got.Status, tt.want.Status)
			}
			if got.Visibility != tt.want.Visibility {
				t.Errorf("Visibility = %v, want %v", got.Visibility, tt.want.Visibility)
			}
			if got.HTMLLink != tt.want.HTMLLink {
				t.Errorf("HTMLLink = %v, want %v", got.HTMLLink, tt.want.HTMLLink)
			}

			// Compare recurrence
			if len(got.Recurrence) != len(tt.want.Recurrence) {
				t.Errorf("len(Recurrence) = %v, want %v", len(got.Recurrence), len(tt.want.Recurrence))
			} else {
				for i, r := range got.Recurrence {
					if r != tt.want.Recurrence[i] {
						t.Errorf("Recurrence[%d] = %v, want %v", i, r, tt.want.Recurrence[i])
					}
				}
			}

			// Compare attendees
			if len(got.Attendees) != len(tt.want.Attendees) {
				t.Errorf("len(Attendees) = %v, want %v", len(got.Attendees), len(tt.want.Attendees))
			} else {
				for i, a := range got.Attendees {
					want := tt.want.Attendees[i]
					if a.Email != want.Email {
						t.Errorf("Attendees[%d].Email = %v, want %v", i, a.Email, want.Email)
					}
					if a.DisplayName != want.DisplayName {
						t.Errorf("Attendees[%d].DisplayName = %v, want %v", i, a.DisplayName, want.DisplayName)
					}
					if a.ResponseStatus != want.ResponseStatus {
						t.Errorf("Attendees[%d].ResponseStatus = %v, want %v", i, a.ResponseStatus, want.ResponseStatus)
					}
					if a.Optional != want.Optional {
						t.Errorf("Attendees[%d].Optional = %v, want %v", i, a.Optional, want.Optional)
					}
					if a.Organizer != want.Organizer {
						t.Errorf("Attendees[%d].Organizer = %v, want %v", i, a.Organizer, want.Organizer)
					}
					if a.Self != want.Self {
						t.Errorf("Attendees[%d].Self = %v, want %v", i, a.Self, want.Self)
					}
				}
			}

			// Compare organizer
			if tt.want.Organizer != nil {
				if got.Organizer == nil {
					t.Error("Organizer = nil, want non-nil")
				} else {
					if got.Organizer.Email != tt.want.Organizer.Email {
						t.Errorf("Organizer.Email = %v, want %v", got.Organizer.Email, tt.want.Organizer.Email)
					}
					if got.Organizer.DisplayName != tt.want.Organizer.DisplayName {
						t.Errorf("Organizer.DisplayName = %v, want %v", got.Organizer.DisplayName, tt.want.Organizer.DisplayName)
					}
				}
			}

			// Compare reminders
			if len(got.Reminders) != len(tt.want.Reminders) {
				t.Errorf("len(Reminders) = %v, want %v", len(got.Reminders), len(tt.want.Reminders))
			} else {
				for i, r := range got.Reminders {
					want := tt.want.Reminders[i]
					if r.Method != want.Method {
						t.Errorf("Reminders[%d].Method = %v, want %v", i, r.Method, want.Method)
					}
					if r.Minutes != want.Minutes {
						t.Errorf("Reminders[%d].Minutes = %v, want %v", i, r.Minutes, want.Minutes)
					}
				}
			}

			// Compare conference data
			if tt.want.ConferenceData != nil {
				if got.ConferenceData == nil {
					t.Error("ConferenceData = nil, want non-nil")
				} else {
					if got.ConferenceData.Type != tt.want.ConferenceData.Type {
						t.Errorf("ConferenceData.Type = %v, want %v", got.ConferenceData.Type, tt.want.ConferenceData.Type)
					}
					if got.ConferenceData.URI != tt.want.ConferenceData.URI {
						t.Errorf("ConferenceData.URI = %v, want %v", got.ConferenceData.URI, tt.want.ConferenceData.URI)
					}
				}
			}
		})
	}
}

// TestDomainEventToGcal tests the domainEventToGcal conversion function.
func TestDomainEventToGcal(t *testing.T) {
	tests := []struct {
		name  string
		input *calendar.Event
		want  *gcal.Event
	}{
		{
			name:  "nil event",
			input: nil,
			want:  nil,
		},
		{
			name: "basic event",
			input: &calendar.Event{
				ID:          "event123",
				Title:       "Test Meeting",
				Description: "A test meeting",
				Location:    "Conference Room",
				Start:       time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC),
				End:         time.Date(2025, 6, 15, 11, 0, 0, 0, time.UTC),
				AllDay:      false,
				Status:      "confirmed",
				Visibility:  "private",
				ColorID:     "5",
			},
			want: &gcal.Event{
				Id:          "event123",
				Summary:     "Test Meeting",
				Description: "A test meeting",
				Location:    "Conference Room",
				Start: &gcal.EventDateTime{
					DateTime: "2025-06-15T10:00:00Z",
				},
				End: &gcal.EventDateTime{
					DateTime: "2025-06-15T11:00:00Z",
				},
				Status:     "confirmed",
				Visibility: "private",
				ColorId:    "5",
			},
		},
		{
			name: "all-day event",
			input: &calendar.Event{
				ID:     "allday123",
				Title:  "Holiday",
				Start:  time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC),
				End:    time.Date(2025, 12, 26, 0, 0, 0, 0, time.UTC),
				AllDay: true,
			},
			want: &gcal.Event{
				Id:      "allday123",
				Summary: "Holiday",
				Start: &gcal.EventDateTime{
					Date: "2025-12-25",
				},
				End: &gcal.EventDateTime{
					Date: "2025-12-26",
				},
			},
		},
		{
			name: "event with attendees",
			input: &calendar.Event{
				ID:    "event456",
				Title: "Team Meeting",
				Start: time.Date(2025, 6, 15, 14, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 6, 15, 15, 0, 0, 0, time.UTC),
				Attendees: []*calendar.Attendee{
					{
						Email:          "alice@example.com",
						DisplayName:    "Alice",
						ResponseStatus: "accepted",
						Optional:       false,
						Organizer:      true,
					},
					{
						Email:          "bob@example.com",
						DisplayName:    "Bob",
						ResponseStatus: "needsAction",
						Optional:       true,
					},
				},
			},
			want: &gcal.Event{
				Id:      "event456",
				Summary: "Team Meeting",
				Start: &gcal.EventDateTime{
					DateTime: "2025-06-15T14:00:00Z",
				},
				End: &gcal.EventDateTime{
					DateTime: "2025-06-15T15:00:00Z",
				},
				Attendees: []*gcal.EventAttendee{
					{
						Email:          "alice@example.com",
						DisplayName:    "Alice",
						ResponseStatus: "accepted",
						Optional:       false,
						Organizer:      true,
					},
					{
						Email:          "bob@example.com",
						DisplayName:    "Bob",
						ResponseStatus: "needsAction",
						Optional:       true,
					},
				},
			},
		},
		{
			name: "event with recurrence",
			input: &calendar.Event{
				ID:    "recurring123",
				Title: "Weekly Standup",
				Start: time.Date(2025, 6, 15, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 6, 15, 9, 30, 0, 0, time.UTC),
				Recurrence: []string{
					"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR",
				},
			},
			want: &gcal.Event{
				Id:      "recurring123",
				Summary: "Weekly Standup",
				Start: &gcal.EventDateTime{
					DateTime: "2025-06-15T09:00:00Z",
				},
				End: &gcal.EventDateTime{
					DateTime: "2025-06-15T09:30:00Z",
				},
				Recurrence: []string{
					"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR",
				},
			},
		},
		{
			name: "event with reminders",
			input: &calendar.Event{
				ID:    "reminder123",
				Title: "Important Meeting",
				Start: time.Date(2025, 6, 15, 15, 0, 0, 0, time.UTC),
				End:   time.Date(2025, 6, 15, 16, 0, 0, 0, time.UTC),
				Reminders: []*calendar.Reminder{
					{Method: "email", Minutes: 60},
					{Method: "popup", Minutes: 10},
				},
			},
			want: &gcal.Event{
				Id:      "reminder123",
				Summary: "Important Meeting",
				Start: &gcal.EventDateTime{
					DateTime: "2025-06-15T15:00:00Z",
				},
				End: &gcal.EventDateTime{
					DateTime: "2025-06-15T16:00:00Z",
				},
				Reminders: &gcal.EventReminders{
					UseDefault: false,
					Overrides: []*gcal.EventReminder{
						{Method: "email", Minutes: 60},
						{Method: "popup", Minutes: 10},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domainEventToGcal(tt.input)

			if tt.want == nil {
				if got != nil {
					t.Errorf("domainEventToGcal() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Error("domainEventToGcal() = nil, want non-nil")
				return
			}

			// Compare basic fields
			if got.Id != tt.want.Id {
				t.Errorf("Id = %v, want %v", got.Id, tt.want.Id)
			}
			if got.Summary != tt.want.Summary {
				t.Errorf("Summary = %v, want %v", got.Summary, tt.want.Summary)
			}
			if got.Description != tt.want.Description {
				t.Errorf("Description = %v, want %v", got.Description, tt.want.Description)
			}
			if got.Location != tt.want.Location {
				t.Errorf("Location = %v, want %v", got.Location, tt.want.Location)
			}
			if got.Status != tt.want.Status {
				t.Errorf("Status = %v, want %v", got.Status, tt.want.Status)
			}
			if got.Visibility != tt.want.Visibility {
				t.Errorf("Visibility = %v, want %v", got.Visibility, tt.want.Visibility)
			}
			if got.ColorId != tt.want.ColorId {
				t.Errorf("ColorId = %v, want %v", got.ColorId, tt.want.ColorId)
			}

			// Compare start/end
			if tt.input.AllDay {
				if got.Start.Date != tt.want.Start.Date {
					t.Errorf("Start.Date = %v, want %v", got.Start.Date, tt.want.Start.Date)
				}
				if got.End.Date != tt.want.End.Date {
					t.Errorf("End.Date = %v, want %v", got.End.Date, tt.want.End.Date)
				}
			} else {
				if got.Start.DateTime != tt.want.Start.DateTime {
					t.Errorf("Start.DateTime = %v, want %v", got.Start.DateTime, tt.want.Start.DateTime)
				}
				if got.End.DateTime != tt.want.End.DateTime {
					t.Errorf("End.DateTime = %v, want %v", got.End.DateTime, tt.want.End.DateTime)
				}
			}

			// Compare recurrence
			if len(got.Recurrence) != len(tt.want.Recurrence) {
				t.Errorf("len(Recurrence) = %v, want %v", len(got.Recurrence), len(tt.want.Recurrence))
			} else {
				for i, r := range got.Recurrence {
					if r != tt.want.Recurrence[i] {
						t.Errorf("Recurrence[%d] = %v, want %v", i, r, tt.want.Recurrence[i])
					}
				}
			}

			// Compare attendees
			if len(got.Attendees) != len(tt.want.Attendees) {
				t.Errorf("len(Attendees) = %v, want %v", len(got.Attendees), len(tt.want.Attendees))
			} else {
				for i, a := range got.Attendees {
					want := tt.want.Attendees[i]
					if a.Email != want.Email {
						t.Errorf("Attendees[%d].Email = %v, want %v", i, a.Email, want.Email)
					}
					if a.DisplayName != want.DisplayName {
						t.Errorf("Attendees[%d].DisplayName = %v, want %v", i, a.DisplayName, want.DisplayName)
					}
					if a.ResponseStatus != want.ResponseStatus {
						t.Errorf("Attendees[%d].ResponseStatus = %v, want %v", i, a.ResponseStatus, want.ResponseStatus)
					}
					if a.Optional != want.Optional {
						t.Errorf("Attendees[%d].Optional = %v, want %v", i, a.Optional, want.Optional)
					}
					if a.Organizer != want.Organizer {
						t.Errorf("Attendees[%d].Organizer = %v, want %v", i, a.Organizer, want.Organizer)
					}
				}
			}

			// Compare reminders
			if tt.want.Reminders != nil {
				if got.Reminders == nil {
					t.Error("Reminders = nil, want non-nil")
				} else {
					if got.Reminders.UseDefault != tt.want.Reminders.UseDefault {
						t.Errorf("Reminders.UseDefault = %v, want %v", got.Reminders.UseDefault, tt.want.Reminders.UseDefault)
					}
					if len(got.Reminders.Overrides) != len(tt.want.Reminders.Overrides) {
						t.Errorf("len(Reminders.Overrides) = %v, want %v", len(got.Reminders.Overrides), len(tt.want.Reminders.Overrides))
					}
				}
			}
		})
	}
}

// TestGcalCalendarToDomain tests the gcalCalendarToDomain conversion function.
func TestGcalCalendarToDomain(t *testing.T) {
	tests := []struct {
		name  string
		input *gcal.CalendarListEntry
		want  *calendar.Calendar
	}{
		{
			name:  "nil calendar",
			input: nil,
			want:  nil,
		},
		{
			name: "basic calendar",
			input: &gcal.CalendarListEntry{
				Id:          "calendar123@group.calendar.google.com",
				Summary:     "Work Calendar",
				Description: "My work calendar",
				TimeZone:    "America/New_York",
				ColorId:     "9",
				Primary:     false,
				Selected:    true,
				AccessRole:  "owner",
			},
			want: &calendar.Calendar{
				ID:          "calendar123@group.calendar.google.com",
				Title:       "Work Calendar",
				Description: "My work calendar",
				TimeZone:    "America/New_York",
				ColorID:     "9",
				Primary:     false,
				Selected:    true,
				AccessRole:  "owner",
			},
		},
		{
			name: "primary calendar",
			input: &gcal.CalendarListEntry{
				Id:         "primary",
				Summary:    "user@gmail.com",
				TimeZone:   "Europe/London",
				Primary:    true,
				Selected:   true,
				AccessRole: "owner",
			},
			want: &calendar.Calendar{
				ID:         "primary",
				Title:      "user@gmail.com",
				TimeZone:   "Europe/London",
				Primary:    true,
				Selected:   true,
				AccessRole: "owner",
			},
		},
		{
			name: "reader access calendar",
			input: &gcal.CalendarListEntry{
				Id:         "shared123@group.calendar.google.com",
				Summary:    "Shared Calendar",
				AccessRole: "reader",
			},
			want: &calendar.Calendar{
				ID:         "shared123@group.calendar.google.com",
				Title:      "Shared Calendar",
				AccessRole: "reader",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gcalCalendarToDomain(tt.input)

			if tt.want == nil {
				if got != nil {
					t.Errorf("gcalCalendarToDomain() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Error("gcalCalendarToDomain() = nil, want non-nil")
				return
			}

			if got.ID != tt.want.ID {
				t.Errorf("ID = %v, want %v", got.ID, tt.want.ID)
			}
			if got.Title != tt.want.Title {
				t.Errorf("Title = %v, want %v", got.Title, tt.want.Title)
			}
			if got.Description != tt.want.Description {
				t.Errorf("Description = %v, want %v", got.Description, tt.want.Description)
			}
			if got.TimeZone != tt.want.TimeZone {
				t.Errorf("TimeZone = %v, want %v", got.TimeZone, tt.want.TimeZone)
			}
			if got.ColorID != tt.want.ColorID {
				t.Errorf("ColorID = %v, want %v", got.ColorID, tt.want.ColorID)
			}
			if got.Primary != tt.want.Primary {
				t.Errorf("Primary = %v, want %v", got.Primary, tt.want.Primary)
			}
			if got.Selected != tt.want.Selected {
				t.Errorf("Selected = %v, want %v", got.Selected, tt.want.Selected)
			}
			if got.AccessRole != tt.want.AccessRole {
				t.Errorf("AccessRole = %v, want %v", got.AccessRole, tt.want.AccessRole)
			}
		})
	}
}

// TestMapAPIError tests the mapAPIError function.
func TestMapAPIError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		resource string
		want     error
	}{
		{
			name:     "nil error",
			err:      nil,
			resource: "event",
			want:     nil,
		},
		{
			name:     "non-google error passthrough",
			err:      context.Canceled,
			resource: "event",
			want:     context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapAPIError(tt.err, tt.resource)
			if got != tt.want {
				t.Errorf("mapAPIError() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestErrInvalidCalendarRequest tests the ErrInvalidCalendarRequest error.
func TestErrInvalidCalendarRequest(t *testing.T) {
	if ErrInvalidCalendarRequest == nil {
		t.Error("ErrInvalidCalendarRequest should not be nil")
	}
	if ErrInvalidCalendarRequest.Error() != "invalid calendar request" {
		t.Errorf("ErrInvalidCalendarRequest.Error() = %v, want %v", ErrInvalidCalendarRequest.Error(), "invalid calendar request")
	}
}

// TestParseRecurrence tests the parseRecurrence function.
func TestParseRecurrence(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "nil input",
			input: nil,
			want:  nil,
		},
		{
			name:  "empty input",
			input: []string{},
			want:  []string{},
		},
		{
			name:  "single rule",
			input: []string{"RRULE:FREQ=DAILY"},
			want:  []string{"RRULE:FREQ=DAILY"},
		},
		{
			name: "multiple rules",
			input: []string{
				"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR",
				"EXDATE:20250615T100000Z",
			},
			want: []string{
				"RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR",
				"EXDATE:20250615T100000Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseRecurrence(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("len(parseRecurrence()) = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseRecurrence()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestDomainCalendarToGcal tests the domainCalendarToGcal conversion function.
func TestDomainCalendarToGcal(t *testing.T) {
	tests := []struct {
		name  string
		input *calendar.Calendar
		want  *gcal.Calendar
	}{
		{
			name:  "nil calendar",
			input: nil,
			want:  nil,
		},
		{
			name: "basic calendar",
			input: &calendar.Calendar{
				ID:          "calendar123",
				Title:       "My Calendar",
				Description: "A test calendar",
				TimeZone:    "America/Los_Angeles",
			},
			want: &gcal.Calendar{
				Id:          "calendar123",
				Summary:     "My Calendar",
				Description: "A test calendar",
				TimeZone:    "America/Los_Angeles",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domainCalendarToGcal(tt.input)

			if tt.want == nil {
				if got != nil {
					t.Errorf("domainCalendarToGcal() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Error("domainCalendarToGcal() = nil, want non-nil")
				return
			}

			if got.Id != tt.want.Id {
				t.Errorf("Id = %v, want %v", got.Id, tt.want.Id)
			}
			if got.Summary != tt.want.Summary {
				t.Errorf("Summary = %v, want %v", got.Summary, tt.want.Summary)
			}
			if got.Description != tt.want.Description {
				t.Errorf("Description = %v, want %v", got.Description, tt.want.Description)
			}
			if got.TimeZone != tt.want.TimeZone {
				t.Errorf("TimeZone = %v, want %v", got.TimeZone, tt.want.TimeZone)
			}
		})
	}
}

// TestGcalACLToDomain tests the gcalACLToDomain conversion function.
func TestGcalACLToDomain(t *testing.T) {
	tests := []struct {
		name  string
		input *gcal.AclRule
		want  *calendar.ACLRule
	}{
		{
			name:  "nil ACL rule",
			input: nil,
			want:  nil,
		},
		{
			name: "user scope ACL",
			input: &gcal.AclRule{
				Id:   "user:alice@example.com",
				Role: "writer",
				Scope: &gcal.AclRuleScope{
					Type:  "user",
					Value: "alice@example.com",
				},
			},
			want: &calendar.ACLRule{
				ID:   "user:alice@example.com",
				Role: "writer",
				Scope: &calendar.ACLScope{
					Type:  "user",
					Value: "alice@example.com",
				},
			},
		},
		{
			name: "domain scope ACL",
			input: &gcal.AclRule{
				Id:   "domain:example.com",
				Role: "reader",
				Scope: &gcal.AclRuleScope{
					Type:  "domain",
					Value: "example.com",
				},
			},
			want: &calendar.ACLRule{
				ID:   "domain:example.com",
				Role: "reader",
				Scope: &calendar.ACLScope{
					Type:  "domain",
					Value: "example.com",
				},
			},
		},
		{
			name: "default scope ACL",
			input: &gcal.AclRule{
				Id:   "default",
				Role: "freeBusyReader",
				Scope: &gcal.AclRuleScope{
					Type: "default",
				},
			},
			want: &calendar.ACLRule{
				ID:   "default",
				Role: "freeBusyReader",
				Scope: &calendar.ACLScope{
					Type: "default",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gcalACLToDomain(tt.input)

			if tt.want == nil {
				if got != nil {
					t.Errorf("gcalACLToDomain() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Error("gcalACLToDomain() = nil, want non-nil")
				return
			}

			if got.ID != tt.want.ID {
				t.Errorf("ID = %v, want %v", got.ID, tt.want.ID)
			}
			if got.Role != tt.want.Role {
				t.Errorf("Role = %v, want %v", got.Role, tt.want.Role)
			}
			if got.Scope == nil {
				t.Error("Scope = nil, want non-nil")
				return
			}
			if got.Scope.Type != tt.want.Scope.Type {
				t.Errorf("Scope.Type = %v, want %v", got.Scope.Type, tt.want.Scope.Type)
			}
			if got.Scope.Value != tt.want.Scope.Value {
				t.Errorf("Scope.Value = %v, want %v", got.Scope.Value, tt.want.Scope.Value)
			}
		})
	}
}

// TestDomainACLToGcal tests the domainACLToGcal conversion function.
func TestDomainACLToGcal(t *testing.T) {
	tests := []struct {
		name  string
		input *calendar.ACLRule
		want  *gcal.AclRule
	}{
		{
			name:  "nil ACL rule",
			input: nil,
			want:  nil,
		},
		{
			name: "user scope ACL",
			input: &calendar.ACLRule{
				ID:   "user:bob@example.com",
				Role: "owner",
				Scope: &calendar.ACLScope{
					Type:  "user",
					Value: "bob@example.com",
				},
			},
			want: &gcal.AclRule{
				Id:   "user:bob@example.com",
				Role: "owner",
				Scope: &gcal.AclRuleScope{
					Type:  "user",
					Value: "bob@example.com",
				},
			},
		},
		{
			name: "group scope ACL",
			input: &calendar.ACLRule{
				ID:   "group:team@example.com",
				Role: "writer",
				Scope: &calendar.ACLScope{
					Type:  "group",
					Value: "team@example.com",
				},
			},
			want: &gcal.AclRule{
				Id:   "group:team@example.com",
				Role: "writer",
				Scope: &gcal.AclRuleScope{
					Type:  "group",
					Value: "team@example.com",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domainACLToGcal(tt.input)

			if tt.want == nil {
				if got != nil {
					t.Errorf("domainACLToGcal() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Error("domainACLToGcal() = nil, want non-nil")
				return
			}

			if got.Id != tt.want.Id {
				t.Errorf("Id = %v, want %v", got.Id, tt.want.Id)
			}
			if got.Role != tt.want.Role {
				t.Errorf("Role = %v, want %v", got.Role, tt.want.Role)
			}
			if got.Scope == nil {
				t.Error("Scope = nil, want non-nil")
				return
			}
			if got.Scope.Type != tt.want.Scope.Type {
				t.Errorf("Scope.Type = %v, want %v", got.Scope.Type, tt.want.Scope.Type)
			}
			if got.Scope.Value != tt.want.Scope.Value {
				t.Errorf("Scope.Value = %v, want %v", got.Scope.Value, tt.want.Scope.Value)
			}
		})
	}
}

// =============================================================================
// Tests Using TestServer Infrastructure
// =============================================================================

// TestGCalEventRepository_ListWithTestServer tests List using the TestServer.
func TestGCalEventRepository_ListWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	now := time.Now()
	event1 := MockEventResponse("event1", "Meeting 1", "Description 1", now, now.Add(time.Hour))
	event2 := MockEventResponse("event2", "Meeting 2", "Description 2", now.Add(2*time.Hour), now.Add(3*time.Hour))

	ts.EventListHandler = func(w http.ResponseWriter, r *http.Request, calendarID string) {
		// Verify query parameters
		timeMin := r.URL.Query().Get("timeMin")
		timeMax := r.URL.Query().Get("timeMax")
		if timeMin == "" || timeMax == "" {
			WriteErrorResponse(w, http.StatusBadRequest, "timeMin and timeMax required")
			return
		}

		WriteJSONResponse(w, MockEventListResponse([]*gcal.Event{event1, event2}, ""))
	}

	service := ts.GCalService(t)
	repo := service.Events()
	ctx := context.Background()

	events, err := repo.List(ctx, "primary", now.Add(-time.Hour), now.Add(4*time.Hour))
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("events count = %d, want %d", len(events), 2)
	}

	if events[0].ID != "event1" {
		t.Errorf("events[0].ID = %q, want %q", events[0].ID, "event1")
	}
	if events[0].Title != "Meeting 1" {
		t.Errorf("events[0].Title = %q, want %q", events[0].Title, "Meeting 1")
	}
	if events[0].CalendarID != "primary" {
		t.Errorf("events[0].CalendarID = %q, want %q", events[0].CalendarID, "primary")
	}
}

// TestGCalEventRepository_GetWithTestServer tests Get using the TestServer.
func TestGCalEventRepository_GetWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	now := time.Now()
	expectedEvent := MockEventResponse("event123", "Team Meeting", "Weekly sync", now, now.Add(time.Hour))
	expectedEvent.Location = "Conference Room A"

	ts.EventGetHandler = func(w http.ResponseWriter, r *http.Request, calendarID, eventID string) {
		if calendarID == "primary" && eventID == "event123" {
			WriteJSONResponse(w, expectedEvent)
		} else {
			WriteErrorResponse(w, http.StatusNotFound, "event not found")
		}
	}

	service := ts.GCalService(t)
	repo := service.Events()
	ctx := context.Background()

	event, err := repo.Get(ctx, "primary", "event123")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if event.ID != "event123" {
		t.Errorf("ID = %q, want %q", event.ID, "event123")
	}
	if event.Title != "Team Meeting" {
		t.Errorf("Title = %q, want %q", event.Title, "Team Meeting")
	}
	if event.Description != "Weekly sync" {
		t.Errorf("Description = %q, want %q", event.Description, "Weekly sync")
	}
	if event.Location != "Conference Room A" {
		t.Errorf("Location = %q, want %q", event.Location, "Conference Room A")
	}
	if event.CalendarID != "primary" {
		t.Errorf("CalendarID = %q, want %q", event.CalendarID, "primary")
	}
}

// TestGCalEventRepository_CreateWithTestServer tests Create using the TestServer.
func TestGCalEventRepository_CreateWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	var createdEvent *gcal.Event

	ts.EventCreateHandler = func(w http.ResponseWriter, r *http.Request, calendarID string) {
		if err := json.NewDecoder(r.Body).Decode(&createdEvent); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}

		// Return the created event with a generated ID
		createdEvent.Id = "new_event_123"
		createdEvent.HtmlLink = "https://calendar.google.com/event?eid=new_event_123"
		createdEvent.Created = time.Now().Format(time.RFC3339)
		createdEvent.Updated = time.Now().Format(time.RFC3339)
		WriteJSONResponse(w, createdEvent)
	}

	service := ts.GCalService(t)
	repo := service.Events()
	ctx := context.Background()

	now := time.Now()
	event := &calendar.Event{
		Title:       "New Meeting",
		Description: "Important discussion",
		Start:       now,
		End:         now.Add(time.Hour),
	}

	created, err := repo.Create(ctx, "primary", event)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if created.ID != "new_event_123" {
		t.Errorf("ID = %q, want %q", created.ID, "new_event_123")
	}
	if created.Title != "New Meeting" {
		t.Errorf("Title = %q, want %q", created.Title, "New Meeting")
	}
	if created.CalendarID != "primary" {
		t.Errorf("CalendarID = %q, want %q", created.CalendarID, "primary")
	}

	// Verify the request was made correctly
	if createdEvent == nil {
		t.Error("expected event to be created, but none was received")
	}
	if createdEvent.Summary != "New Meeting" {
		t.Errorf("sent event Summary = %q, want %q", createdEvent.Summary, "New Meeting")
	}
}

// TestGCalEventRepository_DeleteWithTestServer tests Delete using the TestServer.
func TestGCalEventRepository_DeleteWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	deletedEventID := ""
	ts.EventDeleteHandler = func(w http.ResponseWriter, r *http.Request, calendarID, eventID string) {
		deletedEventID = eventID
		w.WriteHeader(http.StatusNoContent)
	}

	service := ts.GCalService(t)
	repo := service.Events()
	ctx := context.Background()

	err := repo.Delete(ctx, "primary", "event123")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if deletedEventID != "event123" {
		t.Errorf("deletedEventID = %q, want %q", deletedEventID, "event123")
	}
}

// TestGCalCalendarRepository_ListWithTestServer tests List using the TestServer.
func TestGCalCalendarRepository_ListWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.CalendarListHandler = StaticCalendarListHandler([]*gcal.CalendarListEntry{
		MockCalendarListEntryResponse("primary", "user@example.com", "", "America/Los_Angeles", true, "owner"),
		MockCalendarListEntryResponse("work_cal", "Work", "Work calendar", "America/New_York", false, "owner"),
		MockCalendarListEntryResponse("shared_cal", "Shared Calendar", "", "UTC", false, "reader"),
	}, "")

	service := ts.GCalService(t)
	repo := service.Calendars()
	ctx := context.Background()

	calendars, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(calendars) != 3 {
		t.Errorf("calendars count = %d, want %d", len(calendars), 3)
	}

	// Check primary calendar
	if calendars[0].ID != "primary" {
		t.Errorf("calendars[0].ID = %q, want %q", calendars[0].ID, "primary")
	}
	if calendars[0].Title != "user@example.com" {
		t.Errorf("calendars[0].Title = %q, want %q", calendars[0].Title, "user@example.com")
	}
	if !calendars[0].Primary {
		t.Error("calendars[0].Primary = false, want true")
	}
	if calendars[0].AccessRole != "owner" {
		t.Errorf("calendars[0].AccessRole = %q, want %q", calendars[0].AccessRole, "owner")
	}

	// Check shared calendar
	if calendars[2].AccessRole != "reader" {
		t.Errorf("calendars[2].AccessRole = %q, want %q", calendars[2].AccessRole, "reader")
	}
}

// TestGCalEventRepository_NotFoundWithTestServer tests Get for non-existent event.
func TestGCalEventRepository_NotFoundWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.EventGetHandler = func(w http.ResponseWriter, r *http.Request, calendarID, eventID string) {
		WriteErrorResponse(w, http.StatusNotFound, "event not found")
	}

	service := ts.GCalService(t)
	repo := service.Events()
	ctx := context.Background()

	_, err := repo.Get(ctx, "primary", "nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent event, got nil")
	}
	if err != calendar.ErrEventNotFound {
		t.Errorf("error = %v, want %v", err, calendar.ErrEventNotFound)
	}
}

// TestGCalEventRepository_InvalidTimeRangeWithTestServer tests List with invalid time range.
func TestGCalEventRepository_InvalidTimeRangeWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	service := ts.GCalService(t)
	repo := service.Events()
	ctx := context.Background()

	now := time.Now()
	// End time before start time
	_, err := repo.List(ctx, "primary", now, now.Add(-time.Hour))
	if err == nil {
		t.Fatal("expected error for invalid time range, got nil")
	}
	if err != calendar.ErrInvalidTimeRange {
		t.Errorf("error = %v, want %v", err, calendar.ErrInvalidTimeRange)
	}
}

// TestGCalEventRepository_AllDayEventWithTestServer tests handling of all-day events.
func TestGCalEventRepository_AllDayEventWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	holiday := MockAllDayEventResponse("holiday123", "Company Holiday", time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC))

	ts.EventGetHandler = func(w http.ResponseWriter, r *http.Request, calendarID, eventID string) {
		if eventID == "holiday123" {
			WriteJSONResponse(w, holiday)
		} else {
			WriteErrorResponse(w, http.StatusNotFound, "event not found")
		}
	}

	service := ts.GCalService(t)
	repo := service.Events()
	ctx := context.Background()

	event, err := repo.Get(ctx, "primary", "holiday123")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if event.ID != "holiday123" {
		t.Errorf("ID = %q, want %q", event.ID, "holiday123")
	}
	if event.Title != "Company Holiday" {
		t.Errorf("Title = %q, want %q", event.Title, "Company Holiday")
	}
	if !event.AllDay {
		t.Error("AllDay = false, want true")
	}
}

// TestGCalCalendarRepository_GetWithTestServer tests Get calendar by ID.
func TestGCalCalendarRepository_GetWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.CalendarGetHandler = func(w http.ResponseWriter, r *http.Request, calendarID string) {
		if calendarID == "work_cal" {
			WriteJSONResponse(w, MockCalendarListEntryResponse(
				"work_cal", "Work Calendar", "My work calendar", "America/New_York", false, "owner",
			))
		} else {
			WriteErrorResponse(w, http.StatusNotFound, "calendar not found")
		}
	}

	service := ts.GCalService(t)
	repo := service.Calendars()
	ctx := context.Background()

	cal, err := repo.Get(ctx, "work_cal")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if cal.ID != "work_cal" {
		t.Errorf("ID = %q, want %q", cal.ID, "work_cal")
	}
	if cal.Title != "Work Calendar" {
		t.Errorf("Title = %q, want %q", cal.Title, "Work Calendar")
	}
	if cal.Description != "My work calendar" {
		t.Errorf("Description = %q, want %q", cal.Description, "My work calendar")
	}
	if cal.TimeZone != "America/New_York" {
		t.Errorf("TimeZone = %q, want %q", cal.TimeZone, "America/New_York")
	}
}

// =============================================================================
// Additional Event Operations Tests
// =============================================================================

// TestGCalEventRepository_UpdateWithTestServer tests updating an event.
func TestGCalEventRepository_UpdateWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	var updatedEvent *gcal.Event
	ts.EventUpdateHandler = func(w http.ResponseWriter, r *http.Request, calendarID, eventID string) {
		if err := json.NewDecoder(r.Body).Decode(&updatedEvent); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}
		updatedEvent.Id = eventID
		updatedEvent.Updated = time.Now().Format(time.RFC3339)
		WriteJSONResponse(w, updatedEvent)
	}

	service := ts.GCalService(t)
	repo := service.Events()
	ctx := context.Background()

	now := time.Now()
	event := &calendar.Event{
		ID:          "event123",
		Title:       "Updated Meeting",
		Description: "Updated description",
		Start:       now,
		End:         now.Add(2 * time.Hour),
	}

	updated, err := repo.Update(ctx, "primary", event)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.ID != "event123" {
		t.Errorf("ID = %q, want %q", updated.ID, "event123")
	}
	if updated.Title != "Updated Meeting" {
		t.Errorf("Title = %q, want %q", updated.Title, "Updated Meeting")
	}
}

// TestGCalEventRepository_MoveWithTestServer tests moving an event to another calendar.
func TestGCalEventRepository_MoveWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.EventMoveHandler = func(w http.ResponseWriter, r *http.Request, calendarID, eventID string) {
		destination := r.URL.Query().Get("destination")
		if destination == "" {
			WriteErrorResponse(w, http.StatusBadRequest, "destination required")
			return
		}

		event := MockEventResponse(eventID, "Moved Event", "Description", time.Now(), time.Now().Add(time.Hour))
		WriteJSONResponse(w, event)
	}

	service := ts.GCalService(t)
	repo := service.Events()
	ctx := context.Background()

	moved, err := repo.Move(ctx, "primary", "event123", "work_calendar")
	if err != nil {
		t.Fatalf("Move failed: %v", err)
	}

	if moved.ID != "event123" {
		t.Errorf("ID = %q, want %q", moved.ID, "event123")
	}
	if moved.CalendarID != "work_calendar" {
		t.Errorf("CalendarID = %q, want %q", moved.CalendarID, "work_calendar")
	}
}

// TestGCalEventRepository_QuickAddWithTestServer tests creating an event from text.
func TestGCalEventRepository_QuickAddWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.EventQuickAddHandler = func(w http.ResponseWriter, r *http.Request, calendarID string) {
		text := r.URL.Query().Get("text")
		if text == "" {
			WriteErrorResponse(w, http.StatusBadRequest, "text required")
			return
		}

		// Simulate parsing the text into an event
		event := MockEventResponse("quickadd123", text, "", time.Now(), time.Now().Add(time.Hour))
		WriteJSONResponse(w, event)
	}

	service := ts.GCalService(t)
	repo := service.Events()
	ctx := context.Background()

	event, err := repo.QuickAdd(ctx, "primary", "Meeting with Bob tomorrow at 3pm")
	if err != nil {
		t.Fatalf("QuickAdd failed: %v", err)
	}

	if event.ID != "quickadd123" {
		t.Errorf("ID = %q, want %q", event.ID, "quickadd123")
	}
	if event.Title != "Meeting with Bob tomorrow at 3pm" {
		t.Errorf("Title = %q, want %q", event.Title, "Meeting with Bob tomorrow at 3pm")
	}
}

// TestGCalEventRepository_InstancesWithTestServer tests getting instances of a recurring event.
func TestGCalEventRepository_InstancesWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	now := time.Now()
	ts.EventInstancesHandler = func(w http.ResponseWriter, r *http.Request, calendarID, eventID string) {
		instances := []*gcal.Event{
			MockEventResponse(eventID+"_20250615", "Weekly Meeting", "Instance 1", now, now.Add(time.Hour)),
			MockEventResponse(eventID+"_20250622", "Weekly Meeting", "Instance 2", now.Add(7*24*time.Hour), now.Add(7*24*time.Hour+time.Hour)),
			MockEventResponse(eventID+"_20250629", "Weekly Meeting", "Instance 3", now.Add(14*24*time.Hour), now.Add(14*24*time.Hour+time.Hour)),
		}
		WriteJSONResponse(w, MockEventListResponse(instances, ""))
	}

	service := ts.GCalService(t)
	repo := service.Events()
	ctx := context.Background()

	instances, err := repo.Instances(ctx, "primary", "recurring123", now.Add(-time.Hour), now.Add(30*24*time.Hour))
	if err != nil {
		t.Fatalf("Instances failed: %v", err)
	}

	if len(instances) != 3 {
		t.Errorf("instances count = %d, want %d", len(instances), 3)
	}
	if instances[0].Title != "Weekly Meeting" {
		t.Errorf("instances[0].Title = %q, want %q", instances[0].Title, "Weekly Meeting")
	}
}

// TestGCalEventRepository_UpdateNotFound tests updating a non-existent event.
func TestGCalEventRepository_UpdateNotFound(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.EventUpdateHandler = func(w http.ResponseWriter, r *http.Request, calendarID, eventID string) {
		WriteErrorResponse(w, http.StatusNotFound, "event not found")
	}

	service := ts.GCalService(t)
	repo := service.Events()
	ctx := context.Background()

	event := &calendar.Event{
		ID:    "nonexistent",
		Title: "Test Event",
		Start: time.Now(),
		End:   time.Now().Add(time.Hour),
	}

	_, err := repo.Update(ctx, "primary", event)
	if err == nil {
		t.Fatal("expected error for non-existent event, got nil")
	}
	if err != calendar.ErrEventNotFound {
		t.Errorf("error = %v, want %v", err, calendar.ErrEventNotFound)
	}
}

// =============================================================================
// Calendar Repository Tests
// =============================================================================

// TestGCalCalendarRepository_CreateWithTestServer tests creating a calendar.
func TestGCalCalendarRepository_CreateWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	var createdCalendar *gcal.Calendar
	ts.CalendarCreateHandler = func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&createdCalendar); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}
		createdCalendar.Id = "new_calendar_123"
		WriteJSONResponse(w, createdCalendar)
	}

	ts.CalendarGetHandler = func(w http.ResponseWriter, r *http.Request, calendarID string) {
		WriteJSONResponse(w, MockCalendarListEntryResponse(
			calendarID, "New Calendar", "A new calendar", "America/Los_Angeles", false, "owner",
		))
	}

	service := ts.GCalService(t)
	repo := service.Calendars()
	ctx := context.Background()

	cal := &calendar.Calendar{
		Title:       "New Calendar",
		Description: "A new calendar",
		TimeZone:    "America/Los_Angeles",
	}

	created, err := repo.Create(ctx, cal)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if created.ID != "new_calendar_123" {
		t.Errorf("ID = %q, want %q", created.ID, "new_calendar_123")
	}
}

// TestGCalCalendarRepository_UpdateWithTestServer tests updating a calendar.
func TestGCalCalendarRepository_UpdateWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.CalendarUpdateHandler = func(w http.ResponseWriter, r *http.Request, calendarID string) {
		var cal gcal.Calendar
		if err := json.NewDecoder(r.Body).Decode(&cal); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}
		cal.Id = calendarID
		WriteJSONResponse(w, &cal)
	}

	ts.CalendarGetHandler = func(w http.ResponseWriter, r *http.Request, calendarID string) {
		WriteJSONResponse(w, MockCalendarListEntryResponse(
			calendarID, "Updated Calendar", "Updated description", "Europe/London", false, "owner",
		))
	}

	service := ts.GCalService(t)
	repo := service.Calendars()
	ctx := context.Background()

	cal := &calendar.Calendar{
		ID:          "calendar123",
		Title:       "Updated Calendar",
		Description: "Updated description",
		TimeZone:    "Europe/London",
	}

	updated, err := repo.Update(ctx, cal)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.ID != "calendar123" {
		t.Errorf("ID = %q, want %q", updated.ID, "calendar123")
	}
	if updated.Title != "Updated Calendar" {
		t.Errorf("Title = %q, want %q", updated.Title, "Updated Calendar")
	}
}

// TestGCalCalendarRepository_DeleteWithTestServer tests deleting a calendar.
func TestGCalCalendarRepository_DeleteWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	deletedID := ""
	ts.CalendarDeleteHandler = func(w http.ResponseWriter, r *http.Request, calendarID string) {
		deletedID = calendarID
		w.WriteHeader(http.StatusNoContent)
	}

	service := ts.GCalService(t)
	repo := service.Calendars()
	ctx := context.Background()

	err := repo.Delete(ctx, "calendar123")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if deletedID != "calendar123" {
		t.Errorf("deletedID = %q, want %q", deletedID, "calendar123")
	}
}

// TestGCalCalendarRepository_ClearWithTestServer tests clearing all events from a calendar.
func TestGCalCalendarRepository_ClearWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	clearedID := ""
	ts.CalendarClearHandler = func(w http.ResponseWriter, r *http.Request, calendarID string) {
		clearedID = calendarID
		w.WriteHeader(http.StatusNoContent)
	}

	service := ts.GCalService(t)
	repo := service.Calendars()
	ctx := context.Background()

	err := repo.Clear(ctx, "calendar123")
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	if clearedID != "calendar123" {
		t.Errorf("clearedID = %q, want %q", clearedID, "calendar123")
	}
}

// TestGCalCalendarRepository_GetNotFound tests getting a non-existent calendar.
func TestGCalCalendarRepository_GetNotFound(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.CalendarGetHandler = func(w http.ResponseWriter, r *http.Request, calendarID string) {
		WriteErrorResponse(w, http.StatusNotFound, "calendar not found")
	}

	service := ts.GCalService(t)
	repo := service.Calendars()
	ctx := context.Background()

	_, err := repo.Get(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent calendar, got nil")
	}
	if err != calendar.ErrCalendarNotFound {
		t.Errorf("error = %v, want %v", err, calendar.ErrCalendarNotFound)
	}
}

// =============================================================================
// ACL Repository Tests
// =============================================================================

// TestGCalACLRepository_ListWithTestServer tests listing ACL rules.
func TestGCalACLRepository_ListWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.ACLListHandler = func(w http.ResponseWriter, r *http.Request, calendarID string) {
		WriteJSONResponse(w, MockACLListResponse([]*gcal.AclRule{
			MockACLRuleResponse("user:alice@example.com", "writer", "user", "alice@example.com"),
			MockACLRuleResponse("user:bob@example.com", "reader", "user", "bob@example.com"),
			MockACLRuleResponse("domain:example.com", "freeBusyReader", "domain", "example.com"),
		}, ""))
	}

	service := ts.GCalService(t)
	repo := service.ACL()
	ctx := context.Background()

	rules, err := repo.List(ctx, "primary")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(rules) != 3 {
		t.Errorf("rules count = %d, want %d", len(rules), 3)
	}
	if rules[0].ID != "user:alice@example.com" {
		t.Errorf("rules[0].ID = %q, want %q", rules[0].ID, "user:alice@example.com")
	}
	if rules[0].Role != "writer" {
		t.Errorf("rules[0].Role = %q, want %q", rules[0].Role, "writer")
	}
}

// TestGCalACLRepository_GetWithTestServer tests getting an ACL rule.
func TestGCalACLRepository_GetWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.ACLGetHandler = func(w http.ResponseWriter, r *http.Request, calendarID, ruleID string) {
		if ruleID == "user:alice@example.com" {
			WriteJSONResponse(w, MockACLRuleResponse("user:alice@example.com", "writer", "user", "alice@example.com"))
		} else {
			WriteErrorResponse(w, http.StatusNotFound, "ACL rule not found")
		}
	}

	service := ts.GCalService(t)
	repo := service.ACL()
	ctx := context.Background()

	rule, err := repo.Get(ctx, "primary", "user:alice@example.com")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if rule.ID != "user:alice@example.com" {
		t.Errorf("ID = %q, want %q", rule.ID, "user:alice@example.com")
	}
	if rule.Role != "writer" {
		t.Errorf("Role = %q, want %q", rule.Role, "writer")
	}
	if rule.Scope.Type != "user" {
		t.Errorf("Scope.Type = %q, want %q", rule.Scope.Type, "user")
	}
}

// TestGCalACLRepository_InsertWithTestServer tests inserting an ACL rule.
func TestGCalACLRepository_InsertWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	var insertedRule *gcal.AclRule
	ts.ACLInsertHandler = func(w http.ResponseWriter, r *http.Request, calendarID string) {
		if err := json.NewDecoder(r.Body).Decode(&insertedRule); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}
		insertedRule.Id = "user:" + insertedRule.Scope.Value
		WriteJSONResponse(w, insertedRule)
	}

	service := ts.GCalService(t)
	repo := service.ACL()
	ctx := context.Background()

	rule := &calendar.ACLRule{
		Role: "writer",
		Scope: &calendar.ACLScope{
			Type:  "user",
			Value: "charlie@example.com",
		},
	}

	created, err := repo.Insert(ctx, "primary", rule)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	if created.ID != "user:charlie@example.com" {
		t.Errorf("ID = %q, want %q", created.ID, "user:charlie@example.com")
	}
	if created.Role != "writer" {
		t.Errorf("Role = %q, want %q", created.Role, "writer")
	}
}

// TestGCalACLRepository_UpdateWithTestServer tests updating an ACL rule.
func TestGCalACLRepository_UpdateWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.ACLUpdateHandler = func(w http.ResponseWriter, r *http.Request, calendarID, ruleID string) {
		var rule gcal.AclRule
		if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}
		rule.Id = ruleID
		WriteJSONResponse(w, &rule)
	}

	service := ts.GCalService(t)
	repo := service.ACL()
	ctx := context.Background()

	rule := &calendar.ACLRule{
		ID:   "user:alice@example.com",
		Role: "owner",
		Scope: &calendar.ACLScope{
			Type:  "user",
			Value: "alice@example.com",
		},
	}

	updated, err := repo.Update(ctx, "primary", rule)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.Role != "owner" {
		t.Errorf("Role = %q, want %q", updated.Role, "owner")
	}
}

// TestGCalACLRepository_DeleteWithTestServer tests deleting an ACL rule.
func TestGCalACLRepository_DeleteWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	deletedRuleID := ""
	ts.ACLDeleteHandler = func(w http.ResponseWriter, r *http.Request, calendarID, ruleID string) {
		deletedRuleID = ruleID
		w.WriteHeader(http.StatusNoContent)
	}

	service := ts.GCalService(t)
	repo := service.ACL()
	ctx := context.Background()

	err := repo.Delete(ctx, "primary", "user:alice@example.com")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if deletedRuleID != "user:alice@example.com" {
		t.Errorf("deletedRuleID = %q, want %q", deletedRuleID, "user:alice@example.com")
	}
}

// TestGCalACLRepository_GetNotFound tests getting a non-existent ACL rule.
func TestGCalACLRepository_GetNotFound(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.ACLGetHandler = func(w http.ResponseWriter, r *http.Request, calendarID, ruleID string) {
		WriteErrorResponse(w, http.StatusNotFound, "ACL rule not found")
	}

	service := ts.GCalService(t)
	repo := service.ACL()
	ctx := context.Background()

	_, err := repo.Get(ctx, "primary", "nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent ACL rule, got nil")
	}
	if err != calendar.ErrACLNotFound {
		t.Errorf("error = %v, want %v", err, calendar.ErrACLNotFound)
	}
}

// =============================================================================
// FreeBusy Repository Tests
// =============================================================================

// TestGCalFreeBusyRepository_QueryWithTestServer tests querying free/busy information.
func TestGCalFreeBusyRepository_QueryWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.FreeBusyQueryHandler = func(w http.ResponseWriter, r *http.Request) {
		var request gcal.FreeBusyRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}

		// Create mock response with busy periods
		now := time.Now()
		response := MockFreeBusyResponse(map[string][]struct{ Start, End time.Time }{
			"primary": {
				{Start: now, End: now.Add(time.Hour)},
				{Start: now.Add(2 * time.Hour), End: now.Add(3 * time.Hour)},
			},
			"work_calendar": {
				{Start: now.Add(4 * time.Hour), End: now.Add(5 * time.Hour)},
			},
		})
		WriteJSONResponse(w, response)
	}

	service := ts.GCalService(t)
	repo := service.FreeBusy()
	ctx := context.Background()

	now := time.Now()
	request := &calendar.FreeBusyRequest{
		TimeMin:     now,
		TimeMax:     now.Add(24 * time.Hour),
		CalendarIDs: []string{"primary", "work_calendar"},
	}

	response, err := repo.Query(ctx, request)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(response.Calendars) != 2 {
		t.Errorf("calendars count = %d, want %d", len(response.Calendars), 2)
	}

	primaryBusy := response.Calendars["primary"]
	if len(primaryBusy) != 2 {
		t.Errorf("primary busy periods = %d, want %d", len(primaryBusy), 2)
	}

	workBusy := response.Calendars["work_calendar"]
	if len(workBusy) != 1 {
		t.Errorf("work_calendar busy periods = %d, want %d", len(workBusy), 1)
	}
}

// TestGCalFreeBusyRepository_QueryInvalidTimeRange tests query with invalid time range.
func TestGCalFreeBusyRepository_QueryInvalidTimeRange(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	service := ts.GCalService(t)
	repo := service.FreeBusy()
	ctx := context.Background()

	now := time.Now()
	// End time before start time
	request := &calendar.FreeBusyRequest{
		TimeMin:     now,
		TimeMax:     now.Add(-time.Hour),
		CalendarIDs: []string{"primary"},
	}

	_, err := repo.Query(ctx, request)
	if err == nil {
		t.Fatal("expected error for invalid time range, got nil")
	}
	if err != calendar.ErrInvalidTimeRange {
		t.Errorf("error = %v, want %v", err, calendar.ErrInvalidTimeRange)
	}
}

// TestGCalFreeBusyRepository_QueryNilRequest tests query with nil request.
func TestGCalFreeBusyRepository_QueryNilRequest(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	service := ts.GCalService(t)
	repo := service.FreeBusy()
	ctx := context.Background()

	_, err := repo.Query(ctx, nil)
	if err == nil {
		t.Fatal("expected error for nil request, got nil")
	}
}

// =============================================================================
// Error Handling Tests
// =============================================================================

// TestGCalEventRepository_RateLimited tests rate limit handling.
func TestGCalEventRepository_RateLimited(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.EventGetHandler = func(w http.ResponseWriter, r *http.Request, calendarID, eventID string) {
		WriteErrorResponse(w, http.StatusTooManyRequests, "Rate limit exceeded")
	}

	service := ts.GCalService(t)
	repo := service.Events()
	ctx := context.Background()

	_, err := repo.Get(ctx, "primary", "event123")
	if err == nil {
		t.Fatal("expected error for rate limited request, got nil")
	}
}

// TestGCalEventRepository_InternalServerError tests server error handling.
func TestGCalEventRepository_InternalServerError(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.EventGetHandler = func(w http.ResponseWriter, r *http.Request, calendarID, eventID string) {
		WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
	}

	service := ts.GCalService(t)
	repo := service.Events()
	ctx := context.Background()

	_, err := repo.Get(ctx, "primary", "event123")
	if err == nil {
		t.Fatal("expected error for server error, got nil")
	}
}

// TestGCalEventRepository_BadRequest tests bad request handling.
func TestGCalEventRepository_BadRequest(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.EventCreateHandler = func(w http.ResponseWriter, r *http.Request, calendarID string) {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid event data")
	}

	service := ts.GCalService(t)
	repo := service.Events()
	ctx := context.Background()

	event := &calendar.Event{
		Title: "Invalid Event",
		Start: time.Now(),
		End:   time.Now().Add(time.Hour),
	}

	_, err := repo.Create(ctx, "primary", event)
	if err == nil {
		t.Fatal("expected error for bad request, got nil")
	}
}

// TestGCalCalendarRepository_InvalidRequest tests invalid calendar request.
func TestGCalCalendarRepository_InvalidRequest(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	service := ts.GCalService(t)
	repo := service.Calendars()
	ctx := context.Background()

	// Update with empty ID should fail
	cal := &calendar.Calendar{
		ID:    "",
		Title: "Test",
	}

	_, err := repo.Update(ctx, cal)
	if err == nil {
		t.Fatal("expected error for invalid update request, got nil")
	}
}
