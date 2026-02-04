package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	gcal "google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// -----------------------------------------------------------------------------
// Test Server Infrastructure
// -----------------------------------------------------------------------------

// TestServer wraps httptest.Server with API-specific routing and handlers.
// It provides a reusable mock server for testing Gmail and Calendar API methods.
type TestServer struct {
	Server *httptest.Server
	mux    *http.ServeMux

	// Route handlers - assign custom handlers to override default behavior
	mu sync.RWMutex

	// Gmail handlers
	MessageListHandler    func(w http.ResponseWriter, r *http.Request)
	MessageGetHandler     func(w http.ResponseWriter, r *http.Request, msgID string)
	MessageSendHandler    func(w http.ResponseWriter, r *http.Request)
	MessageTrashHandler   func(w http.ResponseWriter, r *http.Request, msgID string)
	MessageUntrashHandler func(w http.ResponseWriter, r *http.Request, msgID string)
	MessageModifyHandler  func(w http.ResponseWriter, r *http.Request, msgID string)
	MessageDeleteHandler  func(w http.ResponseWriter, r *http.Request, msgID string)

	DraftListHandler   func(w http.ResponseWriter, r *http.Request)
	DraftGetHandler    func(w http.ResponseWriter, r *http.Request, draftID string)
	DraftCreateHandler func(w http.ResponseWriter, r *http.Request)
	DraftUpdateHandler func(w http.ResponseWriter, r *http.Request, draftID string)
	DraftSendHandler   func(w http.ResponseWriter, r *http.Request)
	DraftDeleteHandler func(w http.ResponseWriter, r *http.Request, draftID string)

	ThreadListHandler    func(w http.ResponseWriter, r *http.Request)
	ThreadGetHandler     func(w http.ResponseWriter, r *http.Request, threadID string)
	ThreadModifyHandler  func(w http.ResponseWriter, r *http.Request, threadID string)
	ThreadTrashHandler   func(w http.ResponseWriter, r *http.Request, threadID string)
	ThreadUntrashHandler func(w http.ResponseWriter, r *http.Request, threadID string)
	ThreadDeleteHandler  func(w http.ResponseWriter, r *http.Request, threadID string)

	LabelListHandler   func(w http.ResponseWriter, r *http.Request)
	LabelGetHandler    func(w http.ResponseWriter, r *http.Request, labelID string)
	LabelCreateHandler func(w http.ResponseWriter, r *http.Request)
	LabelUpdateHandler func(w http.ResponseWriter, r *http.Request, labelID string)
	LabelDeleteHandler func(w http.ResponseWriter, r *http.Request, labelID string)

	// Calendar handlers
	EventListHandler      func(w http.ResponseWriter, r *http.Request, calendarID string)
	EventGetHandler       func(w http.ResponseWriter, r *http.Request, calendarID, eventID string)
	EventCreateHandler    func(w http.ResponseWriter, r *http.Request, calendarID string)
	EventUpdateHandler    func(w http.ResponseWriter, r *http.Request, calendarID, eventID string)
	EventDeleteHandler    func(w http.ResponseWriter, r *http.Request, calendarID, eventID string)
	EventMoveHandler      func(w http.ResponseWriter, r *http.Request, calendarID, eventID string)
	EventQuickAddHandler  func(w http.ResponseWriter, r *http.Request, calendarID string)
	EventInstancesHandler func(w http.ResponseWriter, r *http.Request, calendarID, eventID string)

	CalendarListHandler   func(w http.ResponseWriter, r *http.Request)
	CalendarGetHandler    func(w http.ResponseWriter, r *http.Request, calendarID string)
	CalendarCreateHandler func(w http.ResponseWriter, r *http.Request)
	CalendarUpdateHandler func(w http.ResponseWriter, r *http.Request, calendarID string)
	CalendarDeleteHandler func(w http.ResponseWriter, r *http.Request, calendarID string)
	CalendarClearHandler  func(w http.ResponseWriter, r *http.Request, calendarID string)

	ACLListHandler   func(w http.ResponseWriter, r *http.Request, calendarID string)
	ACLGetHandler    func(w http.ResponseWriter, r *http.Request, calendarID, ruleID string)
	ACLInsertHandler func(w http.ResponseWriter, r *http.Request, calendarID string)
	ACLUpdateHandler func(w http.ResponseWriter, r *http.Request, calendarID, ruleID string)
	ACLDeleteHandler func(w http.ResponseWriter, r *http.Request, calendarID, ruleID string)

	FreeBusyQueryHandler func(w http.ResponseWriter, r *http.Request)
}

