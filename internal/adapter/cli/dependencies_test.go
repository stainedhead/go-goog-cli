package cli

import (
	"context"
	"testing"
	"time"

	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
	domaincontacts "github.com/stainedhead/go-goog-cli/internal/domain/contacts"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
	domaintasks "github.com/stainedhead/go-goog-cli/internal/domain/tasks"
	"github.com/stainedhead/go-goog-cli/internal/infrastructure/auth"
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
	TokenSource          oauth2.TokenSource
	Err                  error
	GetTokenInfoFunc     func(alias string) (*auth.TokenInfo, error)
	RefreshTokenFunc     func(ctx context.Context, alias string, cfg *oauth2.Config) (*oauth2.Token, error)
	GetGrantedScopesFunc func(alias string) ([]string, error)
	TokenInfo            *auth.TokenInfo
	TokenInfoErr         error
	RefreshTokenRes      *oauth2.Token
	RefreshTokenErr      error
	GrantedScopes        []string
	GrantedScopesErr     error
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

// GetTokenInfo returns mock token information.
func (m *MockTokenManager) GetTokenInfo(alias string) (*auth.TokenInfo, error) {
	if m.GetTokenInfoFunc != nil {
		return m.GetTokenInfoFunc(alias)
	}
	if m.TokenInfoErr != nil {
		return nil, m.TokenInfoErr
	}
	if m.TokenInfo != nil {
		return m.TokenInfo, nil
	}
	// Return a default token info structure
	return &auth.TokenInfo{
		Account:    alias,
		HasToken:   true,
		IsExpired:  false,
		ExpiryTime: time.Now().Add(time.Hour).Format(time.RFC3339),
		Scopes:     []string{"email", "openid"},
		TokenType:  "Bearer",
	}, nil
}

// RefreshToken refreshes the mock OAuth token.
func (m *MockTokenManager) RefreshToken(ctx context.Context, alias string, cfg *oauth2.Config) (*oauth2.Token, error) {
	if m.RefreshTokenFunc != nil {
		return m.RefreshTokenFunc(ctx, alias, cfg)
	}
	if m.RefreshTokenErr != nil {
		return nil, m.RefreshTokenErr
	}
	if m.RefreshTokenRes != nil {
		return m.RefreshTokenRes, nil
	}
	return &oauth2.Token{
		AccessToken: "refreshed-token",
		Expiry:      time.Now().Add(time.Hour),
	}, nil
}

// GetGrantedScopes returns mock granted scopes.
func (m *MockTokenManager) GetGrantedScopes(alias string) ([]string, error) {
	if m.GetGrantedScopesFunc != nil {
		return m.GetGrantedScopesFunc(alias)
	}
	if m.GrantedScopesErr != nil {
		return nil, m.GrantedScopesErr
	}
	if m.GrantedScopes != nil {
		return m.GrantedScopes, nil
	}
	return []string{"email", "openid"}, nil
}

// MockAccountService implements AccountService for testing.
type MockAccountService struct {
	Accounts     []*accountuc.Account
	Account      *accountuc.Account
	ListErr      error
	ResolveErr   error
	TokenManager TokenManager
	AddFunc      func(ctx context.Context, alias string, scopes []string) (*accountuc.Account, error)
	RemoveFunc   func(alias string) error
	SwitchFunc   func(alias string) error
	RenameFunc   func(oldAlias, newAlias string) error
	AddResult    *accountuc.Account
	AddErr       error
	RemoveErr    error
	SwitchErr    error
	RenameErr    error
}

// List returns the mock accounts.
func (m *MockAccountService) List() ([]*accountuc.Account, error) {
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	return m.Accounts, nil
}

// Add adds a new account.
func (m *MockAccountService) Add(ctx context.Context, alias string, scopes []string) (*accountuc.Account, error) {
	if m.AddFunc != nil {
		return m.AddFunc(ctx, alias, scopes)
	}
	if m.AddErr != nil {
		return nil, m.AddErr
	}
	if m.AddResult != nil {
		return m.AddResult, nil
	}
	return &accountuc.Account{
		Alias:     alias,
		Email:     "test@example.com",
		IsDefault: true,
	}, nil
}

// Remove removes an account.
func (m *MockAccountService) Remove(alias string) error {
	if m.RemoveFunc != nil {
		return m.RemoveFunc(alias)
	}
	return m.RemoveErr
}

// Switch switches the default account.
func (m *MockAccountService) Switch(alias string) error {
	if m.SwitchFunc != nil {
		return m.SwitchFunc(alias)
	}
	return m.SwitchErr
}

// Rename renames an account.
func (m *MockAccountService) Rename(oldAlias, newAlias string) error {
	if m.RenameFunc != nil {
		return m.RenameFunc(oldAlias, newAlias)
	}
	return m.RenameErr
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

// MockTaskListRepository implements TaskListRepository for testing.
type MockTaskListRepository struct {
	Lists     []*domaintasks.TaskList
	TaskList  *domaintasks.TaskList
	ListErr   error
	GetErr    error
	CreateErr error
	UpdateErr error
	DeleteErr error
	ListFunc  func(ctx context.Context) ([]*domaintasks.TaskList, error)
	GetFunc   func(ctx context.Context, taskListID string) (*domaintasks.TaskList, error)
}

func (m *MockTaskListRepository) List(ctx context.Context) ([]*domaintasks.TaskList, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	if m.Lists != nil {
		return m.Lists, nil
	}
	return []*domaintasks.TaskList{}, nil
}

func (m *MockTaskListRepository) Get(ctx context.Context, taskListID string) (*domaintasks.TaskList, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, taskListID)
	}
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	if m.TaskList != nil {
		return m.TaskList, nil
	}
	return &domaintasks.TaskList{ID: taskListID, Title: "Test List"}, nil
}

