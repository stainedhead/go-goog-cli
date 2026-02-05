package presenter

import (
	"encoding/json"
	"strings"
	"testing"

	domaincontacts "github.com/stainedhead/go-goog-cli/internal/domain/contacts"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		wantType string
	}{
		{
			name:     "json format returns JSONPresenter",
			format:   FormatJSON,
			wantType: "*presenter.JSONPresenter",
		},
		{
			name:     "table format returns TablePresenter",
			format:   FormatTable,
			wantType: "*presenter.TablePresenter",
		},
		{
			name:     "plain format returns PlainPresenter",
			format:   FormatPlain,
			wantType: "*presenter.PlainPresenter",
		},
		{
			name:     "unknown format returns TablePresenter as default",
			format:   "unknown",
			wantType: "*presenter.TablePresenter",
		},
		{
			name:     "empty format returns TablePresenter as default",
			format:   "",
			wantType: "*presenter.TablePresenter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.format)
			gotType := getTypeName(got)
			if gotType != tt.wantType {
				t.Errorf("New(%q) = %s, want %s", tt.format, gotType, tt.wantType)
			}
		})
	}
}

func getTypeName(p Presenter) string {
	switch p.(type) {
	case *JSONPresenter:
		return "*presenter.JSONPresenter"
	case *TablePresenter:
		return "*presenter.TablePresenter"
	case *PlainPresenter:
		return "*presenter.PlainPresenter"
	default:
		return "unknown"
	}
}

func TestJSONPresenter_RenderContact(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders nil contact as null", func(t *testing.T) {
		got := p.RenderContact(nil)
		if got != "null" {
			t.Errorf("RenderContact(nil) = %q, want %q", got, "null")
		}
	})

	t.Run("renders contact as JSON", func(t *testing.T) {
		contact := createTestContact()
		got := p.RenderContact(contact)

		var result domaincontacts.Contact
		if err := json.Unmarshal([]byte(got), &result); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		if result.ResourceName != contact.ResourceName {
			t.Errorf("ResourceName = %q, want %q", result.ResourceName, contact.ResourceName)
		}
		if len(result.EmailAddresses) != len(contact.EmailAddresses) {
			t.Errorf("EmailAddresses count = %d, want %d", len(result.EmailAddresses), len(contact.EmailAddresses))
		}
	})
}

func TestJSONPresenter_RenderContacts(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders nil contacts as empty array", func(t *testing.T) {
		got := p.RenderContacts(nil)
		if got != "[]" {
			t.Errorf("RenderContacts(nil) = %q, want %q", got, "[]")
		}
	})

	t.Run("renders empty contacts as empty array", func(t *testing.T) {
		got := p.RenderContacts([]*domaincontacts.Contact{})
		if got != "[]" {
			t.Errorf("RenderContacts([]) = %q, want %q", got, "[]")
		}
	})

	t.Run("renders multiple contacts as JSON array", func(t *testing.T) {
		contacts := []*domaincontacts.Contact{createTestContact(), createTestContact()}
		got := p.RenderContacts(contacts)

		var result []*domaincontacts.Contact
		if err := json.Unmarshal([]byte(got), &result); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		if len(result) != len(contacts) {
			t.Errorf("contacts count = %d, want %d", len(result), len(contacts))
		}
	})
}

func TestJSONPresenter_RenderContactGroup(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders nil group as null", func(t *testing.T) {
		got := p.RenderContactGroup(nil)
		if got != "null" {
			t.Errorf("RenderContactGroup(nil) = %q, want %q", got, "null")
		}
	})

	t.Run("renders contact group as JSON", func(t *testing.T) {
		group := createTestContactGroup()
		got := p.RenderContactGroup(group)

		var result domaincontacts.ContactGroup
		if err := json.Unmarshal([]byte(got), &result); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		if result.ResourceName != group.ResourceName {
			t.Errorf("ResourceName = %q, want %q", result.ResourceName, group.ResourceName)
		}
		if result.MemberCount != group.MemberCount {
			t.Errorf("MemberCount = %d, want %d", result.MemberCount, group.MemberCount)
		}
	})
}

