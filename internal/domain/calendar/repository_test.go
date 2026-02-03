package calendar

import (
	"errors"
	"testing"
	"time"
)

func TestDomainErrors(t *testing.T) {
	// Verify error types are distinct and have expected messages
	tests := []struct {
		err     error
		message string
	}{
		{ErrEventNotFound, "event not found"},
		{ErrCalendarNotFound, "calendar not found"},
		{ErrACLNotFound, "ACL rule not found"},
		{ErrInvalidTimeRange, "invalid time range: start must be before end"},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			if tt.err.Error() != tt.message {
				t.Errorf("expected error message %q, got %q", tt.message, tt.err.Error())
			}
		})
	}

	// Verify errors are distinguishable
	if errors.Is(ErrEventNotFound, ErrCalendarNotFound) {
		t.Error("ErrEventNotFound should not be equal to ErrCalendarNotFound")
	}

	if errors.Is(ErrCalendarNotFound, ErrACLNotFound) {
		t.Error("ErrCalendarNotFound should not be equal to ErrACLNotFound")
	}

	if errors.Is(ErrACLNotFound, ErrInvalidTimeRange) {
		t.Error("ErrACLNotFound should not be equal to ErrInvalidTimeRange")
	}
}

func TestNewFreeBusyRequest(t *testing.T) {
	timeMin := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
	timeMax := time.Date(2024, 1, 15, 17, 0, 0, 0, time.UTC)
	calendarIDs := []string{"primary", "work@group.calendar.google.com"}

	req, err := NewFreeBusyRequest(timeMin, timeMax, calendarIDs...)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !req.TimeMin.Equal(timeMin) {
		t.Errorf("expected TimeMin %v, got %v", timeMin, req.TimeMin)
	}

	if !req.TimeMax.Equal(timeMax) {
		t.Errorf("expected TimeMax %v, got %v", timeMax, req.TimeMax)
	}

	if len(req.CalendarIDs) != 2 {
		t.Errorf("expected 2 calendar IDs, got %d", len(req.CalendarIDs))
	}
}

func TestNewFreeBusyRequestInvalidRange(t *testing.T) {
	timeMin := time.Date(2024, 1, 15, 17, 0, 0, 0, time.UTC)
	timeMax := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC) // before timeMin

	_, err := NewFreeBusyRequest(timeMin, timeMax, "primary")

	if !errors.Is(err, ErrInvalidTimeRange) {
		t.Errorf("expected ErrInvalidTimeRange, got %v", err)
	}
}

func TestNewFreeBusyRequestSameTime(t *testing.T) {
	sameTime := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)

	_, err := NewFreeBusyRequest(sameTime, sameTime, "primary")

	if !errors.Is(err, ErrInvalidTimeRange) {
		t.Errorf("expected ErrInvalidTimeRange for same start and end time, got %v", err)
	}
}

func TestNewTimePeriod(t *testing.T) {
	start := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)

	tp, err := NewTimePeriod(start, end)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !tp.Start.Equal(start) {
		t.Errorf("expected Start %v, got %v", start, tp.Start)
	}

	if !tp.End.Equal(end) {
		t.Errorf("expected End %v, got %v", end, tp.End)
	}
}

func TestNewTimePeriodInvalidRange(t *testing.T) {
	start := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC) // before start

	_, err := NewTimePeriod(start, end)

	if !errors.Is(err, ErrInvalidTimeRange) {
		t.Errorf("expected ErrInvalidTimeRange, got %v", err)
	}
}

func TestTimePeriodDuration(t *testing.T) {
	start := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 11, 30, 0, 0, time.UTC)

	tp := &TimePeriod{Start: start, End: end}

	expected := 90 * time.Minute
	if tp.Duration() != expected {
		t.Errorf("expected duration %v, got %v", expected, tp.Duration())
	}
}

func TestTimePeriodOverlaps(t *testing.T) {
	// Base period: 10:00 - 12:00
	base := &TimePeriod{
		Start: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		End:   time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name     string
		other    *TimePeriod
		overlaps bool
	}{
		{
			name: "completely before",
			other: &TimePeriod{
				Start: time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
				End:   time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
			},
			overlaps: false,
		},
		{
			name: "completely after",
			other: &TimePeriod{
				Start: time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC),
				End:   time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC),
			},
			overlaps: false,
		},
		{
			name: "starts before, ends during",
			other: &TimePeriod{
				Start: time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
			},
			overlaps: true,
		},
		{
			name: "starts during, ends after",
			other: &TimePeriod{
				Start: time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
				End:   time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC),
			},
			overlaps: true,
		},
		{
			name: "completely contained",
			other: &TimePeriod{
				Start: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
				End:   time.Date(2024, 1, 15, 11, 30, 0, 0, time.UTC),
			},
			overlaps: true,
		},
		{
			name: "completely contains",
			other: &TimePeriod{
				Start: time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC),
			},
			overlaps: true,
		},
		{
			name: "exact match",
			other: &TimePeriod{
				Start: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
				End:   time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			},
			overlaps: true,
		},
		{
			name: "adjacent before (end equals start)",
			other: &TimePeriod{
				Start: time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
				End:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			},
			overlaps: false,
		},
		{
			name: "adjacent after (start equals end)",
			other: &TimePeriod{
				Start: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
				End:   time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC),
			},
			overlaps: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := base.Overlaps(tt.other)
			if got != tt.overlaps {
				t.Errorf("Overlaps() = %v, want %v", got, tt.overlaps)
			}
		})
	}
}

func TestTimePeriodContains(t *testing.T) {
	// Period: 10:00 - 12:00
	tp := &TimePeriod{
		Start: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		End:   time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name     string
		time     time.Time
		contains bool
	}{
		{
			name:     "before period",
			time:     time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
			contains: false,
		},
		{
			name:     "at start",
			time:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			contains: true,
		},
		{
			name:     "in middle",
			time:     time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
			contains: true,
		},
		{
			name:     "at end (exclusive)",
			time:     time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			contains: false,
		},
		{
			name:     "after period",
			time:     time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC),
			contains: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tp.Contains(tt.time)
			if got != tt.contains {
				t.Errorf("Contains() = %v, want %v", got, tt.contains)
			}
		})
	}
}

func TestFreeBusyResponse(t *testing.T) {
	resp := &FreeBusyResponse{
		Calendars: map[string][]*TimePeriod{
			"primary": {
				{
					Start: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					End:   time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
				},
				{
					Start: time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC),
					End:   time.Date(2024, 1, 15, 15, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	if len(resp.Calendars["primary"]) != 2 {
		t.Errorf("expected 2 busy periods, got %d", len(resp.Calendars["primary"]))
	}

	if resp.Calendars["nonexistent"] != nil {
		t.Error("expected nil for nonexistent calendar")
	}
}
