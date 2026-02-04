# goog

A command-line interface for Google Mail and Calendar, built with Go.

Designed for both human operators and AI agents, `goog` provides programmatic access to Gmail and Google Calendar through a clean, scriptable interface.

## Features

- **Multi-account support** - Manage multiple Google accounts with easy switching
- **Gmail integration** - List, read, send, reply, forward, and manage messages
- **Calendar integration** - Events, calendars, sharing, and availability queries
- **Flexible output** - JSON, table, and plain text formats
- **Secure credentials** - OAuth2 tokens stored in system keyring

## Installation

```bash
# Build from source
go build -o bin/goog ./cmd/goog

# Install to $GOPATH/bin
go install ./cmd/goog
```

## Setup

Before using `goog`, you need to configure Google OAuth2 credentials:

1. Create a Google Cloud project and enable Gmail/Calendar APIs
2. Configure OAuth consent screen
3. Create OAuth credentials (Desktop app)
4. Set environment variables:
   ```bash
   export GOOG_CLIENT_ID="your-client-id.apps.googleusercontent.com"
   export GOOG_CLIENT_SECRET="your-client-secret"
   ```

**See [documentation/SETUP.md](documentation/SETUP.md) for detailed step-by-step instructions with screenshots.**

## Quick Start

```bash
# Authenticate with Google
goog auth login

# List recent emails
goog mail list

# Show today's calendar
goog cal today

# Create an event
goog cal create --title "Meeting" --start "tomorrow 2pm" --end "tomorrow 3pm"
```

## Command Reference

### Authentication

```bash
goog auth login              # Start OAuth flow
goog auth logout             # Remove stored credentials
goog auth status             # Show authentication status
goog auth refresh            # Force token refresh
```

### Account Management

```bash
goog account list            # List all accounts
goog account add [alias]     # Add new account
goog account remove <alias>  # Remove account
goog account switch <alias>  # Set default account
goog account show            # Show current account
goog account rename <old> <new>  # Rename account
```

### Gmail - Messages

```bash
goog mail list               # List inbox messages
goog mail read <id>          # Read message content
goog mail search <query>     # Search messages
goog mail send               # Send new message
goog mail reply <id>         # Reply to message
goog mail forward <id>       # Forward message
goog mail trash <id>         # Move to trash
goog mail untrash <id>       # Restore from trash
goog mail archive <id>       # Archive message
goog mail delete <id>        # Permanently delete (--confirm required)
goog mail modify <id>        # Modify labels
goog mail mark <id>          # Mark read/unread/starred
```

### Gmail - Drafts

```bash
goog draft list              # List all drafts
goog draft show <id>         # Show draft content
goog draft create            # Create new draft
goog draft update <id>       # Update draft
goog draft send <id>         # Send draft
goog draft delete <id>       # Delete draft
```

### Gmail - Labels

```bash
goog label list              # List all labels
goog label show <name>       # Show label details
goog label create <name>     # Create new label
goog label update <name>     # Update label
goog label delete <name>     # Delete label (--confirm required)
```

### Gmail - Threads

```bash
goog thread list             # List threads
goog thread show <id>        # Show thread with all messages
goog thread trash <id>       # Trash entire thread
goog thread modify <id>      # Modify thread labels
```

### Calendar - Events

```bash
goog cal list                # List upcoming events
goog cal show <id>           # Show event details
goog cal today               # Today's events
goog cal week                # This week's events
goog cal create              # Create new event
goog cal update <id>         # Update event
goog cal delete <id>         # Delete event (--confirm required)
goog cal quick <text>        # Create from natural language
goog cal move <id>           # Move to different calendar
goog cal rsvp <id>           # Respond to invitation
goog cal instances <id>      # List recurring event instances
goog cal freebusy            # Check availability
```

### Calendar - Calendars

```bash
goog cal calendars list      # List all calendars
goog cal calendars show <id> # Show calendar details
goog cal calendars create    # Create new calendar
goog cal calendars update <id>   # Update calendar
goog cal calendars delete <id>   # Delete calendar (--confirm required)
goog cal calendars clear <id>    # Clear all events (--confirm required)
```

### Calendar - Sharing (ACL)

```bash
goog cal acl list <calendar-id>           # List sharing rules
goog cal acl add <calendar-id>            # Add sharing rule
goog cal acl remove <calendar-id> <rule>  # Remove sharing rule
goog cal share <calendar-id>              # Share calendar (alias)
goog cal unshare <calendar-id> <rule>     # Unshare calendar (alias)
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--account <alias>` | Use specific account |
| `--format <type>` | Output format: json, table, plain |
| `--quiet` | Suppress non-essential output |
| `--verbose` | Verbose output |
| `--config <path>` | Config file path |

## Examples

### Multi-Account Workflow

```bash
# Add accounts
goog account add personal
goog account add work

# Use specific account
goog mail list --account work

# Switch default account
goog account switch personal
```

### Email Operations

```bash
# Search for unread from a sender
goog mail search "from:boss@company.com is:unread"

# Send an email
goog mail send --to user@example.com --subject "Hello" --body "Message content"

# Reply to a message
goog mail reply abc123 --body "Thanks for your message"

# Forward with intro
goog mail forward abc123 --to colleague@example.com --body "FYI - see below"
```

### Calendar Operations

```bash
# Create event with attendees
goog cal create --title "Team Meeting" \
  --start "2024-01-15 14:00" --end "2024-01-15 15:00" \
  --location "Conference Room A" \
  --attendees user1@example.com,user2@example.com

# Quick add using natural language
goog cal quick "Lunch with Sarah tomorrow at noon"

# Check availability
goog cal freebusy --start "2024-01-15T09:00:00Z" --end "2024-01-15T17:00:00Z"

# Respond to invitation
goog cal rsvp abc123 --accept
```

### JSON Output for Scripting

```bash
# Get messages as JSON for processing
goog mail list --format json | jq '.[] | .subject'

# Get events for today
goog cal today --format json
```

## Project Structure

```
cmd/goog/          # Application entry point
internal/
  domain/          # Business entities (mail, calendar, account)
  usecase/         # Application business logic
  adapter/
    cli/           # Command handlers
    presenter/     # Output formatters
    repository/    # Google API implementations
  infrastructure/  # Auth, config, keyring
```

## Development

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Build
go build -o bin/goog ./cmd/goog

# Lint
golangci-lint run
```

See [AGENTS.md](AGENTS.md) for development guidelines.

## Documentation

- [Setup Guide](documentation/SETUP.md) - OAuth2 configuration and first-time setup
- [Product Summary](documentation/product-summary.md) - Overview and goals
- [Product Details](documentation/product-details.md) - Features and workflows
- [Technical Details](documentation/technical-details.md) - Architecture and APIs

## License

MIT
