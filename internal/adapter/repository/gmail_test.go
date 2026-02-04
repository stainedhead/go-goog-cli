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

// =============================================================================
// Tests Using TestServer Infrastructure
// =============================================================================

// TestGmailRepository_ListWithTestServer tests List using the TestServer.
func TestGmailRepository_ListWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	// Set up mock data
	msg1 := MockMessageResponse("msg1", "thread1", "Subject 1", "alice@example.com", "bob@example.com", "Hello")
	msg2 := MockMessageResponse("msg2", "thread2", "Subject 2", "charlie@example.com", "bob@example.com", "World")

	// Configure handlers
	ts.MessageListHandler = func(w http.ResponseWriter, r *http.Request) {
		WriteJSONResponse(w, MockMessageListResponse(
			[]*gmail.Message{
				{Id: "msg1", ThreadId: "thread1"},
				{Id: "msg2", ThreadId: "thread2"},
			},
			"next_page_token",
			100,
		))
	}

	ts.MessageGetHandler = func(w http.ResponseWriter, r *http.Request, msgID string) {
		switch msgID {
		case "msg1":
			WriteJSONResponse(w, msg1)
		case "msg2":
			WriteJSONResponse(w, msg2)
		default:
			WriteErrorResponse(w, http.StatusNotFound, "message not found")
		}
	}

	repo := ts.GmailRepository(t)
	ctx := context.Background()

	result, err := repo.List(ctx, mail.ListOptions{MaxResults: 10})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if result.NextPageToken != "next_page_token" {
		t.Errorf("NextPageToken = %q, want %q", result.NextPageToken, "next_page_token")
	}
	if result.Total != 100 {
		t.Errorf("Total = %d, want %d", result.Total, 100)
	}
	if len(result.Items) != 2 {
		t.Errorf("Items count = %d, want %d", len(result.Items), 2)
	}
	if result.Items[0].Subject != "Subject 1" {
		t.Errorf("Items[0].Subject = %q, want %q", result.Items[0].Subject, "Subject 1")
	}
}

// TestGmailRepository_GetWithTestServer tests Get using the TestServer.
func TestGmailRepository_GetWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	expectedMsg := MockMessageResponse("msg123", "thread456", "Test Subject", "sender@example.com", "recipient@example.com", "Message body content")
	expectedMsg.LabelIds = []string{"INBOX", "UNREAD"}

	ts.MessageGetHandler = func(w http.ResponseWriter, r *http.Request, msgID string) {
		if msgID == "msg123" {
			WriteJSONResponse(w, expectedMsg)
		} else {
			WriteErrorResponse(w, http.StatusNotFound, "message not found")
		}
	}

	repo := ts.GmailRepository(t)
	ctx := context.Background()

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
	if msg.From != "sender@example.com" {
		t.Errorf("From = %q, want %q", msg.From, "sender@example.com")
	}
	if msg.IsRead {
		t.Error("IsRead = true, want false (message has UNREAD label)")
	}
}

// TestGmailRepository_SendWithTestServer tests Send using the TestServer.
func TestGmailRepository_SendWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	var sentMessageRaw string

	ts.MessageSendHandler = func(w http.ResponseWriter, r *http.Request) {
		// Decode the request body to verify the message
		var msg gmail.Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}
		sentMessageRaw = msg.Raw

		// Return the sent message with an ID
		WriteJSONResponse(w, &gmail.Message{
			Id:       "sent123",
			ThreadId: "thread789",
			LabelIds: []string{"SENT"},
		})
	}

	ts.MessageGetHandler = func(w http.ResponseWriter, r *http.Request, msgID string) {
		if msgID == "sent123" {
			WriteJSONResponse(w, MockMessageResponse(
				"sent123", "thread789", "Test Subject", "sender@example.com", "recipient@example.com", "Test Body",
			))
		} else {
			WriteErrorResponse(w, http.StatusNotFound, "message not found")
		}
	}

	repo := ts.GmailRepository(t)
	ctx := context.Background()

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

	// Verify raw message was actually sent
	if sentMessageRaw == "" {
		t.Error("expected raw message to be sent, but it was empty")
	}
}

