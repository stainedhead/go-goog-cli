// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/stainedhead/go-goog-cli/internal/adapter/repository"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/auth"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/config"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/keyring"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
	"golang.org/x/oauth2"
)

// MessageRepository defines operations for managing email messages.
// This interface mirrors mail.MessageRepository for dependency injection.
type MessageRepository interface {
	List(ctx context.Context, opts mail.ListOptions) (*mail.ListResult[*mail.Message], error)
	Get(ctx context.Context, id string) (*mail.Message, error)
	Send(ctx context.Context, msg *mail.Message) (*mail.Message, error)
	Reply(ctx context.Context, messageID string, reply *mail.Message) (*mail.Message, error)
	Forward(ctx context.Context, messageID string, forward *mail.Message) (*mail.Message, error)
	Trash(ctx context.Context, id string) error
	Untrash(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	Archive(ctx context.Context, id string) error
	Modify(ctx context.Context, id string, req mail.ModifyRequest) (*mail.Message, error)
	Search(ctx context.Context, query string, opts mail.ListOptions) (*mail.ListResult[*mail.Message], error)
}

// DraftRepository defines operations for managing email drafts.
// This interface mirrors mail.DraftRepository for dependency injection.
type DraftRepository interface {
	List(ctx context.Context, opts mail.ListOptions) (*mail.ListResult[*mail.Draft], error)
	Get(ctx context.Context, id string) (*mail.Draft, error)
	Create(ctx context.Context, draft *mail.Draft) (*mail.Draft, error)
	Update(ctx context.Context, draft *mail.Draft) (*mail.Draft, error)
	Send(ctx context.Context, id string) (*mail.Message, error)
	Delete(ctx context.Context, id string) error
}

// ThreadRepository defines operations for managing email threads.
// This interface mirrors mail.ThreadRepository for dependency injection.
type ThreadRepository interface {
	List(ctx context.Context, opts mail.ListOptions) (*mail.ListResult[*mail.Thread], error)
	Get(ctx context.Context, id string) (*mail.Thread, error)
	Modify(ctx context.Context, id string, req mail.ModifyRequest) (*mail.Thread, error)
	Trash(ctx context.Context, id string) error
	Untrash(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

// LabelRepository defines operations for managing email labels.
// This interface mirrors mail.LabelRepository for dependency injection.
type LabelRepository interface {
	List(ctx context.Context) ([]*mail.Label, error)
	Get(ctx context.Context, id string) (*mail.Label, error)
	GetByName(ctx context.Context, name string) (*mail.Label, error)
	Create(ctx context.Context, label *mail.Label) (*mail.Label, error)
	Update(ctx context.Context, label *mail.Label) (*mail.Label, error)
	Delete(ctx context.Context, id string) error
}

// EventRepository defines operations for managing calendar events.
// This interface mirrors calendar.EventRepository for dependency injection.
type EventRepository interface {
	List(ctx context.Context, calendarID string, timeMin, timeMax time.Time) ([]*calendar.Event, error)
	Get(ctx context.Context, calendarID, eventID string) (*calendar.Event, error)
	Create(ctx context.Context, calendarID string, event *calendar.Event) (*calendar.Event, error)
	Update(ctx context.Context, calendarID string, event *calendar.Event) (*calendar.Event, error)
	Delete(ctx context.Context, calendarID, eventID string) error
	Move(ctx context.Context, sourceCalendarID, eventID, destinationCalendarID string) (*calendar.Event, error)
	QuickAdd(ctx context.Context, calendarID, text string) (*calendar.Event, error)
	Instances(ctx context.Context, calendarID, eventID string, timeMin, timeMax time.Time) ([]*calendar.Event, error)
	RSVP(ctx context.Context, calendarID, eventID, response string) error
}

// CalendarRepository defines operations for managing calendars.
// This interface mirrors calendar.CalendarRepository for dependency injection.
type CalendarRepository interface {
	List(ctx context.Context) ([]*calendar.Calendar, error)
	Get(ctx context.Context, calendarID string) (*calendar.Calendar, error)
	Create(ctx context.Context, cal *calendar.Calendar) (*calendar.Calendar, error)
	Update(ctx context.Context, cal *calendar.Calendar) (*calendar.Calendar, error)
	Delete(ctx context.Context, calendarID string) error
	Clear(ctx context.Context, calendarID string) error
}

// ACLRepository defines operations for managing calendar ACL rules.
// This interface mirrors calendar.ACLRepository for dependency injection.
type ACLRepository interface {
	List(ctx context.Context, calendarID string) ([]*calendar.ACLRule, error)
	Get(ctx context.Context, calendarID, ruleID string) (*calendar.ACLRule, error)
	Insert(ctx context.Context, calendarID string, rule *calendar.ACLRule) (*calendar.ACLRule, error)
	Update(ctx context.Context, calendarID string, rule *calendar.ACLRule) (*calendar.ACLRule, error)
	Delete(ctx context.Context, calendarID, ruleID string) error
}

// FreeBusyRepository defines operations for querying calendar availability.
// This interface mirrors calendar.FreeBusyRepository for dependency injection.
type FreeBusyRepository interface {
	Query(ctx context.Context, request *calendar.FreeBusyRequest) (*calendar.FreeBusyResponse, error)
}

// AccountService defines operations for managing user accounts.
type AccountService interface {
	List() ([]*accountuc.Account, error)
	Add(ctx context.Context, alias string, scopes []string) (*accountuc.Account, error)
	Remove(alias string) error
	Switch(alias string) error
	Rename(oldAlias, newAlias string) error
	ResolveAccount(flagValue string) (*accountuc.Account, error)
	GetTokenManager() TokenManager
}

// TokenManager defines operations for managing OAuth tokens.
type TokenManager interface {
	GetTokenSource(ctx context.Context, alias string) (oauth2.TokenSource, error)
	GetTokenInfo(alias string) (*auth.TokenInfo, error)
	RefreshToken(ctx context.Context, alias string, cfg *oauth2.Config) (*oauth2.Token, error)
	GetGrantedScopes(alias string) ([]string, error)
}

// RepositoryFactory creates repository instances from a token source.
type RepositoryFactory interface {
	// Mail repositories
	NewMessageRepository(ctx context.Context, tokenSource oauth2.TokenSource) (MessageRepository, error)
	NewDraftRepository(ctx context.Context, tokenSource oauth2.TokenSource) (DraftRepository, error)
	NewThreadRepository(ctx context.Context, tokenSource oauth2.TokenSource) (ThreadRepository, error)
	NewLabelRepository(ctx context.Context, tokenSource oauth2.TokenSource) (LabelRepository, error)

	// Calendar repositories
	NewEventRepository(ctx context.Context, tokenSource oauth2.TokenSource) (EventRepository, error)
	NewCalendarRepository(ctx context.Context, tokenSource oauth2.TokenSource) (CalendarRepository, error)
	NewACLRepository(ctx context.Context, tokenSource oauth2.TokenSource) (ACLRepository, error)
	NewFreeBusyRepository(ctx context.Context, tokenSource oauth2.TokenSource) (FreeBusyRepository, error)
}

// Dependencies holds all external dependencies required by CLI commands.
// This enables dependency injection for testing.
type Dependencies struct {
	// AccountService provides account management operations.
	AccountService AccountService

	// RepoFactory creates repository instances.
	RepoFactory RepositoryFactory
}

// Global dependencies instance. Use SetDependencies for testing.
var deps *Dependencies

// SetDependencies sets the global dependencies instance.
// This is primarily used for testing to inject mock implementations.
func SetDependencies(d *Dependencies) {
	deps = d
}

// GetDependencies returns the current dependencies instance.
// If no dependencies have been set, it creates and returns default production dependencies.
func GetDependencies() *Dependencies {
	if deps == nil {
		deps = DefaultDependencies()
	}
	return deps
}

// ResetDependencies clears the global dependencies instance.
// This should be called in test cleanup to ensure test isolation.
func ResetDependencies() {
	deps = nil
}

// DefaultDependencies creates production dependencies.
func DefaultDependencies() *Dependencies {
	return &Dependencies{
		AccountService: &defaultAccountService{},
		RepoFactory:    &defaultRepositoryFactory{},
	}
}

// defaultAccountService implements AccountService using production infrastructure.
type defaultAccountService struct {
	svc *accountuc.Service
}

// ensureService initializes the service if not already done.
func (s *defaultAccountService) ensureService() error {
	if s.svc != nil {
		return nil
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	store, err := keyring.NewStore()
	if err != nil {
		return fmt.Errorf("failed to initialize keyring: %w", err)
	}

	s.svc = accountuc.NewService(cfg, store, nil)
	return nil
}

// List returns all configured accounts.
func (s *defaultAccountService) List() ([]*accountuc.Account, error) {
	if err := s.ensureService(); err != nil {
		return nil, err
	}
	return s.svc.List()
}

// Add adds a new account with OAuth authentication.
func (s *defaultAccountService) Add(ctx context.Context, alias string, scopes []string) (*accountuc.Account, error) {
	if err := s.ensureService(); err != nil {
		return nil, err
	}
	// Need to create a service with OAuth flow for Add operation
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	store, err := keyring.NewStore()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize keyring: %w", err)
	}
	flow := accountuc.NewDefaultOAuthFlow()
	svcWithFlow := accountuc.NewService(cfg, store, flow)
	return svcWithFlow.Add(ctx, alias, scopes)
}

// Remove removes an account.
func (s *defaultAccountService) Remove(alias string) error {
	if err := s.ensureService(); err != nil {
		return err
	}
	return s.svc.Remove(alias)
}

// Switch switches the default account.
func (s *defaultAccountService) Switch(alias string) error {
	if err := s.ensureService(); err != nil {
		return err
	}
	return s.svc.Switch(alias)
}

// Rename renames an account.
func (s *defaultAccountService) Rename(oldAlias, newAlias string) error {
	if err := s.ensureService(); err != nil {
		return err
	}
	return s.svc.Rename(oldAlias, newAlias)
}

// ResolveAccount resolves the account to use based on configuration.
func (s *defaultAccountService) ResolveAccount(flagValue string) (*accountuc.Account, error) {
	if err := s.ensureService(); err != nil {
		return nil, err
	}
	return s.svc.ResolveAccount(flagValue)
}

// GetTokenManager returns the token manager.
func (s *defaultAccountService) GetTokenManager() TokenManager {
	if err := s.ensureService(); err != nil {
		return nil
	}
	return &defaultTokenManager{tm: s.svc.GetTokenManager()}
}

// defaultTokenManager wraps the auth.TokenManager.
type defaultTokenManager struct {
	tm *auth.TokenManager
}

// GetTokenSource returns an OAuth2 token source for the given account alias.
func (m *defaultTokenManager) GetTokenSource(ctx context.Context, alias string) (oauth2.TokenSource, error) {
	return m.tm.GetTokenSource(ctx, alias)
}

// GetTokenInfo returns token information for the given account alias.
func (m *defaultTokenManager) GetTokenInfo(alias string) (*auth.TokenInfo, error) {
	return m.tm.GetTokenInfo(alias)
}

// RefreshToken refreshes the OAuth token for the given account alias.
func (m *defaultTokenManager) RefreshToken(ctx context.Context, alias string, cfg *oauth2.Config) (*oauth2.Token, error) {
	return m.tm.RefreshToken(ctx, alias, cfg)
}

// GetGrantedScopes returns the granted scopes for the given account alias.
func (m *defaultTokenManager) GetGrantedScopes(alias string) ([]string, error) {
	return m.tm.GetGrantedScopes(alias)
}

// defaultRepositoryFactory implements RepositoryFactory using production implementations.
type defaultRepositoryFactory struct{}

// NewMessageRepository creates a new message repository.
func (f *defaultRepositoryFactory) NewMessageRepository(ctx context.Context, tokenSource oauth2.TokenSource) (MessageRepository, error) {
	return repository.NewGmailRepository(ctx, tokenSource)
}

// NewDraftRepository creates a new draft repository.
func (f *defaultRepositoryFactory) NewDraftRepository(ctx context.Context, tokenSource oauth2.TokenSource) (DraftRepository, error) {
	gmailRepo, err := repository.NewGmailRepository(ctx, tokenSource)
	if err != nil {
		return nil, err
	}
	return repository.NewGmailDraftRepository(gmailRepo), nil
}

// NewThreadRepository creates a new thread repository.
func (f *defaultRepositoryFactory) NewThreadRepository(ctx context.Context, tokenSource oauth2.TokenSource) (ThreadRepository, error) {
	gmailRepo, err := repository.NewGmailRepository(ctx, tokenSource)
	if err != nil {
		return nil, err
	}
	return repository.NewGmailThreadRepository(gmailRepo), nil
}

// NewLabelRepository creates a new label repository.
func (f *defaultRepositoryFactory) NewLabelRepository(ctx context.Context, tokenSource oauth2.TokenSource) (LabelRepository, error) {
	gmailRepo, err := repository.NewGmailRepository(ctx, tokenSource)
	if err != nil {
		return nil, err
	}
	return repository.NewGmailLabelRepository(gmailRepo), nil
}

// NewEventRepository creates a new event repository.
func (f *defaultRepositoryFactory) NewEventRepository(ctx context.Context, tokenSource oauth2.TokenSource) (EventRepository, error) {
	gcalSvc, err := repository.NewGCalService(ctx, tokenSource)
	if err != nil {
		return nil, err
	}
	return gcalSvc.Events(), nil
}

// NewCalendarRepository creates a new calendar repository.
func (f *defaultRepositoryFactory) NewCalendarRepository(ctx context.Context, tokenSource oauth2.TokenSource) (CalendarRepository, error) {
	gcalSvc, err := repository.NewGCalService(ctx, tokenSource)
	if err != nil {
		return nil, err
	}
	return gcalSvc.Calendars(), nil
}

// NewACLRepository creates a new ACL repository.
func (f *defaultRepositoryFactory) NewACLRepository(ctx context.Context, tokenSource oauth2.TokenSource) (ACLRepository, error) {
	gcalSvc, err := repository.NewGCalService(ctx, tokenSource)
	if err != nil {
		return nil, err
	}
	return gcalSvc.ACL(), nil
}

// NewFreeBusyRepository creates a new free/busy repository.
func (f *defaultRepositoryFactory) NewFreeBusyRepository(ctx context.Context, tokenSource oauth2.TokenSource) (FreeBusyRepository, error) {
	gcalSvc, err := repository.NewGCalService(ctx, tokenSource)
	if err != nil {
		return nil, err
	}
	return gcalSvc.FreeBusy(), nil
}