func (m *MockTaskListRepository) Create(ctx context.Context, taskList *domaintasks.TaskList) (*domaintasks.TaskList, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	return taskList, nil
}

func (m *MockTaskListRepository) Update(ctx context.Context, taskList *domaintasks.TaskList) (*domaintasks.TaskList, error) {
	if m.UpdateErr != nil {
		return nil, m.UpdateErr
	}
	return taskList, nil
}

func (m *MockTaskListRepository) Delete(ctx context.Context, taskListID string) error {
	return m.DeleteErr
}

// MockTaskRepository implements TaskRepository for testing.
type MockTaskRepository struct {
	Tasks     *domaintasks.ListResult[*domaintasks.Task]
	Task      *domaintasks.Task
	ListErr   error
	GetErr    error
	CreateErr error
	UpdateErr error
	DeleteErr error
	MoveErr   error
	ClearErr  error
	ListFunc  func(ctx context.Context, taskListID string, opts domaintasks.ListOptions) (*domaintasks.ListResult[*domaintasks.Task], error)
	GetFunc   func(ctx context.Context, taskListID, taskID string) (*domaintasks.Task, error)
}

func (m *MockTaskRepository) List(ctx context.Context, taskListID string, opts domaintasks.ListOptions) (*domaintasks.ListResult[*domaintasks.Task], error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, taskListID, opts)
	}
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	if m.Tasks != nil {
		return m.Tasks, nil
	}
	return &domaintasks.ListResult[*domaintasks.Task]{Items: []*domaintasks.Task{}}, nil
}

func (m *MockTaskRepository) Get(ctx context.Context, taskListID, taskID string) (*domaintasks.Task, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, taskListID, taskID)
	}
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	if m.Task != nil {
		return m.Task, nil
	}
	return &domaintasks.Task{ID: taskID, TaskListID: taskListID, Title: "Test Task"}, nil
}

func (m *MockTaskRepository) Create(ctx context.Context, taskListID string, task *domaintasks.Task) (*domaintasks.Task, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	return task, nil
}

func (m *MockTaskRepository) Update(ctx context.Context, taskListID string, task *domaintasks.Task) (*domaintasks.Task, error) {
	if m.UpdateErr != nil {
		return nil, m.UpdateErr
	}
	return task, nil
}

func (m *MockTaskRepository) Delete(ctx context.Context, taskListID, taskID string) error {
	return m.DeleteErr
}

func (m *MockTaskRepository) Move(ctx context.Context, taskListID, taskID, parent, previous string) (*domaintasks.Task, error) {
	if m.MoveErr != nil {
		return nil, m.MoveErr
	}
	if m.Task != nil {
		return m.Task, nil
	}
	return &domaintasks.Task{ID: taskID, TaskListID: taskListID, Title: "Test Task"}, nil
}

func (m *MockTaskRepository) Clear(ctx context.Context, taskListID string) error {
	return m.ClearErr
}

// MockContactRepository implements ContactRepository for testing.
type MockContactRepository struct {
	Contacts     *domaincontacts.ListResult[*domaincontacts.Contact]
	Contact      *domaincontacts.Contact
	ListErr      error
	GetErr       error
	CreateErr    error
	UpdateErr    error
	DeleteErr    error
	SearchErr    error
	SearchResult *domaincontacts.ListResult[*domaincontacts.Contact]
}

func (m *MockContactRepository) List(ctx context.Context, opts domaincontacts.ListOptions) (*domaincontacts.ListResult[*domaincontacts.Contact], error) {
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	if m.Contacts != nil {
		return m.Contacts, nil
	}
	return &domaincontacts.ListResult[*domaincontacts.Contact]{Items: []*domaincontacts.Contact{}}, nil
}

func (m *MockContactRepository) Get(ctx context.Context, resourceName string) (*domaincontacts.Contact, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return m.Contact, nil
}

func (m *MockContactRepository) Create(ctx context.Context, contact *domaincontacts.Contact) (*domaincontacts.Contact, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	return contact, nil
}

func (m *MockContactRepository) Update(ctx context.Context, contact *domaincontacts.Contact, updateMask []string) (*domaincontacts.Contact, error) {
	if m.UpdateErr != nil {
		return nil, m.UpdateErr
	}
	return contact, nil
}

func (m *MockContactRepository) Delete(ctx context.Context, resourceName string) error {
	return m.DeleteErr
}