// TestGmailRepository_TrashWithTestServer tests Trash using the TestServer.
func TestGmailRepository_TrashWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	trashedID := ""
	ts.MessageTrashHandler = func(w http.ResponseWriter, r *http.Request, msgID string) {
		trashedID = msgID
		WriteJSONResponse(w, &gmail.Message{
			Id:       msgID,
			LabelIds: []string{"TRASH"},
		})
	}

	repo := ts.GmailRepository(t)
	ctx := context.Background()

	err := repo.Trash(ctx, "msg123")
	if err != nil {
		t.Fatalf("Trash failed: %v", err)
	}

	if trashedID != "msg123" {
		t.Errorf("trashedID = %q, want %q", trashedID, "msg123")
	}
}

// TestGmailRepository_ListLabelsWithTestServer tests ListLabels using the TestServer.
func TestGmailRepository_ListLabelsWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.LabelListHandler = StaticLabelListHandler([]*gmail.Label{
		MockLabelResponse("INBOX", "INBOX", "system"),
		MockLabelResponse("SENT", "SENT", "system"),
		MockLabelResponse("Label_1", "Work", "user"),
		MockLabelResponse("Label_2", "Personal", "user"),
	})

	repo := NewGmailLabelRepository(ts.GmailRepository(t))
	ctx := context.Background()

	labels, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List labels failed: %v", err)
	}

	if len(labels) != 4 {
		t.Errorf("labels count = %d, want %d", len(labels), 4)
	}

	// Verify first label
	if labels[0].ID != "INBOX" {
		t.Errorf("labels[0].ID = %q, want %q", labels[0].ID, "INBOX")
	}
	if labels[0].Name != "INBOX" {
		t.Errorf("labels[0].Name = %q, want %q", labels[0].Name, "INBOX")
	}
	if labels[0].Type != "system" {
		t.Errorf("labels[0].Type = %q, want %q", labels[0].Type, "system")
	}

	// Verify user label
	if labels[2].ID != "Label_1" {
		t.Errorf("labels[2].ID = %q, want %q", labels[2].ID, "Label_1")
	}
	if labels[2].Name != "Work" {
		t.Errorf("labels[2].Name = %q, want %q", labels[2].Name, "Work")
	}
	if labels[2].Type != "user" {
		t.Errorf("labels[2].Type = %q, want %q", labels[2].Type, "user")
	}
}

// TestGmailRepository_GetNotFoundWithTestServer tests Get for non-existent message.
func TestGmailRepository_GetNotFoundWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.MessageGetHandler = func(w http.ResponseWriter, r *http.Request, msgID string) {
		WriteErrorResponse(w, http.StatusNotFound, "Requested entity was not found.")
	}

	repo := ts.GmailRepository(t)
	ctx := context.Background()

	_, err := repo.Get(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent message, got nil")
	}
}

// TestGmailRepository_RateLimitedWithTestServer tests rate limit handling.
func TestGmailRepository_RateLimitedWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.MessageGetHandler = func(w http.ResponseWriter, r *http.Request, msgID string) {
		WriteErrorResponse(w, http.StatusTooManyRequests, "Rate limit exceeded")
	}

	repo := ts.GmailRepository(t)
	ctx := context.Background()

	_, err := repo.Get(ctx, "msg123")
	if err == nil {
		t.Fatal("expected error for rate limited request, got nil")
	}

	if !strings.Contains(err.Error(), ErrRateLimited.Error()) {
		t.Errorf("error = %v, want error containing %v", err, ErrRateLimited)
	}
}

// =============================================================================
// Additional Message Operations Tests
// =============================================================================

