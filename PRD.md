# Product Requirements Document: go-goog-cli

## 1. Product Overview

### 1.1 Purpose
go-goog-cli is a command-line interface for managing Google Mail and Calendar services. Designed for both human operators and AI agents, it provides programmatic access to Gmail and Google Calendar through a clean, scriptable interface.

### 1.2 Target Users
- Developers and power users managing email/calendar via terminal
- AI agents requiring Google Workspace integration
- Automation pipelines needing Gmail/Calendar access

### 1.3 Key Objectives
- Secure OAuth2 authentication with PKCE
- Cross-platform secure credential storage
- Agent-friendly output formats (JSON, plain text)
- Least-privilege scope management

---

## 2. Functional Requirements

### 2.1 Authentication

#### 2.1.1 OAuth2 with PKCE Flow
1. User initiates auth: `goog auth login`
2. CLI generates PKCE code verifier (256-bit entropy)
3. Browser opens to Google consent screen
4. Local callback server receives authorization code
5. Exchange code + verifier for access/refresh tokens
6. Store tokens securely in platform keyring

#### 2.1.2 Token Management
| Command | Description |
|---------|-------------|
| `goog auth login` | Authenticate with Google (browser flow) |
| `goog auth logout` | Remove stored credentials |
| `goog auth status` | Show current authentication state |
| `goog auth refresh` | Force token refresh |

#### 2.1.3 Multi-Account Support

**Use Case:** Support separate accounts for human operators and AI agents.

**Account Management Commands:**
| Command | Description |
|---------|-------------|
| `goog account list` | List all configured accounts |
| `goog account add [alias]` | Add new account (opens browser auth) |
| `goog account remove <alias>` | Remove account and its tokens |
| `goog account switch <alias>` | Set default account |
| `goog account show` | Show current active account |
| `goog account rename <old> <new>` | Rename account alias |

**Account Context Resolution (precedence order):**
1. `--account <alias>` flag (per-command override)
2. `GOOG_ACCOUNT` environment variable
3. `default_account` in config file
4. First authenticated account (fallback)

**Token Storage (per account):**
- Each account stores its own OAuth tokens in platform keyring
- Keyring entries namespaced by account alias: `goog:<alias>:refresh_token`
- Account metadata stored in config: email, scopes granted, last used

**Example Workflow:**
```bash
# Setup: Add human and agent accounts
goog account add human    # Auth as human@gmail.com
goog account add agent    # Auth as agent@company.com

# Set human as default
goog account switch human

# Agent runs commands with explicit account
goog mail list --account agent --format json

# Human uses default
goog mail list
```

**Config Schema (accounts section):**
```yaml
default_account: human
accounts:
  human:
    email: user@gmail.com
    scopes: [gmail.modify, calendar.events]
    added: 2024-01-15T10:00:00Z
  agent:
    email: agent@company.com
    scopes: [gmail.readonly, calendar.readonly]
    added: 2024-01-16T14:30:00Z
```

### 2.2 Gmail Operations

#### 2.2.1 Message Commands
| Command | Description |
|---------|-------------|
| `goog mail list` | List messages (inbox by default) |
| `goog mail read <id>` | Read message content |
| `goog mail send` | Send new message |
| `goog mail reply <id>` | Reply to message |
| `goog mail forward <id>` | Forward message |
| `goog mail search <query>` | Search messages |
| `goog mail trash <id>` | Move to trash |
| `goog mail untrash <id>` | Restore from trash |
| `goog mail delete <id>` | Permanently delete (requires confirmation) |
| `goog mail archive <id>` | Archive message |
| `goog mail modify <id>` | Modify labels on message |
| `goog mail mark <id>` | Mark as read/unread/starred |
| `goog mail move <id>` | Move to label/folder |
| `goog mail attachment <id>` | Download attachment |

#### 2.2.2 Draft Commands
| Command | Description |
|---------|-------------|
| `goog draft list` | List all drafts |
| `goog draft show <id>` | Show draft content |
| `goog draft create` | Create new draft |
| `goog draft update <id>` | Update existing draft |
| `goog draft send <id>` | Send a draft |
| `goog draft delete <id>` | Delete a draft |

