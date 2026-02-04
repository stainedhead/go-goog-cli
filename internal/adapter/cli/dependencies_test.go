package cli

import (
	"context"
	"testing"
	"time"

	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
	"golang.org/x/oauth2"
)

// MockTokenSource implements oauth2.TokenSource for testing.
type MockTokenSource struct {
	token *oauth2.Token
	err   error
}

// Token returns the mock token.
func (m *MockTokenSource) Token() (*oauth2.Token, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.token == nil {
		return &oauth2.Token{
			AccessToken: "mock-access-token",
			TokenType:   "Bearer",
			Expiry:      time.Now().Add(time.Hour),
		}, nil
	}
	return m.token, nil
}

// MockTokenManager implements TokenManager for testing.
type MockTokenManager struct {
	TokenSource oauth2.TokenSource
	Err         error
}

// GetTokenSource returns the mock token source.
func (m *MockTokenManager) GetTokenSource(ctx context.Context, alias string) (oauth2.TokenSource, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if m.TokenSource == nil {
		return &MockTokenSource{}, nil
	}
	return m.TokenSource, nil
}

// MockAccountService implements AccountService for testing.
type MockAccountService struct {
	Accounts     []*accountuc.Account
	Account      *accountuc.Account
	ListErr      error
	ResolveErr   error
	TokenManager TokenManager
}

// List returns the mock accounts.
func (m *MockAccountService) List() ([]*accountuc.Account, error) {
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	return m.Accounts, nil
}

// ResolveAccount returns the mock account.
func (m *MockAccountService) ResolveAccount(flagValue string) (*accountuc.Account, error) {
	if m.ResolveErr != nil {
		return nil, m.ResolveErr
	}
	return m.Account, nil
}

// GetTokenManager returns the mock token manager.
func (m *MockAccountService) GetTokenManager() TokenManager {
	if m.TokenManager == nil {
		return &MockTokenManager{}
	}
	return m.TokenManager
}

// MockMessageRepository implements MessageRepository for testing.
type MockMessageRepository struct {
	Messages      []*mail.Message
	Message       *mail.Message
	ListResult    *mail.ListResult[*mail.Message]
	ListErr       error
	GetErr        error
	SendErr       error
	ReplyErr      error
	ForwardErr    error
	TrashErr      error
	UntrashErr    error
	DeleteErr     error
	ArchiveErr    error
	ModifyErr     error
	SearchErr     error
	ModifyResult  *mail.Message
	SendResult    *mail.Message
	ReplyResult   *mail.Message
	ForwardResult *mail.Message
	SearchResult  *mail.ListResult[*mail.Message]
}

func (m *MockMessageRepository) List(ctx context.Context, opts mail.ListOptions) (*mail.ListResult[*mail.Message], error) {
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	if m.ListResult != nil {
		return m.ListResult, nil
	}
	return &mail.ListResult[*mail.Message]{Items: m.Messages}, nil
}

func (m *MockMessageRepository) Get(ctx context.Context, id string) (*mail.Message, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return m.Message, nil
}

func (m *MockMessageRepository) Send(ctx context.Context, msg *mail.Message) (*mail.Message, error) {
	if m.SendErr != nil {
		return nil, m.SendErr
	}
	if m.SendResult != nil {
		return m.SendResult, nil
	}
	return msg, nil
}

func (m *MockMessageRepository) Reply(ctx context.Context, messageID string, reply *mail.Message) (*mail.Message, error) {
	if m.ReplyErr != nil {
		return nil, m.ReplyErr
	}
	if m.ReplyResult != nil {
		return m.ReplyResult, nil
	}
	return reply, nil
}

func (m *MockMessageRepository) Forward(ctx context.Context, messageID string, forward *mail.Message) (*mail.Message, error) {
	if m.ForwardErr != nil {
		return nil, m.ForwardErr
	}
	if m.ForwardResult != nil {
		return m.ForwardResult, nil
	}
	return forward, nil
}

func (m *MockMessageRepository) Trash(ctx context.Context, id string) error {
	return m.TrashErr
}