// TestGmailRepository_UntrashWithTestServer tests restoring a message from trash.
func TestGmailRepository_UntrashWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	untrashedID := ""
	ts.MessageUntrashHandler = func(w http.ResponseWriter, r *http.Request, msgID string) {
		untrashedID = msgID
		if r.Method != "POST" {
			WriteErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		WriteJSONResponse(w, &gmail.Message{
			Id:       msgID,
			LabelIds: []string{"INBOX"},
		})
	}

	repo := ts.GmailRepository(t)
	ctx := context.Background()

	err := repo.Untrash(ctx, "msg123")
	if err != nil {
		t.Fatalf("Untrash failed: %v", err)
	}

	if untrashedID != "msg123" {
		t.Errorf("untrashedID = %q, want %q", untrashedID, "msg123")
	}
}

// TestGmailRepository_DeleteWithTestServer tests permanently deleting a message.
func TestGmailRepository_DeleteWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	deletedID := ""
	ts.MessageDeleteHandler = func(w http.ResponseWriter, r *http.Request, msgID string) {
		deletedID = msgID
		if r.Method != "DELETE" {
			WriteErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

	repo := ts.GmailRepository(t)
	ctx := context.Background()

	err := repo.Delete(ctx, "msg456")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if deletedID != "msg456" {
		t.Errorf("deletedID = %q, want %q", deletedID, "msg456")
	}
}

// TestGmailRepository_ArchiveWithTestServer tests archiving a message.
func TestGmailRepository_ArchiveWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	var modifyRequest *gmail.ModifyMessageRequest
	ts.MessageModifyHandler = func(w http.ResponseWriter, r *http.Request, msgID string) {
		if err := json.NewDecoder(r.Body).Decode(&modifyRequest); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}
		WriteJSONResponse(w, &gmail.Message{
			Id:       msgID,
			LabelIds: []string{}, // INBOX removed
		})
	}

	repo := ts.GmailRepository(t)
	ctx := context.Background()

	err := repo.Archive(ctx, "msg789")
	if err != nil {
		t.Fatalf("Archive failed: %v", err)
	}

	// Verify the modify request removed INBOX label
	if modifyRequest == nil {
		t.Fatal("expected modify request, got nil")
	}
	foundInbox := false
	for _, label := range modifyRequest.RemoveLabelIds {
		if label == "INBOX" {
			foundInbox = true
			break
		}
	}
	if !foundInbox {
		t.Error("expected INBOX in RemoveLabelIds")
	}
}

// TestGmailRepository_ModifyLabelsWithTestServer tests modifying message labels.
func TestGmailRepository_ModifyLabelsWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	var modifyRequest *gmail.ModifyMessageRequest
	ts.MessageModifyHandler = func(w http.ResponseWriter, r *http.Request, msgID string) {
		if err := json.NewDecoder(r.Body).Decode(&modifyRequest); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}
		WriteJSONResponse(w, &gmail.Message{
			Id:       msgID,
			LabelIds: []string{"INBOX", "STARRED", "Label_123"},
		})
	}

	repo := ts.GmailRepository(t)
	ctx := context.Background()

	msg, err := repo.Modify(ctx, "msg123", mail.ModifyRequest{
		AddLabels:    []string{"STARRED", "Label_123"},
		RemoveLabels: []string{"UNREAD"},
	})
	if err != nil {
		t.Fatalf("Modify failed: %v", err)
	}

	if msg.ID != "msg123" {
		t.Errorf("msg.ID = %q, want %q", msg.ID, "msg123")
	}

	// Verify request was made correctly
	if len(modifyRequest.AddLabelIds) != 2 {
		t.Errorf("AddLabelIds length = %d, want %d", len(modifyRequest.AddLabelIds), 2)
	}
	if len(modifyRequest.RemoveLabelIds) != 1 {
		t.Errorf("RemoveLabelIds length = %d, want %d", len(modifyRequest.RemoveLabelIds), 1)
	}
}

// TestGmailRepository_ModifyLabelsNotFound tests modifying non-existent message.
func TestGmailRepository_ModifyLabelsNotFound(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.MessageModifyHandler = func(w http.ResponseWriter, r *http.Request, msgID string) {
		WriteErrorResponse(w, http.StatusNotFound, "message not found")
	}

	repo := ts.GmailRepository(t)
	ctx := context.Background()

	_, err := repo.Modify(ctx, "nonexistent", mail.ModifyRequest{
		AddLabels: []string{"STARRED"},
	})
	if err == nil {
		t.Fatal("expected error for non-existent message, got nil")
	}
}