#### 2.2.3 Thread Commands
| Command | Description |
|---------|-------------|
| `goog thread list` | List message threads |
| `goog thread show <id>` | Show all messages in thread |
| `goog thread modify <id>` | Modify labels on thread |
| `goog thread trash <id>` | Trash entire thread |
| `goog thread untrash <id>` | Restore thread from trash |
| `goog thread delete <id>` | Permanently delete thread |

#### 2.2.4 Label Commands
| Command | Description |
|---------|-------------|
| `goog label list` | List all labels |
| `goog label show <name>` | Show label details |
| `goog label create <name>` | Create new label |
| `goog label update <name>` | Update label properties |
| `goog label delete <name>` | Delete label |

#### 2.2.5 Settings Commands
| Command | Description |
|---------|-------------|
| `goog settings show` | Show all Gmail settings |
| `goog settings vacation` | Get/set vacation responder |
| `goog settings filters list` | List email filters |
| `goog settings filters create` | Create email filter |
| `goog settings filters delete <id>` | Delete email filter |
| `goog settings forwarding` | Manage forwarding addresses |
| `goog settings send-as list` | List send-as aliases |
| `goog settings send-as create` | Create send-as alias |

#### 2.2.6 Common Mail Flags
| Flag | Description |
|------|-------------|
| `--format json\|plain\|table` | Output format |
| `--limit N` | Max results |
| `--label NAME` | Filter by label |
| `--unread` | Unread only |
| `--from EMAIL` | Filter by sender |
| `--after DATE` | Messages after date |
| `--before DATE` | Messages before date |
| `--include-spam-trash` | Include spam/trash in results |
| `--attachment` | Filter messages with attachments |

### 2.3 Calendar Operations

#### 2.3.1 Event Commands
| Command | Description |
|---------|-------------|
| `goog cal list` | List upcoming events |
| `goog cal show <id>` | Show event details |
| `goog cal create` | Create new event |
| `goog cal quick "<text>"` | Create event from natural language |
| `goog cal update <id>` | Update event |
| `goog cal delete <id>` | Delete event |
| `goog cal move <id>` | Move event to another calendar |
| `goog cal today` | Today's events |
| `goog cal week` | This week's events |
| `goog cal freebusy` | Check availability |
| `goog cal instances <id>` | List recurring event instances |
| `goog cal rsvp <id>` | Respond to event invitation |

#### 2.3.2 Calendar Management Commands
| Command | Description |
|---------|-------------|
| `goog calendar list` | List all calendars |
| `goog calendar show <id>` | Show calendar details |
| `goog calendar create` | Create new calendar |
| `goog calendar update <id>` | Update calendar metadata |
| `goog calendar delete <id>` | Delete calendar |
| `goog calendar clear <id>` | Clear all events from calendar |
| `goog calendar subscribe <id>` | Subscribe to a calendar |
| `goog calendar unsubscribe <id>` | Unsubscribe from calendar |

#### 2.3.3 Sharing/ACL Commands
| Command | Description |
|---------|-------------|
| `goog cal share <calendar>` | Share calendar with user/group |
| `goog cal unshare <calendar>` | Remove sharing |
| `goog cal acl list <calendar>` | List access control rules |
| `goog cal acl set <calendar>` | Set access level for user |

#### 2.3.4 Common Calendar Flags
| Flag | Description |
|------|-------------|
| `--calendar NAME` | Target calendar |
| `--start DATETIME` | Event start |
| `--end DATETIME` | Event end |
| `--duration DURATION` | Event duration (e.g., "1h30m") |
| `--attendees EMAILS` | Comma-separated attendees |
| `--location TEXT` | Event location |
| `--description TEXT` | Event description |
| `--recurrence RRULE` | Recurrence rule (e.g., "FREQ=WEEKLY;COUNT=10") |
| `--reminder MINS` | Set reminder (minutes before) |
| `--notify all\|external\|none` | Attendee notification setting |
| `--visibility public\|private` | Event visibility |
| `--status confirmed\|tentative\|cancelled` | Event status |
| `--color ID` | Event color (1-11) |

