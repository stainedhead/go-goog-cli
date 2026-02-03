# Development Plan: go-goog-cli

This plan enables parallel development with multiple subagents. Tasks are organized by dependency level, allowing independent streams to execute concurrently.

## Dependency Graph Overview

```
Phase 1: Foundation (Sequential - Critical Path)
    ├── 1A: Infrastructure Layer (can parallelize internally)
    │   ├── Config Management
    │   ├── Keyring Storage
    │   └── OAuth2/PKCE Flow
    └── 1B: Domain Entities (parallel)
        ├── Account Domain
        ├── Mail Domain
        └── Calendar Domain

Phase 2: Core Services (parallel streams after Phase 1)
    ├── Stream A: Account & Auth CLI
    ├── Stream B: Output Presenters
    ├── Stream C: Gmail Repository
    └── Stream D: Calendar Repository

Phase 3: Gmail Features (parallel after Stream C)
    ├── Stream E: Messages
    ├── Stream F: Drafts
    ├── Stream G: Labels
    └── Stream H: Attachments

Phase 4: Calendar Features (parallel after Stream D)
    ├── Stream I: Events
    ├── Stream J: Calendar Management
    └── Stream K: Quick Views

Phase 5: Advanced Features (parallel)
    ├── Stream L: Gmail Threads & Settings
    ├── Stream M: Calendar ACL & RSVP
    └── Stream N: Recurring Events

Phase 6: Polish (parallel)
    ├── Stream O: Shell Completions
    ├── Stream P: Cross-Platform Testing
    └── Stream Q: Documentation
```

---

## Phase 1: Foundation

**Goal:** Establish core infrastructure that all other components depend on.

### 1A: Infrastructure Layer

#### Task 1A.1: Configuration Management
**File:** `internal/infrastructure/config/config.go`
**Dependencies:** None
**Parallelizable:** Yes (independent)

```
Implement:
- [ ] Config struct with YAML tags
- [ ] Platform-specific config paths (macOS, Linux, Windows)
- [ ] Load/Save functions
- [ ] Environment variable overrides (GOOG_*)
- [ ] Default values
- [ ] Config file permission validation (600)

Tests:
- [ ] TestConfigLoad
- [ ] TestConfigSave
- [ ] TestConfigPlatformPaths
- [ ] TestConfigEnvOverrides
- [ ] TestConfigDefaults
```

#### Task 1A.2: Keyring Storage
**File:** `internal/infrastructure/keyring/store.go`
**Dependencies:** None
**Parallelizable:** Yes (independent)

```
Implement:
- [ ] KeyringStore interface
- [ ] Per-account token storage (namespaced: goog:<alias>:*)
- [ ] macOS Keychain backend (primary)
- [ ] Encrypted file fallback
- [ ] Store/Retrieve/Delete operations

Tests:
- [ ] TestKeyringStore
- [ ] TestKeyringRetrieve
- [ ] TestKeyringDelete
- [ ] TestKeyringNamespacing
- [ ] TestKeyringFallback
```

#### Task 1A.3: OAuth2/PKCE Authentication
**Files:** `internal/infrastructure/auth/oauth.go`, `internal/infrastructure/auth/token.go`
**Dependencies:** 1A.1 (Config), 1A.2 (Keyring)
**Parallelizable:** After dependencies

```
Implement:
- [ ] PKCE code verifier generation (256-bit entropy)
- [ ] Authorization URL builder with scopes
- [ ] Local callback server (localhost only)
- [ ] Token exchange (code + verifier)
- [ ] Token refresh logic
- [ ] Token caching with expiry
- [ ] Scope management

Tests:
- [ ] TestPKCEVerifierGeneration
- [ ] TestAuthorizationURL
- [ ] TestCallbackServer
- [ ] TestTokenExchange
- [ ] TestTokenRefresh
- [ ] TestScopeValidation
```

### 1B: Domain Entities (Parallel)

#### Task 1B.1: Account Domain
**Files:** `internal/domain/account/account.go`, `internal/domain/account/repository.go`
**Dependencies:** None
**Parallelizable:** Yes

```
Implement:
- [ ] Account entity (alias, email, scopes, timestamps)
- [ ] AccountRepository interface
- [ ] Validation methods

Tests:
- [ ] TestAccountEntity
- [ ] TestAccountValidation
```

