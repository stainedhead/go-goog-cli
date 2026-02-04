// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/adapter/presenter"
	"github.com/stainedhead/go-goog-cli/internal/adapter/repository"
)

// Command flags for calendar event list/show commands.
var (
	calListMaxResults int
)

// getGCalEventRepository creates a GCalEventRepository using the current account's credentials.
// This is a convenience function that wraps getGCalService for commands that only need event operations.
func getGCalEventRepository(ctx context.Context) (*repository.GCalEventRepository, error) {
	gcalSvc, err := getGCalService(ctx)
	if err != nil {
		return nil, err
	}
	return gcalSvc.Events(), nil
}

// calListCmd lists upcoming calendar events.
var calListCmd = &cobra.Command{
	Use:   "list",
	Short: "List upcoming events",
	Long: `List upcoming events from your Google Calendar.

By default, lists events from the primary calendar for the next
30 days. Use --calendar to specify a different calendar and
--max-results to limit the number of events returned.`,
	Example: `  # List upcoming events
  goog cal list

  # List events from a specific calendar
  goog cal list --calendar work@group.calendar.google.com

  # List with JSON output
  goog cal list --format json

  # Limit number of results
  goog cal list --max-results 10`,
	Aliases: []string{"ls"},
	RunE:    runCalList,
}

// calShowCmd shows details of a single event.
var calShowCmd = &cobra.Command{
	Use:   "show <event-id>",
	Short: "Show event details",
	Long: `Show detailed information about a specific calendar event.

Retrieves and displays all available information about the
specified event including attendees, location, and conference data.`,
	Example: `  # Show event details
  goog cal show abc123def456

  # Show with JSON output
  goog cal show abc123def456 --format json

  # Show event from a specific calendar
  goog cal show abc123def456 --calendar work@group.calendar.google.com`,
	Aliases: []string{"get"},
	Args:    cobra.ExactArgs(1),
	RunE:    runCalShow,
}

// calTodayCmd shows today's events.
var calTodayCmd = &cobra.Command{
	Use:   "today",
	Short: "Show today's events",
	Long: `Show all events scheduled for today.

Lists all events from midnight to midnight of the current day
in your local timezone.`,
	Example: `  # Show today's events
  goog cal today

  # Show today's events from a specific calendar
  goog cal today --calendar work@group.calendar.google.com

  # Show with table output
  goog cal today --format table`,
	RunE: runCalToday,
}

// calWeekCmd shows this week's events.
var calWeekCmd = &cobra.Command{
	Use:   "week",
	Short: "Show this week's events",
	Long: `Show all events scheduled for the current week.

Lists all events from the start of the current week (Monday)
to the end of the week (Sunday) in your local timezone.`,
	Example: `  # Show this week's events
  goog cal week

  # Show this week's events from a specific calendar
  goog cal week --calendar work@group.calendar.google.com

  # Show with JSON output
  goog cal week --format json`,
	RunE: runCalWeek,
}

func init() {
	// Add calendar event subcommands to calCmd (defined in cal_utils.go)
	calCmd.AddCommand(calListCmd)
	calCmd.AddCommand(calShowCmd)
	calCmd.AddCommand(calTodayCmd)
	calCmd.AddCommand(calWeekCmd)

	// List command flags
	calListCmd.Flags().StringVar(&calCalendarFlag, "calendar", "primary", "calendar ID to use")
	calListCmd.Flags().IntVar(&calListMaxResults, "max-results", 25, "maximum number of events to return")

	// Show command flags
	calShowCmd.Flags().StringVar(&calCalendarFlag, "calendar", "primary", "calendar ID to use")

	// Today command flags
	calTodayCmd.Flags().StringVar(&calCalendarFlag, "calendar", "primary", "calendar ID to use")

	// Week command flags
	calWeekCmd.Flags().StringVar(&calCalendarFlag, "calendar", "primary", "calendar ID to use")
}

// runCalList handles the cal list command.
func runCalList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get Calendar service (defined in cal_utils.go)
	gcalSvc, err := getGCalService(ctx)
	if err != nil {
		return err
	}

	// Get events repository
	repo := gcalSvc.Events()

	// Calculate time range (now to 30 days from now)
	now := time.Now()
	timeMin := now
	timeMax := now.AddDate(0, 0, 30)

	// List events
	events, err := repo.List(ctx, calCalendarFlag, timeMin, timeMax)
	if err != nil {
		return fmt.Errorf("failed to list events: %w", err)
	}

	// Apply max results limit
	if calListMaxResults > 0 && len(events) > calListMaxResults {
		events = events[:calListMaxResults]
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	// Output result
	output := p.RenderEvents(events)
	cmd.Println(output)

	return nil
}

// runCalShow handles the cal show command.
func runCalShow(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	eventID := args[0]

	// Get Calendar service (defined in cal_utils.go)
	gcalSvc, err := getGCalService(ctx)
	if err != nil {
		return err
	}

	// Get events repository
	repo := gcalSvc.Events()

	// Get the event
	event, err := repo.Get(ctx, calCalendarFlag, eventID)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	// Output result
	output := p.RenderEvent(event)
	cmd.Println(output)

	return nil
}

// runCalToday handles the cal today command.
func runCalToday(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get Calendar service (defined in cal_utils.go)
	gcalSvc, err := getGCalService(ctx)
	if err != nil {
		return err
	}

	// Get events repository
	repo := gcalSvc.Events()

	// Calculate today's time range
	now := time.Now()
	loc := now.Location()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	// List events for today
	events, err := repo.List(ctx, calCalendarFlag, startOfDay, endOfDay)
	if err != nil {
		return fmt.Errorf("failed to list today's events: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	// Output result
	output := p.RenderEvents(events)
	cmd.Println(output)

	// Show event count if not quiet
	if !quietFlag && len(events) > 0 {
		cmd.Printf("\n%d event(s) for %s\n", len(events), startOfDay.Format("Monday, January 2, 2006"))
	}

	return nil
}

// runCalWeek handles the cal week command.
func runCalWeek(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get Calendar service (defined in cal_utils.go)
	gcalSvc, err := getGCalService(ctx)
	if err != nil {
		return err
	}

	// Get events repository
	repo := gcalSvc.Events()

	// Calculate this week's time range (Monday to Sunday)
	now := time.Now()
	loc := now.Location()

	// Find the start of the week (Monday)
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday is day 7
	}
	daysFromMonday := weekday - 1
	startOfWeek := time.Date(now.Year(), now.Month(), now.Day()-daysFromMonday, 0, 0, 0, 0, loc)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	// List events for this week
	events, err := repo.List(ctx, calCalendarFlag, startOfWeek, endOfWeek)
	if err != nil {
		return fmt.Errorf("failed to list this week's events: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	// Output result
	output := p.RenderEvents(events)
	cmd.Println(output)

	// Show event count if not quiet
	if !quietFlag && len(events) > 0 {
		cmd.Printf("\n%d event(s) for week of %s\n", len(events), startOfWeek.Format("January 2, 2006"))
	}

	return nil
}
