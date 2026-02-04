// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

// =============================================================================
// Tests using dependency injection with mocks for mail compose commands
// =============================================================================

func TestRunMailSend_WithMockDependencies(t *testing.T) {
	mockRepo := &MockMessageRepository{
		SendResult: &mail.Message{
			ID:       "sent-msg-id",
			ThreadID: "thread-123",
			Subject:  "Test Subject",
		},
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "sender@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	// Save and restore original flag values
	origTo := mailSendTo
	origSubject := mailSendSubject
	origBody := mailSendBody
	origHTML := mailSendHTML
	origCc := mailSendCc
	origBcc := mailSendBcc
	mailSendTo = []string{"recipient@example.com"}
	mailSendSubject = "Test Subject"
	mailSendBody = "Test body content"
	mailSendHTML = false
	mailSendCc = []string{}
	mailSendBcc = []string{}
	defer func() {
		mailSendTo = origTo
		mailSendSubject = origSubject
		mailSendBody = origBody
		mailSendHTML = origHTML
		mailSendCc = origCc
		mailSendBcc = origBcc
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailSend(cmd, []string{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "sent-msg-id") {
		t.Errorf("expected output to contain sent message ID, got: %s", output)
	}
	if !contains(output, "thread-123") {
		t.Errorf("expected output to contain thread ID, got: %s", output)
	}
}

func TestRunMailReply_WithMockDependencies(t *testing.T) {
	originalMsg := &mail.Message{
		ID:      "original-id",
		From:    "sender@example.com",
		To:      []string{"me@example.com"},
		Subject: "Original Subject",
	}

	mockRepo := &MockMessageRepository{
		Message: originalMsg,
		ReplyResult: &mail.Message{
			ID:       "reply-msg-id",
			ThreadID: "thread-123",
			Subject:  "Re: Original Subject",
		},
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "me@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origBody := mailReplyBody
	origAll := mailReplyAll
	mailReplyBody = "Thanks for your message!"
	mailReplyAll = false
	defer func() {
		mailReplyBody = origBody
		mailReplyAll = origAll
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailReply(cmd, []string{"original-id"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "reply-msg-id") {
		t.Errorf("expected output to contain reply message ID, got: %s", output)
	}
	if !contains(output, "thread-123") {
		t.Errorf("expected output to contain thread ID, got: %s", output)
	}
}

func TestRunMailForward_WithMockDependencies(t *testing.T) {
	mockRepo := &MockMessageRepository{
		ForwardResult: &mail.Message{
			ID:       "forward-msg-id",
			ThreadID: "thread-123",
			Subject:  "Fwd: Original Subject",
		},
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "me@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTo := mailForwardTo
	origBody := mailForwardBody
	mailForwardTo = []string{"colleague@example.com"}
	mailForwardBody = "FYI - see below"
	defer func() {
		mailForwardTo = origTo
		mailForwardBody = origBody
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailForward(cmd, []string{"original-id"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "forward-msg-id") {
		t.Errorf("expected output to contain forward message ID, got: %s", output)
	}
	if !contains(output, "thread-123") {
		t.Errorf("expected output to contain thread ID, got: %s", output)
	}
}

func TestMockMessageRepository_Send(t *testing.T) {
	t.Run("Send success", func(t *testing.T) {
		repo := &MockMessageRepository{
			SendResult: &mail.Message{ID: "sent-id"},
		}

		msg := &mail.Message{To: []string{"user@example.com"}, Subject: "Test"}
		result, err := repo.Send(nil, msg)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "sent-id" {
			t.Errorf("expected ID 'sent-id', got %s", result.ID)
		}
	})

	t.Run("Send error", func(t *testing.T) {
		repo := &MockMessageRepository{SendErr: fmt.Errorf("send error")}

		msg := &mail.Message{To: []string{"user@example.com"}}
		_, err := repo.Send(nil, msg)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockMessageRepository_Reply(t *testing.T) {
	t.Run("Reply success", func(t *testing.T) {
		repo := &MockMessageRepository{
			ReplyResult: &mail.Message{ID: "reply-id"},
		}

		reply := &mail.Message{Body: "Thanks!"}
		result, err := repo.Reply(nil, "msg-id", reply)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "reply-id" {
			t.Errorf("expected ID 'reply-id', got %s", result.ID)
		}
	})

	t.Run("Reply error", func(t *testing.T) {
		repo := &MockMessageRepository{ReplyErr: fmt.Errorf("reply error")}

		reply := &mail.Message{Body: "Thanks!"}
		_, err := repo.Reply(nil, "msg-id", reply)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockMessageRepository_Forward(t *testing.T) {
	t.Run("Forward success", func(t *testing.T) {
		repo := &MockMessageRepository{
			ForwardResult: &mail.Message{ID: "forward-id"},
		}

		forward := &mail.Message{To: []string{"user@example.com"}}
		result, err := repo.Forward(nil, "msg-id", forward)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "forward-id" {
			t.Errorf("expected ID 'forward-id', got %s", result.ID)
		}
	})

	t.Run("Forward error", func(t *testing.T) {
		repo := &MockMessageRepository{ForwardErr: fmt.Errorf("forward error")}

		forward := &mail.Message{To: []string{"user@example.com"}}
		_, err := repo.Forward(nil, "msg-id", forward)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockMessageRepository_ListAndSearch(t *testing.T) {
	t.Run("List success", func(t *testing.T) {
		messages := []*mail.Message{
			{ID: "msg1", Subject: "Test 1"},
			{ID: "msg2", Subject: "Test 2"},
		}
		repo := &MockMessageRepository{Messages: messages}

		result, err := repo.List(nil, mail.ListOptions{MaxResults: 10})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result.Items) != 2 {
			t.Errorf("expected 2 messages, got %d", len(result.Items))
		}
	})

	t.Run("List error", func(t *testing.T) {
		repo := &MockMessageRepository{ListErr: fmt.Errorf("list error")}

		_, err := repo.List(nil, mail.ListOptions{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Search success", func(t *testing.T) {
		messages := []*mail.Message{
			{ID: "msg1", Subject: "Search result"},
		}
		repo := &MockMessageRepository{SearchResult: &mail.ListResult[*mail.Message]{Items: messages}}

		result, err := repo.Search(nil, "is:unread", mail.ListOptions{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result.Items) != 1 {
			t.Errorf("expected 1 message, got %d", len(result.Items))
		}
	})

	t.Run("Search error", func(t *testing.T) {
		repo := &MockMessageRepository{SearchErr: fmt.Errorf("search error")}

		_, err := repo.Search(nil, "query", mail.ListOptions{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestMockMessageRepository_Actions(t *testing.T) {
	t.Run("Trash success", func(t *testing.T) {
		repo := &MockMessageRepository{}
		err := repo.Trash(nil, "msg-id")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Trash error", func(t *testing.T) {
		repo := &MockMessageRepository{TrashErr: fmt.Errorf("trash error")}
		err := repo.Trash(nil, "msg-id")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Untrash success", func(t *testing.T) {
		repo := &MockMessageRepository{}
		err := repo.Untrash(nil, "msg-id")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Untrash error", func(t *testing.T) {
		repo := &MockMessageRepository{UntrashErr: fmt.Errorf("untrash error")}
		err := repo.Untrash(nil, "msg-id")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Archive success", func(t *testing.T) {
		repo := &MockMessageRepository{}
		err := repo.Archive(nil, "msg-id")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Archive error", func(t *testing.T) {
		repo := &MockMessageRepository{ArchiveErr: fmt.Errorf("archive error")}
		err := repo.Archive(nil, "msg-id")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Delete success", func(t *testing.T) {
		repo := &MockMessageRepository{}
		err := repo.Delete(nil, "msg-id")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Delete error", func(t *testing.T) {
		repo := &MockMessageRepository{DeleteErr: fmt.Errorf("delete error")}
		err := repo.Delete(nil, "msg-id")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Modify success", func(t *testing.T) {
		repo := &MockMessageRepository{
			ModifyResult: &mail.Message{ID: "msg-id", Labels: []string{"STARRED"}},
		}
		req := mail.ModifyRequest{AddLabels: []string{"STARRED"}}
		result, err := repo.Modify(nil, "msg-id", req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result == nil {
			t.Error("expected non-nil result")
		}
	})

	t.Run("Modify error", func(t *testing.T) {
		repo := &MockMessageRepository{ModifyErr: fmt.Errorf("modify error")}
		req := mail.ModifyRequest{AddLabels: []string{"STARRED"}}
		_, err := repo.Modify(nil, "msg-id", req)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Get success", func(t *testing.T) {
		msg := &mail.Message{ID: "msg-id", Subject: "Test"}
		repo := &MockMessageRepository{Message: msg}

		result, err := repo.Get(nil, "msg-id")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "msg-id" {
			t.Errorf("expected ID 'msg-id', got %s", result.ID)
		}
	})

	t.Run("Get error", func(t *testing.T) {
		repo := &MockMessageRepository{GetErr: fmt.Errorf("not found")}

		_, err := repo.Get(nil, "nonexistent")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

// =============================================================================
// Comprehensive Error and Edge Case Tests for Mail Compose Commands
// =============================================================================

func TestRunMailSend_InvalidRecipients(t *testing.T) {
	mockRepo := &MockMessageRepository{
		SendResult: &mail.Message{ID: "sent-id"},
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "sender@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	tests := []struct {
		name        string
		to          []string
		cc          []string
		bcc         []string
		expectError bool
	}{
		{
			name:        "invalid to address",
			to:          []string{"notanemail"},
			expectError: true,
		},
		{
			name:        "invalid cc address",
			to:          []string{"valid@example.com"},
			cc:          []string{"invalid"},
			expectError: true,
		},
		{
			name:        "invalid bcc address",
			to:          []string{"valid@example.com"},
			bcc:         []string{"@example.com"},
			expectError: true,
		},
		{
			name:        "all valid addresses",
			to:          []string{"user1@example.com", "user2@example.com"},
			cc:          []string{"cc@example.com"},
			bcc:         []string{"bcc@example.com"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origTo := mailSendTo
			origCc := mailSendCc
			origBcc := mailSendBcc
			origSubject := mailSendSubject
			origBody := mailSendBody

			mailSendTo = tt.to
			mailSendCc = tt.cc
			mailSendBcc = tt.bcc
			mailSendSubject = "Test"
			mailSendBody = "Test body"

			defer func() {
				mailSendTo = origTo
				mailSendCc = origCc
				mailSendBcc = origBcc
				mailSendSubject = origSubject
				mailSendBody = origBody
			}()

			cmd := &cobra.Command{Use: "test"}
			var buf bytes.Buffer
			cmd.SetOut(&buf)

			err := runMailSend(cmd, []string{})

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestRunMailSend_RepositoryError(t *testing.T) {
	mockRepo := &MockMessageRepository{
		SendErr: fmt.Errorf("API error: quota exceeded"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "sender@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTo := mailSendTo
	origSubject := mailSendSubject
	origBody := mailSendBody
	mailSendTo = []string{"recipient@example.com"}
	mailSendSubject = "Test"
	mailSendBody = "Test body"
	defer func() {
		mailSendTo = origTo
		mailSendSubject = origSubject
		mailSendBody = origBody
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailSend(cmd, []string{})
	if err == nil {
		t.Error("expected error from repository, got nil")
	}
	if !contains(err.Error(), "quota exceeded") {
		t.Errorf("expected error to contain 'quota exceeded', got: %v", err)
	}
}

func TestRunMailSend_HTMLContent(t *testing.T) {
	mockRepo := &MockMessageRepository{
		SendResult: &mail.Message{
			ID:       "sent-id",
			ThreadID: "thread-id",
		},
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "sender@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTo := mailSendTo
	origSubject := mailSendSubject
	origBody := mailSendBody
	origHTML := mailSendHTML
	mailSendTo = []string{"recipient@example.com"}
	mailSendSubject = "HTML Test"
	mailSendBody = "<h1>Hello</h1><p>World</p>"
	mailSendHTML = true
	defer func() {
		mailSendTo = origTo
		mailSendSubject = origSubject
		mailSendBody = origBody
		mailSendHTML = origHTML
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailSend(cmd, []string{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunMailReply_GetOriginalError(t *testing.T) {
	mockRepo := &MockMessageRepository{
		GetErr: fmt.Errorf("message not found"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "me@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origBody := mailReplyBody
	mailReplyBody = "Reply text"
	defer func() {
		mailReplyBody = origBody
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailReply(cmd, []string{"nonexistent-id"})
	if err == nil {
		t.Error("expected error when getting original message, got nil")
	}
	if !contains(err.Error(), "message not found") {
		t.Errorf("expected error to contain 'message not found', got: %v", err)
	}
}

func TestRunMailReply_RepositoryError(t *testing.T) {
	originalMsg := &mail.Message{
		ID:      "original-id",
		From:    "sender@example.com",
		To:      []string{"me@example.com"},
		Subject: "Original Subject",
	}

	mockRepo := &MockMessageRepository{
		Message:  originalMsg,
		ReplyErr: fmt.Errorf("network error"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "me@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origBody := mailReplyBody
	mailReplyBody = "Reply text"
	defer func() {
		mailReplyBody = origBody
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailReply(cmd, []string{"original-id"})
	if err == nil {
		t.Error("expected error from repository, got nil")
	}
	if !contains(err.Error(), "network error") {
		t.Errorf("expected error to contain 'network error', got: %v", err)
	}
}

func TestRunMailReply_ReplyAll(t *testing.T) {
	originalMsg := &mail.Message{
		ID:      "original-id",
		From:    "sender@example.com",
		To:      []string{"me@example.com", "other@example.com"},
		Cc:      []string{"cc@example.com"},
		Subject: "Original Subject",
	}

	mockRepo := &MockMessageRepository{
		Message: originalMsg,
		ReplyResult: &mail.Message{
			ID:       "reply-id",
			ThreadID: "thread-id",
		},
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "me@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origBody := mailReplyBody
	origAll := mailReplyAll
	mailReplyBody = "Reply to all"
	mailReplyAll = true
	defer func() {
		mailReplyBody = origBody
		mailReplyAll = origAll
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailReply(cmd, []string{"original-id"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "reply-id") {
		t.Errorf("expected output to contain reply ID, got: %s", output)
	}
}

func TestRunMailForward_InvalidRecipients(t *testing.T) {
	mockRepo := &MockMessageRepository{
		ForwardResult: &mail.Message{ID: "forward-id"},
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "me@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	tests := []struct {
		name        string
		to          []string
		expectError bool
	}{
		{
			name:        "invalid email address",
			to:          []string{"notanemail"},
			expectError: true,
		},
		{
			name:        "valid email addresses",
			to:          []string{"user1@example.com", "user2@example.com"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origTo := mailForwardTo
			origBody := mailForwardBody
			mailForwardTo = tt.to
			mailForwardBody = "FYI"
			defer func() {
				mailForwardTo = origTo
				mailForwardBody = origBody
			}()

			cmd := &cobra.Command{Use: "test"}
			var buf bytes.Buffer
			cmd.SetOut(&buf)

			err := runMailForward(cmd, []string{"msg-id"})

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestRunMailForward_RepositoryError(t *testing.T) {
	mockRepo := &MockMessageRepository{
		ForwardErr: fmt.Errorf("service unavailable"),
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "me@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTo := mailForwardTo
	origBody := mailForwardBody
	mailForwardTo = []string{"colleague@example.com"}
	mailForwardBody = "FYI"
	defer func() {
		mailForwardTo = origTo
		mailForwardBody = origBody
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailForward(cmd, []string{"msg-id"})
	if err == nil {
		t.Error("expected error from repository, got nil")
	}
	if !contains(err.Error(), "service unavailable") {
		t.Errorf("expected error to contain 'service unavailable', got: %v", err)
	}
}

func TestRunMailForward_EmptyBody(t *testing.T) {
	mockRepo := &MockMessageRepository{
		ForwardResult: &mail.Message{
			ID:       "forward-id",
			ThreadID: "thread-id",
		},
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "me@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTo := mailForwardTo
	origBody := mailForwardBody
	mailForwardTo = []string{"colleague@example.com"}
	mailForwardBody = "" // Empty intro message is allowed
	defer func() {
		mailForwardTo = origTo
		mailForwardBody = origBody
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailForward(cmd, []string{"msg-id"})
	if err != nil {
		t.Errorf("unexpected error with empty body: %v", err)
	}
}

func TestRunMailSend_MultipleRecipients(t *testing.T) {
	mockRepo := &MockMessageRepository{
		SendResult: &mail.Message{
			ID:       "sent-id",
			ThreadID: "thread-id",
		},
	}

	deps := &Dependencies{
		AccountService: &MockAccountService{
			Account:      &accountuc.Account{Alias: "test", Email: "sender@example.com"},
			TokenManager: &MockTokenManager{},
		},
		RepoFactory: &MockRepositoryFactory{
			MessageRepo: mockRepo,
		},
	}

	SetDependencies(deps)
	defer ResetDependencies()

	origTo := mailSendTo
	origCc := mailSendCc
	origBcc := mailSendBcc
	origSubject := mailSendSubject
	origBody := mailSendBody
	mailSendTo = []string{"user1@example.com", "user2@example.com", "user3@example.com"}
	mailSendCc = []string{"cc1@example.com", "cc2@example.com"}
	mailSendBcc = []string{"bcc@example.com"}
	mailSendSubject = "Multi-recipient test"
	mailSendBody = "Message to multiple recipients"
	defer func() {
		mailSendTo = origTo
		mailSendCc = origCc
		mailSendBcc = origBcc
		mailSendSubject = origSubject
		mailSendBody = origBody
	}()

	cmd := &cobra.Command{Use: "test"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := runMailSend(cmd, []string{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "sent-id") {
		t.Errorf("expected output to contain sent message ID, got: %s", output)
	}
}

func TestMailSendCmd_HTMLFlag(t *testing.T) {
	flag := mailSendCmd.Flag("html")
	if flag == nil {
		t.Fatal("expected --html flag to be set")
	}

	// Default should be false
	if flag.DefValue != "false" {
		t.Errorf("expected default html to be 'false', got '%s'", flag.DefValue)
	}
}

func TestMailReplyCmd_AllFlag(t *testing.T) {
	flag := mailReplyCmd.Flag("all")
	if flag == nil {
		t.Fatal("expected --all flag to be set")
	}

	// Default should be false
	if flag.DefValue != "false" {
		t.Errorf("expected default all to be 'false', got '%s'", flag.DefValue)
	}
}