#### Task 1B.2: Mail Domain
**Files:** `internal/domain/mail/*.go`
**Dependencies:** None
**Parallelizable:** Yes

```
Implement:
- [ ] Message entity (id, threadId, from, to, subject, body, labels, etc.)
- [ ] Draft entity
- [ ] Thread entity
- [ ] Label entity
- [ ] Attachment entity
- [ ] Filter entity
- [ ] MessageRepository interface
- [ ] DraftRepository interface
- [ ] ThreadRepository interface
- [ ] LabelRepository interface
- [ ] SettingsRepository interface

Tests:
- [ ] TestMessageEntity
- [ ] TestDraftEntity
- [ ] TestThreadEntity
- [ ] TestLabelEntity
- [ ] TestAttachmentEntity
```

#### Task 1B.3: Calendar Domain
**Files:** `internal/domain/calendar/*.go`
**Dependencies:** None
**Parallelizable:** Yes

```
Implement:
- [ ] Event entity (id, title, start, end, attendees, recurrence, etc.)
- [ ] Calendar entity
- [ ] ACL entity
- [ ] Attendee entity
- [ ] EventRepository interface
- [ ] CalendarRepository interface
- [ ] ACLRepository interface

Tests:
- [ ] TestEventEntity
- [ ] TestCalendarEntity
- [ ] TestACLEntity
- [ ] TestAttendeeEntity
```

---

## Phase 2: Core Services

**Goal:** Build repositories, presenters, and account management.
**Precondition:** Phase 1 complete

### Stream A: Account & Auth CLI

#### Task 2A.1: Account Use Cases
**File:** `internal/usecase/account/account.go`
**Dependencies:** 1B.1, 1A.1, 1A.2, 1A.3

```
Implement:
- [ ] AddAccount use case
- [ ] RemoveAccount use case
- [ ] ListAccounts use case
- [ ] SwitchAccount use case
- [ ] ShowAccount use case
- [ ] RenameAccount use case
- [ ] Account context resolution (flag > env > config > first)

Tests:
- [ ] TestAddAccount
- [ ] TestRemoveAccount
- [ ] TestListAccounts
- [ ] TestSwitchAccount
- [ ] TestAccountResolution
```

#### Task 2A.2: Auth CLI Commands
**File:** `internal/adapter/cli/auth.go`
**Dependencies:** 2A.1

```
Implement:
- [ ] goog auth login
- [ ] goog auth logout
- [ ] goog auth status
- [ ] goog auth refresh

Tests:
- [ ] TestAuthLoginCommand
- [ ] TestAuthLogoutCommand
- [ ] TestAuthStatusCommand
```

#### Task 2A.3: Account CLI Commands
**File:** `internal/adapter/cli/account.go`
**Dependencies:** 2A.1

```
Implement:
- [ ] goog account list
- [ ] goog account add [alias]
- [ ] goog account remove <alias>
- [ ] goog account switch <alias>
- [ ] goog account show
- [ ] goog account rename <old> <new>

Tests:
- [ ] TestAccountListCommand
- [ ] TestAccountAddCommand
- [ ] TestAccountRemoveCommand
- [ ] TestAccountSwitchCommand
```

#### Task 2A.4: Config CLI Commands
**File:** `internal/adapter/cli/config.go`
**Dependencies:** 1A.1

```
Implement:
- [ ] goog config show
- [ ] goog config set KEY VALUE
- [ ] goog config get KEY

Tests:
- [ ] TestConfigShowCommand
- [ ] TestConfigSetCommand
- [ ] TestConfigGetCommand
```

### Stream B: Output Presenters (Parallel)

#### Task 2B.1: JSON Presenter
**File:** `internal/adapter/presenter/json.go`
**Dependencies:** 1B.2, 1B.3 (domain entities)

```
Implement:
- [ ] Presenter interface
- [ ] JSON formatting for all entities
- [ ] Pretty print option
- [ ] Error response formatting

Tests:
- [ ] TestJSONPresenterMessage
- [ ] TestJSONPresenterEvent
- [ ] TestJSONPresenterError
```

#### Task 2B.2: Table Presenter
**File:** `internal/adapter/presenter/table.go`
**Dependencies:** 1B.2, 1B.3

