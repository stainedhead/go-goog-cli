// Package repository provides implementations for domain repository interfaces.
package repository

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// Repository errors for Gmail operations.
var (
	ErrBadRequest  = errors.New("bad request")
	ErrRateLimited = errors.New("rate limited")
	ErrTemporary   = errors.New("temporary error")
)

// Default configuration for Gmail repository.
const (
	defaultMaxRetries   = 3
	defaultBaseBackoff  = 100 * time.Millisecond
	gmailLabelInbox     = "INBOX"
	gmailLabelUnread    = "UNREAD"
	gmailLabelStarred   = "STARRED"
	gmailLabelTrash     = "TRASH"
	gmailMessageFormat  = "full"
	gmailMetadataFormat = "metadata"
)

// GmailRepository implements MessageRepository using the Gmail API.
type GmailRepository struct {
	service     *gmail.Service
	userID      string
	maxRetries  int
	baseBackoff time.Duration
}

// Compile-time interface compliance checks.
var (
	_ mail.MessageRepository = (*GmailRepository)(nil)
	_ mail.DraftRepository   = (*GmailDraftRepository)(nil)
	_ mail.LabelRepository   = (*GmailLabelRepository)(nil)
)

// GmailDraftRepository wraps GmailRepository to implement DraftRepository.
type GmailDraftRepository struct {
	*GmailRepository
}

// NewGmailDraftRepository creates a new GmailDraftRepository.
func NewGmailDraftRepository(repo *GmailRepository) *GmailDraftRepository {
	return &GmailDraftRepository{GmailRepository: repo}
}

// GmailLabelRepository wraps GmailRepository to implement LabelRepository.
type GmailLabelRepository struct {
	*GmailRepository
}

// NewGmailLabelRepository creates a new GmailLabelRepository.
func NewGmailLabelRepository(repo *GmailRepository) *GmailLabelRepository {
	return &GmailLabelRepository{GmailRepository: repo}
}

// NewGmailRepository creates a new GmailRepository with the given OAuth2 token source.
func NewGmailRepository(ctx context.Context, tokenSource oauth2.TokenSource) (*GmailRepository, error) {
	httpClient := oauth2.NewClient(ctx, tokenSource)

	service, err := gmail.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gmail service: %w", err)
	}

	return &GmailRepository{
		service:     service,
		userID:      "me",
		maxRetries:  defaultMaxRetries,
		baseBackoff: defaultBaseBackoff,
	}, nil
}

// NewGmailRepositoryWithService creates a GmailRepository with a pre-configured service.
// This is useful for testing with mock servers.
func NewGmailRepositoryWithService(service *gmail.Service, userID string) *GmailRepository {
	return &GmailRepository{
		service:     service,
		userID:      userID,
		maxRetries:  defaultMaxRetries,
		baseBackoff: defaultBaseBackoff,
	}
}

// List retrieves a list of messages matching the given options.
func (r *GmailRepository) List(ctx context.Context, opts mail.ListOptions) (*mail.ListResult[*mail.Message], error) {
	call := r.service.Users.Messages.List(r.userID)

	if opts.MaxResults > 0 {
		call = call.MaxResults(int64(opts.MaxResults))
	}
	if opts.PageToken != "" {
		call = call.PageToken(opts.PageToken)
	}
	if opts.Query != "" {
		call = call.Q(opts.Query)
	}
	if len(opts.LabelIDs) > 0 {
		call = call.LabelIds(opts.LabelIDs...)
	}

	response, err := call.Context(ctx).Do()
	if err != nil {
		return nil, r.handleError(err)
	}

	messages := make([]*mail.Message, 0, len(response.Messages))
	for _, gmailMsg := range response.Messages {
		// Fetch full message details
		fullMsg, err := r.Get(ctx, gmailMsg.Id)
		if err != nil {
			// Log error and continue with partial data
			messages = append(messages, &mail.Message{ID: gmailMsg.Id, ThreadID: gmailMsg.ThreadId})
			continue
		}
		messages = append(messages, fullMsg)
	}

	return &mail.ListResult[*mail.Message]{
		Items:         messages,
		NextPageToken: response.NextPageToken,
		Total:         int(response.ResultSizeEstimate),
	}, nil
}

