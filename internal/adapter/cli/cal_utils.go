// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/adapter/presenter"
	"github.com/stainedhead/go-goog-cli/internal/adapter/repository"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/config"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/keyring"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

// Command flags for calendar utilities.
var (
	// Common calendar flags
	calCalendarFlag string

	// Free/busy flags
	calFreeBusyStart     string
	calFreeBusyEnd       string
	calFreeBusyCalendars []string

	// RSVP flags
	calRSVPAccept    bool
	calRSVPDecline   bool
	calRSVPTentative bool

	// Move flags
	calMoveDestination string
)

// calCmd represents the calendar command group.
var calCmd = &cobra.Command{
	Use:   "cal",
	Short: "Manage Google Calendar",
	Long: `Manage Google Calendar events and schedules.

The cal commands allow you to view, create, and manage
calendar events in your Google Calendar account.`,
	Aliases: []string{"calendar"},
}

// calQuickAddCmd creates an event from natural language.
var calQuickAddCmd = &cobra.Command{
	Use:   "quick <text>",
	Short: "Create event from natural language",
	Long: `Create a new calendar event using natural language text.

Google Calendar will parse the text and create an event
based on the description. The text can include date, time,
location, and other event details.`,
	Example: `  # Create a meeting tomorrow at 3pm
  goog cal quick "Meeting with John tomorrow at 3pm"

  # Create an event with specific date
  goog cal quick "Doctor appointment Friday at 10am"

  # Create an event in a specific calendar
  goog cal quick "Team standup Monday 9am" --calendar work@group.calendar.google.com`,
	Aliases: []string{"quickadd", "add"},
	Args:    cobra.ExactArgs(1),
	RunE:    runCalQuickAdd,
}

// calFreeBusyCmd checks availability.
var calFreeBusyCmd = &cobra.Command{
	Use:   "freebusy",
	Short: "Check calendar availability",
	Long: `Query free/busy information for one or more calendars.

Returns busy time periods within the specified time range.
Use this to check availability before scheduling meetings.`,
	Example: `  # Check availability for the next 24 hours
  goog cal freebusy --start "2024-01-15T09:00:00Z" --end "2024-01-15T17:00:00Z"

  # Check multiple calendars
  goog cal freebusy --start "2024-01-15T09:00:00Z" --end "2024-01-16T17:00:00Z" \
    --calendars "primary,work@group.calendar.google.com"

  # Check with specific calendar
  goog cal freebusy --start "2024-01-15T09:00:00Z" --end "2024-01-15T18:00:00Z" \
    --calendars "team@group.calendar.google.com"`,
	Aliases: []string{"busy", "availability"},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if calFreeBusyStart == "" {
			return fmt.Errorf("--start flag is required")
		}
		if calFreeBusyEnd == "" {
			return fmt.Errorf("--end flag is required")
		}
		return nil
	},
	RunE: runCalFreeBusy,
}

// calRSVPCmd responds to an event invitation.
var calRSVPCmd = &cobra.Command{
	Use:   "rsvp <event-id>",
	Short: "Respond to event invitation",
	Long: `Update your response status for an event invitation.

You must specify exactly one of --accept, --decline, or --tentative.`,
	Example: `  # Accept an invitation
  goog cal rsvp abc123xyz --accept

  # Decline an invitation
  goog cal rsvp abc123xyz --decline

  # Mark as tentative
  goog cal rsvp abc123xyz --tentative

  # RSVP on a specific calendar
  goog cal rsvp abc123xyz --accept --calendar work@group.calendar.google.com`,
	Aliases: []string{"respond"},
	Args:    cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Count how many RSVP options are set
		count := 0
		if calRSVPAccept {
			count++
		}
		if calRSVPDecline {
			count++
		}
		if calRSVPTentative {
			count++
		}

		if count == 0 {
			return fmt.Errorf("one of --accept, --decline, or --tentative is required")
		}
		if count > 1 {
			return fmt.Errorf("only one of --accept, --decline, or --tentative can be specified")
		}
		return nil
	},
	RunE: runCalRSVP,
}

