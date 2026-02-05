package contacts

import (
	"errors"
	"time"
)

// Group type constants
const (
	GroupTypeSystem           = "SYSTEM_CONTACT_GROUP"
	GroupTypeUserContactGroup = "USER_CONTACT_GROUP"
)

// ContactGroup represents a Google Contact Group
type ContactGroup struct {
	ResourceName        string
	ETag                string
	Name                string
	FormattedName       string
	GroupType           string
	MemberCount         int
	MemberResourceNames []string
	Metadata            *GroupMetadata
}

// GroupMetadata contains metadata about the group
type GroupMetadata struct {
	UpdateTime time.Time
	Deleted    bool
}

// NewContactGroup creates a new contact group with validation
func NewContactGroup(name string) (*ContactGroup, error) {
	if name == "" {
		return nil, errors.New("group name cannot be empty")
	}

	return &ContactGroup{
		Name:                name,
		GroupType:           GroupTypeUserContactGroup,
		MemberResourceNames: []string{},
	}, nil
}

// IsSystemGroup returns true if the group is a system-defined group
func (g *ContactGroup) IsSystemGroup() bool {
	return g.GroupType == GroupTypeSystem
}

// CanModify returns true if the group can be modified
func (g *ContactGroup) CanModify() bool {
	return !g.IsSystemGroup()
}
