// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/adapter/presenter"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
)

// Calendar event command flags.
var (
	// Create flags
	calCreateTitle       string
	calCreateStart       string
	calCreateEnd         string
	calCreateLocation    string
	calCreateDescription string
	calCreateAttendees   []string
	calCreateAllDay      bool
	calCreateCalendar    string

	// Update flags
	calUpdateTitle       string
	calUpdateStart       string
	calUpdateEnd         string
	calUpdateLocation    string
	calUpdateDescription string
	calUpdateAttendees   []string
	calUpdateCalendar    string

	// Delete flags
	calDeleteConfirm  bool
	calDeleteCalendar string
)

// calCreateCmd creates a new calendar event.
var calCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new calendar event",
	Long: `Create a new calendar event.

Create a new event in your Google Calendar with the specified details.
The --title and --start flags are required.

Date/time formats supported:
  - "2024-01-15 14:00" (date and time)
  - "2024-01-15" (date only, for all-day events)
  - "2024-01-15T14:00:00Z" (RFC3339)
  - "tomorrow 3pm" (relative date with time)
  - "today 2pm" (relative date with time)`,
	Example: `  # Create a simple event
  goog cal create --title "Team Meeting" --start "2024-01-15 14:00" --end "2024-01-15 15:00"

  # Create an all-day event
  goog cal create --title "Company Holiday" --start "2024-01-15" --all-day

  # Create an event with location and attendees
  goog cal create --title "Sprint Planning" --start "tomorrow 10am" --end "tomorrow 12pm" \
    --location "Conference Room A" --attendees user1@example.com,user2@example.com

  # Create an event in a specific calendar
  goog cal create --title "Personal Errand" --start "today 3pm" --calendar work@example.com`,
	RunE: runCalCreate,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if calCreateTitle == "" {
			return fmt.Errorf("required flag \"title\" not set")
		}
		if calCreateStart == "" {
			return fmt.Errorf("required flag \"start\" not set")
		}
		return nil
	},
}

// calUpdateCmd updates an existing calendar event.
var calUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an existing calendar event",
	Long: `Update an existing calendar event.

Modify the specified event's properties. Only the flags you provide
will be updated; other properties remain unchanged.`,
	Example: `  # Update event title
  goog cal update abc123 --title "Updated Meeting Title"

  # Reschedule an event
  goog cal update abc123 --start "2024-01-16 14:00" --end "2024-01-16 15:00"

  # Update location and description
  goog cal update abc123 --location "Room B" --description "Moved to new room"

  # Add attendees
  goog cal update abc123 --attendees user1@example.com,user2@example.com

  # Update event in a specific calendar
  goog cal update abc123 --title "New Title" --calendar work@example.com`,
	Args: cobra.ExactArgs(1),
	RunE: runCalUpdate,
}

// calDeleteCmd deletes a calendar event.
var calDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a calendar event",
	Long: `Delete a calendar event.

Permanently remove the specified event from the calendar.
The --confirm flag is required to prevent accidental deletion.`,
	Example: `  # Delete an event (requires --confirm)
  goog cal delete abc123 --confirm

  # Delete an event from a specific calendar
  goog cal delete abc123 --confirm --calendar work@example.com`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !calDeleteConfirm {
			cmd.PrintErrln("Error: deletion requires --confirm flag")
			return fmt.Errorf("--confirm flag required for deletion")
		}
		return nil
	},
	RunE: runCalDelete,
}

func init() {
	// Add cal event subcommands
	calCmd.AddCommand(calCreateCmd)
	calCmd.AddCommand(calUpdateCmd)
	calCmd.AddCommand(calDeleteCmd)

	// Create command flags
	calCreateCmd.Flags().StringVar(&calCreateTitle, "title", "", "event title (required)")
	calCreateCmd.Flags().StringVar(&calCreateStart, "start", "", "event start time (required)")
	calCreateCmd.Flags().StringVar(&calCreateEnd, "end", "", "event end time (defaults to 1 hour after start)")
	calCreateCmd.Flags().StringVar(&calCreateLocation, "location", "", "event location")
	calCreateCmd.Flags().StringVar(&calCreateDescription, "description", "", "event description")
	calCreateCmd.Flags().StringSliceVar(&calCreateAttendees, "attendees", nil, "attendee email addresses (comma-separated)")
	calCreateCmd.Flags().BoolVar(&calCreateAllDay, "all-day", false, "create an all-day event")
	calCreateCmd.Flags().StringVar(&calCreateCalendar, "calendar", "primary", "calendar ID to use")

	// Update command flags
	calUpdateCmd.Flags().StringVar(&calUpdateTitle, "title", "", "new event title")
	calUpdateCmd.Flags().StringVar(&calUpdateStart, "start", "", "new event start time")
	calUpdateCmd.Flags().StringVar(&calUpdateEnd, "end", "", "new event end time")
	calUpdateCmd.Flags().StringVar(&calUpdateLocation, "location", "", "new event location")
	calUpdateCmd.Flags().StringVar(&calUpdateDescription, "description", "", "new event description")
	calUpdateCmd.Flags().StringSliceVar(&calUpdateAttendees, "attendees", nil, "new attendee email addresses (comma-separated)")
	calUpdateCmd.Flags().StringVar(&calUpdateCalendar, "calendar", "primary", "calendar ID to use")

	// Delete command flags
	calDeleteCmd.Flags().BoolVar(&calDeleteConfirm, "confirm", false, "confirm deletion")
	calDeleteCmd.Flags().StringVar(&calDeleteCalendar, "calendar", "primary", "calendar ID to use")
}

