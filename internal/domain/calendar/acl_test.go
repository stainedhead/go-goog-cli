package calendar

import "testing"

func TestNewACLRule(t *testing.T) {
	scope := NewUserACLScope("test@example.com")
	role := AccessRoleReader
	rule := NewACLRule(scope, role)

	if rule.Scope != scope {
		t.Error("expected scope to match")
	}

	if rule.Role != role {
		t.Errorf("expected role %q, got %q", role, rule.Role)
	}
}

func TestNewACLScope(t *testing.T) {
	scopeType := ACLScopeTypeUser
	value := "test@example.com"
	scope := NewACLScope(scopeType, value)

	if scope.Type != scopeType {
		t.Errorf("expected type %q, got %q", scopeType, scope.Type)
	}

	if scope.Value != value {
		t.Errorf("expected value %q, got %q", value, scope.Value)
	}
}

func TestNewUserACLScope(t *testing.T) {
	email := "user@example.com"
	scope := NewUserACLScope(email)

	if scope.Type != ACLScopeTypeUser {
		t.Errorf("expected type %q, got %q", ACLScopeTypeUser, scope.Type)
	}

	if scope.Value != email {
		t.Errorf("expected value %q, got %q", email, scope.Value)
	}
}

func TestNewGroupACLScope(t *testing.T) {
	email := "group@example.com"
	scope := NewGroupACLScope(email)

	if scope.Type != ACLScopeTypeGroup {
		t.Errorf("expected type %q, got %q", ACLScopeTypeGroup, scope.Type)
	}

	if scope.Value != email {
		t.Errorf("expected value %q, got %q", email, scope.Value)
	}
}

func TestNewDomainACLScope(t *testing.T) {
	domain := "example.com"
	scope := NewDomainACLScope(domain)

	if scope.Type != ACLScopeTypeDomain {
		t.Errorf("expected type %q, got %q", ACLScopeTypeDomain, scope.Type)
	}

	if scope.Value != domain {
		t.Errorf("expected value %q, got %q", domain, scope.Value)
	}
}

func TestNewDefaultACLScope(t *testing.T) {
	scope := NewDefaultACLScope()

	if scope.Type != ACLScopeTypeDefault {
		t.Errorf("expected type %q, got %q", ACLScopeTypeDefault, scope.Type)
	}

	if scope.Value != "" {
		t.Errorf("expected empty value, got %q", scope.Value)
	}
}

func TestIsValidACLScopeType(t *testing.T) {
	tests := []struct {
		scopeType string
		valid     bool
	}{
		{ACLScopeTypeUser, true},
		{ACLScopeTypeGroup, true},
		{ACLScopeTypeDomain, true},
		{ACLScopeTypeDefault, true},
		{"user", true},
		{"group", true},
		{"domain", true},
		{"default", true},
		{"organization", false},
		{"", false},
		{"USER", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.scopeType, func(t *testing.T) {
			got := IsValidACLScopeType(tt.scopeType)
			if got != tt.valid {
				t.Errorf("IsValidACLScopeType(%q) = %v, want %v", tt.scopeType, got, tt.valid)
			}
		})
	}
}

func TestACLScopeIsUser(t *testing.T) {
	tests := []struct {
		scopeType string
		isUser    bool
	}{
		{ACLScopeTypeUser, true},
		{ACLScopeTypeGroup, false},
		{ACLScopeTypeDomain, false},
		{ACLScopeTypeDefault, false},
	}

	for _, tt := range tests {
		t.Run(tt.scopeType, func(t *testing.T) {
			scope := &ACLScope{Type: tt.scopeType}
			got := scope.IsUser()
			if got != tt.isUser {
				t.Errorf("ACLScope{Type: %q}.IsUser() = %v, want %v", tt.scopeType, got, tt.isUser)
			}
		})
	}
}

func TestACLScopeIsGroup(t *testing.T) {
	tests := []struct {
		scopeType string
		isGroup   bool
	}{
		{ACLScopeTypeUser, false},
		{ACLScopeTypeGroup, true},
		{ACLScopeTypeDomain, false},
		{ACLScopeTypeDefault, false},
	}

	for _, tt := range tests {
		t.Run(tt.scopeType, func(t *testing.T) {
			scope := &ACLScope{Type: tt.scopeType}
			got := scope.IsGroup()
			if got != tt.isGroup {
				t.Errorf("ACLScope{Type: %q}.IsGroup() = %v, want %v", tt.scopeType, got, tt.isGroup)
			}
		})
	}
}