func (m *MockMessageRepository) Untrash(ctx context.Context, id string) error {
	return m.UntrashErr
}

func (m *MockMessageRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteErr
}

func (m *MockMessageRepository) Archive(ctx context.Context, id string) error {
	return m.ArchiveErr
}

func (m *MockMessageRepository) Modify(ctx context.Context, id string, req mail.ModifyRequest) (*mail.Message, error) {
	if m.ModifyErr != nil {
		return nil, m.ModifyErr
	}
	if m.ModifyResult != nil {
		return m.ModifyResult, nil
	}
	return m.Message, nil
}

func (m *MockMessageRepository) Search(ctx context.Context, query string, opts mail.ListOptions) (*mail.ListResult[*mail.Message], error) {
	if m.SearchErr != nil {
		return nil, m.SearchErr
	}
	if m.SearchResult != nil {
		return m.SearchResult, nil
	}
	return &mail.ListResult[*mail.Message]{Items: m.Messages}, nil
}

// MockDraftRepository implements DraftRepository for testing.
type MockDraftRepository struct {
	Drafts       []*mail.Draft
	Draft        *mail.Draft
	ListResult   *mail.ListResult[*mail.Draft]
	ListErr      error
	GetErr       error
	CreateErr    error
	UpdateErr    error
	SendErr      error
	DeleteErr    error
	CreateResult *mail.Draft
	UpdateResult *mail.Draft
	SendResult   *mail.Message
}

func (m *MockDraftRepository) List(ctx context.Context, opts mail.ListOptions) (*mail.ListResult[*mail.Draft], error) {
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	if m.ListResult != nil {
		return m.ListResult, nil
	}
	return &mail.ListResult[*mail.Draft]{Items: m.Drafts}, nil
}

func (m *MockDraftRepository) Get(ctx context.Context, id string) (*mail.Draft, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return m.Draft, nil
}

func (m *MockDraftRepository) Create(ctx context.Context, draft *mail.Draft) (*mail.Draft, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	if m.CreateResult != nil {
		return m.CreateResult, nil
	}
	draft.ID = "mock-draft-id"
	return draft, nil
}

func (m *MockDraftRepository) Update(ctx context.Context, draft *mail.Draft) (*mail.Draft, error) {
	if m.UpdateErr != nil {
		return nil, m.UpdateErr
	}
	if m.UpdateResult != nil {
		return m.UpdateResult, nil
	}
	return draft, nil
}

func (m *MockDraftRepository) Send(ctx context.Context, id string) (*mail.Message, error) {
	if m.SendErr != nil {
		return nil, m.SendErr
	}
	if m.SendResult != nil {
		return m.SendResult, nil
	}
	return &mail.Message{ID: "sent-msg-id"}, nil
}

func (m *MockDraftRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteErr
}

// MockThreadRepository implements ThreadRepository for testing.
type MockThreadRepository struct {
	Threads      []*mail.Thread
	Thread       *mail.Thread
	ListResult   *mail.ListResult[*mail.Thread]
	ListErr      error
	GetErr       error
	ModifyErr    error
	TrashErr     error
	UntrashErr   error
	DeleteErr    error
	ModifyResult *mail.Thread
}

func (m *MockThreadRepository) List(ctx context.Context, opts mail.ListOptions) (*mail.ListResult[*mail.Thread], error) {
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	if m.ListResult != nil {
		return m.ListResult, nil
	}
	return &mail.ListResult[*mail.Thread]{Items: m.Threads}, nil
}

func (m *MockThreadRepository) Get(ctx context.Context, id string) (*mail.Thread, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return m.Thread, nil
}

func (m *MockThreadRepository) Modify(ctx context.Context, id string, req mail.ModifyRequest) (*mail.Thread, error) {
	if m.ModifyErr != nil {
		return nil, m.ModifyErr
	}
	if m.ModifyResult != nil {
		return m.ModifyResult, nil
	}
	return m.Thread, nil
}

func (m *MockThreadRepository) Trash(ctx context.Context, id string) error {
	return m.TrashErr
}

