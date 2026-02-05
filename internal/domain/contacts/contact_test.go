package contacts

import (
	"testing"
)

func TestNewContact(t *testing.T) {
	contact := NewContact()
	if contact == nil {
		t.Error("NewContact() returned nil contact")
	}
	if contact.Names == nil {
		t.Error("NewContact() did not initialize Names slice")
	}
	if contact.EmailAddresses == nil {
		t.Error("NewContact() did not initialize EmailAddresses slice")
	}
}

func TestContact_AddEmail(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		emailType  string
		primary    bool
		wantErr    bool
		errMessage string
	}{
		{
			name:      "adds valid work email",
			email:     "john@example.com",
			emailType: "work",
			primary:   true,
			wantErr:   false,
		},
		{
			name:       "rejects empty email",
			email:      "",
			emailType:  "work",
			primary:    false,
			wantErr:    true,
			errMessage: "email cannot be empty",
		},
		{
			name:       "rejects invalid email format",
			email:      "not-an-email",
			emailType:  "work",
			primary:    false,
			wantErr:    true,
			errMessage: "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contact := NewContact()
			err := contact.AddEmail(tt.email, tt.emailType, tt.primary)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMessage {
				t.Errorf("AddEmail() error message = %v, want %v", err.Error(), tt.errMessage)
			}

			if !tt.wantErr && len(contact.EmailAddresses) != 1 {
				t.Errorf("AddEmail() did not add email, got %d emails", len(contact.EmailAddresses))
			}

			if !tt.wantErr && contact.EmailAddresses[0].Value != tt.email {
				t.Errorf("AddEmail() email = %v, want %v", contact.EmailAddresses[0].Value, tt.email)
			}

			if !tt.wantErr && tt.primary && !contact.EmailAddresses[0].Primary {
				t.Error("AddEmail() did not set primary flag")
			}
		})
	}
}

func TestContact_AddPhone(t *testing.T) {
	tests := []struct {
		name       string
		phone      string
		phoneType  string
		primary    bool
		wantErr    bool
		errMessage string
	}{
		{
			name:      "adds valid mobile phone",
			phone:     "+1-555-0123",
			phoneType: "mobile",
			primary:   true,
			wantErr:   false,
		},
		{
			name:       "rejects empty phone",
			phone:      "",
			phoneType:  "mobile",
			primary:    false,
			wantErr:    true,
			errMessage: "phone cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contact := NewContact()
			err := contact.AddPhone(tt.phone, tt.phoneType, tt.primary)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMessage {
				t.Errorf("AddPhone() error message = %v, want %v", err.Error(), tt.errMessage)
			}

			if !tt.wantErr && len(contact.PhoneNumbers) != 1 {
				t.Errorf("AddPhone() did not add phone, got %d phones", len(contact.PhoneNumbers))
			}
		})
	}
}

func TestContact_SetPrimaryEmail(t *testing.T) {
	tests := []struct {
		name       string
		emails     []string
		setPrimary string
		wantErr    bool
		errMessage string
	}{
		{
			name:       "sets primary email",
			emails:     []string{"john@example.com", "jane@example.com"},
			setPrimary: "jane@example.com",
			wantErr:    false,
		},
		{
			name:       "rejects email not in list",
			emails:     []string{"john@example.com"},
			setPrimary: "jane@example.com",
			wantErr:    true,
			errMessage: "email not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contact := NewContact()
			for _, email := range tt.emails {
				contact.AddEmail(email, "work", false)
			}

			err := contact.SetPrimaryEmail(tt.setPrimary)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetPrimaryEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				found := false
				for _, email := range contact.EmailAddresses {
					if email.Value == tt.setPrimary && email.Primary {
						found = true
					}
				}
				if !found {
					t.Error("SetPrimaryEmail() did not set email as primary")
				}
			}
		})
	}
}

func TestContact_GetPrimaryEmail(t *testing.T) {
	tests := []struct {
		name    string
		emails  []Email
		want    string
		wantErr bool
	}{
		{
			name: "returns primary email",
			emails: []Email{
				{Value: "john@example.com", Primary: false},
				{Value: "jane@example.com", Primary: true},
			},
			want:    "jane@example.com",
			wantErr: false,
		},
		{
			name: "returns first email when no primary",
			emails: []Email{
				{Value: "john@example.com", Primary: false},
				{Value: "jane@example.com", Primary: false},
			},
			want:    "john@example.com",
			wantErr: false,
		},
		{
			name:    "returns error when no emails",
			emails:  []Email{},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contact := &Contact{EmailAddresses: tt.emails}
			got, err := contact.GetPrimaryEmail()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPrimaryEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("GetPrimaryEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContact_GetDisplayName(t *testing.T) {
	tests := []struct {
		name  string
		names []Name
		want  string
	}{
		{
			name: "returns full name",
			names: []Name{
				{GivenName: "John", FamilyName: "Doe"},
			},
			want: "John Doe",
		},
		{
			name: "returns given name only",
			names: []Name{
				{GivenName: "John"},
			},
			want: "John",
		},
		{
			name: "returns display name if set",
			names: []Name{
				{DisplayName: "Johnny", GivenName: "John", FamilyName: "Doe"},
			},
			want: "Johnny",
		},
		{
			name:  "returns empty for no names",
			names: []Name{},
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contact := &Contact{Names: tt.names}
			got := contact.GetDisplayName()

			if got != tt.want {
				t.Errorf("GetDisplayName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContact_IsInGroup(t *testing.T) {
	tests := []struct {
		name        string
		memberships []Membership
		groupName   string
		want        bool
	}{
		{
			name: "returns true for member",
			memberships: []Membership{
				{ContactGroupMembership: &ContactGroupMembership{
					ContactGroupResourceName: "contactGroups/g123",
				}},
			},
			groupName: "contactGroups/g123",
			want:      true,
		},
		{
			name: "returns false for non-member",
			memberships: []Membership{
				{ContactGroupMembership: &ContactGroupMembership{
					ContactGroupResourceName: "contactGroups/g123",
				}},
			},
			groupName: "contactGroups/g456",
			want:      false,
		},
		{
			name:        "returns false for no memberships",
			memberships: []Membership{},
			groupName:   "contactGroups/g123",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contact := &Contact{Memberships: tt.memberships}
			got := contact.IsInGroup(tt.groupName)

			if got != tt.want {
				t.Errorf("IsInGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