func (m *MockContactRepository) Search(ctx context.Context, opts domaincontacts.SearchOptions) (*domaincontacts.ListResult[*domaincontacts.Contact], error) {
	if m.SearchErr != nil {
		return nil, m.SearchErr
	}
	if m.SearchResult != nil {
		return m.SearchResult, nil
	}
	return &domaincontacts.ListResult[*domaincontacts.Contact]{Items: []*domaincontacts.Contact{}}, nil
}

func (m *MockContactRepository) BatchGet(ctx context.Context, resourceNames []string) ([]*domaincontacts.Contact, error) {
	return []*domaincontacts.Contact{}, nil
}

// MockContactGroupRepository implements ContactGroupRepository for testing.
type MockContactGroupRepository struct {
	Groups           []*domaincontacts.ContactGroup
	Group            *domaincontacts.ContactGroup
	Members          *domaincontacts.ListResult[*domaincontacts.Contact]
	ListErr          error
	GetErr           error
	CreateErr        error
	UpdateErr        error
	DeleteErr        error
	ListMembersErr   error
	AddMembersErr    error
	RemoveMembersErr error
}

func (m *MockContactGroupRepository) List(ctx context.Context) ([]*domaincontacts.ContactGroup, error) {
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	return m.Groups, nil
}

func (m *MockContactGroupRepository) Get(ctx context.Context, resourceName string) (*domaincontacts.ContactGroup, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return m.Group, nil
}

func (m *MockContactGroupRepository) Create(ctx context.Context, group *domaincontacts.ContactGroup) (*domaincontacts.ContactGroup, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	return group, nil
}

func (m *MockContactGroupRepository) Update(ctx context.Context, group *domaincontacts.ContactGroup) (*domaincontacts.ContactGroup, error) {
	if m.UpdateErr != nil {
		return nil, m.UpdateErr
	}
	return group, nil
}

func (m *MockContactGroupRepository) Delete(ctx context.Context, resourceName string) error {
	return m.DeleteErr
}

func (m *MockContactGroupRepository) ListMembers(ctx context.Context, resourceName string, opts domaincontacts.ListOptions) (*domaincontacts.ListResult[*domaincontacts.Contact], error) {
	if m.ListMembersErr != nil {
		return nil, m.ListMembersErr
	}
	if m.Members != nil {
		return m.Members, nil
	}
	return &domaincontacts.ListResult[*domaincontacts.Contact]{Items: []*domaincontacts.Contact{}}, nil
}

func (m *MockContactGroupRepository) AddMembers(ctx context.Context, groupResourceName string, contactResourceNames []string) error {
	return m.AddMembersErr
}

func (m *MockContactGroupRepository) RemoveMembers(ctx context.Context, groupResourceName string, contactResourceNames []string) error {
	return m.RemoveMembersErr
}

// MockRepositoryFactory implements RepositoryFactory for testing.
type MockRepositoryFactory struct {
	MessageRepo      MessageRepository
	DraftRepo        DraftRepository
	ThreadRepo       ThreadRepository
	LabelRepo        LabelRepository
	EventRepo        EventRepository
	CalendarRepo     CalendarRepository
	ACLRepo          ACLRepository
	FreeBusyRepo     FreeBusyRepository
	TaskListRepo     TaskListRepository
	TaskRepo         TaskRepository
	ContactRepo      ContactRepository
	ContactGroupRepo ContactGroupRepository
	MessageErr       error
	DraftErr         error
	ThreadErr        error
	LabelErr         error
	EventErr         error
	CalendarErr      error
	ACLErr           error
	FreeBusyErr      error
	TaskListErr      error
	TaskErr          error
	ContactErr       error
	ContactGroupErr  error
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

func (f *MockRepositoryFactory) NewTaskListRepository(ctx context.Context, tokenSource oauth2.TokenSource) (TaskListRepository, error) {
	if f.TaskListErr != nil {
		return nil, f.TaskListErr
	}
	if f.TaskListRepo == nil {
		return &MockTaskListRepository{}, nil
	}
	return f.TaskListRepo, nil
}

func (f *MockRepositoryFactory) NewTaskRepository(ctx context.Context, tokenSource oauth2.TokenSource) (TaskRepository, error) {
	if f.TaskErr != nil {
		return nil, f.TaskErr
	}
	if f.TaskRepo == nil {
		return &MockTaskRepository{}, nil
	}
	return f.TaskRepo, nil
}

func (f *MockRepositoryFactory) NewContactRepository(ctx context.Context, tokenSource oauth2.TokenSource) (ContactRepository, error) {
	if f.ContactErr != nil {
		return nil, f.ContactErr
	}
	if f.ContactRepo == nil {
		return &MockContactRepository{}, nil
	}
	return f.ContactRepo, nil
}

func (f *MockRepositoryFactory) NewContactGroupRepository(ctx context.Context, tokenSource oauth2.TokenSource) (ContactGroupRepository, error) {
	if f.ContactGroupErr != nil {
		return nil, f.ContactGroupErr
	}
	if f.ContactGroupRepo == nil {
		return &MockContactGroupRepository{}, nil
	}
	return f.ContactGroupRepo, nil
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
