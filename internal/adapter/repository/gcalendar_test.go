package repository

import (
	"context"
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
