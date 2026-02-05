package contacts

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Contact represents a Google Contact with business logic
type Contact struct {
	ResourceName   string
	ETag           string
	Names          []Name
	Nicknames      []Nickname
	EmailAddresses []Email
	PhoneNumbers   []Phone
	Addresses      []Address
	Organizations  []Organization
	Birthdays      []Birthday
	Biographies    []Biography
	Photos         []Photo
	URLs           []URL
	Memberships    []Membership
	Metadata       *ResourceMetadata
}

// Name represents a person's name
type Name struct {
	DisplayName        string
	GivenName          string
	FamilyName         string
	MiddleName         string
	HonorificPrefix    string
	HonorificSuffix    string
	PhoneticGivenName  string
	PhoneticFamilyName string
	PhoneticFullName   string
}

// Nickname represents a nickname
type Nickname struct {
	Value string
	Type  string
}

// Email represents an email address
type Email struct {
	Value         string
	Type          string
	DisplayName   string
	Primary       bool
	FormattedType string
}

// Phone represents a phone number
type Phone struct {
	Value         string
	Type          string
	Primary       bool
	FormattedType string
}

// Address represents a physical address
type Address struct {
	FormattedValue  string
	Type            string
	StreetAddress   string
	City            string
	Region          string
	PostalCode      string
	Country         string
	CountryCode     string
	ExtendedAddress string
	PoBox           string
	FormattedType   string
}

// Organization represents an organization affiliation
type Organization struct {
	Name       string
	Title      string
	Department string
	Symbol     string
	Domain     string
	Location   string
	Type       string
	StartDate  *Date
	EndDate    *Date
	Current    bool
}

// Birthday represents a birthday
type Birthday struct {
	Date *Date
	Text string
}

// Biography represents biographical information
type Biography struct {
	Value       string
	ContentType string
}

// Photo represents a photo
type Photo struct {
	URL     string
	Default bool
}

// URL represents a URL
type URL struct {
	Value string
	Type  string
}

// Membership represents group membership
type Membership struct {
	ContactGroupMembership *ContactGroupMembership
	DomainMembership       *DomainMembership
}

// ContactGroupMembership represents membership in a contact group
type ContactGroupMembership struct {
	ContactGroupResourceName string
	ContactGroupID           string
}

// DomainMembership represents membership in a domain
type DomainMembership struct {
	InViewerDomain bool
}

// ResourceMetadata contains metadata about the resource
type ResourceMetadata struct {
	Sources []Source
}

// Source represents the source of contact data
type Source struct {
	Type       string
	ID         string
	ETag       string
	UpdateTime time.Time
}

// Date represents a date
type Date struct {
	Year  int
	Month int
	Day   int
}

// NewContact creates a new contact
func NewContact() *Contact {
	return &Contact{
		Names:          []Name{},
		Nicknames:      []Nickname{},
		EmailAddresses: []Email{},
		PhoneNumbers:   []Phone{},
		Addresses:      []Address{},
		Organizations:  []Organization{},
		Birthdays:      []Birthday{},
		Biographies:    []Biography{},
		Photos:         []Photo{},
		URLs:           []URL{},
		Memberships:    []Membership{},
	}
}

// AddEmail adds an email address to the contact
func (c *Contact) AddEmail(email, emailType string, primary bool) error {
	if err := validateEmail(email); err != nil {
		return err
	}

	c.EmailAddresses = append(c.EmailAddresses, Email{
		Value:   email,
		Type:    emailType,
		Primary: primary,
	})
	return nil
}

// AddPhone adds a phone number to the contact
func (c *Contact) AddPhone(phone, phoneType string, primary bool) error {
	if err := validatePhone(phone); err != nil {
		return err
	}

	c.PhoneNumbers = append(c.PhoneNumbers, Phone{
		Value:   phone,
		Type:    phoneType,
		Primary: primary,
	})
	return nil
}

// SetPrimaryEmail marks the specified email as primary and unmarks others
func (c *Contact) SetPrimaryEmail(email string) error {
	found := false
	for i := range c.EmailAddresses {
		if c.EmailAddresses[i].Value == email {
			c.EmailAddresses[i].Primary = true
			found = true
		} else {
			c.EmailAddresses[i].Primary = false
		}
	}
	if !found {
		return errors.New("email not found")
	}
	return nil
}

// GetPrimaryEmail returns the primary email or the first email if no primary is set
func (c *Contact) GetPrimaryEmail() (string, error) {
	if len(c.EmailAddresses) == 0 {
		return "", errors.New("no email addresses")
	}

	for _, email := range c.EmailAddresses {
		if email.Primary {
			return email.Value, nil
		}
	}

	return c.EmailAddresses[0].Value, nil
}

// GetDisplayName returns the best display name for the contact
func (c *Contact) GetDisplayName() string {
	if len(c.Names) == 0 {
		return ""
	}

	name := c.Names[0]
	if name.DisplayName != "" {
		return name.DisplayName
	}

	parts := []string{}
	if name.GivenName != "" {
		parts = append(parts, name.GivenName)
	}
	if name.FamilyName != "" {
		parts = append(parts, name.FamilyName)
	}

	return strings.Join(parts, " ")
}

// IsInGroup checks if the contact is a member of the specified group
func (c *Contact) IsInGroup(groupResourceName string) bool {
	for _, membership := range c.Memberships {
		if membership.ContactGroupMembership != nil &&
			membership.ContactGroupMembership.ContactGroupResourceName == groupResourceName {
			return true
		}
	}
	return false
}

// validateEmail validates an email address
func validateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(email, "@") || len(email) < 4 {
		return errors.New("invalid email format")
	}
	return nil
}

// validatePhone validates a phone number
func validatePhone(phone string) error {
	if phone == "" {
		return errors.New("phone cannot be empty")
	}
	return nil
}

// FormatDate formats a Date into a string
func (d *Date) FormatDate() string {
	if d == nil {
		return ""
	}
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}