// calMoveCmd moves an event to a different calendar.
var calMoveCmd = &cobra.Command{
	Use:   "move <event-id>",
	Short: "Move event to different calendar",
	Long: `Move an event from one calendar to another.

The event will be removed from the source calendar and
added to the destination calendar.`,
	Example: `  # Move event to a different calendar
  goog cal move abc123xyz --to work@group.calendar.google.com

  # Move from a specific source calendar
  goog cal move abc123xyz --calendar primary --to personal@group.calendar.google.com`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if calMoveDestination == "" {
			return fmt.Errorf("--to flag is required (destination calendar)")
		}
		return nil
	},
	RunE: runCalMove,
}

func init() {
	// Add calendar subcommands
	calCmd.AddCommand(calQuickAddCmd)
	calCmd.AddCommand(calFreeBusyCmd)
	calCmd.AddCommand(calRSVPCmd)
	calCmd.AddCommand(calMoveCmd)

	// Quick add flags
	calQuickAddCmd.Flags().StringVar(&calCalendarFlag, "calendar", "primary", "calendar ID to create event in")

	// Free/busy flags
	calFreeBusyCmd.Flags().StringVar(&calFreeBusyStart, "start", "", "start time (RFC3339 format, required)")
	calFreeBusyCmd.Flags().StringVar(&calFreeBusyEnd, "end", "", "end time (RFC3339 format, required)")
	calFreeBusyCmd.Flags().StringSliceVar(&calFreeBusyCalendars, "calendars", []string{"primary"}, "calendar IDs to check (comma-separated)")

	// RSVP flags
	calRSVPCmd.Flags().BoolVar(&calRSVPAccept, "accept", false, "accept the invitation")
	calRSVPCmd.Flags().BoolVar(&calRSVPDecline, "decline", false, "decline the invitation")
	calRSVPCmd.Flags().BoolVar(&calRSVPTentative, "tentative", false, "mark as tentative")
	calRSVPCmd.Flags().StringVar(&calCalendarFlag, "calendar", "primary", "calendar ID containing the event")

	// Move flags
	calMoveCmd.Flags().StringVar(&calMoveDestination, "to", "", "destination calendar ID (required)")
	calMoveCmd.Flags().StringVar(&calCalendarFlag, "calendar", "primary", "source calendar ID")

	// Add calendar command to root
	rootCmd.AddCommand(calCmd)
}

// getGCalService creates a GCalService using the current account's credentials.
func getGCalService(ctx context.Context) (*repository.GCalService, error) {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Create keyring store
	store, err := keyring.NewStore()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize keyring: %w", err)
	}

	// Create account service
	svc := accountuc.NewService(cfg, store, nil)

	// Resolve account
	acc, err := svc.ResolveAccount(accountFlag)
	if err != nil {
		return nil, fmt.Errorf("no account found: %w (run 'goog auth login' to authenticate)", err)
	}

	// Get token source
	tokenMgr := svc.GetTokenManager()
	tokenSource, err := tokenMgr.GetTokenSource(ctx, acc.Alias)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w (run 'goog auth login' to authenticate)", err)
	}

	// Create GCal service
	gcalSvc, err := repository.NewGCalService(ctx, tokenSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create Calendar client: %w", err)
	}

	return gcalSvc, nil
}

// runCalQuickAdd handles the cal quick command.
func runCalQuickAdd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	text := args[0]

	// Get calendar service
	gcalSvc, err := getGCalService(ctx)
	if err != nil {
		return err
	}

	// Get events repository
	eventRepo := gcalSvc.Events()

	// Quick add the event
	event, err := eventRepo.QuickAdd(ctx, calCalendarFlag, text)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	// Output result
	output := p.RenderEvent(event)
	cmd.Println(output)

	if !quietFlag {
		cmd.Printf("\nEvent created successfully\n")
		if event.HTMLLink != "" {
			cmd.Printf("View at: %s\n", event.HTMLLink)
		}
	}

	return nil
}

