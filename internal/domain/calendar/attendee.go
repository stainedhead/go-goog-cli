// Package calendar provides domain entities for Google Calendar operations.
package calendar

// Attendee represents a participant in a calendar event.
type Attendee struct {
	// Email is the attendee's email address.
	Email string
	// DisplayName is the attendee's display name.
	DisplayName string
	// ResponseStatus indicates the attendee's response: needsAction, declined, tentative, accepted.
	ResponseStatus string
	// Optional indicates whether the attendee is optional.
	Optional bool
	// Organizer indicates whether the attendee is the organizer.
	Organizer bool
	// Self indicates whether this attendee is the current user.
	Self bool
}

// Response status constants.
const (
	ResponseNeedsAction = "needsAction"
	ResponseDeclined    = "declined"
	ResponseTentative   = "tentative"
	ResponseAccepted    = "accepted"
)

// NewAttendee creates a new Attendee with the given email.
func NewAttendee(email string) *Attendee {
	return &Attendee{
		Email:          email,
		ResponseStatus: ResponseNeedsAction,
	}
}

// IsValidResponseStatus checks if the given status is a valid response status.
func IsValidResponseStatus(status string) bool {
	switch status {
	case ResponseNeedsAction, ResponseDeclined, ResponseTentative, ResponseAccepted:
		return true
	default:
		return false
	}
}