func (m *MockThreadRepository) Untrash(ctx context.Context, id string) error {
	return m.UntrashErr
}

func (m *MockThreadRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteErr
}

// MockLabelRepository implements LabelRepository for testing.
type MockLabelRepository struct {
	Labels       []*mail.Label
	Label        *mail.Label
	ListErr      error
	GetErr       error
	GetByNameErr error
	CreateErr    error
	UpdateErr    error
	DeleteErr    error
	CreateResult *mail.Label
	UpdateResult *mail.Label
}

func (m *MockLabelRepository) List(ctx context.Context) ([]*mail.Label, error) {
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	return m.Labels, nil
}

func (m *MockLabelRepository) Get(ctx context.Context, id string) (*mail.Label, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return m.Label, nil
}

func (m *MockLabelRepository) GetByName(ctx context.Context, name string) (*mail.Label, error) {
	if m.GetByNameErr != nil {
		return nil, m.GetByNameErr
	}
	return m.Label, nil
}

func (m *MockLabelRepository) Create(ctx context.Context, label *mail.Label) (*mail.Label, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	if m.CreateResult != nil {
		return m.CreateResult, nil
	}
	label.ID = "mock-label-id"
	return label, nil
}

func (m *MockLabelRepository) Update(ctx context.Context, label *mail.Label) (*mail.Label, error) {
	if m.UpdateErr != nil {
		return nil, m.UpdateErr
	}
	if m.UpdateResult != nil {
		return m.UpdateResult, nil
	}
	return label, nil
}

func (m *MockLabelRepository) Delete(ctx context.Context, id string) error {
	return m.DeleteErr
}

// MockEventRepository implements EventRepository for testing.
type MockEventRepository struct {
	Events         []*calendar.Event
	Event          *calendar.Event
	ListErr        error
	GetErr         error
	CreateErr      error
	UpdateErr      error
	DeleteErr      error
	MoveErr        error
	QuickAddErr    error
	InstancesErr   error
	RSVPErr        error
	CreateResult   *calendar.Event
	UpdateResult   *calendar.Event
	MoveResult     *calendar.Event
	QuickAddResult *calendar.Event
}

func (m *MockEventRepository) List(ctx context.Context, calendarID string, timeMin, timeMax time.Time) ([]*calendar.Event, error) {
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	return m.Events, nil
}

func (m *MockEventRepository) Get(ctx context.Context, calendarID, eventID string) (*calendar.Event, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return m.Event, nil
}

func (m *MockEventRepository) Create(ctx context.Context, calendarID string, event *calendar.Event) (*calendar.Event, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	if m.CreateResult != nil {
		return m.CreateResult, nil
	}
	event.ID = "mock-event-id"
	return event, nil
}

func (m *MockEventRepository) Update(ctx context.Context, calendarID string, event *calendar.Event) (*calendar.Event, error) {
	if m.UpdateErr != nil {
		return nil, m.UpdateErr
	}
	if m.UpdateResult != nil {
		return m.UpdateResult, nil
	}
	return event, nil
}

func (m *MockEventRepository) Delete(ctx context.Context, calendarID, eventID string) error {
	return m.DeleteErr
}

func (m *MockEventRepository) Move(ctx context.Context, sourceCalendarID, eventID, destinationCalendarID string) (*calendar.Event, error) {
	if m.MoveErr != nil {
		return nil, m.MoveErr
	}
	if m.MoveResult != nil {
		return m.MoveResult, nil
	}
	return m.Event, nil
}

func (m *MockEventRepository) QuickAdd(ctx context.Context, calendarID, text string) (*calendar.Event, error) {
	if m.QuickAddErr != nil {
		return nil, m.QuickAddErr
	}
	if m.QuickAddResult != nil {
		return m.QuickAddResult, nil
	}
	return &calendar.Event{ID: "quick-add-id", Title: text}, nil
}

func (m *MockEventRepository) Instances(ctx context.Context, calendarID, eventID string, timeMin, timeMax time.Time) ([]*calendar.Event, error) {
	if m.InstancesErr != nil {
		return nil, m.InstancesErr
	}
	return m.Events, nil
}

