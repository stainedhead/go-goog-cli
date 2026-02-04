// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/domain/mail"
)

// Mail compose command flags.
var (
	// Send flags
	mailSendTo      []string
	mailSendCc      []string
	mailSendBcc     []string
	mailSendSubject string
	mailSendBody    string
	mailSendHTML    bool

	// Reply flags
	mailReplyBody string
	mailReplyAll  bool

	// Forward flags
	mailForwardTo   []string
	mailForwardBody string
)

// mailSendCmd handles sending new messages.
var mailSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a new email message",
	Long: `Send a new email message.

Compose and send a new email to one or more recipients.
The --to flag is required and can be specified multiple times.`,
	Example: `  # Send a simple message
  goog mail send --to user@example.com --subject "Hello" --body "Hi there!"

  # Send to multiple recipients with CC
  goog mail send --to user1@example.com --to user2@example.com \
    --cc manager@example.com --subject "Update" --body "Project update"

  # Send HTML content
  goog mail send --to user@example.com --subject "Report" \
    --body "<h1>Report</h1><p>See attached.</p>" --html

  # Send using a specific account
  goog mail send --to user@example.com --subject "Hello" --body "Hi" --account work`,
	RunE: runMailSend,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(mailSendTo) == 0 {
			return fmt.Errorf("required flag \"to\" not set")
		}
		return nil
	},
}

// mailReplyCmd handles replying to messages.
var mailReplyCmd = &cobra.Command{
	Use:   "reply <id>",
	Short: "Reply to an email message",
	Long: `Reply to an existing email message.

Send a reply to the specified message. Use --all to reply to all
recipients (reply-all). The reply will be part of the same thread
as the original message.`,
	Example: `  # Reply to a message
  goog mail reply abc123 --body "Thanks for your message!"

  # Reply-all
  goog mail reply abc123 --body "I agree with everyone." --all

  # Reply using a specific account
  goog mail reply abc123 --body "Got it!" --account work`,
	Args: cobra.ExactArgs(1),
	RunE: runMailReply,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if mailReplyBody == "" {
			return fmt.Errorf("required flag \"body\" not set")
		}
		return nil
	},
}

// mailForwardCmd handles forwarding messages.
var mailForwardCmd = &cobra.Command{
	Use:   "forward <id>",
	Short: "Forward an email message",
	Long: `Forward an existing email message.

Forward the specified message to one or more recipients.
The --to flag is required. An optional intro message can be
added using --body.`,
	Example: `  # Forward a message
  goog mail forward abc123 --to colleague@example.com

  # Forward with an intro message
  goog mail forward abc123 --to colleague@example.com \
    --body "FYI - see below"

  # Forward to multiple recipients
  goog mail forward abc123 --to user1@example.com --to user2@example.com

  # Forward using a specific account
  goog mail forward abc123 --to user@example.com --account work`,
	Args: cobra.ExactArgs(1),
	RunE: runMailForward,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(mailForwardTo) == 0 {
			return fmt.Errorf("required flag \"to\" not set")
		}
		return nil
	},
}

func init() {
	// Add mail subcommands
	mailCmd.AddCommand(mailSendCmd)
	mailCmd.AddCommand(mailReplyCmd)
	mailCmd.AddCommand(mailForwardCmd)

	// Send command flags
	mailSendCmd.Flags().StringSliceVar(&mailSendTo, "to", nil, "recipient email address(es) (required)")
	mailSendCmd.Flags().StringSliceVar(&mailSendCc, "cc", nil, "CC recipient email address(es)")
	mailSendCmd.Flags().StringSliceVar(&mailSendBcc, "bcc", nil, "BCC recipient email address(es)")
	mailSendCmd.Flags().StringVar(&mailSendSubject, "subject", "", "email subject")
	mailSendCmd.Flags().StringVar(&mailSendBody, "body", "", "email body content")
	mailSendCmd.Flags().BoolVar(&mailSendHTML, "html", false, "treat body as HTML content")

	// Reply command flags
	mailReplyCmd.Flags().StringVar(&mailReplyBody, "body", "", "reply body content (required)")
	mailReplyCmd.Flags().BoolVar(&mailReplyAll, "all", false, "reply to all recipients")

	// Forward command flags
	mailForwardCmd.Flags().StringSliceVar(&mailForwardTo, "to", nil, "recipient email address(es) (required)")
	mailForwardCmd.Flags().StringVar(&mailForwardBody, "body", "", "intro message to add before forwarded content")
}

