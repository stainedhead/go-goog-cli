// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestMailActionsCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "mail") {
		t.Error("expected output to contain 'mail'")
	}
	if !contains(output, "trash") {
		t.Error("expected output to contain 'trash'")
	}
	if !contains(output, "archive") {
		t.Error("expected output to contain 'archive'")
	}
	if !contains(output, "delete") {
		t.Error("expected output to contain 'delete'")
	}
	if !contains(output, "modify") {
		t.Error("expected output to contain 'modify'")
	}
	if !contains(output, "mark") {
		t.Error("expected output to contain 'mark'")
	}
}

func TestMailTrashCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "trash", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "trash") {
		t.Error("expected output to contain 'trash'")
	}
	if !contains(output, "message") {
		t.Error("expected output to contain 'message'")
	}
}

func TestMailUntrashCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "untrash", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "untrash") {
		t.Error("expected output to contain 'untrash'")
	}
	if !contains(output, "restore") || !contains(output, "trash") {
		t.Error("expected output to contain 'restore' or 'trash'")
	}
}

func TestMailArchiveCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "archive", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "archive") {
		t.Error("expected output to contain 'archive'")
	}
	if !contains(output, "INBOX") {
		t.Error("expected output to contain 'INBOX'")
	}
}

func TestMailDeleteCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "delete", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "delete") {
		t.Error("expected output to contain 'delete'")
	}
	if !contains(output, "--confirm") {
		t.Error("expected output to contain '--confirm'")
	}
	if !contains(output, "permanent") {
		t.Error("expected output to contain 'permanent'")
	}
}

func TestMailModifyCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "modify", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "modify") {
		t.Error("expected output to contain 'modify'")
	}
	if !contains(output, "--add-labels") {
		t.Error("expected output to contain '--add-labels'")
	}
	if !contains(output, "--remove-labels") {
		t.Error("expected output to contain '--remove-labels'")
	}
}

func TestMailMarkCmd_Help(t *testing.T) {
	cmd := &cobra.Command{Use: "goog"}
	cmd.AddCommand(mailCmd)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"mail", "mark", "--help"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !contains(output, "mark") {
		t.Error("expected output to contain 'mark'")
	}
	if !contains(output, "--read") {
		t.Error("expected output to contain '--read'")
	}
	if !contains(output, "--unread") {
		t.Error("expected output to contain '--unread'")
	}
	if !contains(output, "--star") {
		t.Error("expected output to contain '--star'")
	}
	if !contains(output, "--unstar") {
		t.Error("expected output to contain '--unstar'")
	}
}

func TestMailDeleteCmd_RequiresConfirmFlag(t *testing.T) {
	// Test that PreRunE validates the --confirm flag
	// We test the validation logic directly since Cobra flag parsing
	// behavior varies in test contexts
	mailDeleteConfirm = false

	mockCmd := &cobra.Command{Use: "test"}

	if mailDeleteCmd.PreRunE != nil {
		err := mailDeleteCmd.PreRunE(mockCmd, []string{"msg123"})
		if err == nil {
			t.Error("expected error when --confirm flag is not set")
		}
	} else {
		t.Error("mailDeleteCmd should have PreRunE defined")
	}
}

func TestMailMarkCmd_RequiresAtLeastOneFlag(t *testing.T) {
	// Test that PreRunE validates at least one flag is set
	mailMarkRead = false
	mailMarkUnread = false
	mailMarkStar = false
	mailMarkUnstar = false

	mockCmd := &cobra.Command{Use: "test"}

	if mailMarkCmd.PreRunE != nil {
		err := mailMarkCmd.PreRunE(mockCmd, []string{"msg123"})
		if err == nil {
			t.Error("expected error when no flags are set")
		}
	} else {
		t.Error("mailMarkCmd should have PreRunE defined")
	}
}

func TestMailModifyCmd_RequiresAtLeastOneFlag(t *testing.T) {
	// Test that PreRunE validates at least one flag is set
	mailModifyAddLabels = nil
	mailModifyRemoveLabels = nil

	mockCmd := &cobra.Command{Use: "test"}

	if mailModifyCmd.PreRunE != nil {
		err := mailModifyCmd.PreRunE(mockCmd, []string{"msg123"})
		if err == nil {
			t.Error("expected error when no flags are set")
		}
	} else {
		t.Error("mailModifyCmd should have PreRunE defined")
	}
}

func TestMailMarkCmd_ConflictingReadFlags(t *testing.T) {
	// Test that --read and --unread cannot be used together
	mailMarkRead = true
	mailMarkUnread = true
	mailMarkStar = false
	mailMarkUnstar = false

	mockCmd := &cobra.Command{Use: "test"}

	if mailMarkCmd.PreRunE != nil {
		err := mailMarkCmd.PreRunE(mockCmd, []string{"msg123"})
		if err == nil {
			t.Error("expected error when both --read and --unread are set")
		}
	} else {
		t.Error("mailMarkCmd should have PreRunE defined")
	}
}

func TestMailMarkCmd_ConflictingStarFlags(t *testing.T) {
	// Test that --star and --unstar cannot be used together
	mailMarkRead = false
	mailMarkUnread = false
	mailMarkStar = true
	mailMarkUnstar = true

	mockCmd := &cobra.Command{Use: "test"}

	if mailMarkCmd.PreRunE != nil {
		err := mailMarkCmd.PreRunE(mockCmd, []string{"msg123"})
		if err == nil {
			t.Error("expected error when both --star and --unstar are set")
		}
	} else {
		t.Error("mailMarkCmd should have PreRunE defined")
	}
}

func TestMailTrashCmd_RequiresMessageID(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if mailTrashCmd.Args == nil {
		t.Error("mailTrashCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := mailTrashCmd.Args(mailTrashCmd, []string{})
	if err == nil {
		t.Error("expected error when message ID is missing")
	}
}

func TestMailUntrashCmd_RequiresMessageID(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if mailUntrashCmd.Args == nil {
		t.Error("mailUntrashCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := mailUntrashCmd.Args(mailUntrashCmd, []string{})
	if err == nil {
		t.Error("expected error when message ID is missing")
	}
}

func TestMailArchiveCmd_RequiresMessageID(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if mailArchiveCmd.Args == nil {
		t.Error("mailArchiveCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := mailArchiveCmd.Args(mailArchiveCmd, []string{})
	if err == nil {
		t.Error("expected error when message ID is missing")
	}
}

func TestMailDeleteCmd_RequiresMessageID(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if mailDeleteCmd.Args == nil {
		t.Error("mailDeleteCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := mailDeleteCmd.Args(mailDeleteCmd, []string{})
	if err == nil {
		t.Error("expected error when message ID is missing")
	}
}

func TestMailModifyCmd_RequiresMessageID(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if mailModifyCmd.Args == nil {
		t.Error("mailModifyCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := mailModifyCmd.Args(mailModifyCmd, []string{})
	if err == nil {
		t.Error("expected error when message ID is missing")
	}
}

func TestMailMarkCmd_RequiresMessageID(t *testing.T) {
	// Test that Args validator requires exactly one argument
	if mailMarkCmd.Args == nil {
		t.Error("mailMarkCmd should have Args validator defined")
		return
	}

	// Test with no args - should fail
	err := mailMarkCmd.Args(mailMarkCmd, []string{})
	if err == nil {
		t.Error("expected error when message ID is missing")
	}
}
