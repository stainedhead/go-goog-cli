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
- Pure business entities: `Mail`, `Calendar`, `Account`
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
| Messages | list, get, send, reply, forward, trash, untrash, delete, modify |
| Drafts | list, get, create, update, send, delete |
| Labels | list, get, create, update, delete |
| Threads | list, get, trash, modify |

### Calendar API

| Category | Operations |
|----------|------------|
| Events | list, get, create, update, delete, quickAdd, move, instances |
| Calendars | list, get, insert, update, delete, clear |
| ACL | list, get, insert, delete |
| FreeBusy | query |

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

Tests follow TDD with the standard Go testing package:

```bash
go test ./...                      # All tests
go test -cover ./...               # With coverage
go test -v ./internal/adapter/cli  # Specific package
```

Test organization:
- Unit tests in same package (`*_test.go`)
- Table-driven tests for command variations
- Mock interfaces for external dependencies

## Project Structure

```
.
├── cmd/goog/
│   └── main.go                    # Entry point
├── internal/
│   ├── domain/
│   │   ├── account/               # Account entity
│   │   ├── mail/                  # Message, Draft, Label, Thread
│   │   └── calendar/              # Event, Calendar, ACL, FreeBusy
│   ├── usecase/
│   │   └── account/               # Account service, OAuth flow
│   ├── adapter/
│   │   ├── cli/                   # Command handlers
│   │   ├── presenter/             # JSON, Table, Plain formatters
│   │   └── repository/            # Gmail, Calendar repositories
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