```
Implement:
- [ ] Table formatting for lists
- [ ] Column width handling
- [ ] Truncation for long fields
- [ ] Color support (optional)

Tests:
- [ ] TestTablePresenterMessages
- [ ] TestTablePresenterEvents
- [ ] TestTablePresenterTruncation
```

#### Task 2B.3: Plain Presenter
**File:** `internal/adapter/presenter/plain.go`
**Dependencies:** 1B.2, 1B.3

```
Implement:
- [ ] Plain text formatting
- [ ] Single item detail view
- [ ] List view

Tests:
- [ ] TestPlainPresenterMessage
- [ ] TestPlainPresenterEvent
```

### Stream C: Gmail Repository (Parallel)

#### Task 2C.1: Gmail API Client Setup
**File:** `internal/adapter/repository/gmail.go`
**Dependencies:** 1A.3, 1B.2

```
Implement:
- [ ] Gmail service initialization
- [ ] API client with auth token
- [ ] Rate limiting / retry logic
- [ ] Error mapping to domain errors

Tests:
- [ ] TestGmailClientInit
- [ ] TestGmailErrorMapping
```

### Stream D: Calendar Repository (Parallel)

#### Task 2D.1: Calendar API Client Setup
**File:** `internal/adapter/repository/gcalendar.go`
**Dependencies:** 1A.3, 1B.3

```
Implement:
- [ ] Calendar service initialization
- [ ] API client with auth token
- [ ] Rate limiting / retry logic
- [ ] Error mapping to domain errors

Tests:
- [ ] TestCalendarClientInit
- [ ] TestCalendarErrorMapping
```

---

## Phase 3: Gmail Features

**Goal:** Implement all Gmail commands.
**Precondition:** Phase 2 Stream A, B, C complete

### Stream E: Messages (Can split further)

#### Task 3E.1: Message Use Cases - Read Operations
**File:** `internal/usecase/mail/message.go`
**Dependencies:** 2C.1

```
Implement:
- [ ] ListMessages use case
- [ ] ReadMessage use case
- [ ] SearchMessages use case

Tests:
- [ ] TestListMessages
- [ ] TestReadMessage
- [ ] TestSearchMessages
```

#### Task 3E.2: Message Use Cases - Write Operations
**File:** `internal/usecase/mail/message.go`
**Dependencies:** 2C.1

```
Implement:
- [ ] SendMessage use case
- [ ] ReplyMessage use case
- [ ] ForwardMessage use case
- [ ] TrashMessage use case
- [ ] UntrashMessage use case
- [ ] DeleteMessage use case
- [ ] ArchiveMessage use case
- [ ] ModifyMessage use case
- [ ] MarkMessage use case
- [ ] MoveMessage use case

Tests:
- [ ] TestSendMessage
- [ ] TestReplyMessage
- [ ] TestTrashMessage
- [ ] TestDeleteMessage
```

#### Task 3E.3: Message CLI Commands
**File:** `internal/adapter/cli/mail.go`
**Dependencies:** 3E.1, 3E.2, 2B.*

```
Implement:
- [ ] goog mail list (with all flags)
- [ ] goog mail read <id>
- [ ] goog mail send
- [ ] goog mail reply <id>
- [ ] goog mail forward <id>
- [ ] goog mail search <query>
- [ ] goog mail trash <id>
- [ ] goog mail untrash <id>
- [ ] goog mail delete <id>
- [ ] goog mail archive <id>
- [ ] goog mail modify <id>
- [ ] goog mail mark <id>
- [ ] goog mail move <id>

Tests:
- [ ] TestMailListCommand
- [ ] TestMailReadCommand
- [ ] TestMailSendCommand
- [ ] TestMailSearchCommand
```

### Stream F: Drafts (Parallel with E)

#### Task 3F.1: Draft Use Cases
**File:** `internal/usecase/mail/draft.go`
**Dependencies:** 2C.1

```
Implement:
- [ ] ListDrafts use case
- [ ] ShowDraft use case
- [ ] CreateDraft use case
- [ ] UpdateDraft use case
- [ ] SendDraft use case
- [ ] DeleteDraft use case

Tests:
- [ ] TestListDrafts
- [ ] TestCreateDraft
- [ ] TestSendDraft
```

#### Task 3F.2: Draft CLI Commands
**File:** `internal/adapter/cli/draft.go`
**Dependencies:** 3F.1