func (m *MockEventRepository) RSVP(ctx context.Context, calendarID, eventID, response string) error {
	return m.RSVPErr
}

// MockCalendarRepository implements CalendarRepository for testing.
type MockCalendarRepository struct {
	Calendars    []*calendar.Calendar
	Calendar     *calendar.Calendar
	ListErr      error
	GetErr       error
	CreateErr    error
	UpdateErr    error
	DeleteErr    error
	ClearErr     error
	CreateResult *calendar.Calendar
	UpdateResult *calendar.Calendar
}

func (m *MockCalendarRepository) List(ctx context.Context) ([]*calendar.Calendar, error) {
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	return m.Calendars, nil
}

func (m *MockCalendarRepository) Get(ctx context.Context, calendarID string) (*calendar.Calendar, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return m.Calendar, nil
}

func (m *MockCalendarRepository) Create(ctx context.Context, cal *calendar.Calendar) (*calendar.Calendar, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	if m.CreateResult != nil {
		return m.CreateResult, nil
	}
	cal.ID = "mock-calendar-id"
	return cal, nil
}

func (m *MockCalendarRepository) Update(ctx context.Context, cal *calendar.Calendar) (*calendar.Calendar, error) {
	if m.UpdateErr != nil {
		return nil, m.UpdateErr
	}
	if m.UpdateResult != nil {
		return m.UpdateResult, nil
	}
	return cal, nil
}

func (m *MockCalendarRepository) Delete(ctx context.Context, calendarID string) error {
	return m.DeleteErr
}

func (m *MockCalendarRepository) Clear(ctx context.Context, calendarID string) error {
	return m.ClearErr
}

// MockACLRepository implements ACLRepository for testing.
type MockACLRepository struct {
	Rules        []*calendar.ACLRule
	Rule         *calendar.ACLRule
	ListErr      error
	GetErr       error
	InsertErr    error
	UpdateErr    error
	DeleteErr    error
	InsertResult *calendar.ACLRule
	UpdateResult *calendar.ACLRule
}

func (m *MockACLRepository) List(ctx context.Context, calendarID string) ([]*calendar.ACLRule, error) {
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	return m.Rules, nil
}

func (m *MockACLRepository) Get(ctx context.Context, calendarID, ruleID string) (*calendar.ACLRule, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return m.Rule, nil
}

func (m *MockACLRepository) Insert(ctx context.Context, calendarID string, rule *calendar.ACLRule) (*calendar.ACLRule, error) {
	if m.InsertErr != nil {
		return nil, m.InsertErr
	}
	if m.InsertResult != nil {
		return m.InsertResult, nil
	}
	rule.ID = "mock-acl-id"
	return rule, nil
}

func (m *MockACLRepository) Update(ctx context.Context, calendarID string, rule *calendar.ACLRule) (*calendar.ACLRule, error) {
	if m.UpdateErr != nil {
		return nil, m.UpdateErr
	}
	if m.UpdateResult != nil {
		return m.UpdateResult, nil
	}
	return rule, nil
}

func (m *MockACLRepository) Delete(ctx context.Context, calendarID, ruleID string) error {
	return m.DeleteErr
}

// MockFreeBusyRepository implements FreeBusyRepository for testing.
type MockFreeBusyRepository struct {
	Response *calendar.FreeBusyResponse
	QueryErr error
}

func (m *MockFreeBusyRepository) Query(ctx context.Context, request *calendar.FreeBusyRequest) (*calendar.FreeBusyResponse, error) {
	if m.QueryErr != nil {
		return nil, m.QueryErr
	}
	if m.Response != nil {
		return m.Response, nil
	}
	return &calendar.FreeBusyResponse{Calendars: make(map[string][]*calendar.TimePeriod)}, nil
}

// MockRepositoryFactory implements RepositoryFactory for testing.
type MockRepositoryFactory struct {
	MessageRepo  MessageRepository
	DraftRepo    DraftRepository
	ThreadRepo   ThreadRepository
	LabelRepo    LabelRepository
	EventRepo    EventRepository
	CalendarRepo CalendarRepository
	ACLRepo      ACLRepository
	FreeBusyRepo FreeBusyRepository
	MessageErr   error
	DraftErr     error
	ThreadErr    error
	LabelErr     error
	EventErr     error
	CalendarErr  error
	ACLErr       error
	FreeBusyErr  error
}

