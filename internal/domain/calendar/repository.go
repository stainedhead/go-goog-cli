package calendar

import (
	"context"
	"errors"
	"time"
)

// Domain errors.
var (
	// ErrEventNotFound is returned when an event cannot be found.
	ErrEventNotFound = errors.New("event not found")
	// ErrCalendarNotFound is returned when a calendar cannot be found.
	ErrCalendarNotFound = errors.New("calendar not found")
	// ErrACLNotFound is returned when an ACL rule cannot be found.
	ErrACLNotFound = errors.New("ACL rule not found")
	// ErrInvalidTimeRange is returned when an invalid time range is provided.
	ErrInvalidTimeRange = errors.New("invalid time range: start must be before end")
)

// EventRepository defines the interface for event persistence operations.
type EventRepository interface {
	// List returns events from a calendar within the specified time range.
	List(ctx context.Context, calendarID string, timeMin, timeMax time.Time) ([]*Event, error)
	// Get retrieves a single event by ID.
	Get(ctx context.Context, calendarID, eventID string) (*Event, error)
	// Create creates a new event in the specified calendar.
	Create(ctx context.Context, calendarID string, event *Event) (*Event, error)
	// Update updates an existing event.
	Update(ctx context.Context, calendarID string, event *Event) (*Event, error)
	// Delete removes an event from a calendar.
	Delete(ctx context.Context, calendarID, eventID string) error
	// Move moves an event to a different calendar.
	Move(ctx context.Context, sourceCalendarID, eventID, destinationCalendarID string) (*Event, error)
	// QuickAdd creates an event based on a simple text string (e.g., "Meeting tomorrow 3pm").
	QuickAdd(ctx context.Context, calendarID, text string) (*Event, error)
	// Instances returns instances of a recurring event.
	Instances(ctx context.Context, calendarID, eventID string, timeMin, timeMax time.Time) ([]*Event, error)
	// RSVP updates the current user's response to an event.
	RSVP(ctx context.Context, calendarID, eventID, response string) error
}

// CalendarRepository defines the interface for calendar persistence operations.
type CalendarRepository interface {
	// List returns all calendars accessible to the user.
	List(ctx context.Context) ([]*Calendar, error)
	// Get retrieves a single calendar by ID.
	Get(ctx context.Context, calendarID string) (*Calendar, error)
	// Create creates a new calendar.
	Create(ctx context.Context, calendar *Calendar) (*Calendar, error)
	// Update updates an existing calendar.
	Update(ctx context.Context, calendar *Calendar) (*Calendar, error)
	// Delete removes a calendar.
	Delete(ctx context.Context, calendarID string) error
	// Clear clears all events from a calendar.
	Clear(ctx context.Context, calendarID string) error
}

// ACLRepository defines the interface for calendar ACL operations.
type ACLRepository interface {
	// List returns all ACL rules for a calendar.
	List(ctx context.Context, calendarID string) ([]*ACLRule, error)
	// Get retrieves a single ACL rule by ID.
	Get(ctx context.Context, calendarID, ruleID string) (*ACLRule, error)
	// Insert creates a new ACL rule for a calendar.
	Insert(ctx context.Context, calendarID string, rule *ACLRule) (*ACLRule, error)
	// Update updates an existing ACL rule.
	Update(ctx context.Context, calendarID string, rule *ACLRule) (*ACLRule, error)
	// Delete removes an ACL rule from a calendar.
	Delete(ctx context.Context, calendarID, ruleID string) error
}

// FreeBusyRequest represents a request for free/busy information.
type FreeBusyRequest struct {
	// TimeMin is the start of the time range to query.
	TimeMin time.Time
	// TimeMax is the end of the time range to query.
	TimeMax time.Time
	// CalendarIDs is the list of calendar IDs to query.
	CalendarIDs []string
}

// FreeBusyResponse represents the response from a free/busy query.
type FreeBusyResponse struct {
	// Calendars maps calendar IDs to their busy periods.
	Calendars map[string][]*TimePeriod
}

// TimePeriod represents a period of time.
type TimePeriod struct {
	// Start is the start of the time period.
	Start time.Time
	// End is the end of the time period.
	End time.Time
}

// FreeBusyRepository defines the interface for free/busy queries.
type FreeBusyRepository interface {
	// Query returns free/busy information for the specified calendars and time range.
	Query(ctx context.Context, request *FreeBusyRequest) (*FreeBusyResponse, error)
}

// NewFreeBusyRequest creates a new FreeBusyRequest.
func NewFreeBusyRequest(timeMin, timeMax time.Time, calendarIDs ...string) (*FreeBusyRequest, error) {
	if !timeMin.Before(timeMax) {
		return nil, ErrInvalidTimeRange
	}
	return &FreeBusyRequest{
		TimeMin:     timeMin,
		TimeMax:     timeMax,
		CalendarIDs: calendarIDs,
	}, nil
}

// NewTimePeriod creates a new TimePeriod.
func NewTimePeriod(start, end time.Time) (*TimePeriod, error) {
	if !start.Before(end) {
		return nil, ErrInvalidTimeRange
	}
	return &TimePeriod{
		Start: start,
		End:   end,
	}, nil
}

// Duration returns the duration of the time period.
func (tp *TimePeriod) Duration() time.Duration {
	return tp.End.Sub(tp.Start)
}

// Overlaps returns true if this time period overlaps with another.
func (tp *TimePeriod) Overlaps(other *TimePeriod) bool {
	return tp.Start.Before(other.End) && other.Start.Before(tp.End)
}

// Contains returns true if the given time is within this time period.
func (tp *TimePeriod) Contains(t time.Time) bool {
	return !t.Before(tp.Start) && t.Before(tp.End)
}
