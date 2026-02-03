package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// TestParseHeaders tests the header parsing function.
func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name            string
		headers         []*gmail.MessagePartHeader
		expectedFrom    string
		expectedTo      string
		expectedSubject string
		expectedDate    string
	}{
		{
			name: "basic headers",
			headers: []*gmail.MessagePartHeader{
				{Name: "From", Value: "sender@example.com"},
				{Name: "To", Value: "recipient@example.com"},
				{Name: "Subject", Value: "Test Subject"},
				{Name: "Date", Value: "Mon, 2 Jan 2006 15:04:05 -0700"},
			},
			expectedFrom:    "sender@example.com",
			expectedTo:      "recipient@example.com",
			expectedSubject: "Test Subject",
			expectedDate:    "2006-01-02",
		},
		{
			name: "multiple recipients",
			headers: []*gmail.MessagePartHeader{
				{Name: "From", Value: "sender@example.com"},
				{Name: "To", Value: "one@example.com, two@example.com"},
				{Name: "Subject", Value: "Multi Recipient"},
			},
			expectedFrom:    "sender@example.com",
			expectedTo:      "one@example.com, two@example.com",
			expectedSubject: "Multi Recipient",
		},
		{
			name: "case insensitive headers",
			headers: []*gmail.MessagePartHeader{
				{Name: "from", Value: "lower@example.com"},
				{Name: "TO", Value: "upper@example.com"},
				{Name: "SUBJECT", Value: "Caps Subject"},
			},
			expectedFrom:    "lower@example.com",
			expectedTo:      "upper@example.com",
			expectedSubject: "Caps Subject",
		},
		{
			name:            "empty headers",
			headers:         []*gmail.MessagePartHeader{},
			expectedFrom:    "",
			expectedTo:      "",
			expectedSubject: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, to, subject, date := parseHeaders(tt.headers)

			if from != tt.expectedFrom {
				t.Errorf("from = %q, want %q", from, tt.expectedFrom)
			}
			if to != tt.expectedTo {
				t.Errorf("to = %q, want %q", to, tt.expectedTo)
			}
			if subject != tt.expectedSubject {
				t.Errorf("subject = %q, want %q", subject, tt.expectedSubject)
			}
			if tt.expectedDate != "" {
				if date.Format("2006-01-02") != tt.expectedDate {
					t.Errorf("date = %v, want date string containing %q", date, tt.expectedDate)
				}
			}
		})
	}
}