### 2.4 Global Flags
| Flag | Description |
|------|-------------|
| `--account NAME` | Use specific account |
| `--format json\|plain\|table` | Output format (default: table) |
| `--quiet` | Suppress non-essential output |
| `--verbose` | Verbose output |
| `--config PATH` | Config file path |

### 2.5 Configuration
| Command | Description |
|---------|-------------|
| `goog config show` | Display current config |
| `goog config set KEY VALUE` | Set config value |
| `goog config get KEY` | Get config value |

---

## 3. Technical Architecture

### 3.1 Clean Architecture Layers

```
cmd/goog/
└── main.go                 # Entry point, DI setup

internal/
├── domain/
│   ├── mail/
│   │   ├── message.go      # Message entity
│   │   ├── draft.go        # Draft entity
│   │   ├── thread.go       # Thread entity
│   │   ├── label.go        # Label entity
│   │   ├── attachment.go   # Attachment entity
│   │   ├── filter.go       # Filter entity
│   │   └── repository.go   # Repository interfaces
│   ├── calendar/
│   │   ├── event.go        # Event entity
│   │   ├── calendar.go     # Calendar entity
│   │   ├── acl.go          # ACL entity
│   │   ├── attendee.go     # Attendee entity
│   │   └── repository.go   # Repository interfaces
│   └── account/
│       ├── account.go      # Account entity
│       └── repository.go   # Repository interface
│
├── usecase/
│   ├── mail/
│   │   ├── message.go      # Message use cases (list, read, send, etc.)
│   │   ├── draft.go        # Draft use cases
│   │   ├── thread.go       # Thread use cases
│   │   ├── label.go        # Label use cases
│   │   └── settings.go     # Settings use cases
│   ├── calendar/
│   │   ├── event.go        # Event use cases
│   │   ├── calendar.go     # Calendar management use cases
│   │   └── acl.go          # ACL use cases
│   └── account/
│       └── account.go      # Account management use cases
│
├── adapter/
│   ├── cli/
│   │   ├── root.go         # Root command
│   │   ├── auth.go         # Auth commands
│   │   ├── account.go      # Account commands
│   │   ├── mail.go         # Mail commands
│   │   ├── draft.go        # Draft commands
│   │   ├── thread.go       # Thread commands
│   │   ├── label.go        # Label commands
│   │   ├── settings.go     # Settings commands
│   │   ├── calendar.go     # Calendar event commands
│   │   ├── calendar_mgmt.go # Calendar management commands
│   │   └── acl.go          # ACL commands
│   ├── repository/
│   │   ├── gmail.go        # Gmail API repository
│   │   └── gcalendar.go    # Calendar API repository
│   └── presenter/
│       ├── json.go         # JSON output
│       ├── table.go        # Table output
│       └── plain.go        # Plain text output
│
└── infrastructure/
    ├── auth/
    │   ├── oauth.go        # OAuth2/PKCE flow
    │   └── token.go        # Token management
    ├── keyring/
    │   └── store.go        # Credential storage (per-account)
    └── config/
        └── config.go       # Configuration management
```

### 3.2 Technology Stack

| Component | Package | Purpose |
|-----------|---------|---------|
| CLI Framework | `github.com/spf13/cobra` | Command parsing, help generation |
| Configuration | `github.com/spf13/viper` | Config file, env vars, flags |
| OAuth2 | `golang.org/x/oauth2` | OAuth2 client with PKCE |
| OAuth CLI | `github.com/int128/oauth2cli` | Browser-based auth flow |
| Gmail API | `google.golang.org/api/gmail/v1` | Gmail operations |
| Calendar API | `google.golang.org/api/calendar/v3` | Calendar operations |
| Keyring | `github.com/99designs/keyring` | Secure credential storage |
| Table Output | `github.com/olekukonko/tablewriter` | CLI table formatting |
| JSON | `encoding/json` (stdlib) | JSON output |

