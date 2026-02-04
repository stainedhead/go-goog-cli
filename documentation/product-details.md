# Product Details

## Features

### Authentication

OAuth2/PKCE flow with browser-based consent:

```bash
goog auth login                    # Start auth flow
goog auth login --account work     # Add as named account
goog auth login --scopes gmail.modify,calendar  # Specific scopes
goog auth status                   # Check token status
goog auth refresh                  # Force token refresh
goog auth logout                   # Remove credentials
```

Scope shorthand supported: `gmail`, `gmail.modify`, `calendar`, `calendar.full`, `tasks`, `tasks.readonly`, etc.

### Multi-Account Management

```bash
goog account add personal          # Add account with alias
goog account list                  # Show all accounts
goog account switch work           # Set default
goog account show                  # Current account info
goog account rename old new        # Rename alias
goog account remove work           # Remove account
```

Account resolution order:
1. `--account` flag
2. `GOOG_ACCOUNT` environment variable
3. Default account in config
4. First authenticated account

### Gmail - Messages

List and search:
```bash
goog mail list                     # Inbox messages
goog mail list --labels STARRED    # By label
goog mail list --unread-only       # Unread only
goog mail search "is:unread from:boss@company.com"
```

Read and actions:
```bash
goog mail read <id>                # Full message content
goog mail trash <id>               # Move to trash
goog mail untrash <id>             # Restore from trash
goog mail archive <id>             # Remove from inbox
goog mail delete <id> --confirm    # Permanent delete
goog mail mark <id> --read         # Mark as read
goog mail mark <id> --star         # Add star
goog mail modify <id> --add-labels IMPORTANT
goog mail move <id> --to "Work"    # Move to specific label
```

Compose:
```bash
goog mail send --to user@example.com --subject "Hello" --body "Content"
goog mail send --to a@ex.com --cc b@ex.com --bcc c@ex.com --html
goog mail reply <id> --body "Thanks"
goog mail reply <id> --body "Thanks" --all    # Reply all
goog mail forward <id> --to user@example.com --body "FYI"
```

### Gmail - Drafts

```bash
goog draft list                    # List all drafts
goog draft show <id>               # View draft content
goog draft create --to user@example.com --subject "Draft" --body "Content"
goog draft update <id> --subject "New Subject"
goog draft send <id>               # Send draft
goog draft delete <id>             # Delete draft
```

### Gmail - Labels

```bash
goog label list                    # All labels
goog label show "Work"             # Label details
goog label create "Projects" --background "#4285f4" --text "#ffffff"
goog label update "Projects" --background "#ff0000"
goog label delete "Old" --confirm
```

### Gmail - Threads

```bash
goog thread list                   # List conversations
goog thread list --labels INBOX --max-results 50
goog thread show <id>              # Full conversation
goog thread trash <id>             # Trash entire thread
goog thread untrash <id>           # Restore thread from trash
goog thread delete <id> --confirm  # Permanently delete thread
goog thread modify <id> --add-labels Archive --remove-labels INBOX
```

### Calendar - Events

List and view:
```bash
goog cal list                      # Upcoming 30 days
goog cal today                     # Today's events
goog cal week                      # This week's events
goog cal show <id>                 # Event details
goog cal instances <id>            # Recurring event instances
```

Create and update:
```bash
goog cal create --title "Meeting" --start "2024-01-15 14:00" --end "2024-01-15 15:00"
goog cal create --title "Holiday" --start "2024-01-15" --all-day
goog cal create --title "Team Sync" --start "tomorrow 10am" --attendees a@ex.com,b@ex.com
goog cal quick "Lunch with John tomorrow at noon"
goog cal update <id> --title "New Title" --location "Room B"
goog cal delete <id> --confirm
goog cal move <id> --to work@group.calendar.google.com
```

RSVP:
```bash
goog cal rsvp <id> --accept
goog cal rsvp <id> --decline
goog cal rsvp <id> --tentative
```

Availability:
```bash
goog cal freebusy --start "2024-01-15T09:00:00Z" --end "2024-01-15T17:00:00Z"
goog cal freebusy --calendars "primary,team@group.calendar.google.com"
```

### Calendar - Calendars

```bash
goog cal calendars list            # All calendars
goog cal calendars show primary    # Calendar details
goog cal calendars create --title "Projects" --description "Work projects"
goog cal calendars update <id> --title "New Name"
goog cal calendars delete <id> --confirm
goog cal calendars clear <id> --confirm    # Remove all events
```

### Calendar - Sharing (ACL)

```bash
goog cal acl list primary                              # List rules
goog cal acl add primary --email user@ex.com --role reader
goog cal acl add primary --email user@ex.com --role writer
goog cal acl remove primary "user:user@ex.com" --confirm

# Aliases
goog cal share primary --email user@ex.com --role reader
goog cal unshare primary "user:user@ex.com" --confirm
```

Roles: `reader`, `writer`, `owner`, `freeBusyReader`

### Tasks - Task Lists

```bash
goog tasks lists                       # List all task lists
goog tasks create-list "Work Tasks"    # Create new list
goog tasks update-list <id> --title "New Name"
goog tasks delete-list <id> --confirm  # Delete list
```

### Tasks - Tasks

List and view:
```bash
goog tasks list                        # List tasks in default list
goog tasks list --list <list-id>       # Tasks from specific list
goog tasks list --show-completed       # Include completed
goog tasks list --max-results 50       # Limit results
goog tasks get <id>                    # Task details
```

Create and update:
```bash
goog tasks create "Buy groceries"
goog tasks create "Submit report" --due "2024-12-31" --notes "Q4 report"
goog tasks create "Review slides" --parent <parent-id>    # Create subtask
goog tasks update <id> --title "New title"
goog tasks update <id> --notes "Updated notes" --due "2024-12-31"
```

Actions:
```bash
goog tasks complete <id>               # Mark as completed
goog tasks reopen <id>                 # Reopen completed task
goog tasks move <id> --parent <parent-id>      # Move under parent
goog tasks move <id> --previous <previous-id>  # Reorder tasks
goog tasks delete <id> --confirm       # Delete task
goog tasks clear --confirm             # Clear completed tasks
```

Task Lists:
- Default list ID: `@default`
- Each account can have multiple task lists
- Tasks belong to one list

Subtasks:
- Tasks can have one level of subtasks (parent-child only)
- Use `--parent` flag when creating to make a subtask
- Maximum 2000 subtasks per parent task

## Output Formats

| Format | Flag | Use Case |
|--------|------|----------|
| Table | `--format table` | Human-readable (default) |
| JSON | `--format json` | Scripting, parsing |
| Plain | `--format plain` | Simple text output |

## User Workflows

### Daily Email Routine
```bash
goog mail list --unread-only       # Check unread
goog mail read abc123              # Read message
goog mail reply abc123 --body "On it"
goog mail archive abc123           # Archive when done
```

### Meeting Scheduling
```bash
goog cal freebusy --start "..." --end "..."   # Check availability
goog cal create --title "Meeting" --start "..." --attendees ...
```

### Task Management
```bash
goog tasks list                        # Check today's tasks
goog tasks create "New task" --due "2024-12-31"
goog tasks complete abc123             # Mark done
goog tasks clear --confirm             # Clean up completed
```

### Automated Notifications (AI Agent)
```bash
goog mail search "from:alerts@system.com" --format json | jq ...
goog cal today --format json | jq ...
goog tasks list --format json | jq ...
```
