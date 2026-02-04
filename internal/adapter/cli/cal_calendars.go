// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/adapter/presenter"
	"github.com/stainedhead/go-goog-cli/internal/adapter/repository"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
)

// Calendar management command flags.
var (
	calendarsTitle       string
	calendarsDescription string
	calendarsTimezone    string
	calendarsConfirm     bool
)

// calendarsCmd represents the calendars command group.
var calendarsCmd = &cobra.Command{
	Use:   "calendars",
	Short: "Manage calendars",
	Long: `Manage Google Calendar calendars.

Calendar management commands allow you to list, view, create, update,
and delete calendars, as well as clear all events from a calendar.

This is different from the event commands (list, show, today, week)
which operate on events within calendars.`,
	Aliases: []string{"calendar"},
}

// calendarsListCmd lists all calendars.
var calendarsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all calendars",
	Long: `List all calendars accessible to the user.

Displays both primary and secondary calendars, along with
their access roles and time zones.`,
	Aliases: []string{"ls"},
	Example: `  # List all calendars
  goog cal calendars list

  # List calendars with JSON output
  goog cal calendars list --format json`,
	RunE: runCalendarsList,
}

// calendarsShowCmd shows a single calendar.
var calendarsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show calendar details",
	Long: `Show details of a specific calendar.

Displays the calendar ID, title, description, time zone,
primary status, and access role.`,
	Aliases: []string{"get", "info"},
	Example: `  # Show calendar by ID
  goog cal calendars show primary

  # Show a secondary calendar
  goog cal calendars show "example@group.calendar.google.com"

  # Show with JSON output
  goog cal calendars show primary --format json`,
	Args: cobra.ExactArgs(1),
	RunE: runCalendarsShow,
}

// calendarsCreateCmd creates a new calendar.
var calendarsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new calendar",
	Long: `Create a new Google Calendar.

You must specify a title for the calendar. Optionally, you can
set a description and time zone.`,
	Example: `  # Create a simple calendar
  goog cal calendars create --title "Work Projects"

  # Create a calendar with description and timezone
  goog cal calendars create --title "Team Meetings" --description "Team sync meetings" --timezone "America/New_York"`,
	RunE: runCalendarsCreate,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if calendarsTitle == "" {
			return fmt.Errorf("required flag \"title\" not set")
		}
		return nil
	},
}

// calendarsUpdateCmd updates an existing calendar.
var calendarsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a calendar",
	Long: `Update an existing calendar's metadata.

You can update the title, description, and time zone.
Note that you can only update calendars you own.`,
	Example: `  # Update calendar title
  goog cal calendars update "example@group.calendar.google.com" --title "New Title"

  # Update multiple fields
  goog cal calendars update "example@group.calendar.google.com" --title "New Title" --description "New description" --timezone "Europe/London"`,
	Args: cobra.ExactArgs(1),
	RunE: runCalendarsUpdate,
}

// calendarsDeleteCmd deletes a calendar.
var calendarsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a calendar",
	Long: `Delete a calendar permanently.

This action cannot be undone. All events in the calendar
will be permanently deleted.

You can only delete calendars you own. The primary calendar
cannot be deleted.

Requires --confirm flag for safety.`,
	Aliases: []string{"rm", "remove"},
	Example: `  # Delete a calendar (requires confirmation)
  goog cal calendars delete "example@group.calendar.google.com" --confirm`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !calendarsConfirm {
			return fmt.Errorf("deletion requires --confirm flag")
		}
		return nil
	},
	RunE: runCalendarsDelete,
}

// calendarsClearCmd clears all events from a calendar.
var calendarsClearCmd = &cobra.Command{
	Use:   "clear <id>",
	Short: "Clear all events from a calendar",
	Long: `Clear all events from a calendar.

This action cannot be undone. All events in the calendar
will be permanently deleted, but the calendar itself remains.

Note: Only primary calendars can be cleared.

Requires --confirm flag for safety.`,
	Example: `  # Clear all events from primary calendar (requires confirmation)
  goog cal calendars clear primary --confirm`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !calendarsConfirm {
			return fmt.Errorf("clearing all events requires --confirm flag")
		}
		return nil
	},
	RunE: runCalendarsClear,
}

func init() {
	// Add calendars subcommands
	calendarsCmd.AddCommand(calendarsListCmd)
	calendarsCmd.AddCommand(calendarsShowCmd)
	calendarsCmd.AddCommand(calendarsCreateCmd)
	calendarsCmd.AddCommand(calendarsUpdateCmd)
	calendarsCmd.AddCommand(calendarsDeleteCmd)
	calendarsCmd.AddCommand(calendarsClearCmd)

	// Create flags
	calendarsCreateCmd.Flags().StringVar(&calendarsTitle, "title", "", "calendar title (required)")
	calendarsCreateCmd.Flags().StringVar(&calendarsDescription, "description", "", "calendar description")
	calendarsCreateCmd.Flags().StringVar(&calendarsTimezone, "timezone", "", "calendar time zone (e.g., America/New_York)")
	_ = calendarsCreateCmd.MarkFlagRequired("title")

	// Update flags
	calendarsUpdateCmd.Flags().StringVar(&calendarsTitle, "title", "", "calendar title")
	calendarsUpdateCmd.Flags().StringVar(&calendarsDescription, "description", "", "calendar description")
	calendarsUpdateCmd.Flags().StringVar(&calendarsTimezone, "timezone", "", "calendar time zone (e.g., America/New_York)")

	// Delete flags
	calendarsDeleteCmd.Flags().BoolVar(&calendarsConfirm, "confirm", false, "confirm deletion")

	// Clear flags
	calendarsClearCmd.Flags().BoolVar(&calendarsConfirm, "confirm", false, "confirm clearing all events")

	// Add to cal command
	calCmd.AddCommand(calendarsCmd)
}