### 3.3 Authentication Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    OAuth2 + PKCE Flow                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  User                CLI                  Browser    Google     │
│   │                   │                      │          │       │
│   │──goog auth login──>│                      │          │       │
│   │                   │──generate verifier───│          │       │
│   │                   │──start local server──│          │       │
│   │                   │──────────────────────>│          │       │
│   │                   │      open auth URL    │          │       │
│   │                   │                      │──login──>│       │
│   │                   │                      │<─consent─│       │
│   │                   │<────callback + code───│          │       │
│   │                   │──exchange code+verifier────────>│       │
│   │                   │<─────access + refresh tokens────│       │
│   │                   │──store in keyring────│          │       │
│   │<──authenticated───│                      │          │       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 3.4 Secure Credential Storage

**Per-Account Token Storage:**

| Platform | Keyring Entry Pattern |
|----------|----------------------|
| macOS | Keychain: service="go-goog-cli", account="<alias>" |
| Windows | Credential Manager: target="go-goog-cli:<alias>" |
| Linux | Secret Service: label="go-goog-cli:<alias>" |
| Fallback | `~/.config/goog/tokens/<alias>.enc` |

**Stored per account:**
- OAuth2 refresh token (encrypted)
- Access token (cached, with expiry)
- Granted scopes
- Account email (for display/verification)

### 3.5 Google API Scopes

#### Gmail Scopes (Full Coverage)
| Scope | Level | Use Case |
|-------|-------|----------|
| `gmail.readonly` | Restricted | Read-only mail access |
| `gmail.send` | Sensitive | Send emails |
| `gmail.compose` | Sensitive | Create/modify drafts |
| `gmail.modify` | Restricted | Read/write (excludes permanent delete) |
| `gmail.labels` | Non-sensitive | Label management only |
| `gmail.settings.basic` | Sensitive | Manage basic mail settings |
| `gmail.settings.sharing` | Sensitive | Manage send-as, forwarding |
| `mail.google.com` | Restricted | Full mailbox access (includes delete) |

#### Calendar Scopes (Full Coverage)
| Scope | Use Case |
|-------|----------|
| `calendar.readonly` | Read-only calendar access |
| `calendar` | Full calendar access |
| `calendar.events` | View and edit events only |
| `calendar.events.readonly` | Read events only |
| `calendar.settings.readonly` | Read calendar settings |
| `calendar.freebusy` | Check availability only |
| `calendar.acls` | Manage calendar sharing |

**Scope Strategy:** Request minimal scopes by default. Prompt for additional scopes when operations require them. Use incremental authorization to request additional scopes as needed.

---

## 4. Configuration

### 4.1 Config File Location
| Platform | Path |
|----------|------|
| macOS | `~/Library/Application Support/goog/config.yaml` |
| Linux | `~/.config/goog/config.yaml` |
| Windows | `%APPDATA%\goog\config.yaml` |

### 4.2 Config Schema
```yaml
default_account: personal
default_format: table
timezone: America/New_York

accounts:
  personal:
    email: user@gmail.com
  work:
    email: user@company.com

mail:
  default_label: INBOX
  page_size: 25

calendar:
  default_calendar: primary
  week_start: sunday
```

### 4.3 Environment Variables
| Variable | Description |
|----------|-------------|
| `GOOG_ACCOUNT` | Override default account |
| `GOOG_FORMAT` | Override output format |
| `GOOG_CONFIG` | Custom config file path |
| `GOOG_CLIENT_ID` | OAuth client ID |
| `GOOG_CLIENT_SECRET` | OAuth client secret |

---

## 5. Agent Integration

### 5.1 Design Principles
- All commands support `--format json` for machine parsing
- Exit codes indicate success (0) or failure (non-zero)
- Error messages written to stderr, data to stdout
- No interactive prompts in non-TTY mode
- Idempotent operations where possible

### 5.2 Example Agent Workflows

**Check for important emails:**
```bash
goog mail search "is:unread is:important" --format json --limit 10
```

**Schedule a meeting:**
```bash
goog cal create --title "Team Sync" \
    --start "2024-01-15T10:00:00" \
    --duration "30m" \
    --attendees "team@company.com" \
    --format json
```