// TestGmailMessageToDomain tests conversion from Gmail API message to domain message.
func TestGmailMessageToDomain(t *testing.T) {
	tests := []struct {
		name     string
		gmailMsg *gmail.Message
		want     *mail.Message
	}{
		{
			name: "basic message",
			gmailMsg: &gmail.Message{
				Id:       "msg123",
				ThreadId: "thread456",
				Snippet:  "This is a preview...",
				LabelIds: []string{"INBOX", "UNREAD"},
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{Name: "From", Value: "sender@example.com"},
						{Name: "To", Value: "recipient@example.com"},
						{Name: "Subject", Value: "Test Subject"},
					},
					Body: &gmail.MessagePartBody{
						Data: base64.URLEncoding.EncodeToString([]byte("Hello, World!")),
					},
				},
			},
			want: &mail.Message{
				ID:        "msg123",
				ThreadID:  "thread456",
				From:      "sender@example.com",
				To:        []string{"recipient@example.com"},
				Subject:   "Test Subject",
				Body:      "Hello, World!",
				Snippet:   "This is a preview...",
				Labels:    []string{"INBOX", "UNREAD"},
				IsRead:    false,
				IsStarred: false,
			},
		},
		{
			name: "read and starred message",
			gmailMsg: &gmail.Message{
				Id:       "msg789",
				ThreadId: "thread101",
				LabelIds: []string{"INBOX", "STARRED"},
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{Name: "From", Value: "another@example.com"},
						{Name: "To", Value: "me@example.com"},
						{Name: "Subject", Value: "Starred Message"},
					},
				},
			},
			want: &mail.Message{
				ID:        "msg789",
				ThreadID:  "thread101",
				From:      "another@example.com",
				To:        []string{"me@example.com"},
				Subject:   "Starred Message",
				Labels:    []string{"INBOX", "STARRED"},
				IsRead:    true,
				IsStarred: true,
			},
		},
		{
			name: "multipart message",
			gmailMsg: &gmail.Message{
				Id:       "multipart123",
				ThreadId: "thread789",
				LabelIds: []string{"UNREAD"}, // Include UNREAD label so IsRead is false
				Payload: &gmail.MessagePart{
					MimeType: "multipart/alternative",
					Headers: []*gmail.MessagePartHeader{
						{Name: "From", Value: "html@example.com"},
						{Name: "To", Value: "reader@example.com"},
						{Name: "Subject", Value: "HTML Email"},
					},
					Parts: []*gmail.MessagePart{
						{
							MimeType: "text/plain",
							Body: &gmail.MessagePartBody{
								Data: base64.URLEncoding.EncodeToString([]byte("Plain text content")),
							},
						},
						{
							MimeType: "text/html",
							Body: &gmail.MessagePartBody{
								Data: base64.URLEncoding.EncodeToString([]byte("<p>HTML content</p>")),
							},
						},
					},
				},
			},
			want: &mail.Message{
				ID:       "multipart123",
				ThreadID: "thread789",
				From:     "html@example.com",
				To:       []string{"reader@example.com"},
				Subject:  "HTML Email",
				Body:     "Plain text content",
				BodyHTML: "<p>HTML content</p>",
				IsRead:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gmailMessageToDomain(tt.gmailMsg)

			if got.ID != tt.want.ID {
				t.Errorf("ID = %q, want %q", got.ID, tt.want.ID)
			}
			if got.ThreadID != tt.want.ThreadID {
				t.Errorf("ThreadID = %q, want %q", got.ThreadID, tt.want.ThreadID)
			}
			if got.From != tt.want.From {
				t.Errorf("From = %q, want %q", got.From, tt.want.From)
			}
			if got.Subject != tt.want.Subject {
				t.Errorf("Subject = %q, want %q", got.Subject, tt.want.Subject)
			}
			if got.IsRead != tt.want.IsRead {
				t.Errorf("IsRead = %v, want %v", got.IsRead, tt.want.IsRead)
			}
			if got.IsStarred != tt.want.IsStarred {
				t.Errorf("IsStarred = %v, want %v", got.IsStarred, tt.want.IsStarred)
			}
			if tt.want.Body != "" && got.Body != tt.want.Body {
				t.Errorf("Body = %q, want %q", got.Body, tt.want.Body)
			}
			if tt.want.BodyHTML != "" && got.BodyHTML != tt.want.BodyHTML {
				t.Errorf("BodyHTML = %q, want %q", got.BodyHTML, tt.want.BodyHTML)
			}
		})
	}
}

// TestBuildMimeMessage tests MIME message building.
func TestBuildMimeMessage(t *testing.T) {
	tests := []struct {
		name        string
		msg         *mail.Message
		wantHeaders []string
	}{
		{
			name: "basic message",
			msg: &mail.Message{
				From:    "sender@example.com",
				To:      []string{"recipient@example.com"},
				Subject: "Test Subject",
				Body:    "Hello, World!",
			},
			wantHeaders: []string{
				"From: sender@example.com",
				"To: recipient@example.com",
				"Subject: Test Subject",
				"MIME-Version: 1.0",
				"Content-Type: text/plain; charset=\"utf-8\"",
			},
		},
		{
			name: "multiple recipients with cc and bcc",
			msg: &mail.Message{
				From:    "sender@example.com",
				To:      []string{"one@example.com", "two@example.com"},
				Cc:      []string{"cc@example.com"},
				Bcc:     []string{"bcc@example.com"},
				Subject: "Multi Recipient",
				Body:    "Content",
			},
			wantHeaders: []string{
				"From: sender@example.com",
				"To: one@example.com, two@example.com",
				"Cc: cc@example.com",
				"Bcc: bcc@example.com",
				"Subject: Multi Recipient",
			},
		},
		{
			name: "html message",
			msg: &mail.Message{
				From:     "sender@example.com",
				To:       []string{"recipient@example.com"},
				Subject:  "HTML Email",
				BodyHTML: "<p>Hello, World!</p>",
			},
			wantHeaders: []string{
				"Content-Type: text/html; charset=\"utf-8\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildMimeMessage(tt.msg)
			gotStr := string(got)

			for _, header := range tt.wantHeaders {
				if !strings.Contains(gotStr, header) {
					t.Errorf("MIME message missing header %q\nGot:\n%s", header, gotStr)
				}
			}
		})
	}
}