// getCalendarRepository creates a calendar repository for the current account.
func getCalendarRepository(ctx context.Context) (*repository.GCalCalendarRepository, error) {
	tokenSource, err := getTokenSource(ctx)
	if err != nil {
		return nil, err
	}

	// Create GCal service
	gcalService, err := repository.NewGCalService(ctx, tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	return gcalService.Calendars(), nil
}

// runCalendarsList handles the calendars list command.
func runCalendarsList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get repository using dependency injection
	repo, err := getCalendarRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	calendars, err := repo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list calendars: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	output := p.RenderCalendars(calendars)
	cmd.Println(output)

	return nil
}

// runCalendarsShow handles the calendars show command.
func runCalendarsShow(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	calendarID := args[0]

	// Get repository using dependency injection
	repo, err := getCalendarRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	cal, err := repo.Get(ctx, calendarID)
	if err != nil {
		return fmt.Errorf("calendar not found: %s", calendarID)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	output := p.RenderCalendar(cal)
	cmd.Println(output)

	return nil
}

// runCalendarsCreate handles the calendars create command.
func runCalendarsCreate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get repository using dependency injection
	repo, err := getCalendarRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	// Create new calendar
	cal := calendar.NewCalendar(calendarsTitle)
	cal.Description = calendarsDescription
	cal.TimeZone = calendarsTimezone

	created, err := repo.Create(ctx, cal)
	if err != nil {
		return fmt.Errorf("failed to create calendar: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	if formatFlag == "json" {
		cmd.Println(p.RenderCalendar(created))
	} else {
		cmd.Printf("Calendar created successfully.\n")
		cmd.Printf("ID: %s\n", created.ID)
		cmd.Printf("Title: %s\n", created.Title)
		if created.Description != "" {
			cmd.Printf("Description: %s\n", created.Description)
		}
		if created.TimeZone != "" {
			cmd.Printf("Time Zone: %s\n", created.TimeZone)
		}
	}

	return nil
}

// runCalendarsUpdate handles the calendars update command.
func runCalendarsUpdate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	calendarID := args[0]

	// Get repository using dependency injection
	repo, err := getCalendarRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	// Get existing calendar
	cal, err := repo.Get(ctx, calendarID)
	if err != nil {
		return fmt.Errorf("calendar not found: %s", calendarID)
	}

	// Check if calendar can be modified
	if !cal.IsOwner() {
		return fmt.Errorf("cannot modify calendar: insufficient permissions (access role: %s)", cal.AccessRole)
	}

	// Update fields if provided
	if calendarsTitle != "" {
		cal.Title = calendarsTitle
	}
	if calendarsDescription != "" {
		cal.Description = calendarsDescription
	}
	if calendarsTimezone != "" {
		cal.TimeZone = calendarsTimezone
	}

	updated, err := repo.Update(ctx, cal)
	if err != nil {
		return fmt.Errorf("failed to update calendar: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	if formatFlag == "json" {
		cmd.Println(p.RenderCalendar(updated))
	} else {
		cmd.Printf("Calendar updated successfully.\n")
		cmd.Printf("ID: %s\n", updated.ID)
		cmd.Printf("Title: %s\n", updated.Title)
	}

	return nil
}

// runCalendarsDelete handles the calendars delete command.
func runCalendarsDelete(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	calendarID := args[0]

	// Get repository using dependency injection
	repo, err := getCalendarRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	// Get calendar to verify it exists and check permissions
	cal, err := repo.Get(ctx, calendarID)
	if err != nil {
		return fmt.Errorf("calendar not found: %s", calendarID)
	}

	// Check if calendar can be deleted
	if !cal.IsOwner() {
		return fmt.Errorf("cannot delete calendar: insufficient permissions (access role: %s)", cal.AccessRole)
	}

	// Cannot delete primary calendar
	if cal.Primary {
		return fmt.Errorf("cannot delete primary calendar")
	}

	if err := repo.Delete(ctx, calendarID); err != nil {
		return fmt.Errorf("failed to delete calendar: %w", err)
	}

	if !quietFlag {
		cmd.Printf("Calendar '%s' deleted successfully.\n", cal.Title)
	}

	return nil
}

// runCalendarsClear handles the calendars clear command.
func runCalendarsClear(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	calendarID := args[0]

	// Get repository using dependency injection
	repo, err := getCalendarRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	// Get calendar to verify it exists
	cal, err := repo.Get(ctx, calendarID)
	if err != nil {
		return fmt.Errorf("calendar not found: %s", calendarID)
	}

	if err := repo.Clear(ctx, calendarID); err != nil {
		return fmt.Errorf("failed to clear calendar: %w", err)
	}

	if !quietFlag {
		cmd.Printf("All events cleared from calendar '%s'.\n", cal.Title)
	}

	return nil
}