**Send an email:**
```bash
goog mail send --to "user@example.com" \
    --subject "Meeting Notes" \
    --body "Please find attached..." \
    --format json
```

---

## 6. Security Requirements

### 6.1 Authentication
- OAuth2 Authorization Code flow with PKCE (required)
- No storage of user passwords
- Refresh tokens stored in platform keyring
- Access tokens kept in memory only

### 6.2 Token Security
- PKCE code verifier: 256-bit minimum entropy
- State parameter for CSRF protection
- Localhost-only callback server
- Automatic token refresh before expiry

### 6.3 Data Handling
- No logging of sensitive data (tokens, email content)
- Secure memory handling for credentials
- Config file permissions: 600 (owner read/write only)

---

## 7. Non-Functional Requirements

### 7.1 Platform Support
| Platform | Version |
|----------|---------|
| macOS | 10.15+ (Catalina) |
| Linux | Kernel 4.4+ with D-Bus |
| Windows | 10/11 |

### 7.2 Performance
- Command startup: < 200ms
- API response display: < 500ms (network dependent)
- Token refresh: Background, non-blocking

### 7.3 Distribution
- Single binary (no runtime dependencies)
- Homebrew formula (macOS/Linux)
- Scoop/Chocolatey (Windows)
- GitHub releases (all platforms)

---

## 8. Development Phases

### Phase 1: Foundation
- [ ] Project structure (Clean Architecture)
- [ ] OAuth2/PKCE authentication flow
- [ ] Multi-account token storage in keyring
- [ ] Account management commands (add, remove, list, switch)
- [ ] Basic config management with account context
- [ ] Root command with global flags (including --account)

### Phase 2: Gmail Core
- [ ] Messages: list, read, search, send, reply
- [ ] Messages: trash, untrash, delete, modify, archive
- [ ] Drafts: full CRUD (list, show, create, update, send, delete)
- [ ] Labels: full CRUD (list, show, create, update, delete)
- [ ] Attachments: download
- [ ] Output formatters (JSON, table, plain)

### Phase 3: Gmail Advanced
- [ ] Threads: full CRUD (list, show, modify, trash, untrash, delete)
- [ ] Settings: vacation responder
- [ ] Settings: filters (list, create, delete)
- [ ] Settings: forwarding, send-as
- [ ] Batch operations

### Phase 4: Calendar Core
- [ ] Events: list, show, create, update, delete
- [ ] Events: quickAdd, move, instances
- [ ] Calendar management: list, create, update, delete
- [ ] Today/week quick views

### Phase 5: Calendar Advanced
- [ ] Recurring events with instances
- [ ] Attendee management and RSVP
- [ ] ACL/sharing management (share, unshare, acl list/set)
- [ ] Freebusy queries

### Phase 6: Polish
- [ ] Windows/Linux keyring backends (macOS in Phase 1)
- [ ] Shell completions (bash, zsh, fish)
- [ ] Batch operations for both APIs
- [ ] Interactive mode (optional)

---

## 9. References

### Google APIs
- [Gmail API Go Package](https://pkg.go.dev/google.golang.org/api/gmail/v1)
- [Calendar API Go Package](https://pkg.go.dev/google.golang.org/api/calendar/v3)
- [Gmail API Scopes](https://developers.google.com/workspace/gmail/api/auth/scopes)
- [Calendar API Scopes](https://developers.google.com/workspace/calendar/api/auth)
- [OAuth2 Scopes Reference](https://developers.google.com/identity/protocols/oauth2/scopes)

### Go Packages
- [spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [spf13/viper](https://github.com/spf13/viper) - Configuration
- [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2) - OAuth2 with PKCE
- [int128/oauth2cli](https://pkg.go.dev/github.com/int128/oauth2cli) - CLI OAuth flow
- [99designs/keyring](https://github.com/99designs/keyring) - Credential storage

### Reference Implementations
- [gogcli](https://github.com/steipete/gogcli) - Comprehensive Google Suite CLI
- [golang-calendar-cli](https://github.com/Sup3r-Us3r/golang-calendar-cli) - Calendar CLI example
