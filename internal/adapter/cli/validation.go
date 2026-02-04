// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"fmt"
	"regexp"
	"time"
)

// emailRegex is a simple but effective regex pattern for validating email addresses.
// It matches most common email formats while keeping the pattern readable.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// isValidEmail validates an email address using a simple regex pattern.
// It returns true if the email address has a valid format.
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	return emailRegex.MatchString(email)
}

// isStartTimeValid validates that the start time is not in the past.
// It allows events starting on the same day for usability.
func isStartTimeValid(startTime time.Time) bool {
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return !startTime.Before(startOfToday)
}

// isEventDurationValid validates that the event duration is at least the minimum required.
// Returns true if the duration is at least 1 minute.
func isEventDurationValid(startTime, endTime time.Time) bool {
	duration := endTime.Sub(startTime)
	return duration >= time.Minute
}

// validateEventTime validates both start time and duration for a calendar event.
// It returns an error if validation fails.
func validateEventTime(startTime, endTime time.Time, isAllDay bool) error {
	// Validate start time is not in the past
	if !isStartTimeValid(startTime) {
		return fmt.Errorf("start time cannot be in the past")
	}

	// Validate minimum duration for non-all-day events
	if !isAllDay && !isEventDurationValid(startTime, endTime) {
		return fmt.Errorf("event duration must be at least 1 minute")
	}

	return nil
}