```
Implement:
- [ ] goog draft list
- [ ] goog draft show <id>
- [ ] goog draft create
- [ ] goog draft update <id>
- [ ] goog draft send <id>
- [ ] goog draft delete <id>

Tests:
- [ ] TestDraftListCommand
- [ ] TestDraftCreateCommand
```

### Stream G: Labels (Parallel with E, F)

#### Task 3G.1: Label Use Cases
**File:** `internal/usecase/mail/label.go`
**Dependencies:** 2C.1

```
Implement:
- [ ] ListLabels use case
- [ ] ShowLabel use case
- [ ] CreateLabel use case
- [ ] UpdateLabel use case
- [ ] DeleteLabel use case

Tests:
- [ ] TestListLabels
- [ ] TestCreateLabel
- [ ] TestDeleteLabel
```

#### Task 3G.2: Label CLI Commands
**File:** `internal/adapter/cli/label.go`
**Dependencies:** 3G.1

```
Implement:
- [ ] goog label list
- [ ] goog label show <name>
- [ ] goog label create <name>
- [ ] goog label update <name>
- [ ] goog label delete <name>

Tests:
- [ ] TestLabelListCommand
- [ ] TestLabelCreateCommand
```

### Stream H: Attachments (Parallel with E, F, G)

#### Task 3H.1: Attachment Use Cases
**File:** `internal/usecase/mail/message.go` (extend)
**Dependencies:** 2C.1

```
Implement:
- [ ] DownloadAttachment use case
- [ ] List attachments in message

Tests:
- [ ] TestDownloadAttachment
```

#### Task 3H.2: Attachment CLI
**File:** `internal/adapter/cli/mail.go` (extend)
**Dependencies:** 3H.1

```
Implement:
- [ ] goog mail attachment <id> [--output PATH]

Tests:
- [ ] TestAttachmentDownloadCommand
```

---

## Phase 4: Calendar Features

**Goal:** Implement core Calendar commands.
**Precondition:** Phase 2 Stream A, B, D complete

### Stream I: Events

#### Task 4I.1: Event Use Cases - Read Operations
**File:** `internal/usecase/calendar/event.go`
**Dependencies:** 2D.1

```
Implement:
- [ ] ListEvents use case
- [ ] ShowEvent use case
- [ ] TodayEvents use case
- [ ] WeekEvents use case
- [ ] FreeBusy use case

Tests:
- [ ] TestListEvents
- [ ] TestShowEvent
- [ ] TestTodayEvents
- [ ] TestFreeBusy
```

#### Task 4I.2: Event Use Cases - Write Operations
**File:** `internal/usecase/calendar/event.go`
**Dependencies:** 2D.1

```
Implement:
- [ ] CreateEvent use case
- [ ] QuickAddEvent use case
- [ ] UpdateEvent use case
- [ ] DeleteEvent use case
- [ ] MoveEvent use case

Tests:
- [ ] TestCreateEvent
- [ ] TestQuickAddEvent
- [ ] TestUpdateEvent
- [ ] TestDeleteEvent
```

#### Task 4I.3: Event CLI Commands
**File:** `internal/adapter/cli/calendar.go`
**Dependencies:** 4I.1, 4I.2

```
Implement:
- [ ] goog cal list (with all flags)
- [ ] goog cal show <id>
- [ ] goog cal create
- [ ] goog cal quick "<text>"
- [ ] goog cal update <id>
- [ ] goog cal delete <id>
- [ ] goog cal move <id>
- [ ] goog cal today
- [ ] goog cal week
- [ ] goog cal freebusy

Tests:
- [ ] TestCalListCommand
- [ ] TestCalCreateCommand
- [ ] TestCalTodayCommand
```

### Stream J: Calendar Management (Parallel with I)

#### Task 4J.1: Calendar Management Use Cases
**File:** `internal/usecase/calendar/calendar.go`
**Dependencies:** 2D.1

```
Implement:
- [ ] ListCalendars use case
- [ ] ShowCalendar use case
- [ ] CreateCalendar use case
- [ ] UpdateCalendar use case
- [ ] DeleteCalendar use case
- [ ] ClearCalendar use case
- [ ] SubscribeCalendar use case
- [ ] UnsubscribeCalendar use case

Tests:
- [ ] TestListCalendars
- [ ] TestCreateCalendar
- [ ] TestDeleteCalendar
```

#### Task 4J.2: Calendar Management CLI
**File:** `internal/adapter/cli/calendar_mgmt.go`
**Dependencies:** 4J.1