// TestMapGmailError tests error mapping from Gmail API errors to domain errors.
func TestMapGmailError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    error
	}{
		{
			name:       "404 maps to message not found",
			statusCode: http.StatusNotFound,
			wantErr:    mail.ErrMessageNotFound,
		},
		{
			name:       "400 returns validation error",
			statusCode: http.StatusBadRequest,
			wantErr:    ErrBadRequest,
		},
		{
			name:       "429 returns rate limit error",
			statusCode: http.StatusTooManyRequests,
			wantErr:    ErrRateLimited,
		},
		{
			name:       "500 returns temporary error",
			statusCode: http.StatusInternalServerError,
			wantErr:    ErrTemporary,
		},
		{
			name:       "503 returns temporary error",
			statusCode: http.StatusServiceUnavailable,
			wantErr:    ErrTemporary,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mapGmailError(tt.statusCode, "test error")

			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr.Error()) {
				t.Errorf("error = %v, want error containing %v", err, tt.wantErr)
			}
		})
	}
}

// TestGmailRepository_List tests the List method with a mock server.
func TestGmailRepository_List(t *testing.T) {
	// Create a mock Gmail API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/gmail/v1/users/me/messages" {
			response := gmail.ListMessagesResponse{
				Messages: []*gmail.Message{
					{Id: "msg1", ThreadId: "thread1"},
					{Id: "msg2", ThreadId: "thread2"},
				},
				NextPageToken:      "next_token",
				ResultSizeEstimate: 100,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		// Handle individual message fetches
		if strings.HasPrefix(r.URL.Path, "/gmail/v1/users/me/messages/") {
			msgID := strings.TrimPrefix(r.URL.Path, "/gmail/v1/users/me/messages/")
			response := gmail.Message{
				Id:       msgID,
				ThreadId: "thread1",
				Snippet:  "Test snippet",
				LabelIds: []string{"INBOX"},
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{Name: "From", Value: "test@example.com"},
						{Name: "To", Value: "me@example.com"},
						{Name: "Subject", Value: "Test Message"},
					},
				},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	ctx := context.Background()
	service, err := gmail.NewService(ctx, option.WithEndpoint(server.URL), option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("failed to create Gmail service: %v", err)
	}

	repo := &GmailRepository{
		service: service,
		userID:  "me",
	}

	result, err := repo.List(ctx, mail.ListOptions{MaxResults: 10})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if result.NextPageToken != "next_token" {
		t.Errorf("NextPageToken = %q, want %q", result.NextPageToken, "next_token")
	}
	if result.Total != 100 {
		t.Errorf("Total = %d, want %d", result.Total, 100)
	}
	if len(result.Items) != 2 {
		t.Errorf("Items count = %d, want %d", len(result.Items), 2)
	}
}

// TestGmailRepository_Get tests the Get method.
func TestGmailRepository_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/gmail/v1/users/me/messages/msg123" {
			response := gmail.Message{
				Id:       "msg123",
				ThreadId: "thread456",
				Snippet:  "Message snippet",
				LabelIds: []string{"INBOX", "UNREAD"},
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{Name: "From", Value: "sender@example.com"},
						{Name: "To", Value: "recipient@example.com"},
						{Name: "Subject", Value: "Test Subject"},
					},
					Body: &gmail.MessagePartBody{
						Data: base64.URLEncoding.EncodeToString([]byte("Message body")),
					},
				},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	ctx := context.Background()
	service, err := gmail.NewService(ctx, option.WithEndpoint(server.URL), option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("failed to create Gmail service: %v", err)
	}

	repo := &GmailRepository{
		service: service,
		userID:  "me",
	}

	msg, err := repo.Get(ctx, "msg123")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if msg.ID != "msg123" {
		t.Errorf("ID = %q, want %q", msg.ID, "msg123")
	}
	if msg.ThreadID != "thread456" {
		t.Errorf("ThreadID = %q, want %q", msg.ThreadID, "thread456")
	}
	if msg.Subject != "Test Subject" {
		t.Errorf("Subject = %q, want %q", msg.Subject, "Test Subject")
	}
}

// TestGmailRepository_GetNotFound tests Get with non-existent message.
func TestGmailRepository_GetNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    404,
				"message": "Requested entity was not found.",
			},
		})
	}))
	defer server.Close()

	ctx := context.Background()
	service, err := gmail.NewService(ctx, option.WithEndpoint(server.URL), option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("failed to create Gmail service: %v", err)
	}

	repo := &GmailRepository{
		service: service,
		userID:  "me",
	}

	_, err = repo.Get(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent message, got nil")
	}
}