func TestACLScopeIsDomain(t *testing.T) {
	tests := []struct {
		scopeType string
		isDomain  bool
	}{
		{ACLScopeTypeUser, false},
		{ACLScopeTypeGroup, false},
		{ACLScopeTypeDomain, true},
		{ACLScopeTypeDefault, false},
	}

	for _, tt := range tests {
		t.Run(tt.scopeType, func(t *testing.T) {
			scope := &ACLScope{Type: tt.scopeType}
			got := scope.IsDomain()
			if got != tt.isDomain {
				t.Errorf("ACLScope{Type: %q}.IsDomain() = %v, want %v", tt.scopeType, got, tt.isDomain)
			}
		})
	}
}

func TestACLScopeIsDefault(t *testing.T) {
	tests := []struct {
		scopeType string
		isDefault bool
	}{
		{ACLScopeTypeUser, false},
		{ACLScopeTypeGroup, false},
		{ACLScopeTypeDomain, false},
		{ACLScopeTypeDefault, true},
	}

	for _, tt := range tests {
		t.Run(tt.scopeType, func(t *testing.T) {
			scope := &ACLScope{Type: tt.scopeType}
			got := scope.IsDefault()
			if got != tt.isDefault {
				t.Errorf("ACLScope{Type: %q}.IsDefault() = %v, want %v", tt.scopeType, got, tt.isDefault)
			}
		})
	}
}

func TestACLRuleGrantsOwnerAccess(t *testing.T) {
	tests := []struct {
		role   string
		grants bool
	}{
		{AccessRoleOwner, true},
		{AccessRoleWriter, false},
		{AccessRoleReader, false},
		{AccessRoleFreeBusyReader, false},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			rule := &ACLRule{Role: tt.role}
			got := rule.GrantsOwnerAccess()
			if got != tt.grants {
				t.Errorf("ACLRule{Role: %q}.GrantsOwnerAccess() = %v, want %v", tt.role, got, tt.grants)
			}
		})
	}
}

func TestACLRuleGrantsWriteAccess(t *testing.T) {
	tests := []struct {
		role   string
		grants bool
	}{
		{AccessRoleOwner, true},
		{AccessRoleWriter, true},
		{AccessRoleReader, false},
		{AccessRoleFreeBusyReader, false},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			rule := &ACLRule{Role: tt.role}
			got := rule.GrantsWriteAccess()
			if got != tt.grants {
				t.Errorf("ACLRule{Role: %q}.GrantsWriteAccess() = %v, want %v", tt.role, got, tt.grants)
			}
		})
	}
}

func TestACLRuleGrantsReadAccess(t *testing.T) {
	tests := []struct {
		role   string
		grants bool
	}{
		{AccessRoleOwner, true},
		{AccessRoleWriter, true},
		{AccessRoleReader, true},
		{AccessRoleFreeBusyReader, false},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			rule := &ACLRule{Role: tt.role}
			got := rule.GrantsReadAccess()
			if got != tt.grants {
				t.Errorf("ACLRule{Role: %q}.GrantsReadAccess() = %v, want %v", tt.role, got, tt.grants)
			}
		})
	}
}

func TestACLRuleFields(t *testing.T) {
	scope := NewUserACLScope("test@example.com")
	rule := &ACLRule{
		ID:    "rule123",
		Scope: scope,
		Role:  AccessRoleWriter,
	}

	if rule.ID != "rule123" {
		t.Errorf("unexpected ID: %s", rule.ID)
	}

	if rule.Scope.Type != ACLScopeTypeUser {
		t.Errorf("unexpected Scope.Type: %s", rule.Scope.Type)
	}

	if rule.Scope.Value != "test@example.com" {
		t.Errorf("unexpected Scope.Value: %s", rule.Scope.Value)
	}

	if rule.Role != AccessRoleWriter {
		t.Errorf("unexpected Role: %s", rule.Role)
	}
}