```
Implement:
- [ ] goog calendar list
- [ ] goog calendar show <id>
- [ ] goog calendar create
- [ ] goog calendar update <id>
- [ ] goog calendar delete <id>
- [ ] goog calendar clear <id>
- [ ] goog calendar subscribe <id>
- [ ] goog calendar unsubscribe <id>

Tests:
- [ ] TestCalendarListCommand
- [ ] TestCalendarCreateCommand
```

---

## Phase 5: Advanced Features

**Goal:** Implement advanced Gmail and Calendar features.
**Precondition:** Phases 3 and 4 complete

### Stream L: Gmail Threads & Settings

#### Task 5L.1: Thread Use Cases
**File:** `internal/usecase/mail/thread.go`
**Dependencies:** 2C.1

```
Implement:
- [ ] ListThreads use case
- [ ] ShowThread use case
- [ ] ModifyThread use case
- [ ] TrashThread use case
- [ ] UntrashThread use case
- [ ] DeleteThread use case

Tests:
- [ ] TestListThreads
- [ ] TestShowThread
```

#### Task 5L.2: Thread CLI Commands
**File:** `internal/adapter/cli/thread.go`
**Dependencies:** 5L.1

```
Implement:
- [ ] goog thread list
- [ ] goog thread show <id>
- [ ] goog thread modify <id>
- [ ] goog thread trash <id>
- [ ] goog thread untrash <id>
- [ ] goog thread delete <id>

Tests:
- [ ] TestThreadListCommand
```

#### Task 5L.3: Settings Use Cases
**File:** `internal/usecase/mail/settings.go`
**Dependencies:** 2C.1

```
Implement:
- [ ] ShowSettings use case
- [ ] GetVacation use case
- [ ] SetVacation use case
- [ ] ListFilters use case
- [ ] CreateFilter use case
- [ ] DeleteFilter use case
- [ ] ListSendAs use case
- [ ] ManageForwarding use case

Tests:
- [ ] TestGetVacation
- [ ] TestListFilters
```

#### Task 5L.4: Settings CLI Commands
**File:** `internal/adapter/cli/settings.go`
**Dependencies:** 5L.3

```
Implement:
- [ ] goog settings show
- [ ] goog settings vacation
- [ ] goog settings filters list
- [ ] goog settings filters create
- [ ] goog settings filters delete <id>
- [ ] goog settings forwarding
- [ ] goog settings send-as list
- [ ] goog settings send-as create

Tests:
- [ ] TestSettingsShowCommand
- [ ] TestSettingsVacationCommand
```

### Stream M: Calendar ACL & RSVP (Parallel with L)

#### Task 5M.1: ACL Use Cases
**File:** `internal/usecase/calendar/acl.go`
**Dependencies:** 2D.1

```
Implement:
- [ ] ListACL use case
- [ ] SetACL use case
- [ ] ShareCalendar use case
- [ ] UnshareCalendar use case

Tests:
- [ ] TestListACL
- [ ] TestShareCalendar
```

#### Task 5M.2: RSVP Use Cases
**File:** `internal/usecase/calendar/event.go` (extend)
**Dependencies:** 2D.1

```
Implement:
- [ ] RSVPEvent use case (accept/decline/tentative)

Tests:
- [ ] TestRSVPEvent
```

#### Task 5M.3: ACL CLI Commands
**File:** `internal/adapter/cli/acl.go`
**Dependencies:** 5M.1, 5M.2

```
Implement:
- [ ] goog cal share <calendar>
- [ ] goog cal unshare <calendar>
- [ ] goog cal acl list <calendar>
- [ ] goog cal acl set <calendar>
- [ ] goog cal rsvp <id>

Tests:
- [ ] TestCalShareCommand
- [ ] TestCalRSVPCommand
```

### Stream N: Recurring Events (Parallel with L, M)

#### Task 5N.1: Recurring Event Use Cases
**File:** `internal/usecase/calendar/event.go` (extend)
**Dependencies:** 2D.1

```
Implement:
- [ ] ListEventInstances use case
- [ ] Recurrence rule parsing
- [ ] Exception handling

Tests:
- [ ] TestListEventInstances
- [ ] TestRecurrenceParsing
```