func TestJSONPresenter_RenderContactGroups(t *testing.T) {
	p := NewJSONPresenter()

	t.Run("renders nil groups as empty array", func(t *testing.T) {
		got := p.RenderContactGroups(nil)
		if got != "[]" {
			t.Errorf("RenderContactGroups(nil) = %q, want %q", got, "[]")
		}
	})

	t.Run("renders multiple groups as JSON array", func(t *testing.T) {
		groups := []*domaincontacts.ContactGroup{createTestContactGroup(), createTestContactGroup()}
		got := p.RenderContactGroups(groups)

		var result []*domaincontacts.ContactGroup
		if err := json.Unmarshal([]byte(got), &result); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		if len(result) != len(groups) {
			t.Errorf("groups count = %d, want %d", len(result), len(groups))
		}
	})
}

func TestTablePresenter_RenderContact(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders nil contact with message", func(t *testing.T) {
		got := p.RenderContact(nil)
		if got != "No contact found" {
			t.Errorf("RenderContact(nil) = %q, want %q", got, "No contact found")
		}
	})

	t.Run("renders contact as table", func(t *testing.T) {
		contact := createTestContact()
		got := p.RenderContact(contact)

		if !strings.Contains(got, contact.ResourceName) {
			t.Errorf("output missing ResourceName: %s", got)
		}
		if !strings.Contains(got, "John Doe") {
			t.Errorf("output missing display name: %s", got)
		}
		if !strings.Contains(got, "john.doe@example.com") {
			t.Errorf("output missing email: %s", got)
		}
	})
}

func TestTablePresenter_RenderContacts(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders empty contacts with message", func(t *testing.T) {
		got := p.RenderContacts([]*domaincontacts.Contact{})
		if got != "No contacts found" {
			t.Errorf("RenderContacts([]) = %q, want %q", got, "No contacts found")
		}
	})

	t.Run("renders contacts as table", func(t *testing.T) {
		contacts := []*domaincontacts.Contact{createTestContact()}
		got := p.RenderContacts(contacts)

		upperResult := strings.ToUpper(got)
		if !strings.Contains(upperResult, "NAME") || !strings.Contains(upperResult, "EMAIL") {
			t.Errorf("output missing headers: %s", got)
		}
		if !strings.Contains(got, "John Doe") {
			t.Errorf("output missing name: %s", got)
		}
		if !strings.Contains(got, "john.doe@example.com") {
			t.Errorf("output missing email: %s", got)
		}
	})
}

func TestTablePresenter_RenderContactGroup(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders nil group with message", func(t *testing.T) {
		got := p.RenderContactGroup(nil)
		if got != "No contact group found" {
			t.Errorf("RenderContactGroup(nil) = %q, want %q", got, "No contact group found")
		}
	})

	t.Run("renders contact group as table", func(t *testing.T) {
		group := createTestContactGroup()
		got := p.RenderContactGroup(group)

		if !strings.Contains(got, group.ResourceName) {
			t.Errorf("output missing ResourceName: %s", got)
		}
		if !strings.Contains(got, group.Name) {
			t.Errorf("output missing name: %s", got)
		}
	})
}

func TestTablePresenter_RenderContactGroups(t *testing.T) {
	p := NewTablePresenter()

	t.Run("renders empty groups with message", func(t *testing.T) {
		got := p.RenderContactGroups([]*domaincontacts.ContactGroup{})
		if got != "No contact groups found" {
			t.Errorf("RenderContactGroups([]) = %q, want %q", got, "No contact groups found")
		}
	})

	t.Run("renders groups as table", func(t *testing.T) {
		groups := []*domaincontacts.ContactGroup{createTestContactGroup()}
		got := p.RenderContactGroups(groups)

		upperResult := strings.ToUpper(got)
		if !strings.Contains(upperResult, "NAME") || !strings.Contains(upperResult, "TYPE") {
			t.Errorf("output missing headers: %s", got)
		}
		if !strings.Contains(got, "Friends") {
			t.Errorf("output missing name: %s", got)
		}
		if !strings.Contains(got, "USER_CONTACT_GROUP") {
			t.Errorf("output missing group type: %s", got)
		}
	})
}

