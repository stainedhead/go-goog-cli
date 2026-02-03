package calendar

import "time"

// Event represents a Google Calendar event.
type Event struct {
	// ID is the unique identifier for the event.
	ID string
	// CalendarID is the ID of the calendar containing this event.
	CalendarID string
	// Title is the event summary/title.
	Title string
	// Description is the event description.
	Description string
	// Location is the geographic location of the event.
	Location string
	// Start is the event start time.
	Start time.Time
	// End is the event end time.
	End time.Time
	// AllDay indicates whether this is an all-day event.
	AllDay bool
	// Recurrence contains RRULE strings for recurring events.
	Recurrence []string
	// Attendees is the list of event attendees.
	Attendees []*Attendee
	// Organizer is the event organizer.
	Organizer *Attendee
	// Status is the event status: confirmed, tentative, cancelled.
	Status string
	// Visibility is the event visibility: public, private.
	Visibility string
	// ColorID is the color ID for the event.
	ColorID string
	// Reminders is the list of reminders for the event.
	Reminders []*Reminder
	// ConferenceData contains video conference information.
	ConferenceData *ConferenceData
	// Created is when the event was created.
	Created time.Time
	// Updated is when the event was last updated.
	Updated time.Time
	// HTMLLink is the URL to the event in Google Calendar.
	HTMLLink string
}

// Event status constants.
const (
	StatusConfirmed = "confirmed"
	StatusTentative = "tentative"
	StatusCancelled = "cancelled"
)

// Event visibility constants.
const (
	VisibilityPublic  = "public"
	VisibilityPrivate = "private"
)

// Reminder represents an event reminder.
type Reminder struct {
	// Method is the reminder delivery method: email, popup.
	Method string
	// Minutes is how many minutes before the event the reminder should trigger.
	Minutes int
}

// Reminder method constants.
const (
	ReminderMethodEmail = "email"
	ReminderMethodPopup = "popup"
)

// ConferenceData represents video conference information.
type ConferenceData struct {
	// Type is the conference type (e.g., "hangoutsMeet").
	Type string
	// URI is the conference join URI.
	URI string
}

// NewEvent creates a new Event with the given title and time range.
func NewEvent(title string, start, end time.Time) *Event {
	return &Event{
		Title:      title,
		Start:      start,
		End:        end,
		Status:     StatusConfirmed,
		Visibility: VisibilityPrivate,
		Attendees:  make([]*Attendee, 0),
		Reminders:  make([]*Reminder, 0),
	}
}

// NewAllDayEvent creates a new all-day Event.
func NewAllDayEvent(title string, date time.Time) *Event {
	// For all-day events, use the start of the day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	return &Event{
		Title:      title,
		Start:      startOfDay,
		End:        endOfDay,
		AllDay:     true,
		Status:     StatusConfirmed,
		Visibility: VisibilityPrivate,
		Attendees:  make([]*Attendee, 0),
		Reminders:  make([]*Reminder, 0),
	}
}

// NewReminder creates a new Reminder with the given method and minutes.
func NewReminder(method string, minutes int) *Reminder {
	return &Reminder{
		Method:  method,
		Minutes: minutes,
	}
}

// IsValidStatus checks if the given status is a valid event status.
func IsValidStatus(status string) bool {
	switch status {
	case StatusConfirmed, StatusTentative, StatusCancelled:
		return true
	default:
		return false
	}
}

// IsValidVisibility checks if the given visibility is a valid event visibility.
func IsValidVisibility(visibility string) bool {
	switch visibility {
	case VisibilityPublic, VisibilityPrivate:
		return true
	default:
		return false
	}
}

// IsValidReminderMethod checks if the given method is a valid reminder method.
func IsValidReminderMethod(method string) bool {
	switch method {
	case ReminderMethodEmail, ReminderMethodPopup:
		return true
	default:
		return false
	}
}

// Duration returns the duration of the event.
func (e *Event) Duration() time.Duration {
	return e.End.Sub(e.Start)
}

// IsRecurring returns true if the event has recurrence rules.
func (e *Event) IsRecurring() bool {
	return len(e.Recurrence) > 0
}

// HasConference returns true if the event has conference data.
func (e *Event) HasConference() bool {
	return e.ConferenceData != nil && e.ConferenceData.URI != ""
}

// AddAttendee adds an attendee to the event.
func (e *Event) AddAttendee(attendee *Attendee) {
	if e.Attendees == nil {
		e.Attendees = make([]*Attendee, 0)
	}
	e.Attendees = append(e.Attendees, attendee)
}

// AddReminder adds a reminder to the event.
func (e *Event) AddReminder(reminder *Reminder) {
	if e.Reminders == nil {
		e.Reminders = make([]*Reminder, 0)
	}
	e.Reminders = append(e.Reminders, reminder)
}