// NewTestServer creates a new TestServer with default routing.
func NewTestServer() *TestServer {
	ts := &TestServer{
		mux: http.NewServeMux(),
	}
	ts.setupRoutes()
	ts.Server = httptest.NewServer(ts.mux)
	return ts
}

// Close shuts down the test server.
func (ts *TestServer) Close() {
	ts.Server.Close()
}

// URL returns the server's URL.
func (ts *TestServer) URL() string {
	return ts.Server.URL
}

// GmailService creates a Gmail service configured to use this test server.
func (ts *TestServer) GmailService(t *testing.T) *gmail.Service {
	t.Helper()
	ctx := context.Background()
	service, err := gmail.NewService(ctx,
		option.WithEndpoint(ts.Server.URL),
		option.WithoutAuthentication(),
	)
	if err != nil {
		t.Fatalf("failed to create Gmail service: %v", err)
	}
	return service
}

// CalendarService creates a Calendar service configured to use this test server.
func (ts *TestServer) CalendarService(t *testing.T) *gcal.Service {
	t.Helper()
	ctx := context.Background()
	service, err := gcal.NewService(ctx,
		option.WithEndpoint(ts.Server.URL),
		option.WithoutAuthentication(),
	)
	if err != nil {
		t.Fatalf("failed to create Calendar service: %v", err)
	}
	return service
}

// GmailRepository creates a GmailRepository configured to use this test server.
func (ts *TestServer) GmailRepository(t *testing.T) *GmailRepository {
	t.Helper()
	return NewGmailRepositoryWithService(ts.GmailService(t), "me")
}

// CalendarService creates a GCalService configured to use this test server.
func (ts *TestServer) GCalService(t *testing.T) *GCalService {
	t.Helper()
	return NewGCalServiceWithService(ts.CalendarService(t))
}

// setupRoutes configures all API routes with default handlers.
func (ts *TestServer) setupRoutes() {
	// Gmail API routes - more specific routes first
	ts.mux.HandleFunc("/gmail/v1/users/me/messages/send", ts.handleGmailMessageSend)
	ts.mux.HandleFunc("/gmail/v1/users/me/messages", ts.handleGmailMessages)
	ts.mux.HandleFunc("/gmail/v1/users/me/messages/", ts.handleGmailMessage)
	ts.mux.HandleFunc("/gmail/v1/users/me/drafts/send", ts.handleGmailDraftSend)
	ts.mux.HandleFunc("/gmail/v1/users/me/drafts", ts.handleGmailDrafts)
	ts.mux.HandleFunc("/gmail/v1/users/me/drafts/", ts.handleGmailDraft)
	ts.mux.HandleFunc("/gmail/v1/users/me/threads", ts.handleGmailThreads)
	ts.mux.HandleFunc("/gmail/v1/users/me/threads/", ts.handleGmailThread)
	ts.mux.HandleFunc("/gmail/v1/users/me/labels", ts.handleGmailLabels)
	ts.mux.HandleFunc("/gmail/v1/users/me/labels/", ts.handleGmailLabel)

	// Calendar API routes - the Google API client strips the /calendar/v3 prefix
	ts.mux.HandleFunc("/calendars", ts.handleCalendarCreate)
	ts.mux.HandleFunc("/calendars/", ts.handleCalendar)
	ts.mux.HandleFunc("/users/me/calendarList", ts.handleCalendarList)
	ts.mux.HandleFunc("/users/me/calendarList/", ts.handleCalendarListEntry)
	ts.mux.HandleFunc("/freeBusy", ts.handleFreeBusy)
}

// -----------------------------------------------------------------------------
// Gmail Route Handlers
// -----------------------------------------------------------------------------