// runCalFreeBusy handles the cal freebusy command.
func runCalFreeBusy(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Parse start time
	startTime, err := time.Parse(time.RFC3339, calFreeBusyStart)
	if err != nil {
		return fmt.Errorf("invalid start time format (use RFC3339, e.g., 2024-01-15T09:00:00Z): %w", err)
	}

	// Parse end time
	endTime, err := time.Parse(time.RFC3339, calFreeBusyEnd)
	if err != nil {
		return fmt.Errorf("invalid end time format (use RFC3339, e.g., 2024-01-15T17:00:00Z): %w", err)
	}

	// Validate time range
	if !startTime.Before(endTime) {
		return calendar.ErrInvalidTimeRange
	}

	// Get calendar service
	gcalSvc, err := getGCalService(ctx)
	if err != nil {
		return err
	}

	// Get free/busy repository
	freeBusyRepo := gcalSvc.FreeBusy()

	// Create request
	request, err := calendar.NewFreeBusyRequest(startTime, endTime, calFreeBusyCalendars...)
	if err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	// Query free/busy
	response, err := freeBusyRepo.Query(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to query free/busy: %w", err)
	}

	// Output result based on format
	if formatFlag == presenter.FormatJSON {
		// Render custom JSON structure for free/busy response
		cmd.Println(renderFreeBusyJSON(response))
	} else {
		// Render as table or plain text
		output := renderFreeBusyTable(response, startTime, endTime)
		cmd.Println(output)
	}

	return nil
}

// renderFreeBusyTable renders free/busy response as a table.
func renderFreeBusyTable(response *calendar.FreeBusyResponse, start, end time.Time) string {
	if response == nil || len(response.Calendars) == 0 {
		return "No busy periods found"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Free/Busy Information (%s to %s)\n\n",
		start.Format("2006-01-02 15:04"),
		end.Format("2006-01-02 15:04")))

	for calID, periods := range response.Calendars {
		sb.WriteString(fmt.Sprintf("Calendar: %s\n", calID))
		if len(periods) == 0 {
			sb.WriteString("  No busy periods (free)\n")
		} else {
			for _, period := range periods {
				sb.WriteString(fmt.Sprintf("  BUSY: %s - %s (%s)\n",
					period.Start.Format("2006-01-02 15:04"),
					period.End.Format("2006-01-02 15:04"),
					period.Duration().String()))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderFreeBusyJSON renders free/busy response as JSON string.
func renderFreeBusyJSON(response *calendar.FreeBusyResponse) string {
	if response == nil {
		return "{}"
	}

	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString("  \"calendars\": {\n")

	calIndex := 0
	for calID, periods := range response.Calendars {
		if calIndex > 0 {
			sb.WriteString(",\n")
		}
		sb.WriteString(fmt.Sprintf("    \"%s\": [\n", calID))

		for i, period := range periods {
			if i > 0 {
				sb.WriteString(",\n")
			}
			sb.WriteString(fmt.Sprintf("      {\"start\": \"%s\", \"end\": \"%s\"}",
				period.Start.Format(time.RFC3339),
				period.End.Format(time.RFC3339)))
		}

		sb.WriteString("\n    ]")
		calIndex++
	}

	sb.WriteString("\n  }\n")
	sb.WriteString("}")

	return sb.String()
}

// runCalRSVP handles the cal rsvp command.
func runCalRSVP(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	eventID := args[0]

	// Determine response status
	var response string
	switch {
	case calRSVPAccept:
		response = calendar.ResponseAccepted
	case calRSVPDecline:
		response = calendar.ResponseDeclined
	case calRSVPTentative:
		response = calendar.ResponseTentative
	}

	// Get calendar service
	gcalSvc, err := getGCalService(ctx)
	if err != nil {
		return err
	}

	// Get events repository
	eventRepo := gcalSvc.Events()

	// RSVP to the event
	err = eventRepo.RSVP(ctx, calCalendarFlag, eventID, response)
	if err != nil {
		return fmt.Errorf("failed to update RSVP: %w", err)
	}

	if !quietFlag {
		responseText := map[string]string{
			calendar.ResponseAccepted:  "accepted",
			calendar.ResponseDeclined:  "declined",
			calendar.ResponseTentative: "tentative",
		}
		cmd.Printf("RSVP updated: %s\n", responseText[response])
	}

	return nil
}

// runCalMove handles the cal move command.
func runCalMove(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	eventID := args[0]

	// Get calendar service
	gcalSvc, err := getGCalService(ctx)
	if err != nil {
		return err
	}

	// Get events repository
	eventRepo := gcalSvc.Events()

	// Move the event
	event, err := eventRepo.Move(ctx, calCalendarFlag, eventID, calMoveDestination)
	if err != nil {
		return fmt.Errorf("failed to move event: %w", err)
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	// Output result
	output := p.RenderEvent(event)
	cmd.Println(output)

	if !quietFlag {
		cmd.Printf("\nEvent moved to calendar: %s\n", calMoveDestination)
	}

	return nil
}