// Get retrieves a single message by ID.
func (r *GmailRepository) Get(ctx context.Context, id string) (*mail.Message, error) {
	gmailMsg, err := r.service.Users.Messages.Get(r.userID, id).
		Format(gmailMessageFormat).
		Context(ctx).
		Do()
	if err != nil {
		return nil, r.handleError(err)
	}

	return gmailMessageToDomain(gmailMsg), nil
}

// Send sends a new message.
func (r *GmailRepository) Send(ctx context.Context, msg *mail.Message) (*mail.Message, error) {
	raw := buildMimeMessage(msg)
	encodedRaw := base64.URLEncoding.EncodeToString(raw)

	gmailMsg := &gmail.Message{
		Raw: encodedRaw,
	}

	sent, err := r.service.Users.Messages.Send(r.userID, gmailMsg).
		Context(ctx).
		Do()
	if err != nil {
		return nil, r.handleError(err)
	}

	// Fetch the sent message to get full details
	return r.Get(ctx, sent.Id)
}

// Reply sends a reply to an existing message.
func (r *GmailRepository) Reply(ctx context.Context, messageID string, reply *mail.Message) (*mail.Message, error) {
	// Get the original message to find the thread ID
	original, err := r.Get(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get original message: %w", err)
	}

	// Set the thread ID for the reply
	reply.ThreadID = original.ThreadID

	// Build the MIME message with References and In-Reply-To headers
	raw := buildReplyMimeMessage(reply, messageID)
	encodedRaw := base64.URLEncoding.EncodeToString(raw)

	gmailMsg := &gmail.Message{
		Raw:      encodedRaw,
		ThreadId: original.ThreadID,
	}

	sent, err := r.service.Users.Messages.Send(r.userID, gmailMsg).
		Context(ctx).
		Do()
	if err != nil {
		return nil, r.handleError(err)
	}

	return r.Get(ctx, sent.Id)
}

// Forward forwards an existing message.
func (r *GmailRepository) Forward(ctx context.Context, messageID string, forward *mail.Message) (*mail.Message, error) {
	// Get the original message to include in the forward body
	original, err := r.Get(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get original message: %w", err)
	}

	// Append the original message content to the forward body
	if forward.Body != "" {
		forward.Body = forward.Body + buildForwardBody(original)
	} else {
		forward.Body = buildForwardBody(original)
	}

	// Set subject if not provided
	if forward.Subject == "" {
		forward.Subject = "Fwd: " + original.Subject
	}

	return r.Send(ctx, forward)
}

// Trash moves a message to trash.
func (r *GmailRepository) Trash(ctx context.Context, id string) error {
	_, err := r.service.Users.Messages.Trash(r.userID, id).
		Context(ctx).
		Do()
	if err != nil {
		return r.handleError(err)
	}
	return nil
}

// Untrash removes a message from trash.
func (r *GmailRepository) Untrash(ctx context.Context, id string) error {
	_, err := r.service.Users.Messages.Untrash(r.userID, id).
		Context(ctx).
		Do()
	if err != nil {
		return r.handleError(err)
	}
	return nil
}

// Delete permanently deletes a message.
func (r *GmailRepository) Delete(ctx context.Context, id string) error {
	err := r.service.Users.Messages.Delete(r.userID, id).
		Context(ctx).
		Do()
	if err != nil {
		return r.handleError(err)
	}
	return nil
}

// Archive archives a message by removing the INBOX label.
func (r *GmailRepository) Archive(ctx context.Context, id string) error {
	_, err := r.Modify(ctx, id, mail.ModifyRequest{
		RemoveLabels: []string{gmailLabelInbox},
	})
	return err
}