// runMailSend handles the mail send command.
func runMailSend(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get Gmail repository
	repo, senderEmail, err := getGmailRepository(ctx)
	if err != nil {
		return err
	}

	// Parse and validate recipients
	toRecipients, err := parseEmailRecipients(mailSendTo)
	if err != nil {
		return fmt.Errorf("invalid 'to' recipient: %w", err)
	}

	ccRecipients, err := parseEmailRecipients(mailSendCc)
	if err != nil {
		return fmt.Errorf("invalid 'cc' recipient: %w", err)
	}

	bccRecipients, err := parseEmailRecipients(mailSendBcc)
	if err != nil {
		return fmt.Errorf("invalid 'bcc' recipient: %w", err)
	}

	// Build message
	msg := &mail.Message{
		From:    senderEmail,
		To:      toRecipients,
		Cc:      ccRecipients,
		Bcc:     bccRecipients,
		Subject: mailSendSubject,
	}

	if mailSendHTML {
		msg.BodyHTML = mailSendBody
	} else {
		msg.Body = mailSendBody
	}

	// Send message
	sent, err := repo.Send(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	cmd.Printf("Message sent successfully.\n")
	cmd.Printf("Message ID: %s\n", sent.ID)
	cmd.Printf("Thread ID: %s\n", sent.ThreadID)

	return nil
}

// runMailReply handles the mail reply command.
func runMailReply(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	messageID := args[0]

	// Get Gmail repository
	repo, senderEmail, err := getGmailRepository(ctx)
	if err != nil {
		return err
	}

	// Get original message to determine recipients
	original, err := repo.Get(ctx, messageID)
	if err != nil {
		return fmt.Errorf("failed to get original message: %w", err)
	}

	// Build reply message
	reply := &mail.Message{
		From:    senderEmail,
		Body:    mailReplyBody,
		Subject: buildReplySubject(original.Subject),
	}

	// Set recipients based on reply-all flag
	if mailReplyAll {
		// Reply to sender and all original recipients (except ourselves)
		reply.To = []string{original.From}
		for _, to := range original.To {
			if to != senderEmail {
				reply.To = append(reply.To, to)
			}
		}
		reply.Cc = original.Cc
	} else {
		// Reply only to sender
		reply.To = []string{original.From}
	}

	// Send reply
	sent, err := repo.Reply(ctx, messageID, reply)
	if err != nil {
		return fmt.Errorf("failed to send reply: %w", err)
	}

	cmd.Printf("Reply sent successfully.\n")
	cmd.Printf("Message ID: %s\n", sent.ID)
	cmd.Printf("Thread ID: %s\n", sent.ThreadID)

	return nil
}

// runMailForward handles the mail forward command.
func runMailForward(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	messageID := args[0]

	// Get Gmail repository
	repo, senderEmail, err := getGmailRepository(ctx)
	if err != nil {
		return err
	}

	// Parse and validate recipients
	toRecipients, err := parseEmailRecipients(mailForwardTo)
	if err != nil {
		return fmt.Errorf("invalid 'to' recipient: %w", err)
	}

	// Build forward message
	forward := &mail.Message{
		From: senderEmail,
		To:   toRecipients,
		Body: mailForwardBody,
	}

	// Send forward
	sent, err := repo.Forward(ctx, messageID, forward)
	if err != nil {
		return fmt.Errorf("failed to forward message: %w", err)
	}

	cmd.Printf("Message forwarded successfully.\n")
	cmd.Printf("Message ID: %s\n", sent.ID)
	cmd.Printf("Thread ID: %s\n", sent.ThreadID)

	return nil
}

// parseEmailRecipients cleans, validates, and returns email recipients.
// Returns an error if any email address is invalid.
func parseEmailRecipients(recipients []string) ([]string, error) {
	if recipients == nil {
		return []string{}, nil
	}

	result := make([]string, 0, len(recipients))
	for _, r := range recipients {
		trimmed := strings.TrimSpace(r)
		if trimmed != "" {
			if !isValidEmail(trimmed) {
				return nil, fmt.Errorf("invalid email address: %q", trimmed)
			}
			result = append(result, trimmed)
		}
	}
	return result, nil
}

// buildReplySubject prepends "Re: " to the subject if not already present.
func buildReplySubject(subject string) string {
	if strings.HasPrefix(strings.ToLower(subject), "re:") {
		return subject
	}
	return "Re: " + subject
}