// =============================================================================
// Draft Repository Tests
// =============================================================================

// TestGmailDraftRepository_ListWithTestServer tests listing drafts.
func TestGmailDraftRepository_ListWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.DraftListHandler = func(w http.ResponseWriter, r *http.Request) {
		WriteJSONResponse(w, &gmail.ListDraftsResponse{
			Drafts: []*gmail.Draft{
				{Id: "draft1"},
				{Id: "draft2"},
			},
			NextPageToken:      "next_token",
			ResultSizeEstimate: 10,
		})
	}

	ts.DraftGetHandler = func(w http.ResponseWriter, r *http.Request, draftID string) {
		WriteJSONResponse(w, MockDraftResponse(draftID, "msg_"+draftID, "Draft Subject", "me@example.com", "you@example.com", "Draft body"))
	}

	repo := NewGmailDraftRepository(ts.GmailRepository(t))
	ctx := context.Background()

	result, err := repo.List(ctx, mail.ListOptions{MaxResults: 10})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(result.Items) != 2 {
		t.Errorf("drafts count = %d, want %d", len(result.Items), 2)
	}
	if result.NextPageToken != "next_token" {
		t.Errorf("NextPageToken = %q, want %q", result.NextPageToken, "next_token")
	}
}

// TestGmailDraftRepository_GetWithTestServer tests getting a single draft.
func TestGmailDraftRepository_GetWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.DraftGetHandler = func(w http.ResponseWriter, r *http.Request, draftID string) {
		if draftID == "draft123" {
			WriteJSONResponse(w, MockDraftResponse("draft123", "msg123", "Test Draft", "sender@example.com", "recipient@example.com", "Draft content"))
		} else {
			WriteErrorResponse(w, http.StatusNotFound, "draft not found")
		}
	}

	repo := NewGmailDraftRepository(ts.GmailRepository(t))
	ctx := context.Background()

	draft, err := repo.Get(ctx, "draft123")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if draft.ID != "draft123" {
		t.Errorf("draft.ID = %q, want %q", draft.ID, "draft123")
	}
	if draft.Message == nil {
		t.Fatal("draft.Message should not be nil")
	}
	if draft.Message.Subject != "Test Draft" {
		t.Errorf("draft.Message.Subject = %q, want %q", draft.Message.Subject, "Test Draft")
	}
}

// TestGmailDraftRepository_CreateWithTestServer tests creating a draft.
func TestGmailDraftRepository_CreateWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	var createdDraft *gmail.Draft
	ts.DraftCreateHandler = func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&createdDraft); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}
		WriteJSONResponse(w, &gmail.Draft{
			Id: "new_draft_123",
			Message: &gmail.Message{
				Id: "msg_new_draft",
			},
		})
	}

	ts.DraftGetHandler = func(w http.ResponseWriter, r *http.Request, draftID string) {
		WriteJSONResponse(w, MockDraftResponse("new_draft_123", "msg_new_draft", "New Draft", "me@example.com", "you@example.com", "Draft content"))
	}

	repo := NewGmailDraftRepository(ts.GmailRepository(t))
	ctx := context.Background()

	draft := &mail.Draft{
		Message: &mail.Message{
			From:    "me@example.com",
			To:      []string{"you@example.com"},
			Subject: "New Draft",
			Body:    "Draft content",
		},
	}

	created, err := repo.Create(ctx, draft)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if created.ID != "new_draft_123" {
		t.Errorf("created.ID = %q, want %q", created.ID, "new_draft_123")
	}
}

