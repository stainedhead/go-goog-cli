package calendar

// ACLRule represents an access control rule for a calendar.
type ACLRule struct {
	// ID is the unique identifier for the ACL rule.
	ID string
	// Scope defines who the rule applies to.
	Scope *ACLScope
	// Role is the access level: owner, writer, reader, freeBusyReader.
	Role string
}

// ACLScope defines the scope of an ACL rule.
type ACLScope struct {
	// Type is the scope type: user, group, domain, default.
	Type string
	// Value is the email address or domain for user, group, or domain scopes.
	Value string
}

// ACL scope type constants.
const (
	ACLScopeTypeUser    = "user"
	ACLScopeTypeGroup   = "group"
	ACLScopeTypeDomain  = "domain"
	ACLScopeTypeDefault = "default"
)

// NewACLRule creates a new ACLRule with the given scope and role.
func NewACLRule(scope *ACLScope, role string) *ACLRule {
	return &ACLRule{
		Scope: scope,
		Role:  role,
	}
}

// NewACLScope creates a new ACLScope with the given type and value.
func NewACLScope(scopeType, value string) *ACLScope {
	return &ACLScope{
		Type:  scopeType,
		Value: value,
	}
}

// NewUserACLScope creates a new ACLScope for a user.
func NewUserACLScope(email string) *ACLScope {
	return &ACLScope{
		Type:  ACLScopeTypeUser,
		Value: email,
	}
}

// NewGroupACLScope creates a new ACLScope for a group.
func NewGroupACLScope(email string) *ACLScope {
	return &ACLScope{
		Type:  ACLScopeTypeGroup,
		Value: email,
	}
}

// NewDomainACLScope creates a new ACLScope for a domain.
func NewDomainACLScope(domain string) *ACLScope {
	return &ACLScope{
		Type:  ACLScopeTypeDomain,
		Value: domain,
	}
}

// NewDefaultACLScope creates a new ACLScope for the default (public) scope.
func NewDefaultACLScope() *ACLScope {
	return &ACLScope{
		Type: ACLScopeTypeDefault,
	}
}

// IsValidACLScopeType checks if the given type is a valid ACL scope type.
func IsValidACLScopeType(scopeType string) bool {
	switch scopeType {
	case ACLScopeTypeUser, ACLScopeTypeGroup, ACLScopeTypeDomain, ACLScopeTypeDefault:
		return true
	default:
		return false
	}
}

// IsUser returns true if the scope is for a specific user.
func (s *ACLScope) IsUser() bool {
	return s.Type == ACLScopeTypeUser
}

// IsGroup returns true if the scope is for a group.
func (s *ACLScope) IsGroup() bool {
	return s.Type == ACLScopeTypeGroup
}

// IsDomain returns true if the scope is for a domain.
func (s *ACLScope) IsDomain() bool {
	return s.Type == ACLScopeTypeDomain
}

// IsDefault returns true if the scope is the default (public) scope.
func (s *ACLScope) IsDefault() bool {
	return s.Type == ACLScopeTypeDefault
}

// GrantsOwnerAccess returns true if the rule grants owner access.
func (r *ACLRule) GrantsOwnerAccess() bool {
	return r.Role == AccessRoleOwner
}

// GrantsWriteAccess returns true if the rule grants write access or higher.
func (r *ACLRule) GrantsWriteAccess() bool {
	return r.Role == AccessRoleOwner || r.Role == AccessRoleWriter
}

// GrantsReadAccess returns true if the rule grants read access or higher.
func (r *ACLRule) GrantsReadAccess() bool {
	switch r.Role {
	case AccessRoleOwner, AccessRoleWriter, AccessRoleReader:
		return true
	default:
		return false
	}
}