func TestPlainPresenter_RenderContact(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders nil contact as empty", func(t *testing.T) {
		got := p.RenderContact(nil)
		if got != "" {
			t.Errorf("RenderContact(nil) = %q, want empty string", got)
		}
	})

	t.Run("renders contact as plain text", func(t *testing.T) {
		contact := createTestContact()
		got := p.RenderContact(contact)

		if !strings.Contains(got, contact.ResourceName) {
			t.Errorf("output missing ResourceName: %s", got)
		}
		if !strings.Contains(got, "John Doe") {
			t.Errorf("output missing name: %s", got)
		}
		if !strings.Contains(got, "john.doe@example.com") {
			t.Errorf("output missing email: %s", got)
		}
	})
}

func TestPlainPresenter_RenderContacts(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders empty contacts as empty", func(t *testing.T) {
		got := p.RenderContacts([]*domaincontacts.Contact{})
		if got != "" {
			t.Errorf("RenderContacts([]) = %q, want empty string", got)
		}
	})

	t.Run("renders contacts as plain text lines", func(t *testing.T) {
		contacts := []*domaincontacts.Contact{createTestContact()}
		got := p.RenderContacts(contacts)

		if !strings.Contains(got, "John Doe") {
			t.Errorf("output missing name: %s", got)
		}
	})
}

func TestPlainPresenter_RenderContactGroup(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders nil group as empty", func(t *testing.T) {
		got := p.RenderContactGroup(nil)
		if got != "" {
			t.Errorf("RenderContactGroup(nil) = %q, want empty string", got)
		}
	})

	t.Run("renders contact group as plain text", func(t *testing.T) {
		group := createTestContactGroup()
		got := p.RenderContactGroup(group)

		if !strings.Contains(got, group.ResourceName) {
			t.Errorf("output missing ResourceName: %s", got)
		}
		if !strings.Contains(got, group.Name) {
			t.Errorf("output missing name: %s", got)
		}
	})
}

func TestPlainPresenter_RenderContactGroups(t *testing.T) {
	p := NewPlainPresenter()

	t.Run("renders empty groups as empty", func(t *testing.T) {
		got := p.RenderContactGroups([]*domaincontacts.ContactGroup{})
		if got != "" {
			t.Errorf("RenderContactGroups([]) = %q, want empty string", got)
		}
	})

	t.Run("renders groups as plain text lines", func(t *testing.T) {
		groups := []*domaincontacts.ContactGroup{createTestContactGroup()}
		got := p.RenderContactGroups(groups)

		if !strings.Contains(got, "Friends") {
			t.Errorf("output missing name: %s", got)
		}
	})
}

// Helper functions

func createTestContact() *domaincontacts.Contact {
	contact := &domaincontacts.Contact{
		ResourceName: "people/c123456",
		ETag:         "etag123",
		Names: []domaincontacts.Name{
			{
				DisplayName: "John Doe",
				GivenName:   "John",
				FamilyName:  "Doe",
			},
		},
		EmailAddresses: []domaincontacts.Email{
			{
				Value:   "john.doe@example.com",
				Type:    "work",
				Primary: true,
			},
		},
		PhoneNumbers: []domaincontacts.Phone{
			{
				Value:   "+1-555-1234",
				Type:    "mobile",
				Primary: true,
			},
		},
		Addresses: []domaincontacts.Address{
			{
				FormattedValue: "123 Main St, Anytown, CA 12345",
				Type:           "home",
				StreetAddress:  "123 Main St",
				City:           "Anytown",
				Region:         "CA",
				PostalCode:     "12345",
			},
		},
		Organizations: []domaincontacts.Organization{
			{
				Name:  "ACME Corp",
				Title: "Software Engineer",
			},
		},
	}
	return contact
}

func createTestContactGroup() *domaincontacts.ContactGroup {
	return &domaincontacts.ContactGroup{
		ResourceName:  "contactGroups/myContacts",
		Name:          "Friends",
		FormattedName: "Friends",
		GroupType:     domaincontacts.GroupTypeUserContactGroup,
		MemberCount:   5,
	}
}
