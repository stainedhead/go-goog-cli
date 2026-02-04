// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"github.com/spf13/cobra"
)

var (
	// Global flags
	accountFlag string
	formatFlag  string
	quietFlag   bool
	verboseFlag bool
	configFlag  string
)

// Version information set at build time.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "goog",
	Short: "CLI for Google Mail, Calendar, and Tasks",
	Long: `goog is a command-line interface for managing Google Mail, Calendar, and Tasks services.

Designed for both human operators and AI agents, it provides programmatic
access to Gmail, Google Calendar, and Google Tasks through a clean, scriptable interface.

Examples:
  goog auth login                    # Authenticate with Google
  goog mail list                     # List inbox messages
  goog mail search "is:unread"       # Search for unread messages
  goog cal today                     # Show today's events
  goog cal create --title "Meeting"  # Create a calendar event
  goog tasks list                    # List tasks
  goog tasks create "Buy groceries"  # Create a task`,
}

// versionCmd prints the version information.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("goog %s (commit: %s, built: %s)\n", version, commit, date)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags available to all commands
	rootCmd.PersistentFlags().StringVar(&accountFlag, "account", "", "use specific account")
	rootCmd.PersistentFlags().StringVar(&formatFlag, "format", "table", "output format (json|plain|table)")
	rootCmd.PersistentFlags().BoolVar(&quietFlag, "quiet", false, "suppress non-essential output")
	rootCmd.PersistentFlags().BoolVar(&verboseFlag, "verbose", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&configFlag, "config", "", "config file path")

	// Add subcommands
	rootCmd.AddCommand(versionCmd)
}