// Modify modifies the labels on a message.
func (r *GmailRepository) Modify(ctx context.Context, id string, req mail.ModifyRequest) (*mail.Message, error) {
	modifyReq := &gmail.ModifyMessageRequest{
		AddLabelIds:    req.AddLabels,
		RemoveLabelIds: req.RemoveLabels,
	}

	gmailMsg, err := r.service.Users.Messages.Modify(r.userID, id, modifyReq).
		Context(ctx).
		Do()
	if err != nil {
		return nil, r.handleError(err)
	}

	return gmailMessageToDomain(gmailMsg), nil
}

// Search searches for messages matching the query.
func (r *GmailRepository) Search(ctx context.Context, query string, opts mail.ListOptions) (*mail.ListResult[*mail.Message], error) {
	opts.Query = query
	return r.List(ctx, opts)
}

// handleError maps Gmail API errors to domain errors.
func (r *GmailRepository) handleError(err error) error {
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) {
		return mapGmailError(apiErr.Code, apiErr.Message)
	}
	return fmt.Errorf("gmail error: %w", err)
}

// mapGmailError maps HTTP status codes to domain errors.
func mapGmailError(statusCode int, message string) error {
	switch statusCode {
	case http.StatusNotFound:
		return fmt.Errorf("%w: %s", mail.ErrMessageNotFound, message)
	case http.StatusBadRequest:
		return fmt.Errorf("%w: %s", ErrBadRequest, message)
	case http.StatusTooManyRequests:
		return fmt.Errorf("%w: %s", ErrRateLimited, message)
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return fmt.Errorf("%w: %s", ErrTemporary, message)
	default:
		return fmt.Errorf("gmail API error (status %d): %s", statusCode, message)
	}
}

// gmailMessageToDomain converts a Gmail API message to a domain Message.
func gmailMessageToDomain(msg *gmail.Message) *mail.Message {
	if msg == nil {
		return nil
	}

	result := &mail.Message{
		ID:       msg.Id,
		ThreadID: msg.ThreadId,
		Snippet:  msg.Snippet,
		Labels:   msg.LabelIds,
	}

	// Initialize slices
	result.To = []string{}
	result.Cc = []string{}
	result.Bcc = []string{}
	if result.Labels == nil {
		result.Labels = []string{}
	}

	// Determine read and starred status from labels
	result.IsRead = !hasLabel(msg.LabelIds, gmailLabelUnread)
	result.IsStarred = hasLabel(msg.LabelIds, gmailLabelStarred)

	// Parse headers and body from payload
	if msg.Payload != nil {
		from, to, subject, date := parseHeaders(msg.Payload.Headers)
		result.From = from
		result.Subject = subject
		result.Date = date

		// Parse recipients
		if to != "" {
			result.To = parseRecipients(to)
		}

		// Get Cc from headers
		for _, header := range msg.Payload.Headers {
			if strings.EqualFold(header.Name, "Cc") {
				result.Cc = parseRecipients(header.Value)
				break
			}
		}

		// Extract body content
		result.Body, result.BodyHTML = extractBody(msg.Payload)
	}

	return result
}

// domainMessageToGmail converts a domain Message to a Gmail API message.
func domainMessageToGmail(msg *mail.Message) *gmail.Message {
	if msg == nil {
		return nil
	}

	raw := buildMimeMessage(msg)
	encodedRaw := base64.URLEncoding.EncodeToString(raw)

	return &gmail.Message{
		Id:       msg.ID,
		ThreadId: msg.ThreadID,
		Raw:      encodedRaw,
		LabelIds: msg.Labels,
	}
}

