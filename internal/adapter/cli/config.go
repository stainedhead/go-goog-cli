// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/config"
)

// configCmd represents the config command group.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long: `Manage goog configuration.

The config commands allow you to view and modify configuration
settings such as default account, output format, and service-specific
options.`,
}

// configShowCmd shows the current configuration.
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	Long: `Display the current configuration in YAML format.

This shows all configuration values including accounts,
mail settings, and calendar settings.`,
	Aliases: []string{"view", "list"},
	RunE:    runConfigShow,
}

// configSetCmd sets a configuration value.
var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value.

Available keys:
  default_account          - Default account alias
  default_format           - Default output format (json|plain|table)
  timezone                 - Timezone for date/time display
  mail.default_label       - Default mail label to list
  mail.page_size           - Default number of messages per page
  calendar.default_calendar - Default calendar ID
  calendar.week_start      - First day of week (sunday|monday)`,
	Example: `  # Set default format to JSON
  goog config set default_format json

  # Set mail page size
  goog config set mail.page_size 50

  # Set default calendar
  goog config set calendar.default_calendar primary`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigSet,
}

// configGetCmd gets a configuration value.
var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Long: `Get a configuration value by key.

Available keys:
  default_account          - Default account alias
  default_format           - Default output format
  timezone                 - Timezone for date/time display
  mail.default_label       - Default mail label
  mail.page_size           - Messages per page
  calendar.default_calendar - Default calendar ID
  calendar.week_start      - First day of week`,
	Example: `  # Get default format
  goog config get default_format

  # Get mail page size
  goog config get mail.page_size`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigGet,
}

// configPathCmd shows the config file path.
var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	Long:  `Show the path to the configuration file.`,
	RunE:  runConfigPath,
}

func init() {
	// Add config subcommands
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configPathCmd)

	// Add to root
	rootCmd.AddCommand(configCmd)
}

// runConfigShow handles the config show command.
func runConfigShow(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Output config in YAML-like format
	cmd.Printf("default_account: %s\n", cfg.DefaultAccount)
	cmd.Printf("default_format: %s\n", cfg.DefaultFormat)
	cmd.Printf("timezone: %s\n", cfg.Timezone)

	cmd.Println()
	cmd.Println("mail:")
	cmd.Printf("  default_label: %s\n", cfg.Mail.DefaultLabel)
	cmd.Printf("  page_size: %d\n", cfg.Mail.PageSize)

	cmd.Println()
	cmd.Println("calendar:")
	cmd.Printf("  default_calendar: %s\n", cfg.Calendar.DefaultCalendar)
	cmd.Printf("  week_start: %s\n", cfg.Calendar.WeekStart)

	if len(cfg.Accounts) > 0 {
		cmd.Println()
		cmd.Println("accounts:")
		// Sort aliases for deterministic output
		aliases := make([]string, 0, len(cfg.Accounts))
		for alias := range cfg.Accounts {
			aliases = append(aliases, alias)
		}
		sort.Strings(aliases)

		for _, alias := range aliases {
			acc := cfg.Accounts[alias]
			cmd.Printf("  %s:\n", alias)
			cmd.Printf("    email: %s\n", acc.Email)
			if !acc.AddedAt.IsZero() {
				cmd.Printf("    added_at: %s\n", acc.AddedAt.Format("2006-01-02T15:04:05Z07:00"))
			}
			if len(acc.Scopes) > 0 {
				cmd.Println("    scopes:")
				for _, scope := range acc.Scopes {
					cmd.Printf("      - %s\n", scope)
				}
			}
		}
	}

	return nil
}

// runConfigSet handles the config set command.
func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Set value
	if err := cfg.SetValue(key, value); err != nil {
		return fmt.Errorf("failed to set config value: %w", err)
	}

	// Save config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	cmd.Printf("Set %s = %s\n", key, value)
	return nil
}

// runConfigGet handles the config get command.
func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get value
	value, err := cfg.GetValue(key)
	if err != nil {
		return fmt.Errorf("failed to get config value: %w", err)
	}

	cmd.Println(value)
	return nil
}

// runConfigPath handles the config path command.
func runConfigPath(cmd *cobra.Command, args []string) error {
	cmd.Println(config.GetConfigPath())
	return nil
}