func (f *MockRepositoryFactory) NewMessageRepository(ctx context.Context, tokenSource oauth2.TokenSource) (MessageRepository, error) {
	if f.MessageErr != nil {
		return nil, f.MessageErr
	}
	if f.MessageRepo == nil {
		return &MockMessageRepository{}, nil
	}
	return f.MessageRepo, nil
}

func (f *MockRepositoryFactory) NewDraftRepository(ctx context.Context, tokenSource oauth2.TokenSource) (DraftRepository, error) {
	if f.DraftErr != nil {
		return nil, f.DraftErr
	}
	if f.DraftRepo == nil {
		return &MockDraftRepository{}, nil
	}
	return f.DraftRepo, nil
}

func (f *MockRepositoryFactory) NewThreadRepository(ctx context.Context, tokenSource oauth2.TokenSource) (ThreadRepository, error) {
	if f.ThreadErr != nil {
		return nil, f.ThreadErr
	}
	if f.ThreadRepo == nil {
		return &MockThreadRepository{}, nil
	}
	return f.ThreadRepo, nil
}

func (f *MockRepositoryFactory) NewLabelRepository(ctx context.Context, tokenSource oauth2.TokenSource) (LabelRepository, error) {
	if f.LabelErr != nil {
		return nil, f.LabelErr
	}
	if f.LabelRepo == nil {
		return &MockLabelRepository{}, nil
	}
	return f.LabelRepo, nil
}

func (f *MockRepositoryFactory) NewEventRepository(ctx context.Context, tokenSource oauth2.TokenSource) (EventRepository, error) {
	if f.EventErr != nil {
		return nil, f.EventErr
	}
	if f.EventRepo == nil {
		return &MockEventRepository{}, nil
	}
	return f.EventRepo, nil
}

func (f *MockRepositoryFactory) NewCalendarRepository(ctx context.Context, tokenSource oauth2.TokenSource) (CalendarRepository, error) {
	if f.CalendarErr != nil {
		return nil, f.CalendarErr
	}
	if f.CalendarRepo == nil {
		return &MockCalendarRepository{}, nil
	}
	return f.CalendarRepo, nil
}

func (f *MockRepositoryFactory) NewACLRepository(ctx context.Context, tokenSource oauth2.TokenSource) (ACLRepository, error) {
	if f.ACLErr != nil {
		return nil, f.ACLErr
	}
	if f.ACLRepo == nil {
		return &MockACLRepository{}, nil
	}
	return f.ACLRepo, nil
}

func (f *MockRepositoryFactory) NewFreeBusyRepository(ctx context.Context, tokenSource oauth2.TokenSource) (FreeBusyRepository, error) {
	if f.FreeBusyErr != nil {
		return nil, f.FreeBusyErr
	}
	if f.FreeBusyRepo == nil {
		return &MockFreeBusyRepository{}, nil
	}
	return f.FreeBusyRepo, nil
}