// parseHeaders extracts common headers from Gmail message headers.
func parseHeaders(headers []*gmail.MessagePartHeader) (from, to, subject string, date time.Time) {
	for _, header := range headers {
		switch strings.ToLower(header.Name) {
		case "from":
			from = header.Value
		case "to":
			to = header.Value
		case "subject":
			subject = header.Value
		case "date":
			// Try parsing RFC 2822 date format
			parsed, err := time.Parse(time.RFC1123Z, header.Value)
			if err == nil {
				date = parsed
			} else {
				// Try alternative formats
				parsed, err = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", header.Value)
				if err == nil {
					date = parsed
				}
			}
		}
	}
	return
}

// parseRecipients parses a comma-separated list of email addresses.
func parseRecipients(addresses string) []string {
	if addresses == "" {
		return []string{}
	}

	parts := strings.Split(addresses, ",")
	recipients := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			recipients = append(recipients, trimmed)
		}
	}
	return recipients
}

// hasLabel checks if a label exists in the label list.
func hasLabel(labels []string, target string) bool {
	for _, label := range labels {
		if label == target {
			return true
		}
	}
	return false
}

// extractBody extracts plain text and HTML body from message payload.
func extractBody(payload *gmail.MessagePart) (plain, html string) {
	if payload == nil {
		return "", ""
	}

	// Single part message
	if payload.Body != nil && payload.Body.Data != "" {
		decoded, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			if strings.HasPrefix(payload.MimeType, "text/html") {
				return "", string(decoded)
			}
			return string(decoded), ""
		}
	}

	// Multipart message
	if len(payload.Parts) > 0 {
		for _, part := range payload.Parts {
			partPlain, partHTML := extractBodyFromPart(part)
			if partPlain != "" && plain == "" {
				plain = partPlain
			}
			if partHTML != "" && html == "" {
				html = partHTML
			}
		}
	}

	return plain, html
}

// extractBodyFromPart recursively extracts body content from a message part.
func extractBodyFromPart(part *gmail.MessagePart) (plain, html string) {
	if part == nil {
		return "", ""
	}

	// Recursively handle nested multipart
	if strings.HasPrefix(part.MimeType, "multipart/") && len(part.Parts) > 0 {
		for _, subpart := range part.Parts {
			subPlain, subHTML := extractBodyFromPart(subpart)
			if subPlain != "" && plain == "" {
				plain = subPlain
			}
			if subHTML != "" && html == "" {
				html = subHTML
			}
		}
		return plain, html
	}

	// Extract content from leaf parts
	if part.Body != nil && part.Body.Data != "" {
		decoded, err := base64.URLEncoding.DecodeString(part.Body.Data)
		if err == nil {
			content := string(decoded)
			switch part.MimeType {
			case "text/plain":
				return content, ""
			case "text/html":
				return "", content
			}
		}
	}

	return "", ""
}

// buildMimeMessage constructs a MIME message from a domain Message.
func buildMimeMessage(msg *mail.Message) []byte {
	var builder strings.Builder

	// Write headers
	builder.WriteString(fmt.Sprintf("From: %s\r\n", msg.From))
	builder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(msg.To, ", ")))
	if len(msg.Cc) > 0 {
		builder.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(msg.Cc, ", ")))
	}
	if len(msg.Bcc) > 0 {
		builder.WriteString(fmt.Sprintf("Bcc: %s\r\n", strings.Join(msg.Bcc, ", ")))
	}
	builder.WriteString(fmt.Sprintf("Subject: %s\r\n", msg.Subject))
	builder.WriteString("MIME-Version: 1.0\r\n")

	// Determine content type
	if msg.BodyHTML != "" {
		builder.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
		builder.WriteString("\r\n")
		builder.WriteString(msg.BodyHTML)
	} else {
		builder.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
		builder.WriteString("\r\n")
		builder.WriteString(msg.Body)
	}

	return []byte(builder.String())
}