// TestGmailDraftRepository_UpdateWithTestServer tests updating a draft.
func TestGmailDraftRepository_UpdateWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.DraftUpdateHandler = func(w http.ResponseWriter, r *http.Request, draftID string) {
		WriteJSONResponse(w, &gmail.Draft{
			Id: draftID,
			Message: &gmail.Message{
				Id: "msg_" + draftID,
			},
		})
	}

	ts.DraftGetHandler = func(w http.ResponseWriter, r *http.Request, draftID string) {
		WriteJSONResponse(w, MockDraftResponse(draftID, "msg_"+draftID, "Updated Subject", "me@example.com", "you@example.com", "Updated body"))
	}

	repo := NewGmailDraftRepository(ts.GmailRepository(t))
	ctx := context.Background()

	draft := &mail.Draft{
		ID: "draft123",
		Message: &mail.Message{
			From:    "me@example.com",
			To:      []string{"you@example.com"},
			Subject: "Updated Subject",
			Body:    "Updated body",
		},
	}

	updated, err := repo.Update(ctx, draft)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.ID != "draft123" {
		t.Errorf("updated.ID = %q, want %q", updated.ID, "draft123")
	}
}

// TestGmailDraftRepository_SendWithTestServer tests sending a draft.
func TestGmailDraftRepository_SendWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.DraftSendHandler = func(w http.ResponseWriter, r *http.Request) {
		WriteJSONResponse(w, &gmail.Message{
			Id:       "sent_msg_123",
			ThreadId: "thread_456",
			LabelIds: []string{"SENT"},
		})
	}

	ts.MessageGetHandler = func(w http.ResponseWriter, r *http.Request, msgID string) {
		WriteJSONResponse(w, MockMessageResponse(msgID, "thread_456", "Sent Subject", "me@example.com", "you@example.com", "Sent body"))
	}

	repo := NewGmailDraftRepository(ts.GmailRepository(t))
	ctx := context.Background()

	msg, err := repo.Send(ctx, "draft123")
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	if msg.ID != "sent_msg_123" {
		t.Errorf("msg.ID = %q, want %q", msg.ID, "sent_msg_123")
	}
}

// TestGmailDraftRepository_DeleteWithTestServer tests deleting a draft.
func TestGmailDraftRepository_DeleteWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	deletedID := ""
	ts.DraftDeleteHandler = func(w http.ResponseWriter, r *http.Request, draftID string) {
		deletedID = draftID
		w.WriteHeader(http.StatusNoContent)
	}

	repo := NewGmailDraftRepository(ts.GmailRepository(t))
	ctx := context.Background()

	err := repo.Delete(ctx, "draft123")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if deletedID != "draft123" {
		t.Errorf("deletedID = %q, want %q", deletedID, "draft123")
	}
}

// TestGmailDraftRepository_GetNotFound tests getting a non-existent draft.
func TestGmailDraftRepository_GetNotFound(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.DraftGetHandler = func(w http.ResponseWriter, r *http.Request, draftID string) {
		WriteErrorResponse(w, http.StatusNotFound, "draft not found")
	}

	repo := NewGmailDraftRepository(ts.GmailRepository(t))
	ctx := context.Background()

	_, err := repo.Get(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent draft, got nil")
	}
}

// =============================================================================
// Thread Repository Tests
// =============================================================================

// TestGmailThreadRepository_ListWithTestServer tests listing threads.
func TestGmailThreadRepository_ListWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.ThreadListHandler = func(w http.ResponseWriter, r *http.Request) {
		WriteJSONResponse(w, &gmail.ListThreadsResponse{
			Threads: []*gmail.Thread{
				{Id: "thread1", Snippet: "Thread 1 snippet"},
				{Id: "thread2", Snippet: "Thread 2 snippet"},
			},
			NextPageToken:      "next_token",
			ResultSizeEstimate: 50,
		})
	}

	repo := NewGmailThreadRepository(ts.GmailRepository(t))
	ctx := context.Background()

	result, err := repo.List(ctx, mail.ListOptions{MaxResults: 10})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(result.Items) != 2 {
		t.Errorf("threads count = %d, want %d", len(result.Items), 2)
	}
	if result.Items[0].ID != "thread1" {
		t.Errorf("result.Items[0].ID = %q, want %q", result.Items[0].ID, "thread1")
	}
	if result.Items[0].Snippet != "Thread 1 snippet" {
		t.Errorf("result.Items[0].Snippet = %q, want %q", result.Items[0].Snippet, "Thread 1 snippet")
	}
}

