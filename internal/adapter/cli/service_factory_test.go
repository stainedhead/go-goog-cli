// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"context"
	"fmt"
	"strings"
	"testing"

	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

// =============================================================================
// Tests for service factory functions using dependency injection
// =============================================================================

func TestGetTokenSourceFromDeps_Success(t *testing.T) {
	mockAccount := &accountuc.Account{
		Alias: "test",
		Email: "test@example.com",
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      mockAccount,
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	tokenSource, err := getTokenSourceFromDeps(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokenSource == nil {
		t.Error("expected non-nil token source")
	}
}

func TestGetTokenSourceFromDeps_AccountResolveError(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			ResolveErr: fmt.Errorf("no account found"),
		},
		RepoFactory: &MockRepositoryFactory{},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	_, err := getTokenSourceFromDeps(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "no account found") {
		t.Errorf("expected error to contain 'no account found', got: %v", err)
	}
}

func TestGetTokenSourceFromDeps_TokenManagerError(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account: &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{
				Err: fmt.Errorf("token error"),
			},
		},
		RepoFactory: &MockRepositoryFactory{},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	_, err := getTokenSourceFromDeps(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to get token") {
		t.Errorf("expected error to contain 'failed to get token', got: %v", err)
	}
}

func TestGetTokenSourceWithEmailFromDeps_Success(t *testing.T) {
	mockAccount := &accountuc.Account{
		Alias: "test",
		Email: "test@example.com",
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      mockAccount,
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	tokenSource, email, err := getTokenSourceWithEmailFromDeps(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokenSource == nil {
		t.Error("expected non-nil token source")
	}
	if email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", email)
	}
}

func TestGetTokenSourceWithEmailFromDeps_Error(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			ResolveErr: fmt.Errorf("no account found"),
		},
		RepoFactory: &MockRepositoryFactory{},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	_, _, err := getTokenSourceWithEmailFromDeps(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetMessageRepositoryFromDeps_Success(t *testing.T) {
	mockAccount := &accountuc.Account{
		Alias: "test",
		Email: "test@example.com",
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      mockAccount,
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: &MockMessageRepository{},
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	repo, email, err := getMessageRepositoryFromDeps(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo == nil {
		t.Error("expected non-nil repository")
	}
	if email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", email)
	}
}

func TestGetMessageRepositoryFromDeps_FactoryError(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageErr: fmt.Errorf("factory error"),
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	_, _, err := getMessageRepositoryFromDeps(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to create message repository") {
		t.Errorf("expected error to contain 'failed to create message repository', got: %v", err)
	}
}

func TestGetDraftRepositoryFromDeps_Success(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftRepo: &MockDraftRepository{},
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	repo, err := getDraftRepositoryFromDeps(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo == nil {
		t.Error("expected non-nil repository")
	}
}

func TestGetDraftRepositoryFromDeps_Error(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			DraftErr: fmt.Errorf("draft factory error"),
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	_, err := getDraftRepositoryFromDeps(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetThreadRepositoryFromDeps_Success(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadRepo: &MockThreadRepository{},
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	repo, err := getThreadRepositoryFromDeps(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo == nil {
		t.Error("expected non-nil repository")
	}
}

func TestGetThreadRepositoryFromDeps_Error(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ThreadErr: fmt.Errorf("thread factory error"),
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	_, err := getThreadRepositoryFromDeps(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetLabelRepositoryFromDeps_Success(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelRepo: &MockLabelRepository{},
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	repo, err := getLabelRepositoryFromDeps(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo == nil {
		t.Error("expected non-nil repository")
	}
}

func TestGetLabelRepositoryFromDeps_Error(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			LabelErr: fmt.Errorf("label factory error"),
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	_, err := getLabelRepositoryFromDeps(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetEventRepositoryFromDeps_Success(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			EventRepo: &MockEventRepository{},
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	repo, err := getEventRepositoryFromDeps(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo == nil {
		t.Error("expected non-nil repository")
	}
}

func TestGetEventRepositoryFromDeps_Error(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			EventErr: fmt.Errorf("event factory error"),
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	_, err := getEventRepositoryFromDeps(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetCalendarRepositoryFromDeps_Success(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			CalendarRepo: &MockCalendarRepository{},
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	repo, err := getCalendarRepositoryFromDeps(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo == nil {
		t.Error("expected non-nil repository")
	}
}

func TestGetCalendarRepositoryFromDeps_Error(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			CalendarErr: fmt.Errorf("calendar factory error"),
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	_, err := getCalendarRepositoryFromDeps(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetACLRepositoryFromDeps_Success(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ACLRepo: &MockACLRepository{},
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	repo, err := getACLRepositoryFromDeps(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo == nil {
		t.Error("expected non-nil repository")
	}
}

func TestGetACLRepositoryFromDeps_Error(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			ACLErr: fmt.Errorf("ACL factory error"),
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	_, err := getACLRepositoryFromDeps(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetFreeBusyRepositoryFromDeps_Success(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			FreeBusyRepo: &MockFreeBusyRepository{},
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	repo, err := getFreeBusyRepositoryFromDeps(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo == nil {
		t.Error("expected non-nil repository")
	}
}

func TestGetFreeBusyRepositoryFromDeps_Error(t *testing.T) {
	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "test@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			FreeBusyErr: fmt.Errorf("free/busy factory error"),
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	ctx := context.Background()
	_, err := getFreeBusyRepositoryFromDeps(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestMockTokenManager_GetTokenSource(t *testing.T) {
	t.Run("success with default token source", func(t *testing.T) {
		tm := &MockTokenManager{}
		ctx := context.Background()

		ts, err := tm.GetTokenSource(ctx, "test")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if ts == nil {
			t.Error("expected non-nil token source")
		}
	})

	t.Run("success with custom token source", func(t *testing.T) {
		customTS := &MockTokenSource{}
		tm := &MockTokenManager{TokenSource: customTS}
		ctx := context.Background()

		ts, err := tm.GetTokenSource(ctx, "test")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if ts != customTS {
			t.Error("expected custom token source")
		}
	})

	t.Run("error", func(t *testing.T) {
		tm := &MockTokenManager{Err: fmt.Errorf("token error")}
		ctx := context.Background()

		_, err := tm.GetTokenSource(ctx, "test")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockTokenSource_Token(t *testing.T) {
	t.Run("success with default token", func(t *testing.T) {
		ts := &MockTokenSource{}

		token, err := ts.Token()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if token == nil {
			t.Error("expected non-nil token")
		}
		if token.AccessToken != "mock-access-token" {
			t.Errorf("expected 'mock-access-token', got %s", token.AccessToken)
		}
	})

	t.Run("error", func(t *testing.T) {
		ts := &MockTokenSource{err: fmt.Errorf("token error")}

		_, err := ts.Token()
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