// buildReplyMimeMessage constructs a MIME message for a reply.
func buildReplyMimeMessage(msg *mail.Message, originalMessageID string) []byte {
	var builder strings.Builder

	// Write headers
	builder.WriteString(fmt.Sprintf("From: %s\r\n", msg.From))
	builder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(msg.To, ", ")))
	if len(msg.Cc) > 0 {
		builder.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(msg.Cc, ", ")))
	}
	builder.WriteString(fmt.Sprintf("Subject: %s\r\n", msg.Subject))
	builder.WriteString(fmt.Sprintf("In-Reply-To: <%s>\r\n", originalMessageID))
	builder.WriteString(fmt.Sprintf("References: <%s>\r\n", originalMessageID))
	builder.WriteString("MIME-Version: 1.0\r\n")

	// Determine content type
	if msg.BodyHTML != "" {
		builder.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
		builder.WriteString("\r\n")
		builder.WriteString(msg.BodyHTML)
	} else {
		builder.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
		builder.WriteString("\r\n")
		builder.WriteString(msg.Body)
	}

	return []byte(builder.String())
}

// buildForwardBody creates the body text for a forwarded message.
func buildForwardBody(original *mail.Message) string {
	var builder strings.Builder

	builder.WriteString("\r\n\r\n---------- Forwarded message ---------\r\n")
	builder.WriteString(fmt.Sprintf("From: %s\r\n", original.From))
	builder.WriteString(fmt.Sprintf("Date: %s\r\n", original.Date.Format(time.RFC1123Z)))
	builder.WriteString(fmt.Sprintf("Subject: %s\r\n", original.Subject))
	builder.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(original.To, ", ")))
	builder.WriteString("\r\n")
	builder.WriteString(original.Body)

	return builder.String()
}

// retryWithBackoff executes a function with exponential backoff retry.
func retryWithBackoff[T any](ctx context.Context, maxRetries int, baseBackoff time.Duration, fn func() (T, error)) (T, error) {
	var zero T
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}

		// Check if error is retryable
		if !isRetryableError(err) {
			return zero, err
		}

		lastErr = err

		// Calculate backoff duration with exponential increase
		backoff := baseBackoff * time.Duration(1<<attempt)

		// Wait for backoff or context cancellation
		select {
		case <-ctx.Done():
			return zero, fmt.Errorf("context cancelled during retry: %w", ctx.Err())
		case <-time.After(backoff):
			// Continue to next attempt
		}
	}

	return zero, fmt.Errorf("max retries (%d) exceeded: %w", maxRetries, lastErr)
}

// isRetryableError determines if an error should trigger a retry.
func isRetryableError(err error) bool {
	return errors.Is(err, ErrTemporary) || errors.Is(err, ErrRateLimited)
}

// =============================================================================
// DraftRepository Implementation (GmailDraftRepository)
// =============================================================================

// List retrieves a list of drafts.
func (r *GmailDraftRepository) List(ctx context.Context, opts mail.ListOptions) (*mail.ListResult[*mail.Draft], error) {
	call := r.service.Users.Drafts.List(r.userID)

	if opts.MaxResults > 0 {
		call = call.MaxResults(int64(opts.MaxResults))
	}
	if opts.PageToken != "" {
		call = call.PageToken(opts.PageToken)
	}
	if opts.Query != "" {
		call = call.Q(opts.Query)
	}

	response, err := call.Context(ctx).Do()
	if err != nil {
		return nil, r.handleError(err)
	}

	drafts := make([]*mail.Draft, 0, len(response.Drafts))
	for _, gmailDraft := range response.Drafts {
		// Fetch full draft details
		fullDraft, err := r.Get(ctx, gmailDraft.Id)
		if err != nil {
			// Create minimal draft on error
			drafts = append(drafts, &mail.Draft{ID: gmailDraft.Id})
			continue
		}
		drafts = append(drafts, fullDraft)
	}

	return &mail.ListResult[*mail.Draft]{
		Items:         drafts,
		NextPageToken: response.NextPageToken,
		Total:         int(response.ResultSizeEstimate),
	}, nil
}