// NewTestDependencies creates a Dependencies instance with all mock implementations.
// This is a convenience function for setting up tests.
func NewTestDependencies() *Dependencies {
	return &Dependencies{
		AccountService: &MockAccountService{
			Account: &accountuc.Account{
				Alias: "test",
				Email: "test@example.com",
			},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{},
	}
}

// TestDependencies tests the dependency injection infrastructure.
func TestDependencies(t *testing.T) {
	t.Run("SetDependencies and GetDependencies", func(t *testing.T) {
		// Reset to clean state
		ResetDependencies()
		defer ResetDependencies()

		testDeps := NewTestDependencies()
		SetDependencies(testDeps)

		got := GetDependencies()
		if got != testDeps {
			t.Errorf("GetDependencies() = %v, want %v", got, testDeps)
		}
	})

	t.Run("GetDependencies returns default when not set", func(t *testing.T) {
		ResetDependencies()
		defer ResetDependencies()

		got := GetDependencies()
		if got == nil {
			t.Error("GetDependencies() returned nil, want default dependencies")
		}
		if got.AccountService == nil {
			t.Error("Default dependencies has nil AccountService")
		}
		if got.RepoFactory == nil {
			t.Error("Default dependencies has nil RepoFactory")
		}
	})

	t.Run("ResetDependencies clears dependencies", func(t *testing.T) {
		testDeps := NewTestDependencies()
		SetDependencies(testDeps)
		ResetDependencies()

		// GetDependencies should now return default
		got := GetDependencies()
		if got == testDeps {
			t.Error("ResetDependencies did not clear dependencies")
		}
		// Clean up
		ResetDependencies()
	})
}

// TestMockAccountService tests the mock account service.
func TestMockAccountService(t *testing.T) {
	t.Run("List returns accounts", func(t *testing.T) {
		accounts := []*accountuc.Account{
			{Alias: "test1", Email: "test1@example.com"},
			{Alias: "test2", Email: "test2@example.com"},
		}
		svc := &MockAccountService{Accounts: accounts}

		got, err := svc.List()
		if err != nil {
			t.Errorf("List() error = %v", err)
		}
		if len(got) != len(accounts) {
			t.Errorf("List() returned %d accounts, want %d", len(got), len(accounts))
		}
	})

	t.Run("ResolveAccount returns account", func(t *testing.T) {
		account := &accountuc.Account{Alias: "test", Email: "test@example.com"}
		svc := &MockAccountService{Account: account}

		got, err := svc.ResolveAccount("test")
		if err != nil {
			t.Errorf("ResolveAccount() error = %v", err)
		}
		if got != account {
			t.Errorf("ResolveAccount() = %v, want %v", got, account)
		}
	})
}

// TestMockRepositoryFactory tests the mock repository factory.
func TestMockRepositoryFactory(t *testing.T) {
	ctx := context.Background()
	tokenSource := &MockTokenSource{}
	factory := &MockRepositoryFactory{}

	t.Run("NewMessageRepository", func(t *testing.T) {
		repo, err := factory.NewMessageRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("NewMessageRepository() error = %v", err)
		}
		if repo == nil {
			t.Error("NewMessageRepository() returned nil")
		}
	})

	t.Run("NewDraftRepository", func(t *testing.T) {
		repo, err := factory.NewDraftRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("NewDraftRepository() error = %v", err)
		}
		if repo == nil {
			t.Error("NewDraftRepository() returned nil")
		}
	})

	t.Run("NewThreadRepository", func(t *testing.T) {
		repo, err := factory.NewThreadRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("NewThreadRepository() error = %v", err)
		}
		if repo == nil {
			t.Error("NewThreadRepository() returned nil")
		}
	})

	t.Run("NewLabelRepository", func(t *testing.T) {
		repo, err := factory.NewLabelRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("NewLabelRepository() error = %v", err)
		}
		if repo == nil {
			t.Error("NewLabelRepository() returned nil")
		}
	})

	t.Run("NewEventRepository", func(t *testing.T) {
		repo, err := factory.NewEventRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("NewEventRepository() error = %v", err)
		}
		if repo == nil {
			t.Error("NewEventRepository() returned nil")
		}
	})

	t.Run("NewCalendarRepository", func(t *testing.T) {
		repo, err := factory.NewCalendarRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("NewCalendarRepository() error = %v", err)
		}
		if repo == nil {
			t.Error("NewCalendarRepository() returned nil")
		}
	})

	t.Run("NewACLRepository", func(t *testing.T) {
		repo, err := factory.NewACLRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("NewACLRepository() error = %v", err)
		}
		if repo == nil {
			t.Error("NewACLRepository() returned nil")
		}
	})

	t.Run("NewFreeBusyRepository", func(t *testing.T) {
		repo, err := factory.NewFreeBusyRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("NewFreeBusyRepository() error = %v", err)
		}
		if repo == nil {
			t.Error("NewFreeBusyRepository() returned nil")
		}
	})
}