// TestGmailThreadRepository_GetWithTestServer tests getting a thread.
func TestGmailThreadRepository_GetWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.ThreadGetHandler = func(w http.ResponseWriter, r *http.Request, threadID string) {
		if threadID == "thread123" {
			WriteJSONResponse(w, MockThreadResponse("thread123", []*gmail.Message{
				MockMessageResponse("msg1", "thread123", "First message", "alice@example.com", "bob@example.com", "Hello"),
				MockMessageResponse("msg2", "thread123", "Re: First message", "bob@example.com", "alice@example.com", "Hi there"),
			}))
		} else {
			WriteErrorResponse(w, http.StatusNotFound, "thread not found")
		}
	}

	repo := NewGmailThreadRepository(ts.GmailRepository(t))
	ctx := context.Background()

	thread, err := repo.Get(ctx, "thread123")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if thread.ID != "thread123" {
		t.Errorf("thread.ID = %q, want %q", thread.ID, "thread123")
	}
	if len(thread.Messages) != 2 {
		t.Errorf("thread.Messages length = %d, want %d", len(thread.Messages), 2)
	}
	if thread.Messages[0].Subject != "First message" {
		t.Errorf("thread.Messages[0].Subject = %q, want %q", thread.Messages[0].Subject, "First message")
	}
}

// TestGmailThreadRepository_ModifyWithTestServer tests modifying thread labels.
func TestGmailThreadRepository_ModifyWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	var modifyRequest *gmail.ModifyThreadRequest
	ts.ThreadModifyHandler = func(w http.ResponseWriter, r *http.Request, threadID string) {
		if err := json.NewDecoder(r.Body).Decode(&modifyRequest); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}
		WriteJSONResponse(w, &gmail.Thread{
			Id:      threadID,
			Snippet: "Modified thread",
		})
	}

	repo := NewGmailThreadRepository(ts.GmailRepository(t))
	ctx := context.Background()

	thread, err := repo.Modify(ctx, "thread123", mail.ModifyRequest{
		AddLabels:    []string{"STARRED"},
		RemoveLabels: []string{"UNREAD"},
	})
	if err != nil {
		t.Fatalf("Modify failed: %v", err)
	}

	if thread.ID != "thread123" {
		t.Errorf("thread.ID = %q, want %q", thread.ID, "thread123")
	}

	// Verify request
	if len(modifyRequest.AddLabelIds) != 1 || modifyRequest.AddLabelIds[0] != "STARRED" {
		t.Errorf("expected STARRED in AddLabelIds")
	}
}

// TestGmailThreadRepository_TrashWithTestServer tests trashing a thread.
func TestGmailThreadRepository_TrashWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	trashedID := ""
	ts.ThreadTrashHandler = func(w http.ResponseWriter, r *http.Request, threadID string) {
		trashedID = threadID
		WriteJSONResponse(w, &gmail.Thread{Id: threadID})
	}

	repo := NewGmailThreadRepository(ts.GmailRepository(t))
	ctx := context.Background()

	err := repo.Trash(ctx, "thread123")
	if err != nil {
		t.Fatalf("Trash failed: %v", err)
	}

	if trashedID != "thread123" {
		t.Errorf("trashedID = %q, want %q", trashedID, "thread123")
	}
}

// TestGmailThreadRepository_UntrashWithTestServer tests untrashing a thread.
func TestGmailThreadRepository_UntrashWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	untrashedID := ""
	ts.ThreadUntrashHandler = func(w http.ResponseWriter, r *http.Request, threadID string) {
		untrashedID = threadID
		WriteJSONResponse(w, &gmail.Thread{Id: threadID})
	}

	repo := NewGmailThreadRepository(ts.GmailRepository(t))
	ctx := context.Background()

	err := repo.Untrash(ctx, "thread123")
	if err != nil {
		t.Fatalf("Untrash failed: %v", err)
	}

	if untrashedID != "thread123" {
		t.Errorf("untrashedID = %q, want %q", untrashedID, "thread123")
	}
}

