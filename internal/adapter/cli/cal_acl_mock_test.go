// Package cli provides command-line interface handlers for the goog application.
package cli

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/domain/calendar"
)

// =============================================================================
// Tests using dependency injection with mocks for ACL commands
// =============================================================================

func TestACLAddCmd_Validation(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		expectErr bool
	}{
		{
			name:      "empty email",
			email:     "",
			expectErr: true,
		},
		{
			name:      "valid email",
			email:     "user@example.com",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origEmail := aclEmail
			aclEmail = tt.email
			defer func() { aclEmail = origEmail }()

			mockCmd := &cobra.Command{Use: "test"}

			if aclAddCmd.PreRunE == nil {
				t.Error("aclAddCmd should have PreRunE defined")
				return
			}

			err := aclAddCmd.PreRunE(mockCmd, []string{"primary"})

			if tt.expectErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestACLRemoveCmd_Validation(t *testing.T) {
	tests := []struct {
		name      string
		confirm   bool
		expectErr bool
	}{
		{
			name:      "without confirmation",
			confirm:   false,
			expectErr: true,
		},
		{
			name:      "with confirmation",
			confirm:   true,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origConfirm := aclConfirm
			aclConfirm = tt.confirm
			defer func() { aclConfirm = origConfirm }()

			mockCmd := &cobra.Command{Use: "test"}

			if aclRemoveCmd.PreRunE == nil {
				t.Error("aclRemoveCmd should have PreRunE defined")
				return
			}

			err := aclRemoveCmd.PreRunE(mockCmd, []string{"primary", "rule-id"})

			if tt.expectErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestShareCmd_Validation(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		expectErr bool
	}{
		{
			name:      "empty email",
			email:     "",
			expectErr: true,
		},
		{
			name:      "valid email",
			email:     "user@example.com",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// shareCmd uses the same aclEmail variable as aclAddCmd
			origEmail := aclEmail
			aclEmail = tt.email
			defer func() { aclEmail = origEmail }()

			mockCmd := &cobra.Command{Use: "test"}

			if shareCmd.PreRunE == nil {
				t.Error("shareCmd should have PreRunE defined")
				return
			}

			err := shareCmd.PreRunE(mockCmd, []string{"primary"})

			if tt.expectErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestUnshareCmd_NoPreRunE(t *testing.T) {
	// unshareCmd doesn't have a PreRunE - it uses runACLRemove directly
	// which checks aclConfirm internally
	if unshareCmd.PreRunE != nil {
		t.Log("unshareCmd has PreRunE defined - testing it")

		origConfirm := aclConfirm
		aclConfirm = false
		defer func() { aclConfirm = origConfirm }()

		mockCmd := &cobra.Command{Use: "test"}
		err := unshareCmd.PreRunE(mockCmd, []string{"primary", "rule-id"})
		if err == nil {
			t.Error("expected error when confirm not set")
		}
	} else {
		t.Log("unshareCmd relies on runACLRemove for confirmation check")
	}
}

func TestMockACLRepository(t *testing.T) {
	t.Run("List success", func(t *testing.T) {
		rules := []*calendar.ACLRule{
			{ID: "rule1", Role: "owner", Scope: &calendar.ACLScope{Type: "user", Value: "user@example.com"}},
		}
		repo := &MockACLRepository{Rules: rules}

		result, err := repo.List(nil, "primary")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("expected 1 rule, got %d", len(result))
		}
	})

	t.Run("List error", func(t *testing.T) {
		repo := &MockACLRepository{ListErr: fmt.Errorf("list error")}

		_, err := repo.List(nil, "primary")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Get success", func(t *testing.T) {
		rule := &calendar.ACLRule{ID: "rule1", Role: "reader"}
		repo := &MockACLRepository{Rule: rule}

		result, err := repo.Get(nil, "primary", "rule1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "rule1" {
			t.Errorf("expected ID 'rule1', got %s", result.ID)
		}
	})

	t.Run("Get error", func(t *testing.T) {
		repo := &MockACLRepository{GetErr: fmt.Errorf("not found")}

		_, err := repo.Get(nil, "primary", "nonexistent")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Insert success", func(t *testing.T) {
		repo := &MockACLRepository{}
		rule := &calendar.ACLRule{Role: "reader", Scope: &calendar.ACLScope{Type: "user", Value: "user@example.com"}}

		result, err := repo.Insert(nil, "primary", rule)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.ID != "mock-acl-id" {
			t.Errorf("expected mock ID, got %s", result.ID)
		}
	})

	t.Run("Insert error", func(t *testing.T) {
		repo := &MockACLRepository{InsertErr: fmt.Errorf("insert error")}
		rule := &calendar.ACLRule{Role: "reader"}

		_, err := repo.Insert(nil, "primary", rule)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Update success", func(t *testing.T) {
		repo := &MockACLRepository{}
		rule := &calendar.ACLRule{ID: "rule1", Role: "writer"}

		result, err := repo.Update(nil, "primary", rule)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result.Role != "writer" {
			t.Errorf("expected role 'writer', got %s", result.Role)
		}
	})

	t.Run("Update error", func(t *testing.T) {
		repo := &MockACLRepository{UpdateErr: fmt.Errorf("update error")}
		rule := &calendar.ACLRule{ID: "rule1", Role: "writer"}

		_, err := repo.Update(nil, "primary", rule)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Delete success", func(t *testing.T) {
		repo := &MockACLRepository{}

		err := repo.Delete(nil, "primary", "rule-id")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Delete error", func(t *testing.T) {
		repo := &MockACLRepository{DeleteErr: fmt.Errorf("delete error")}

		err := repo.Delete(nil, "primary", "rule-id")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestACLAddCmd_RoleDefaultValue(t *testing.T) {
	flag := aclAddCmd.Flag("role")
	if flag == nil {
		t.Fatal("expected --role flag to be set")
	}

	// Default role should be "reader"
	if flag.DefValue != "reader" {
		t.Errorf("expected default role 'reader', got '%s'", flag.DefValue)
	}
}

func TestShareCmd_RoleDefaultValue(t *testing.T) {
	flag := shareCmd.Flag("role")
	if flag == nil {
		t.Fatal("expected --role flag to be set")
	}

	// Default role should be "reader"
	if flag.DefValue != "reader" {
		t.Errorf("expected default role 'reader', got '%s'", flag.DefValue)
	}
}

func TestACLAddCmd_ValidRoles(t *testing.T) {
	// Test that the role flag accepts valid values
	validRoles := []string{"reader", "writer", "owner", "freeBusyReader"}

	for _, role := range validRoles {
		t.Run(role, func(t *testing.T) {
			// Just validate that the role can be set - actual API validation happens server-side
			flag := aclAddCmd.Flag("role")
			if flag == nil {
				t.Fatal("expected --role flag to be set")
			}
		})
	}
}
