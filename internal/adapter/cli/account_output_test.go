// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	accountuc "github.com/stainedhead/go-goog-cli/internal/usecase/account"
)

// =============================================================================
// Tests for account output functions
// =============================================================================

func TestOutputAccountsJSON(t *testing.T) {
	accounts := []*accountuc.Account{
		{
			Alias:     "default",
			Email:     "user@example.com",
			IsDefault: true,
			Added:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Scopes:    []string{"gmail.readonly", "calendar.readonly"},
		},
		{
			Alias:     "work",
			Email:     "work@company.com",
			IsDefault: false,
			Added:     time.Date(2024, 2, 20, 14, 30, 0, 0, time.UTC),
			Scopes:    []string{},
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := outputAccountsJSON(cmd, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var output bytes.Buffer
	_, err = output.ReadFrom(r)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	outputStr := output.String()

	// Check JSON structure
	if !strings.Contains(outputStr, "\"alias\": \"default\"") {
		t.Errorf("expected output to contain alias 'default', got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "\"email\": \"user@example.com\"") {
		t.Errorf("expected output to contain email 'user@example.com', got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "\"is_default\": true") {
		t.Errorf("expected output to contain is_default: true, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "\"scopes\":") {
		t.Errorf("expected output to contain scopes for default account, got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "\"alias\": \"work\"") {
		t.Errorf("expected output to contain alias 'work', got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "\"is_default\": false") {
		t.Errorf("expected output to contain is_default: false for work account, got: %s", outputStr)
	}
}

func TestOutputAccountsJSON_EmptyList(t *testing.T) {
	accounts := []*accountuc.Account{}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := outputAccountsJSON(cmd, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var output bytes.Buffer
	_, err = output.ReadFrom(r)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	outputStr := strings.TrimSpace(output.String())

	// Should output empty array
	if outputStr != "[\n]" {
		t.Errorf("expected empty array '[]', got: %s", outputStr)
	}
}

func TestOutputAccountsJSON_SingleAccount(t *testing.T) {
	accounts := []*accountuc.Account{
		{
			Alias:     "test",
			Email:     "test@example.com",
			IsDefault: true,
			Added:     time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC),
			Scopes:    []string{"gmail.readonly"},
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := outputAccountsJSON(cmd, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var output bytes.Buffer
	_, err = output.ReadFrom(r)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	outputStr := output.String()

	// Should contain valid JSON structure with one account
	if !strings.Contains(outputStr, "\"alias\": \"test\"") {
		t.Errorf("expected output to contain test alias, got: %s", outputStr)
	}
	// Should not have comma at the end (single item)
	if strings.Contains(outputStr, "},\n]") {
		t.Errorf("expected no trailing comma for single item, got: %s", outputStr)
	}
}

func TestOutputAccountsJSON_NoScopes(t *testing.T) {
	accounts := []*accountuc.Account{
		{
			Alias:     "test",
			Email:     "test@example.com",
			IsDefault: false,
			Added:     time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC),
			Scopes:    nil,
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := outputAccountsJSON(cmd, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var output bytes.Buffer
	_, err = output.ReadFrom(r)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	outputStr := output.String()

	// Should NOT contain scopes field when empty/nil
	if strings.Contains(outputStr, "\"scopes\":") {
		t.Errorf("expected output to not contain scopes when nil, got: %s", outputStr)
	}
}

func TestOutputAccountsTable_SingleAccount(t *testing.T) {
	accounts := []*accountuc.Account{
		{
			Alias:     "default",
			Email:     "user@example.com",
			IsDefault: true,
			Added:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
	}

	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := outputAccountsTable(cmd, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// Check table structure
	if !strings.Contains(output, "ALIAS") {
		t.Error("expected output to contain 'ALIAS' header")
	}
	if !strings.Contains(output, "EMAIL") {
		t.Error("expected output to contain 'EMAIL' header")
	}
	if !strings.Contains(output, "DEFAULT") {
		t.Error("expected output to contain 'DEFAULT' header")
	}
	if !strings.Contains(output, "ADDED") {
		t.Error("expected output to contain 'ADDED' header")
	}
	if !strings.Contains(output, "default") {
		t.Error("expected output to contain 'default' alias")
	}
	if !strings.Contains(output, "user@example.com") {
		t.Error("expected output to contain email")
	}
	if !strings.Contains(output, "*") {
		t.Error("expected output to contain '*' for default account")
	}
}

func TestOutputAccountsPlain_MultipleAccounts(t *testing.T) {
	accounts := []*accountuc.Account{
		{
			Alias:     "personal",
			Email:     "personal@example.com",
			IsDefault: false,
		},
		{
			Alias:     "work",
			Email:     "work@company.com",
			IsDefault: true,
		},
		{
			Alias:     "test",
			Email:     "test@test.com",
			IsDefault: false,
		},
	}

	cmd := &cobra.Command{Use: "test"}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := outputAccountsPlain(cmd, accounts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// Check all accounts present
	if !strings.Contains(output, "personal: personal@example.com") {
		t.Error("expected output to contain personal account")
	}
	if !strings.Contains(output, "work: work@company.com (default)") {
		t.Error("expected output to contain work account with (default) marker")
	}
	if !strings.Contains(output, "test: test@test.com") {
		t.Error("expected output to contain test account")
	}

	// Ensure only work has (default) marker
	lines := strings.Split(output, "\n")
	defaultCount := 0
	for _, line := range lines {
		if strings.Contains(line, "(default)") {
			defaultCount++
		}
	}
	if defaultCount != 1 {
		t.Errorf("expected exactly 1 account with (default), got %d", defaultCount)
	}
}
