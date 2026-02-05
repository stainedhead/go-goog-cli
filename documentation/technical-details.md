# Technical Details

## Architecture

This application follows Clean Architecture with four layers:

```
internal/
├── domain/           # Business entities and rules (innermost)
├── usecase/          # Application-specific business logic
├── adapter/          # Interface adapters (CLI, repositories)
└── infrastructure/   # External concerns (APIs, config, keyring)
```

Dependencies flow inward only. Inner layers define interfaces; outer layers implement them.

### Layer Responsibilities

**Domain** (`internal/domain/`)
- Pure business entities: `Mail`, `Calendar`, `Tasks`, `Contacts`, `Account`
- Business rules and validation
- No external dependencies

**Use Case** (`internal/usecase/`)
- Application orchestration
- Account management service
- OAuth flow coordination

**Adapter** (`internal/adapter/`)
- CLI command handlers (`cli/`)
- Output formatters (`presenter/`)
- Google API implementations (`repository/`)

**Infrastructure** (`internal/infrastructure/`)
- OAuth2/PKCE authentication (`auth/`)
- Configuration management (`config/`)
- Keyring integration (`keyring/`)

## Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.24+ |
| CLI Framework | [spf13/cobra](https://github.com/spf13/cobra) |
| Configuration | [spf13/viper](https://github.com/spf13/viper) |
| Keyring | [99designs/keyring](https://github.com/99designs/keyring) |
| Google APIs | [google.golang.org/api](https://google.golang.org/api) |
| OAuth2 | [golang.org/x/oauth2](https://golang.org/x/oauth2) |

## Data Flow

### Command Execution
```
User Input → Cobra CLI → Command Handler → Use Case → Repository → Google API
                                              ↓
                                          Presenter → Output
```

### Authentication Flow
```
goog auth login
    → Open browser for Google OAuth consent
    → Start local HTTP server for callback
    → Exchange authorization code for tokens
    → Store tokens in system keyring
    → Save account metadata to config
```

## API Coverage

### Gmail API

| Category | Operations |
|----------|------------|
| Messages | list, get, send, reply, forward, trash, untrash, delete, modify, move |
| Drafts | list, get, create, update, send, delete |
| Labels | list, get, create, update, delete |
| Threads | list, get, trash, untrash, delete, modify |

### Calendar API

| Category | Operations |
|----------|------------|
| Events | list, get, create, update, delete, quickAdd, move, instances |
| Calendars | list, get, insert, update, delete, clear |
| ACL | list, get, insert, delete |
| FreeBusy | query |

### Tasks API

| Category | Operations |
|----------|------------|
| TaskLists | list, get, insert, update, delete |
| Tasks | list, get, insert, update, delete, move, clear |

### People API (Contacts)

| Category | Operations |
|----------|------------|
| Contacts | list, get, create, update, delete, search |
| ContactGroups | list, get, create, update, delete, members |
| Memberships | modify (add/remove contacts to/from groups) |

**Architecture Pattern**: Contacts follows the Tasks pattern (simple CRUD, no use case layer).

**API Details**:
- Uses Google People API v1
- Person resource fields: names, emailAddresses, phoneNumbers, addresses, organizations, birthdays, biographies, photos, urls, memberships, metadata
- ContactGroup resource fields: name, groupType, memberCount, memberResourceNames, metadata
- Resource names: `people/c<id>` for contacts, `contactGroups/<id>` for groups

**Data Transformations**:
- API to Domain: `PeopleRepository` transforms `people.Person` to `contacts.Contact`
- Domain to API: Transform functions convert domain entities to People API structures
- Handles nested structures (names, emails, phones, addresses, etc.)

## Configuration

Config file location: `~/.config/goog/config.yaml`

```yaml
default_account: personal
accounts:
  personal:
    email: user@gmail.com
    scopes:
      - https://www.googleapis.com/auth/gmail.modify
      - https://www.googleapis.com/auth/calendar
    added: 2024-01-15T10:00:00Z
  work:
    email: user@company.com
    scopes:
      - https://www.googleapis.com/auth/gmail.readonly
    added: 2024-01-16T14:30:00Z
```

## Credential Storage

OAuth tokens stored securely in system keyring:

| Platform | Backend |
|----------|---------|
| macOS | Keychain |
| Windows | Credential Manager |
| Linux | Secret Service (GNOME Keyring, KWallet) |

Keyring entries per account:
- `oauth_token` - Serialized OAuth2 token (access, refresh, expiry)
- `oauth_scopes` - Granted scopes list

## OAuth Scopes

### Gmail Scopes

| Scope | Description |
|-------|-------------|
| `gmail.readonly` | Read-only mail access |
| `gmail.send` | Send emails |
| `gmail.compose` | Create/modify drafts |
| `gmail.modify` | Read/write (excludes delete) |
| `gmail.labels` | Label management |

### Calendar Scopes

| Scope | Description |
|-------|-------------|
| `calendar.readonly` | Read-only calendar |
| `calendar` | Full calendar access |
| `calendar.events` | Events only |

### Tasks Scopes

| Scope | Description |
|-------|-------------|
| `tasks.readonly` | Read-only tasks access |
| `tasks` | Full tasks access |

### Contacts Scopes

| Scope | Description |
|-------|-------------|
| `contacts.readonly` | Read-only contacts access |
| `contacts.other.readonly` | Read access to other contacts data |
| `contacts` | Full contacts access (read/write) |

## Error Handling

Errors are wrapped with context using `fmt.Errorf("context: %w", err)`:

```go
if err := repo.Send(ctx, msg); err != nil {
    return fmt.Errorf("failed to send message: %w", err)
}
```

User-facing errors include actionable guidance:
```
Error: no account found (run 'goog auth login' to authenticate)
```

## Testing

Tests follow Test-Driven Development (TDD) with comprehensive coverage:

```bash
go test ./...                      # All tests
go test -cover ./...               # With coverage
go test -v ./internal/adapter/cli  # Specific package
```

### Test Coverage

| Package | Coverage | Tests | Status |
|---------|----------|-------|--------|
| **domain/account** | 100% | 15 | ✅ Perfect |
| **domain/calendar** | 100% | 45 | ✅ Perfect |
| **domain/mail** | 100% | 38 | ✅ Perfect |
| **domain/tasks** | 100% | 14 | ✅ Perfect |
| **domain/contacts** | 94.5% | 12 | ✅ Excellent |
| **infrastructure/auth** | 93.3% | 25 | ✅ Excellent |
| **infrastructure/keyring** | 91.3% | 18 | ✅ Excellent |
| **infrastructure/config** | 80.1% | 12 | ✅ Good |
| **adapter/presenter** | 93.7% | 42 | ✅ Excellent |
| **adapter/repository** | 84.1% | 73 | ✅ Excellent |
| **adapter/cli** | 83.5% | 174+ | ✅ Very Good |
| **usecase/account** | 90.6% | 20 | ✅ Excellent |

**Overall Project:** 80.6% coverage | 436+ total tests

### Test Infrastructure

**Dependency Injection Framework:**
- All CLI commands use DI for full testability
- Mock implementations for repositories and services
- Enables isolated unit testing without real API calls

**HTTP Test Server:**
- Comprehensive mock Google API server
- Handles all Gmail and Calendar API endpoints
- Returns realistic test data for integration tests

**Test Organization:**
- Unit tests in same package (`*_test.go`)
- Table-driven tests for command variations
- Execution tests with mocked dependencies
- Edge case and error path testing
- Helper function tests at 100% coverage

## Project Structure

```
.
├── cmd/goog/
│   └── main.go                    # Entry point
├── internal/
│   ├── domain/
│   │   ├── account/               # Account entity
│   │   ├── mail/                  # Message, Draft, Label, Thread
│   │   ├── calendar/              # Event, Calendar, ACL, FreeBusy
│   │   ├── tasks/                 # Task, TaskList
│   │   └── contacts/              # Contact, ContactGroup
│   ├── usecase/
│   │   └── account/               # Account service, OAuth flow
│   ├── adapter/
│   │   ├── cli/                   # Command handlers
│   │   ├── presenter/             # JSON, Table, Plain formatters
│   │   └── repository/            # Gmail, Calendar, Tasks, People repositories
│   └── infrastructure/
│       ├── auth/                  # OAuth2/PKCE, token management
│       ├── config/                # Viper configuration
│       └── keyring/               # Secure credential storage
├── documentation/
│   ├── product-summary.md
│   ├── product-details.md
│   └── technical-details.md
├── AGENTS.md                      # Development guidelines
├── CLAUDE.md                      # AI agent instructions
└── README.md                      # Project documentation
```
