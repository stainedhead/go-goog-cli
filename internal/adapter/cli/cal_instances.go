// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/adapter/presenter"
)

// Command flags for instances command.
var (
	calInstancesStart      string
	calInstancesEnd        string
	calInstancesMaxResults int
	calInstancesCalendar   string
)

// calInstancesCmd lists instances of a recurring event.
var calInstancesCmd = &cobra.Command{
	Use:   "instances <event-id>",
	Short: "List instances of a recurring event",
	Long: `List instances of a recurring event.

Returns the individual occurrences of a recurring calendar event
within the specified time range. If no time range is specified,
returns instances from now onwards.

Use this command to see when recurring events are scheduled,
modify individual instances, or check for conflicts.`,
	Example: `  # List instances of a recurring event
  goog cal instances abc123def456

  # List instances within a specific time range
  goog cal instances abc123def456 --start "2024-01-01T00:00:00Z" --end "2024-03-01T00:00:00Z"

  # List instances from a specific calendar
  goog cal instances abc123def456 --calendar work@group.calendar.google.com

  # Limit number of results
  goog cal instances abc123def456 --max-results 10

  # Output as JSON
  goog cal instances abc123def456 --format json`,
	Aliases: []string{"recurring"},
	Args:    cobra.ExactArgs(1),
	RunE:    runCalInstances,
}

func init() {
	// Add instances subcommand to calCmd (defined in cal_utils.go)
	calCmd.AddCommand(calInstancesCmd)

	// Instances command flags
	calInstancesCmd.Flags().StringVar(&calInstancesCalendar, "calendar", "primary", "calendar ID to use")
	calInstancesCmd.Flags().StringVar(&calInstancesStart, "start", "", "start time for instances (RFC3339 format)")
	calInstancesCmd.Flags().StringVar(&calInstancesEnd, "end", "", "end time for instances (RFC3339 format)")
	calInstancesCmd.Flags().IntVar(&calInstancesMaxResults, "max-results", 25, "maximum number of instances to return")
}

// runCalInstances handles the cal instances command.
func runCalInstances(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	eventID := args[0]

	// Parse start time if provided
	var timeMin time.Time
	if calInstancesStart != "" {
		var err error
		timeMin, err = time.Parse(time.RFC3339, calInstancesStart)
		if err != nil {
			return fmt.Errorf("invalid start time format (use RFC3339, e.g., 2024-01-15T09:00:00Z): %w", err)
		}
	}

	// Parse end time if provided
	var timeMax time.Time
	if calInstancesEnd != "" {
		var err error
		timeMax, err = time.Parse(time.RFC3339, calInstancesEnd)
		if err != nil {
			return fmt.Errorf("invalid end time format (use RFC3339, e.g., 2024-01-15T17:00:00Z): %w", err)
		}
	}

	// Validate time range if both are provided
	if !timeMin.IsZero() && !timeMax.IsZero() && !timeMin.Before(timeMax) {
		return fmt.Errorf("start time must be before end time")
	}

	// Get event repository using dependency injection
	repo, err := getEventRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	// Get instances of the recurring event
	instances, err := repo.Instances(ctx, calInstancesCalendar, eventID, timeMin, timeMax)
	if err != nil {
		return fmt.Errorf("failed to get event instances: %w", err)
	}

	// Apply max results limit
	if calInstancesMaxResults > 0 && len(instances) > calInstancesMaxResults {
		instances = instances[:calInstancesMaxResults]
	}

	// Create presenter based on format flag
	p := presenter.New(formatFlag)

	// Output result
	output := p.RenderEvents(instances)
	cmd.Println(output)

	// Show instance count if not quiet
	if !quietFlag && len(instances) > 0 {
		cmd.Printf("\n%d instance(s) of recurring event\n", len(instances))
	}

	return nil
}