func (ts *TestServer) handleGmailMessageSend(w http.ResponseWriter, r *http.Request) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	if r.Method == http.MethodPost {
		if ts.MessageSendHandler != nil {
			ts.MessageSendHandler(w, r)
		} else {
			http.Error(w, "send handler not configured", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleGmailMessages(w http.ResponseWriter, r *http.Request) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	switch r.Method {
	case http.MethodGet:
		if ts.MessageListHandler != nil {
			ts.MessageListHandler(w, r)
		} else {
			WriteJSONResponse(w, &gmail.ListMessagesResponse{Messages: []*gmail.Message{}})
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleGmailMessage(w http.ResponseWriter, r *http.Request) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	path := strings.TrimPrefix(r.URL.Path, "/gmail/v1/users/me/messages/")
	parts := strings.Split(path, "/")
	msgID := parts[0]

	// Check for sub-resources
	if len(parts) > 1 {
		switch parts[1] {
		case "send":
			if ts.MessageSendHandler != nil {
				ts.MessageSendHandler(w, r)
			}
			return
		case "trash":
			if ts.MessageTrashHandler != nil {
				ts.MessageTrashHandler(w, r, msgID)
			} else {
				WriteJSONResponse(w, &gmail.Message{Id: msgID, LabelIds: []string{"TRASH"}})
			}
			return
		case "untrash":
			if ts.MessageUntrashHandler != nil {
				ts.MessageUntrashHandler(w, r, msgID)
			} else {
				WriteJSONResponse(w, &gmail.Message{Id: msgID, LabelIds: []string{"INBOX"}})
			}
			return
		case "modify":
			if ts.MessageModifyHandler != nil {
				ts.MessageModifyHandler(w, r, msgID)
			}
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		if ts.MessageGetHandler != nil {
			ts.MessageGetHandler(w, r, msgID)
		} else {
			http.Error(w, "message not found", http.StatusNotFound)
		}
	case http.MethodDelete:
		if ts.MessageDeleteHandler != nil {
			ts.MessageDeleteHandler(w, r, msgID)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleGmailDraftSend(w http.ResponseWriter, r *http.Request) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	if r.Method == http.MethodPost {
		if ts.DraftSendHandler != nil {
			ts.DraftSendHandler(w, r)
		} else {
			http.Error(w, "draft send handler not configured", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleGmailDrafts(w http.ResponseWriter, r *http.Request) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	switch r.Method {
	case http.MethodGet:
		if ts.DraftListHandler != nil {
			ts.DraftListHandler(w, r)
		} else {
			WriteJSONResponse(w, &gmail.ListDraftsResponse{Drafts: []*gmail.Draft{}})
		}
	case http.MethodPost:
		if ts.DraftCreateHandler != nil {
			ts.DraftCreateHandler(w, r)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleGmailDraft(w http.ResponseWriter, r *http.Request) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	path := strings.TrimPrefix(r.URL.Path, "/gmail/v1/users/me/drafts/")
	parts := strings.Split(path, "/")
	draftID := parts[0]

	// Check for send sub-resource
	if len(parts) > 1 && parts[1] == "send" {
		if ts.DraftSendHandler != nil {
			ts.DraftSendHandler(w, r)
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		if ts.DraftGetHandler != nil {
			ts.DraftGetHandler(w, r, draftID)
		} else {
			http.Error(w, "draft not found", http.StatusNotFound)
		}
	case http.MethodPut:
		if ts.DraftUpdateHandler != nil {
			ts.DraftUpdateHandler(w, r, draftID)
		}
	case http.MethodDelete:
		if ts.DraftDeleteHandler != nil {
			ts.DraftDeleteHandler(w, r, draftID)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleGmailThreads(w http.ResponseWriter, r *http.Request) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	switch r.Method {
	case http.MethodGet:
		if ts.ThreadListHandler != nil {
			ts.ThreadListHandler(w, r)
		} else {
			WriteJSONResponse(w, &gmail.ListThreadsResponse{Threads: []*gmail.Thread{}})
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleGmailThread(w http.ResponseWriter, r *http.Request) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	path := strings.TrimPrefix(r.URL.Path, "/gmail/v1/users/me/threads/")
	parts := strings.Split(path, "/")
	threadID := parts[0]

	// Check for sub-resources
	if len(parts) > 1 {
		switch parts[1] {
		case "trash":
			if ts.ThreadTrashHandler != nil {
				ts.ThreadTrashHandler(w, r, threadID)
			} else {
				WriteJSONResponse(w, &gmail.Thread{Id: threadID})
			}
			return
		case "untrash":
			if ts.ThreadUntrashHandler != nil {
				ts.ThreadUntrashHandler(w, r, threadID)
			} else {
				WriteJSONResponse(w, &gmail.Thread{Id: threadID})
			}
			return
		case "modify":
			if ts.ThreadModifyHandler != nil {
				ts.ThreadModifyHandler(w, r, threadID)
			}
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		if ts.ThreadGetHandler != nil {
			ts.ThreadGetHandler(w, r, threadID)
		} else {
			http.Error(w, "thread not found", http.StatusNotFound)
		}
	case http.MethodDelete:
		if ts.ThreadDeleteHandler != nil {
			ts.ThreadDeleteHandler(w, r, threadID)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleGmailLabels(w http.ResponseWriter, r *http.Request) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	switch r.Method {
	case http.MethodGet:
		if ts.LabelListHandler != nil {
			ts.LabelListHandler(w, r)
		} else {
			WriteJSONResponse(w, &gmail.ListLabelsResponse{Labels: []*gmail.Label{}})
		}
	case http.MethodPost:
		if ts.LabelCreateHandler != nil {
			ts.LabelCreateHandler(w, r)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleGmailLabel(w http.ResponseWriter, r *http.Request) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	labelID := strings.TrimPrefix(r.URL.Path, "/gmail/v1/users/me/labels/")

	switch r.Method {
	case http.MethodGet:
		if ts.LabelGetHandler != nil {
			ts.LabelGetHandler(w, r, labelID)
		} else {
			http.Error(w, "label not found", http.StatusNotFound)
		}
	case http.MethodPut, http.MethodPatch:
		if ts.LabelUpdateHandler != nil {
			ts.LabelUpdateHandler(w, r, labelID)
		}
	case http.MethodDelete:
		if ts.LabelDeleteHandler != nil {
			ts.LabelDeleteHandler(w, r, labelID)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// -----------------------------------------------------------------------------
// Calendar Route Handlers
// -----------------------------------------------------------------------------

func (ts *TestServer) handleCalendarCreate(w http.ResponseWriter, r *http.Request) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	if r.Method == http.MethodPost {
		if ts.CalendarCreateHandler != nil {
			ts.CalendarCreateHandler(w, r)
		} else {
			http.Error(w, "calendar create handler not configured", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleCalendar(w http.ResponseWriter, r *http.Request) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	path := strings.TrimPrefix(r.URL.Path, "/calendars/")
	parts := strings.Split(path, "/")
	calendarID := parts[0]

	// Route based on path structure
	if len(parts) == 1 {
		// Calendar operations
		ts.handleCalendarResource(w, r, calendarID)
		return
	}

	switch parts[1] {
	case "events":
		if len(parts) == 2 {
			// List or create events
			ts.handleEventsList(w, r, calendarID)
		} else if len(parts) >= 3 {
			eventID := parts[2]
			// Handle quickAdd endpoint: /calendars/{calendarId}/events/quickAdd
			if eventID == "quickAdd" {
				if ts.EventQuickAddHandler != nil {
					ts.EventQuickAddHandler(w, r, calendarID)
				} else {
					http.Error(w, "quick add handler not configured", http.StatusInternalServerError)
				}
				return
			}
			if len(parts) == 4 {
				switch parts[3] {
				case "move":
					if ts.EventMoveHandler != nil {
						ts.EventMoveHandler(w, r, calendarID, eventID)
					}
				case "instances":
					if ts.EventInstancesHandler != nil {
						ts.EventInstancesHandler(w, r, calendarID, eventID)
					} else {
						WriteJSONResponse(w, &gcal.Events{Items: []*gcal.Event{}})
					}
				default:
					http.Error(w, "not found", http.StatusNotFound)
				}
			} else {
				ts.handleEventResource(w, r, calendarID, eventID)
			}
		}
	case "acl":
		if len(parts) == 2 {
			ts.handleACLList(w, r, calendarID)
		} else {
			ruleID := parts[2]
			ts.handleACLResource(w, r, calendarID, ruleID)
		}
	case "clear":
		if ts.CalendarClearHandler != nil {
			ts.CalendarClearHandler(w, r, calendarID)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func (ts *TestServer) handleCalendarResource(w http.ResponseWriter, r *http.Request, calendarID string) {
	switch r.Method {
	case http.MethodGet:
		if ts.CalendarGetHandler != nil {
			ts.CalendarGetHandler(w, r, calendarID)
		} else {
			http.Error(w, "calendar not found", http.StatusNotFound)
		}
	case http.MethodPut:
		if ts.CalendarUpdateHandler != nil {
			ts.CalendarUpdateHandler(w, r, calendarID)
		}
	case http.MethodDelete:
		if ts.CalendarDeleteHandler != nil {
			ts.CalendarDeleteHandler(w, r, calendarID)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	case http.MethodPost:
		// Calendar insert
		if ts.CalendarCreateHandler != nil {
			ts.CalendarCreateHandler(w, r)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleEventsList(w http.ResponseWriter, r *http.Request, calendarID string) {
	switch r.Method {
	case http.MethodGet:
		// Check for quickAdd parameter
		if r.URL.Query().Get("text") != "" {
			if ts.EventQuickAddHandler != nil {
				ts.EventQuickAddHandler(w, r, calendarID)
			}
			return
		}
		if ts.EventListHandler != nil {
			ts.EventListHandler(w, r, calendarID)
		} else {
			WriteJSONResponse(w, &gcal.Events{Items: []*gcal.Event{}})
		}
	case http.MethodPost:
		// Check if this is a quickAdd
		if r.URL.Query().Get("text") != "" {
			if ts.EventQuickAddHandler != nil {
				ts.EventQuickAddHandler(w, r, calendarID)
			}
		} else if ts.EventCreateHandler != nil {
			ts.EventCreateHandler(w, r, calendarID)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleEventResource(w http.ResponseWriter, r *http.Request, calendarID, eventID string) {
	switch r.Method {
	case http.MethodGet:
		if ts.EventGetHandler != nil {
			ts.EventGetHandler(w, r, calendarID, eventID)
		} else {
			http.Error(w, "event not found", http.StatusNotFound)
		}
	case http.MethodPut, http.MethodPatch:
		if ts.EventUpdateHandler != nil {
			ts.EventUpdateHandler(w, r, calendarID, eventID)
		}
	case http.MethodDelete:
		if ts.EventDeleteHandler != nil {
			ts.EventDeleteHandler(w, r, calendarID, eventID)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleCalendarList(w http.ResponseWriter, r *http.Request) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	switch r.Method {
	case http.MethodGet:
		if ts.CalendarListHandler != nil {
			ts.CalendarListHandler(w, r)
		} else {
			WriteJSONResponse(w, &gcal.CalendarList{Items: []*gcal.CalendarListEntry{}})
		}
	case http.MethodPost:
		if ts.CalendarCreateHandler != nil {
			ts.CalendarCreateHandler(w, r)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleCalendarListEntry(w http.ResponseWriter, r *http.Request) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	calendarID := strings.TrimPrefix(r.URL.Path, "/users/me/calendarList/")

	switch r.Method {
	case http.MethodGet:
		if ts.CalendarGetHandler != nil {
			ts.CalendarGetHandler(w, r, calendarID)
		} else {
			http.Error(w, "calendar not found", http.StatusNotFound)
		}
	case http.MethodDelete:
		if ts.CalendarDeleteHandler != nil {
			ts.CalendarDeleteHandler(w, r, calendarID)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleACLList(w http.ResponseWriter, r *http.Request, calendarID string) {
	switch r.Method {
	case http.MethodGet:
		if ts.ACLListHandler != nil {
			ts.ACLListHandler(w, r, calendarID)
		} else {
			WriteJSONResponse(w, &gcal.Acl{Items: []*gcal.AclRule{}})
		}
	case http.MethodPost:
		if ts.ACLInsertHandler != nil {
			ts.ACLInsertHandler(w, r, calendarID)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleACLResource(w http.ResponseWriter, r *http.Request, calendarID, ruleID string) {
	switch r.Method {
	case http.MethodGet:
		if ts.ACLGetHandler != nil {
			ts.ACLGetHandler(w, r, calendarID, ruleID)
		} else {
			http.Error(w, "ACL rule not found", http.StatusNotFound)
		}
	case http.MethodPut, http.MethodPatch:
		if ts.ACLUpdateHandler != nil {
			ts.ACLUpdateHandler(w, r, calendarID, ruleID)
		}
	case http.MethodDelete:
		if ts.ACLDeleteHandler != nil {
			ts.ACLDeleteHandler(w, r, calendarID, ruleID)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ts *TestServer) handleFreeBusy(w http.ResponseWriter, r *http.Request) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	if r.Method == http.MethodPost {
		if ts.FreeBusyQueryHandler != nil {
			ts.FreeBusyQueryHandler(w, r)
		} else {
			WriteJSONResponse(w, &gcal.FreeBusyResponse{
				Calendars: map[string]gcal.FreeBusyCalendar{},
			})
		}
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// -----------------------------------------------------------------------------
// Mock Response Helpers
// -----------------------------------------------------------------------------

// WriteJSONResponse writes a JSON response with proper headers.
func WriteJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// WriteErrorResponse writes an error response in Google API format.
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    statusCode,
			"message": message,
		},
	})
}

// MockMessageResponse creates a mock Gmail message response.
func MockMessageResponse(id, threadID, subject, from, to, body string) *gmail.Message {
	encodedBody := base64.URLEncoding.EncodeToString([]byte(body))
	return &gmail.Message{
		Id:       id,
		ThreadId: threadID,
		Snippet:  truncateSnippet(body, 100),
		LabelIds: []string{"INBOX"},
		Payload: &gmail.MessagePart{
			Headers: []*gmail.MessagePartHeader{
				{Name: "From", Value: from},
				{Name: "To", Value: to},
				{Name: "Subject", Value: subject},
				{Name: "Date", Value: time.Now().Format(time.RFC1123Z)},
			},
			Body: &gmail.MessagePartBody{
				Data: encodedBody,
			},
		},
	}
}

// MockMessageListResponse creates a mock Gmail message list response.
func MockMessageListResponse(messages []*gmail.Message, nextPageToken string, total int64) *gmail.ListMessagesResponse {
	return &gmail.ListMessagesResponse{
		Messages:           messages,
		NextPageToken:      nextPageToken,
		ResultSizeEstimate: total,
	}
}

// MockDraftResponse creates a mock Gmail draft response.
func MockDraftResponse(id, messageID, subject, from, to, body string) *gmail.Draft {
	return &gmail.Draft{
		Id:      id,
		Message: MockMessageResponse(messageID, "thread1", subject, from, to, body),
	}
}

// MockLabelResponse creates a mock Gmail label response.
func MockLabelResponse(id, name, labelType string) *gmail.Label {
	return &gmail.Label{
		Id:   id,
		Name: name,
		Type: labelType,
	}
}

// MockLabelListResponse creates a mock Gmail label list response.
func MockLabelListResponse(labels []*gmail.Label) *gmail.ListLabelsResponse {
	return &gmail.ListLabelsResponse{
		Labels: labels,
	}
}

// MockThreadResponse creates a mock Gmail thread response.
func MockThreadResponse(id string, messages []*gmail.Message) *gmail.Thread {
	var snippet string
	if len(messages) > 0 {
		snippet = messages[0].Snippet
	}
	return &gmail.Thread{
		Id:       id,
		Snippet:  snippet,
		Messages: messages,
	}
}

// MockEventResponse creates a mock Calendar event response.
func MockEventResponse(id, summary, description string, start, end time.Time) *gcal.Event {
	return &gcal.Event{
		Id:          id,
		Summary:     summary,
		Description: description,
		Start: &gcal.EventDateTime{
			DateTime: start.Format(time.RFC3339),
		},
		End: &gcal.EventDateTime{
			DateTime: end.Format(time.RFC3339),
		},
		Status:   "confirmed",
		HtmlLink: fmt.Sprintf("https://calendar.google.com/event?eid=%s", id),
		Created:  time.Now().Format(time.RFC3339),
		Updated:  time.Now().Format(time.RFC3339),
	}
}

// MockAllDayEventResponse creates a mock all-day Calendar event response.
func MockAllDayEventResponse(id, summary string, date time.Time) *gcal.Event {
	return &gcal.Event{
		Id:      id,
		Summary: summary,
		Start: &gcal.EventDateTime{
			Date: date.Format("2006-01-02"),
		},
		End: &gcal.EventDateTime{
			Date: date.AddDate(0, 0, 1).Format("2006-01-02"),
		},
		Status:  "confirmed",
		Created: time.Now().Format(time.RFC3339),
		Updated: time.Now().Format(time.RFC3339),
	}
}

// MockEventListResponse creates a mock Calendar events list response.
func MockEventListResponse(events []*gcal.Event, nextPageToken string) *gcal.Events {
	return &gcal.Events{
		Items:         events,
		NextPageToken: nextPageToken,
	}
}

// MockCalendarResponse creates a mock Calendar response.
func MockCalendarResponse(id, summary, description, timeZone string) *gcal.Calendar {
	return &gcal.Calendar{
		Id:          id,
		Summary:     summary,
		Description: description,
		TimeZone:    timeZone,
	}
}

// MockCalendarListEntryResponse creates a mock CalendarListEntry response.
func MockCalendarListEntryResponse(id, summary, description, timeZone string, primary bool, accessRole string) *gcal.CalendarListEntry {
	return &gcal.CalendarListEntry{
		Id:          id,
		Summary:     summary,
		Description: description,
		TimeZone:    timeZone,
		Primary:     primary,
		AccessRole:  accessRole,
		Selected:    true,
	}
}

// MockCalendarListResponse creates a mock CalendarList response.
func MockCalendarListResponse(calendars []*gcal.CalendarListEntry, nextPageToken string) *gcal.CalendarList {
	return &gcal.CalendarList{
		Items:         calendars,
		NextPageToken: nextPageToken,
	}
}

// MockACLRuleResponse creates a mock ACL rule response.
func MockACLRuleResponse(id, role, scopeType, scopeValue string) *gcal.AclRule {
	return &gcal.AclRule{
		Id:   id,
		Role: role,
		Scope: &gcal.AclRuleScope{
			Type:  scopeType,
			Value: scopeValue,
		},
	}
}

// MockACLListResponse creates a mock ACL list response.
func MockACLListResponse(rules []*gcal.AclRule, nextPageToken string) *gcal.Acl {
	return &gcal.Acl{
		Items:         rules,
		NextPageToken: nextPageToken,
	}
}

// MockFreeBusyResponse creates a mock FreeBusy response.
func MockFreeBusyResponse(calendars map[string][]struct{ Start, End time.Time }) *gcal.FreeBusyResponse {
	result := &gcal.FreeBusyResponse{
		Calendars: make(map[string]gcal.FreeBusyCalendar),
	}
	for calID, periods := range calendars {
		busyPeriods := make([]*gcal.TimePeriod, len(periods))
		for i, p := range periods {
			busyPeriods[i] = &gcal.TimePeriod{
				Start: p.Start.Format(time.RFC3339),
				End:   p.End.Format(time.RFC3339),
			}
		}
		result.Calendars[calID] = gcal.FreeBusyCalendar{
			Busy: busyPeriods,
		}
	}
	return result
}

// truncateSnippet truncates a string to the specified length.
func truncateSnippet(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// -----------------------------------------------------------------------------
// Handler Builder Helpers
// -----------------------------------------------------------------------------

// StaticMessageHandler returns a handler that always returns the given message.
func StaticMessageHandler(msg *gmail.Message) func(w http.ResponseWriter, r *http.Request, msgID string) {
	return func(w http.ResponseWriter, r *http.Request, msgID string) {
		WriteJSONResponse(w, msg)
	}
}

// StaticMessageListHandler returns a handler that always returns the given messages.
func StaticMessageListHandler(messages []*gmail.Message, nextPageToken string, total int64) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteJSONResponse(w, MockMessageListResponse(messages, nextPageToken, total))
	}
}

// StaticLabelListHandler returns a handler that always returns the given labels.
func StaticLabelListHandler(labels []*gmail.Label) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteJSONResponse(w, MockLabelListResponse(labels))
	}
}

// StaticEventHandler returns a handler that always returns the given event.
func StaticEventHandler(event *gcal.Event) func(w http.ResponseWriter, r *http.Request, calendarID, eventID string) {
	return func(w http.ResponseWriter, r *http.Request, calendarID, eventID string) {
		WriteJSONResponse(w, event)
	}
}

// StaticEventListHandler returns a handler that always returns the given events.
func StaticEventListHandler(events []*gcal.Event, nextPageToken string) func(w http.ResponseWriter, r *http.Request, calendarID string) {
	return func(w http.ResponseWriter, r *http.Request, calendarID string) {
		WriteJSONResponse(w, MockEventListResponse(events, nextPageToken))
	}
}

// StaticCalendarListHandler returns a handler that always returns the given calendars.
func StaticCalendarListHandler(calendars []*gcal.CalendarListEntry, nextPageToken string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteJSONResponse(w, MockCalendarListResponse(calendars, nextPageToken))
	}
}

// NotFoundHandler returns a handler that returns a 404 error.
func NotFoundHandler(resource string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("%s not found", resource))
	}
}

// ErrorHandler returns a handler that returns the specified error.
func ErrorHandler(statusCode int, message string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteErrorResponse(w, statusCode, message)
	}
}