// TestGmailThreadRepository_DeleteWithTestServer tests permanently deleting a thread.
func TestGmailThreadRepository_DeleteWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	deletedID := ""
	ts.ThreadDeleteHandler = func(w http.ResponseWriter, r *http.Request, threadID string) {
		deletedID = threadID
		w.WriteHeader(http.StatusNoContent)
	}

	repo := NewGmailThreadRepository(ts.GmailRepository(t))
	ctx := context.Background()

	err := repo.Delete(ctx, "thread123")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if deletedID != "thread123" {
		t.Errorf("deletedID = %q, want %q", deletedID, "thread123")
	}
}

// TestGmailThreadRepository_GetNotFound tests getting a non-existent thread.
func TestGmailThreadRepository_GetNotFound(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.ThreadGetHandler = func(w http.ResponseWriter, r *http.Request, threadID string) {
		WriteErrorResponse(w, http.StatusNotFound, "thread not found")
	}

	repo := NewGmailThreadRepository(ts.GmailRepository(t))
	ctx := context.Background()

	_, err := repo.Get(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent thread, got nil")
	}
}

// =============================================================================
// Label Repository Tests
// =============================================================================

// TestGmailLabelRepository_GetWithTestServer tests getting a label by ID.
func TestGmailLabelRepository_GetWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.LabelGetHandler = func(w http.ResponseWriter, r *http.Request, labelID string) {
		if labelID == "Label_1" {
			WriteJSONResponse(w, MockLabelResponse("Label_1", "Work", "user"))
		} else {
			WriteErrorResponse(w, http.StatusNotFound, "label not found")
		}
	}

	repo := NewGmailLabelRepository(ts.GmailRepository(t))
	ctx := context.Background()

	label, err := repo.Get(ctx, "Label_1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if label.ID != "Label_1" {
		t.Errorf("label.ID = %q, want %q", label.ID, "Label_1")
	}
	if label.Name != "Work" {
		t.Errorf("label.Name = %q, want %q", label.Name, "Work")
	}
}

// TestGmailLabelRepository_CreateWithTestServer tests creating a label.
func TestGmailLabelRepository_CreateWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	var createdLabel *gmail.Label
	ts.LabelCreateHandler = func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&createdLabel); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}
		createdLabel.Id = "Label_new_123"
		WriteJSONResponse(w, createdLabel)
	}

	repo := NewGmailLabelRepository(ts.GmailRepository(t))
	ctx := context.Background()

	label := &mail.Label{
		Name: "New Label",
		Type: "user",
	}

	created, err := repo.Create(ctx, label)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if created.ID != "Label_new_123" {
		t.Errorf("created.ID = %q, want %q", created.ID, "Label_new_123")
	}
	if created.Name != "New Label" {
		t.Errorf("created.Name = %q, want %q", created.Name, "New Label")
	}
}

// TestGmailLabelRepository_UpdateWithTestServer tests updating a label.
func TestGmailLabelRepository_UpdateWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.LabelUpdateHandler = func(w http.ResponseWriter, r *http.Request, labelID string) {
		var label gmail.Label
		if err := json.NewDecoder(r.Body).Decode(&label); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
			return
		}
		label.Id = labelID
		WriteJSONResponse(w, &label)
	}

	repo := NewGmailLabelRepository(ts.GmailRepository(t))
	ctx := context.Background()

	label := &mail.Label{
		ID:   "Label_1",
		Name: "Updated Label",
	}

	updated, err := repo.Update(ctx, label)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.ID != "Label_1" {
		t.Errorf("updated.ID = %q, want %q", updated.ID, "Label_1")
	}
	if updated.Name != "Updated Label" {
		t.Errorf("updated.Name = %q, want %q", updated.Name, "Updated Label")
	}
}