#### Task 5N.2: Recurring Event CLI
**File:** `internal/adapter/cli/calendar.go` (extend)
**Dependencies:** 5N.1

```
Implement:
- [ ] goog cal instances <id>
- [ ] --recurrence flag for create/update

Tests:
- [ ] TestCalInstancesCommand
```

---

## Phase 6: Polish

**Goal:** Cross-platform support, completions, documentation.
**Precondition:** Phases 1-5 complete

### Stream O: Shell Completions

#### Task 6O.1: Shell Completion Generation
**File:** `internal/adapter/cli/completion.go`
**Dependencies:** All CLI commands

```
Implement:
- [ ] Bash completion
- [ ] Zsh completion
- [ ] Fish completion
- [ ] PowerShell completion

Tests:
- [ ] TestBashCompletion
- [ ] TestZshCompletion
```

### Stream P: Cross-Platform (Parallel with O)

#### Task 6P.1: Windows Keyring Backend
**File:** `internal/infrastructure/keyring/windows.go`
**Dependencies:** 1A.2

```
Implement:
- [ ] Windows Credential Manager integration

Tests:
- [ ] TestWindowsKeyring (CI with Windows runner)
```

#### Task 6P.2: Linux Keyring Backend
**File:** `internal/infrastructure/keyring/linux.go`
**Dependencies:** 1A.2

```
Implement:
- [ ] Secret Service (D-Bus) integration
- [ ] KWallet fallback

Tests:
- [ ] TestLinuxKeyring (CI with Linux runner)
```

### Stream Q: Documentation (Parallel with O, P)

#### Task 6Q.1: User Documentation
```
Update:
- [ ] README.md with full usage examples
- [ ] Man page generation
- [ ] --help text review for all commands
```

#### Task 6Q.2: API Documentation
```
Update:
- [ ] documentation/technical-details.md
- [ ] GoDoc comments on all exported types
```

---

## Parallel Execution Matrix

This matrix shows which tasks can run simultaneously:

| Phase | Parallel Streams | Max Concurrent Agents |
|-------|------------------|----------------------|
| 1A    | Config, Keyring (then OAuth) | 2 → 1 |
| 1B    | Account, Mail, Calendar domains | 3 |
| 2     | Auth CLI, Presenters, Gmail Repo, Calendar Repo | 4 |
| 3     | Messages, Drafts, Labels, Attachments | 4 |
| 4     | Events, Calendar Mgmt | 2 |
| 5     | Threads/Settings, ACL/RSVP, Recurring | 3 |
| 6     | Completions, Windows, Linux, Docs | 4 |

**Maximum theoretical parallelism:** 4 agents working simultaneously

---

## Task Assignment for Subagents

When spawning subagents, assign by stream:

```
Agent 1 (Foundation):     Phase 1A (sequential), then Stream A
Agent 2 (Mail Domain):    Task 1B.2, then Stream E
Agent 3 (Calendar Domain): Task 1B.3, then Stream I
Agent 4 (Presenters):     Stream B, then Stream F or G
```

After Phase 2:
```
Agent 1: Stream E (Messages)
Agent 2: Stream F (Drafts)
Agent 3: Stream G (Labels) + H (Attachments)
Agent 4: Stream I (Events) + J (Calendar Mgmt)
```

---

## Quality Gates (Per Task)

Each task must pass before marking complete:

```bash
# Run for every task completion
go fmt ./...
go mod tidy
go vet ./...
golangci-lint run
go test ./...
go build -o bin/goog ./cmd/goog
./bin/goog --help
```

---

## Integration Checkpoints

After each phase, run integration verification:

### After Phase 1
```bash
./bin/goog config show
# Should display default config
```

### After Phase 2
```bash
./bin/goog auth login
./bin/goog account list
./bin/goog auth status
```

### After Phase 3
```bash
./bin/goog mail list --format json
./bin/goog mail search "is:unread" --limit 5
./bin/goog label list
./bin/goog draft list
```

### After Phase 4
```bash
./bin/goog cal today
./bin/goog cal week
./bin/goog calendar list
./bin/goog cal create --title "Test" --start "tomorrow 10am" --duration "1h"
```

### After Phase 5
```bash
./bin/goog thread list
./bin/goog settings show
./bin/goog cal acl list primary
```

### After Phase 6
```bash
./bin/goog completion bash > /tmp/goog.bash
./bin/goog --help  # Verify all commands documented
```
