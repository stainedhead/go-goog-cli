// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"testing"
	"time"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		// Valid emails
		{
			name:     "simple email",
			email:    "user@example.com",
			expected: true,
		},
		{
			name:     "email with dots",
			email:    "first.last@example.com",
			expected: true,
		},
		{
			name:     "email with plus",
			email:    "user+tag@example.com",
			expected: true,
		},
		{
			name:     "email with underscore",
			email:    "user_name@example.com",
			expected: true,
		},
		{
			name:     "email with hyphen",
			email:    "user-name@example.com",
			expected: true,
		},
		{
			name:     "email with subdomain",
			email:    "user@mail.example.com",
			expected: true,
		},
		{
			name:     "email with numbers",
			email:    "user123@example123.com",
			expected: true,
		},
		{
			name:     "email with percent",
			email:    "user%tag@example.com",
			expected: true,
		},
		{
			name:     "email with long TLD",
			email:    "user@example.museum",
			expected: true,
		},
		{
			name:     "email with two letter TLD",
			email:    "user@example.co",
			expected: true,
		},

		// Invalid emails
		{
			name:     "empty string",
			email:    "",
			expected: false,
		},
		{
			name:     "no at sign",
			email:    "userexample.com",
			expected: false,
		},
		{
			name:     "no domain",
			email:    "user@",
			expected: false,
		},
		{
			name:     "no local part",
			email:    "@example.com",
			expected: false,
		},
		{
			name:     "no TLD",
			email:    "user@example",
			expected: false,
		},
		{
			name:     "single letter TLD",
			email:    "user@example.c",
			expected: false,
		},
		{
			name:     "double at sign",
			email:    "user@@example.com",
			expected: false,
		},
		{
			name:     "spaces in email",
			email:    "user @example.com",
			expected: false,
		},
		{
			name:     "just text",
			email:    "notanemail",
			expected: false,
		},
		{
			name:     "missing dot in domain",
			email:    "user@examplecom",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("isValidEmail(%q) = %v, expected %v", tt.email, result, tt.expected)
			}
		})
	}
}

func TestIsStartTimeValid(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		startTime time.Time
		expected  bool
	}{
		{
			name:      "today",
			startTime: now,
			expected:  true,
		},
		{
			name:      "tomorrow",
			startTime: now.AddDate(0, 0, 1),
			expected:  true,
		},
		{
			name:      "next week",
			startTime: now.AddDate(0, 0, 7),
			expected:  true,
		},
		{
			name:      "yesterday",
			startTime: now.AddDate(0, 0, -1),
			expected:  false,
		},
		{
			name:      "last week",
			startTime: now.AddDate(0, 0, -7),
			expected:  false,
		},
		{
			name:      "last month",
			startTime: now.AddDate(0, -1, 0),
			expected:  false,
		},
		{
			name:      "start of today",
			startTime: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
			expected:  true,
		},
		{
			name:      "end of yesterday",
			startTime: time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, now.Location()),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isStartTimeValid(tt.startTime)
			if result != tt.expected {
				t.Errorf("isStartTimeValid(%v) = %v, expected %v", tt.startTime, result, tt.expected)
			}
		})
	}
}

func TestIsEventDurationValid(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		expected  bool
	}{
		{
			name:      "1 hour duration",
			startTime: now,
			endTime:   now.Add(time.Hour),
			expected:  true,
		},
		{
			name:      "30 minute duration",
			startTime: now,
			endTime:   now.Add(30 * time.Minute),
			expected:  true,
		},
		{
			name:      "exactly 1 minute",
			startTime: now,
			endTime:   now.Add(time.Minute),
			expected:  true,
		},
		{
			name:      "59 seconds - invalid",
			startTime: now,
			endTime:   now.Add(59 * time.Second),
			expected:  false,
		},
		{
			name:      "0 duration - invalid",
			startTime: now,
			endTime:   now,
			expected:  false,
		},
		{
			name:      "negative duration - invalid",
			startTime: now,
			endTime:   now.Add(-time.Hour),
			expected:  false,
		},
		{
			name:      "all day event (24 hours)",
			startTime: now,
			endTime:   now.AddDate(0, 0, 1),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isEventDurationValid(tt.startTime, tt.endTime)
			if result != tt.expected {
				t.Errorf("isEventDurationValid(%v, %v) = %v, expected %v", tt.startTime, tt.endTime, result, tt.expected)
			}
		})
	}
}

func TestValidateEventTime(t *testing.T) {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	yesterday := now.AddDate(0, 0, -1)

	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		isAllDay  bool
		expectErr bool
	}{
		{
			name:      "valid future event",
			startTime: tomorrow,
			endTime:   tomorrow.Add(time.Hour),
			isAllDay:  false,
			expectErr: false,
		},
		{
			name:      "valid all-day event",
			startTime: tomorrow,
			endTime:   tomorrow.AddDate(0, 0, 1),
			isAllDay:  true,
			expectErr: false,
		},
		{
			name:      "past start time",
			startTime: yesterday,
			endTime:   yesterday.Add(time.Hour),
			isAllDay:  false,
			expectErr: true,
		},
		{
			name:      "duration too short",
			startTime: tomorrow,
			endTime:   tomorrow.Add(30 * time.Second),
			isAllDay:  false,
			expectErr: true,
		},
		{
			name:      "all-day event ignores duration check",
			startTime: tomorrow,
			endTime:   tomorrow.Add(30 * time.Second),
			isAllDay:  true,
			expectErr: false,
		},
		{
			name:      "today event is valid",
			startTime: now,
			endTime:   now.Add(time.Hour),
			isAllDay:  false,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEventTime(tt.startTime, tt.endTime, tt.isAllDay)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
