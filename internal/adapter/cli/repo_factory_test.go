// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"context"
	"fmt"
	"testing"

	"golang.org/x/oauth2"
)

// =============================================================================
// Tests for mock repository factory
// =============================================================================

func TestMockRepositoryFactory_NewMessageRepository(t *testing.T) {
	ctx := context.Background()
	tokenSource := &MockTokenSource{}

	t.Run("returns default mock when not set", func(t *testing.T) {
		factory := &MockRepositoryFactory{}

		repo, err := factory.NewMessageRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if repo == nil {
			t.Error("expected non-nil repository")
		}
	})

	t.Run("returns custom repo when set", func(t *testing.T) {
		customRepo := &MockMessageRepository{}
		factory := &MockRepositoryFactory{MessageRepo: customRepo}

		repo, err := factory.NewMessageRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if repo != customRepo {
			t.Error("expected custom repository")
		}
	})

	t.Run("returns error when set", func(t *testing.T) {
		factory := &MockRepositoryFactory{MessageErr: fmt.Errorf("message error")}

		_, err := factory.NewMessageRepository(ctx, tokenSource)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockRepositoryFactory_NewDraftRepository(t *testing.T) {
	ctx := context.Background()
	tokenSource := &MockTokenSource{}

	t.Run("returns default mock when not set", func(t *testing.T) {
		factory := &MockRepositoryFactory{}

		repo, err := factory.NewDraftRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if repo == nil {
			t.Error("expected non-nil repository")
		}
	})

	t.Run("returns custom repo when set", func(t *testing.T) {
		customRepo := &MockDraftRepository{}
		factory := &MockRepositoryFactory{DraftRepo: customRepo}

		repo, err := factory.NewDraftRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if repo != customRepo {
			t.Error("expected custom repository")
		}
	})

	t.Run("returns error when set", func(t *testing.T) {
		factory := &MockRepositoryFactory{DraftErr: fmt.Errorf("draft error")}

		_, err := factory.NewDraftRepository(ctx, tokenSource)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockRepositoryFactory_NewThreadRepository(t *testing.T) {
	ctx := context.Background()
	tokenSource := &MockTokenSource{}

	t.Run("returns default mock when not set", func(t *testing.T) {
		factory := &MockRepositoryFactory{}

		repo, err := factory.NewThreadRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if repo == nil {
			t.Error("expected non-nil repository")
		}
	})

	t.Run("returns error when set", func(t *testing.T) {
		factory := &MockRepositoryFactory{ThreadErr: fmt.Errorf("thread error")}

		_, err := factory.NewThreadRepository(ctx, tokenSource)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockRepositoryFactory_NewLabelRepository(t *testing.T) {
	ctx := context.Background()
	tokenSource := &MockTokenSource{}

	t.Run("returns default mock when not set", func(t *testing.T) {
		factory := &MockRepositoryFactory{}

		repo, err := factory.NewLabelRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if repo == nil {
			t.Error("expected non-nil repository")
		}
	})

	t.Run("returns error when set", func(t *testing.T) {
		factory := &MockRepositoryFactory{LabelErr: fmt.Errorf("label error")}

		_, err := factory.NewLabelRepository(ctx, tokenSource)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockRepositoryFactory_NewEventRepository(t *testing.T) {
	ctx := context.Background()
	tokenSource := &MockTokenSource{}

	t.Run("returns default mock when not set", func(t *testing.T) {
		factory := &MockRepositoryFactory{}

		repo, err := factory.NewEventRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if repo == nil {
			t.Error("expected non-nil repository")
		}
	})

	t.Run("returns error when set", func(t *testing.T) {
		factory := &MockRepositoryFactory{EventErr: fmt.Errorf("event error")}

		_, err := factory.NewEventRepository(ctx, tokenSource)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockRepositoryFactory_NewCalendarRepository(t *testing.T) {
	ctx := context.Background()
	tokenSource := &MockTokenSource{}

	t.Run("returns default mock when not set", func(t *testing.T) {
		factory := &MockRepositoryFactory{}

		repo, err := factory.NewCalendarRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if repo == nil {
			t.Error("expected non-nil repository")
		}
	})

	t.Run("returns error when set", func(t *testing.T) {
		factory := &MockRepositoryFactory{CalendarErr: fmt.Errorf("calendar error")}

		_, err := factory.NewCalendarRepository(ctx, tokenSource)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockRepositoryFactory_NewACLRepository(t *testing.T) {
	ctx := context.Background()
	tokenSource := &MockTokenSource{}

	t.Run("returns default mock when not set", func(t *testing.T) {
		factory := &MockRepositoryFactory{}

		repo, err := factory.NewACLRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if repo == nil {
			t.Error("expected non-nil repository")
		}
	})

	t.Run("returns error when set", func(t *testing.T) {
		factory := &MockRepositoryFactory{ACLErr: fmt.Errorf("ACL error")}

		_, err := factory.NewACLRepository(ctx, tokenSource)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockRepositoryFactory_NewFreeBusyRepository(t *testing.T) {
	ctx := context.Background()
	tokenSource := &MockTokenSource{}

	t.Run("returns default mock when not set", func(t *testing.T) {
		factory := &MockRepositoryFactory{}

		repo, err := factory.NewFreeBusyRepository(ctx, tokenSource)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if repo == nil {
			t.Error("expected non-nil repository")
		}
	})

	t.Run("returns error when set", func(t *testing.T) {
		factory := &MockRepositoryFactory{FreeBusyErr: fmt.Errorf("freebusy error")}

		_, err := factory.NewFreeBusyRepository(ctx, tokenSource)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockTokenSource_CustomToken(t *testing.T) {
	customToken := &oauth2.Token{
		AccessToken: "custom-token",
		TokenType:   "Bearer",
	}
	ts := &MockTokenSource{token: customToken}

	token, err := ts.Token()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if token.AccessToken != "custom-token" {
		t.Errorf("expected 'custom-token', got %s", token.AccessToken)
	}
}