// runCalCreate handles the cal create command.
func runCalCreate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Parse start time
	startTime, err := parseDateTime(calCreateStart)
	if err != nil {
		return fmt.Errorf("invalid start time: %w", err)
	}

	// Determine end time
	var endTime time.Time
	if calCreateEnd != "" {
		endTime, err = parseDateTime(calCreateEnd)
		if err != nil {
			return fmt.Errorf("invalid end time: %w", err)
		}
	} else if calCreateAllDay {
		// For all-day events without explicit end, use next day
		endTime = startTime.AddDate(0, 0, 1)
	} else {
		// Default to 1 hour duration
		endTime = startTime.Add(time.Hour)
	}

	// Validate time range
	if !startTime.Before(endTime) {
		return fmt.Errorf("start time must be before end time")
	}

	// Validate start time and duration
	if err := validateEventTime(startTime, endTime, calCreateAllDay); err != nil {
		return err
	}

	// Get repository
	repo, err := getGCalEventRepository(ctx)
	if err != nil {
		return err
	}

	// Build event
	var event *calendar.Event
	if calCreateAllDay {
		event = calendar.NewAllDayEvent(calCreateTitle, startTime)
		// Set the end date for multi-day all-day events
		if calCreateEnd != "" {
			event.End = endTime
		}
	} else {
		event = calendar.NewEvent(calCreateTitle, startTime, endTime)
	}

	// Set optional fields
	if calCreateLocation != "" {
		event.Location = calCreateLocation
	}
	if calCreateDescription != "" {
		event.Description = calCreateDescription
	}

	// Add attendees
	attendees, err := parseAttendees(calCreateAttendees)
	if err != nil {
		return err
	}
	for _, email := range attendees {
		event.AddAttendee(calendar.NewAttendee(email))
	}

	// Create event
	created, err := repo.Create(ctx, calCreateCalendar, event)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	// Output result
	p := presenter.New(formatFlag)
	output := p.RenderEvent(created)
	cmd.Println(output)

	if !quietFlag {
		cmd.Printf("\nEvent created successfully.\n")
		cmd.Printf("Event ID: %s\n", created.ID)
		if created.HTMLLink != "" {
			cmd.Printf("Link: %s\n", created.HTMLLink)
		}
	}

	return nil
}

// runCalUpdate handles the cal update command.
func runCalUpdate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	eventID := args[0]

	// Get repository
	repo, err := getGCalEventRepository(ctx)
	if err != nil {
		return err
	}

	// Fetch existing event
	existing, err := repo.Get(ctx, calUpdateCalendar, eventID)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	// Update fields if flags were provided
	if calUpdateTitle != "" {
		existing.Title = calUpdateTitle
	}

	if calUpdateStart != "" {
		startTime, err := parseDateTime(calUpdateStart)
		if err != nil {
			return fmt.Errorf("invalid start time: %w", err)
		}
		existing.Start = startTime
	}

	if calUpdateEnd != "" {
		endTime, err := parseDateTime(calUpdateEnd)
		if err != nil {
			return fmt.Errorf("invalid end time: %w", err)
		}
		existing.End = endTime
	}

	// Validate time range
	if !existing.Start.Before(existing.End) {
		return fmt.Errorf("start time must be before end time")
	}

	if calUpdateLocation != "" {
		existing.Location = calUpdateLocation
	}

	if calUpdateDescription != "" {
		existing.Description = calUpdateDescription
	}

	// Update attendees if provided
	if len(calUpdateAttendees) > 0 {
		attendees, err := parseAttendees(calUpdateAttendees)
		if err != nil {
			return err
		}
		existing.Attendees = make([]*calendar.Attendee, 0, len(attendees))
		for _, email := range attendees {
			existing.AddAttendee(calendar.NewAttendee(email))
		}
	}

	// Update event
	updated, err := repo.Update(ctx, calUpdateCalendar, existing)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	// Output result
	p := presenter.New(formatFlag)
	output := p.RenderEvent(updated)
	cmd.Println(output)

	if !quietFlag {
		cmd.Printf("\nEvent updated successfully.\n")
		cmd.Printf("Event ID: %s\n", updated.ID)
		if updated.HTMLLink != "" {
			cmd.Printf("Link: %s\n", updated.HTMLLink)
		}
	}

	return nil
}