// Get retrieves a single draft by ID.
func (r *GmailDraftRepository) Get(ctx context.Context, id string) (*mail.Draft, error) {
	gmailDraft, err := r.service.Users.Drafts.Get(r.userID, id).
		Format(gmailMessageFormat).
		Context(ctx).
		Do()
	if err != nil {
		return nil, r.handleDraftError(err)
	}

	return gmailDraftToDomain(gmailDraft), nil
}

// Create creates a new draft.
func (r *GmailDraftRepository) Create(ctx context.Context, draft *mail.Draft) (*mail.Draft, error) {
	if draft.Message == nil {
		return nil, fmt.Errorf("%w: draft must have a message", ErrBadRequest)
	}

	raw := buildMimeMessage(draft.Message)
	encodedRaw := base64.URLEncoding.EncodeToString(raw)

	gmailDraft := &gmail.Draft{
		Message: &gmail.Message{
			Raw: encodedRaw,
		},
	}

	created, err := r.service.Users.Drafts.Create(r.userID, gmailDraft).
		Context(ctx).
		Do()
	if err != nil {
		return nil, r.handleError(err)
	}

	// Fetch the created draft to get full details
	return r.Get(ctx, created.Id)
}

// Update updates an existing draft.
func (r *GmailDraftRepository) Update(ctx context.Context, draft *mail.Draft) (*mail.Draft, error) {
	if draft.Message == nil {
		return nil, fmt.Errorf("%w: draft must have a message", ErrBadRequest)
	}

	raw := buildMimeMessage(draft.Message)
	encodedRaw := base64.URLEncoding.EncodeToString(raw)

	gmailDraft := &gmail.Draft{
		Message: &gmail.Message{
			Raw: encodedRaw,
		},
	}

	updated, err := r.service.Users.Drafts.Update(r.userID, draft.ID, gmailDraft).
		Context(ctx).
		Do()
	if err != nil {
		return nil, r.handleDraftError(err)
	}

	// Fetch the updated draft to get full details
	return r.Get(ctx, updated.Id)
}

// Send sends a draft.
func (r *GmailDraftRepository) Send(ctx context.Context, id string) (*mail.Message, error) {
	gmailDraft := &gmail.Draft{
		Id: id,
	}

	sent, err := r.service.Users.Drafts.Send(r.userID, gmailDraft).
		Context(ctx).
		Do()
	if err != nil {
		return nil, r.handleDraftError(err)
	}

	// Fetch the sent message to get full details using the message API
	gmailMsg, err := r.service.Users.Messages.Get(r.userID, sent.Id).
		Format(gmailMessageFormat).
		Context(ctx).
		Do()
	if err != nil {
		return nil, r.handleError(err)
	}

	return gmailMessageToDomain(gmailMsg), nil
}

// Delete deletes a draft.
func (r *GmailDraftRepository) Delete(ctx context.Context, id string) error {
	err := r.service.Users.Drafts.Delete(r.userID, id).
		Context(ctx).
		Do()
	if err != nil {
		return r.handleDraftError(err)
	}
	return nil
}

// handleDraftError maps Gmail API errors to domain draft errors.
func (r *GmailRepository) handleDraftError(err error) error {
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) {
		if apiErr.Code == http.StatusNotFound {
			return fmt.Errorf("%w: %s", mail.ErrDraftNotFound, apiErr.Message)
		}
		return mapGmailError(apiErr.Code, apiErr.Message)
	}
	return fmt.Errorf("gmail error: %w", err)
}

// gmailDraftToDomain converts a Gmail API draft to a domain Draft.
func gmailDraftToDomain(draft *gmail.Draft) *mail.Draft {
	if draft == nil {
		return nil
	}

	result := &mail.Draft{
		ID: draft.Id,
	}

	if draft.Message != nil {
		result.Message = gmailMessageToDomain(draft.Message)
	}

	return result
}

// =============================================================================
// LabelRepository Implementation (GmailLabelRepository)
// =============================================================================