// TestGmailRepository_Send tests the Send method.
func TestGmailRepository_Send(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/gmail/v1/users/me/messages/send" {
			response := gmail.Message{
				Id:       "sent123",
				ThreadId: "thread789",
				LabelIds: []string{"SENT"},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		// Handle the subsequent Get request to fetch the sent message details
		if r.Method == "GET" && r.URL.Path == "/gmail/v1/users/me/messages/sent123" {
			response := gmail.Message{
				Id:       "sent123",
				ThreadId: "thread789",
				LabelIds: []string{"SENT"},
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{Name: "From", Value: "sender@example.com"},
						{Name: "To", Value: "recipient@example.com"},
						{Name: "Subject", Value: "Test Subject"},
					},
					Body: &gmail.MessagePartBody{
						Data: base64.URLEncoding.EncodeToString([]byte("Test Body")),
					},
				},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	ctx := context.Background()
	service, err := gmail.NewService(ctx, option.WithEndpoint(server.URL), option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("failed to create Gmail service: %v", err)
	}

	repo := &GmailRepository{
		service: service,
		userID:  "me",
	}

	msg := &mail.Message{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	sent, err := repo.Send(ctx, msg)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	if sent.ID != "sent123" {
		t.Errorf("sent ID = %q, want %q", sent.ID, "sent123")
	}
}

// TestGmailRepository_Trash tests the Trash method.
func TestGmailRepository_Trash(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/gmail/v1/users/me/messages/msg123/trash" {
			response := gmail.Message{
				Id:       "msg123",
				LabelIds: []string{"TRASH"},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	ctx := context.Background()
	service, err := gmail.NewService(ctx, option.WithEndpoint(server.URL), option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("failed to create Gmail service: %v", err)
	}

	repo := &GmailRepository{
		service: service,
		userID:  "me",
	}

	err = repo.Trash(ctx, "msg123")
	if err != nil {
		t.Fatalf("Trash failed: %v", err)
	}
}

// TestGmailRepository_Modify tests the Modify method.
func TestGmailRepository_Modify(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/gmail/v1/users/me/messages/msg123/modify" {
			response := gmail.Message{
				Id:       "msg123",
				LabelIds: []string{"INBOX", "STARRED"},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	ctx := context.Background()
	service, err := gmail.NewService(ctx, option.WithEndpoint(server.URL), option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("failed to create Gmail service: %v", err)
	}

	repo := &GmailRepository{
		service: service,
		userID:  "me",
	}

	msg, err := repo.Modify(ctx, "msg123", mail.ModifyRequest{
		AddLabels:    []string{"STARRED"},
		RemoveLabels: []string{"UNREAD"},
	})
	if err != nil {
		t.Fatalf("Modify failed: %v", err)
	}

	if msg.ID != "msg123" {
		t.Errorf("modified message ID = %q, want %q", msg.ID, "msg123")
	}
}

// TestRetryWithBackoff tests the retry mechanism.
func TestRetryWithBackoff(t *testing.T) {
	attempts := 0
	ctx := context.Background()

	result, err := retryWithBackoff(ctx, 3, 10*time.Millisecond, func() (string, error) {
		attempts++
		if attempts < 3 {
			return "", ErrTemporary
		}
		return "success", nil
	})

	if err != nil {
		t.Fatalf("retryWithBackoff failed: %v", err)
	}
	if result != "success" {
		t.Errorf("result = %q, want %q", result, "success")
	}
	if attempts != 3 {
		t.Errorf("attempts = %d, want %d", attempts, 3)
	}
}

// TestRetryWithBackoffExhausted tests retry exhaustion.
func TestRetryWithBackoffExhausted(t *testing.T) {
	ctx := context.Background()

	_, err := retryWithBackoff(ctx, 3, 10*time.Millisecond, func() (string, error) {
		return "", ErrTemporary
	})

	if err == nil {
		t.Fatal("expected error after retries exhausted, got nil")
	}
}

// TestRetryWithBackoffContextCancelled tests context cancellation during retry.
func TestRetryWithBackoffContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	attempts := 0
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err := retryWithBackoff(ctx, 10, 30*time.Millisecond, func() (string, error) {
		attempts++
		return "", ErrTemporary
	})

	if err == nil {
		t.Fatal("expected error after context cancelled, got nil")
	}
}