// runCalDelete handles the cal delete command.
func runCalDelete(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	eventID := args[0]

	// Get repository
	repo, err := getGCalEventRepository(ctx)
	if err != nil {
		return err
	}

	// Delete event
	if err := repo.Delete(ctx, calDeleteCalendar, eventID); err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	if !quietFlag {
		cmd.Printf("Event %s deleted successfully.\n", eventID)
	}

	return nil
}

// parseDateTime parses a date/time string into a time.Time.
// Supports various formats including:
// - "2024-01-15 14:00" (date and time)
// - "2024-01-15" (date only)
// - "2024-01-15 14:30:45" (date with seconds)
// - "2024-01-15T14:00:00Z" (RFC3339)
// - "tomorrow" or "tomorrow 3pm" (relative)
// - "today" or "today 2pm" (relative)
func parseDateTime(input string) (time.Time, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return time.Time{}, fmt.Errorf("empty date/time string")
	}

	// Try RFC3339 first
	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return t, nil
	}

	// Try common formats
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(format, input, time.Local); err == nil {
			return t, nil
		}
	}

	// Handle relative dates
	lower := strings.ToLower(input)

	if strings.HasPrefix(lower, "tomorrow") {
		return parseRelativeDate(input, 1)
	}

	if strings.HasPrefix(lower, "today") {
		return parseRelativeDate(input, 0)
	}

	return time.Time{}, fmt.Errorf("unable to parse date/time: %q", input)
}

// parseRelativeDate parses a relative date string like "tomorrow 3pm".
func parseRelativeDate(input string, daysOffset int) (time.Time, error) {
	now := time.Now()
	baseDate := now.AddDate(0, 0, daysOffset)

	// Extract time component if present
	parts := strings.Fields(strings.ToLower(input))
	if len(parts) == 1 {
		// Just "tomorrow" or "today" - return start of day
		return time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 0, 0, 0, 0, time.Local), nil
	}

	// Try to parse time component
	timeStr := parts[1]
	hour, minute, err := parseTimeOfDay(timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time in relative date: %w", err)
	}

	return time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), hour, minute, 0, 0, time.Local), nil
}

// parseTimeOfDay parses a time string like "3pm", "3:30pm", or "14:00".
func parseTimeOfDay(input string) (hour, minute int, err error) {
	input = strings.ToLower(strings.TrimSpace(input))

	// Try 24-hour format first (e.g., "14:00")
	if matched, _ := regexp.MatchString(`^\d{1,2}:\d{2}$`, input); matched {
		parts := strings.Split(input, ":")
		hour, _ = strconv.Atoi(parts[0])
		minute, _ = strconv.Atoi(parts[1])
		if hour >= 0 && hour <= 23 && minute >= 0 && minute <= 59 {
			return hour, minute, nil
		}
		return 0, 0, fmt.Errorf("invalid time: %s", input)
	}

	// Try 12-hour format (e.g., "3pm", "3:30pm")
	ampmRegex := regexp.MustCompile(`^(\d{1,2})(?::(\d{2}))?(am|pm)$`)
	matches := ampmRegex.FindStringSubmatch(input)
	if matches != nil {
		hour, _ = strconv.Atoi(matches[1])
		if matches[2] != "" {
			minute, _ = strconv.Atoi(matches[2])
		}
		isPM := matches[3] == "pm"

		// Convert to 24-hour format
		if hour == 12 {
			if !isPM {
				hour = 0 // 12am = midnight
			}
		} else if isPM {
			hour += 12
		}

		if hour >= 0 && hour <= 23 && minute >= 0 && minute <= 59 {
			return hour, minute, nil
		}
	}

	return 0, 0, fmt.Errorf("invalid time format: %s", input)
}

// parseAttendees cleans, validates, and returns attendee email addresses.
// Returns an error if any email address is invalid.
func parseAttendees(attendees []string) ([]string, error) {
	if attendees == nil {
		return []string{}, nil
	}

	result := make([]string, 0, len(attendees))
	for _, a := range attendees {
		trimmed := strings.TrimSpace(a)
		if trimmed != "" {
			if !isValidEmail(trimmed) {
				return nil, fmt.Errorf("invalid attendee email: %q", trimmed)
			}
			result = append(result, trimmed)
		}
	}
	return result, nil
}