// List retrieves all labels.
func (r *GmailLabelRepository) List(ctx context.Context) ([]*mail.Label, error) {
	response, err := r.service.Users.Labels.List(r.userID).
		Context(ctx).
		Do()
	if err != nil {
		return nil, r.handleError(err)
	}

	labels := make([]*mail.Label, 0, len(response.Labels))
	for _, gmailLabel := range response.Labels {
		labels = append(labels, gmailLabelToDomain(gmailLabel))
	}

	return labels, nil
}

// Get retrieves a single label by ID.
func (r *GmailLabelRepository) Get(ctx context.Context, id string) (*mail.Label, error) {
	gmailLabel, err := r.service.Users.Labels.Get(r.userID, id).
		Context(ctx).
		Do()
	if err != nil {
		return nil, r.handleLabelError(err)
	}

	return gmailLabelToDomain(gmailLabel), nil
}

// GetByName retrieves a label by name.
func (r *GmailLabelRepository) GetByName(ctx context.Context, name string) (*mail.Label, error) {
	labels, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, label := range labels {
		if label.Name == name {
			return label, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", mail.ErrLabelNotFound, name)
}

// Create creates a new label.
func (r *GmailLabelRepository) Create(ctx context.Context, label *mail.Label) (*mail.Label, error) {
	gmailLabel := domainLabelToGmail(label)

	created, err := r.service.Users.Labels.Create(r.userID, gmailLabel).
		Context(ctx).
		Do()
	if err != nil {
		return nil, r.handleError(err)
	}

	return gmailLabelToDomain(created), nil
}

// Update updates an existing label.
func (r *GmailLabelRepository) Update(ctx context.Context, label *mail.Label) (*mail.Label, error) {
	gmailLabel := domainLabelToGmail(label)

	updated, err := r.service.Users.Labels.Update(r.userID, label.ID, gmailLabel).
		Context(ctx).
		Do()
	if err != nil {
		return nil, r.handleLabelError(err)
	}

	return gmailLabelToDomain(updated), nil
}

// Delete deletes a label.
func (r *GmailLabelRepository) Delete(ctx context.Context, id string) error {
	err := r.service.Users.Labels.Delete(r.userID, id).
		Context(ctx).
		Do()
	if err != nil {
		return r.handleLabelError(err)
	}
	return nil
}

// handleLabelError maps Gmail API errors to domain label errors.
func (r *GmailRepository) handleLabelError(err error) error {
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) {
		if apiErr.Code == http.StatusNotFound {
			return fmt.Errorf("%w: %s", mail.ErrLabelNotFound, apiErr.Message)
		}
		return mapGmailError(apiErr.Code, apiErr.Message)
	}
	return fmt.Errorf("gmail error: %w", err)
}

// gmailLabelToDomain converts a Gmail API label to a domain Label.
func gmailLabelToDomain(label *gmail.Label) *mail.Label {
	if label == nil {
		return nil
	}

	result := &mail.Label{
		ID:                    label.Id,
		Name:                  label.Name,
		Type:                  label.Type,
		MessageListVisibility: label.MessageListVisibility,
		LabelListVisibility:   label.LabelListVisibility,
	}

	if label.Color != nil {
		result.Color = &mail.LabelColor{
			Background: label.Color.BackgroundColor,
			Text:       label.Color.TextColor,
		}
	}

	return result
}

// domainLabelToGmail converts a domain Label to a Gmail API label.
func domainLabelToGmail(label *mail.Label) *gmail.Label {
	if label == nil {
		return nil
	}

	result := &gmail.Label{
		Id:                    label.ID,
		Name:                  label.Name,
		Type:                  label.Type,
		MessageListVisibility: label.MessageListVisibility,
		LabelListVisibility:   label.LabelListVisibility,
	}

	if label.Color != nil {
		result.Color = &gmail.LabelColor{
			BackgroundColor: label.Color.Background,
			TextColor:       label.Color.Text,
		}
	}

	return result
}