// TestGmailLabelRepository_DeleteWithTestServer tests deleting a label.
func TestGmailLabelRepository_DeleteWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	deletedID := ""
	ts.LabelDeleteHandler = func(w http.ResponseWriter, r *http.Request, labelID string) {
		deletedID = labelID
		w.WriteHeader(http.StatusNoContent)
	}

	repo := NewGmailLabelRepository(ts.GmailRepository(t))
	ctx := context.Background()

	err := repo.Delete(ctx, "Label_1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if deletedID != "Label_1" {
		t.Errorf("deletedID = %q, want %q", deletedID, "Label_1")
	}
}

// TestGmailLabelRepository_GetByNameWithTestServer tests getting a label by name.
func TestGmailLabelRepository_GetByNameWithTestServer(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.LabelListHandler = StaticLabelListHandler([]*gmail.Label{
		MockLabelResponse("INBOX", "INBOX", "system"),
		MockLabelResponse("Label_1", "Work", "user"),
		MockLabelResponse("Label_2", "Personal", "user"),
	})

	repo := NewGmailLabelRepository(ts.GmailRepository(t))
	ctx := context.Background()

	label, err := repo.GetByName(ctx, "Work")
	if err != nil {
		t.Fatalf("GetByName failed: %v", err)
	}

	if label.ID != "Label_1" {
		t.Errorf("label.ID = %q, want %q", label.ID, "Label_1")
	}
	if label.Name != "Work" {
		t.Errorf("label.Name = %q, want %q", label.Name, "Work")
	}
}

// TestGmailLabelRepository_GetByNameNotFound tests GetByName for non-existent label.
func TestGmailLabelRepository_GetByNameNotFound(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.LabelListHandler = StaticLabelListHandler([]*gmail.Label{
		MockLabelResponse("INBOX", "INBOX", "system"),
	})

	repo := NewGmailLabelRepository(ts.GmailRepository(t))
	ctx := context.Background()

	_, err := repo.GetByName(ctx, "Nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent label name, got nil")
	}
}

// =============================================================================
// Error Handling Tests
// =============================================================================

// TestGmailRepository_ForbiddenError tests 403 Forbidden handling.
func TestGmailRepository_ForbiddenError(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.MessageGetHandler = func(w http.ResponseWriter, r *http.Request, msgID string) {
		WriteErrorResponse(w, http.StatusForbidden, "Access denied")
	}

	repo := ts.GmailRepository(t)
	ctx := context.Background()

	_, err := repo.Get(ctx, "msg123")
	if err == nil {
		t.Fatal("expected error for forbidden request, got nil")
	}

	if !strings.Contains(err.Error(), "403") {
		t.Errorf("expected error to contain 403, got %v", err)
	}
}

// TestGmailRepository_InternalServerError tests 500 error handling.
func TestGmailRepository_InternalServerError(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.MessageGetHandler = func(w http.ResponseWriter, r *http.Request, msgID string) {
		WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
	}

	repo := ts.GmailRepository(t)
	ctx := context.Background()

	_, err := repo.Get(ctx, "msg123")
	if err == nil {
		t.Fatal("expected error for server error, got nil")
	}

	if !strings.Contains(err.Error(), ErrTemporary.Error()) {
		t.Errorf("expected temporary error, got %v", err)
	}
}

// TestGmailRepository_BadRequestError tests 400 error handling.
func TestGmailRepository_BadRequestError(t *testing.T) {
	ts := NewTestServer()
	defer ts.Close()

	ts.MessageModifyHandler = func(w http.ResponseWriter, r *http.Request, msgID string) {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid label ID")
	}

	repo := ts.GmailRepository(t)
	ctx := context.Background()

	_, err := repo.Modify(ctx, "msg123", mail.ModifyRequest{
		AddLabels: []string{"invalid_label"},
	})
	if err == nil {
		t.Fatal("expected error for bad request, got nil")
	}

	if !strings.Contains(err.Error(), ErrBadRequest.Error()) {
		t.Errorf("expected bad request error, got %v", err)
	}
}
