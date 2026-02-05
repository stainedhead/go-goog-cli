package repository

import (
	"testing"
	"time"

	"github.com/stainedhead/go-goog-cli/internal/domain/contacts"
	"google.golang.org/api/people/v1"
)

// TestApiPersonToDomain_BasicFields tests basic field mapping.
func TestApiPersonToDomain_BasicFields(t *testing.T) {
	apiPerson := &people.Person{
		ResourceName: "people/c123",
		Etag:         "etag456",
	}

	contact := apiPersonToDomain(apiPerson)

	if contact.ResourceName != apiPerson.ResourceName {
		t.Errorf("ResourceName = %s, want %s", contact.ResourceName, apiPerson.ResourceName)
	}
	if contact.ETag != apiPerson.Etag {
		t.Errorf("ETag = %s, want %s", contact.ETag, apiPerson.Etag)
	}
}

// TestApiPersonToDomain_Names tests name field mapping.
func TestApiPersonToDomain_Names(t *testing.T) {
	apiPerson := &people.Person{
		ResourceName: "people/c123",
		Names: []*people.Name{
			{
				DisplayName:        "John Doe",
				GivenName:          "John",
				FamilyName:         "Doe",
				MiddleName:         "A",
				HonorificPrefix:    "Mr",
				HonorificSuffix:    "Jr",
				PhoneticGivenName:  "Jon",
				PhoneticFamilyName: "Do",
				PhoneticFullName:   "Jon Do",
			},
		},
	}

	contact := apiPersonToDomain(apiPerson)

	if len(contact.Names) != 1 {
		t.Fatalf("expected 1 name, got %d", len(contact.Names))
	}

	name := contact.Names[0]
	if name.DisplayName != "John Doe" {
		t.Errorf("DisplayName = %s, want John Doe", name.DisplayName)
	}
	if name.GivenName != "John" {
		t.Errorf("GivenName = %s, want John", name.GivenName)
	}
	if name.FamilyName != "Doe" {
		t.Errorf("FamilyName = %s, want Doe", name.FamilyName)
	}
}

// TestApiPersonToDomain_Emails tests email field mapping.
func TestApiPersonToDomain_Emails(t *testing.T) {
	apiPerson := &people.Person{
		ResourceName: "people/c123",
		EmailAddresses: []*people.EmailAddress{
			{
				Value:         "john@example.com",
				Type:          "work",
				DisplayName:   "Work Email",
				FormattedType: "Work",
			},
			{
				Value:         "john.personal@example.com",
				Type:          "home",
				FormattedType: "Home",
			},
		},
	}

	// Mark first email as primary using metadata
	apiPerson.EmailAddresses[0].Metadata = &people.FieldMetadata{Primary: true}

	contact := apiPersonToDomain(apiPerson)

	if len(contact.EmailAddresses) != 2 {
		t.Fatalf("expected 2 emails, got %d", len(contact.EmailAddresses))
	}

	email := contact.EmailAddresses[0]
	if email.Value != "john@example.com" {
		t.Errorf("Email = %s, want john@example.com", email.Value)
	}
	if email.Type != "work" {
		t.Errorf("Type = %s, want work", email.Type)
	}
	if !email.Primary {
		t.Error("expected email to be primary")
	}
}

// TestApiPersonToDomain_PhoneNumbers tests phone number mapping.
func TestApiPersonToDomain_PhoneNumbers(t *testing.T) {
	apiPerson := &people.Person{
		ResourceName: "people/c123",
		PhoneNumbers: []*people.PhoneNumber{
			{
				Value:         "+1234567890",
				Type:          "mobile",
				FormattedType: "Mobile",
			},
		},
	}

	contact := apiPersonToDomain(apiPerson)

	if len(contact.PhoneNumbers) != 1 {
		t.Fatalf("expected 1 phone, got %d", len(contact.PhoneNumbers))
	}

	phone := contact.PhoneNumbers[0]
	if phone.Value != "+1234567890" {
		t.Errorf("Phone = %s, want +1234567890", phone.Value)
	}
	if phone.Type != "mobile" {
		t.Errorf("Type = %s, want mobile", phone.Type)
	}
}

// TestApiPersonToDomain_Addresses tests address mapping.
func TestApiPersonToDomain_Addresses(t *testing.T) {
	apiPerson := &people.Person{
		ResourceName: "people/c123",
		Addresses: []*people.Address{
			{
				FormattedValue:  "123 Main St, City, ST 12345",
				Type:            "home",
				StreetAddress:   "123 Main St",
				City:            "City",
				Region:          "ST",
				PostalCode:      "12345",
				Country:         "USA",
				CountryCode:     "US",
				ExtendedAddress: "Apt 4",
				PoBox:           "PO Box 123",
				FormattedType:   "Home",
			},
		},
	}

	contact := apiPersonToDomain(apiPerson)

	if len(contact.Addresses) != 1 {
		t.Fatalf("expected 1 address, got %d", len(contact.Addresses))
	}

	addr := contact.Addresses[0]
	if addr.City != "City" {
		t.Errorf("City = %s, want City", addr.City)
	}
	if addr.PostalCode != "12345" {
		t.Errorf("PostalCode = %s, want 12345", addr.PostalCode)
	}
}

// TestApiPersonToDomain_Organizations tests organization mapping.
func TestApiPersonToDomain_Organizations(t *testing.T) {
	apiPerson := &people.Person{
		ResourceName: "people/c123",
		Organizations: []*people.Organization{
			{
				Name:       "Acme Corp",
				Title:      "Engineer",
				Department: "IT",
				Symbol:     "ACME",
				Domain:     "acme.com",
				Location:   "New York",
				Type:       "work",
				StartDate: &people.Date{
					Year:  2020,
					Month: 1,
					Day:   15,
				},
				Current: true,
			},
		},
	}

	contact := apiPersonToDomain(apiPerson)

	if len(contact.Organizations) != 1 {
		t.Fatalf("expected 1 organization, got %d", len(contact.Organizations))
	}

	org := contact.Organizations[0]
	if org.Name != "Acme Corp" {
		t.Errorf("Name = %s, want Acme Corp", org.Name)
	}
	if org.Title != "Engineer" {
		t.Errorf("Title = %s, want Engineer", org.Title)
	}
	if !org.Current {
		t.Error("expected organization to be current")
	}
}

// TestApiPersonToDomain_Birthdays tests birthday mapping.
func TestApiPersonToDomain_Birthdays(t *testing.T) {
	apiPerson := &people.Person{
		ResourceName: "people/c123",
		Birthdays: []*people.Birthday{
			{
				Date: &people.Date{
					Year:  1990,
					Month: 6,
					Day:   15,
				},
				Text: "June 15, 1990",
			},
		},
	}

	contact := apiPersonToDomain(apiPerson)

	if len(contact.Birthdays) != 1 {
		t.Fatalf("expected 1 birthday, got %d", len(contact.Birthdays))
	}

	birthday := contact.Birthdays[0]
	if birthday.Date.Year != 1990 {
		t.Errorf("Year = %d, want 1990", birthday.Date.Year)
	}
	if birthday.Date.Month != 6 {
		t.Errorf("Month = %d, want 6", birthday.Date.Month)
	}
}

// TestApiPersonToDomain_NilFields tests handling of nil fields.
func TestApiPersonToDomain_NilFields(t *testing.T) {
	apiPerson := &people.Person{
		ResourceName: "people/c123",
	}

	contact := apiPersonToDomain(apiPerson)

	if contact == nil {
		t.Fatal("expected non-nil contact")
	}
	if len(contact.Names) != 0 {
		t.Errorf("expected 0 names, got %d", len(contact.Names))
	}
	if len(contact.EmailAddresses) != 0 {
		t.Errorf("expected 0 emails, got %d", len(contact.EmailAddresses))
	}
}

// TestDomainToApiPerson_BasicFields tests basic field conversion.
func TestDomainToApiPerson_BasicFields(t *testing.T) {
	contact := &contacts.Contact{
		ResourceName: "people/c123",
		ETag:         "etag456",
	}

	apiPerson := domainToApiPerson(contact)

	if apiPerson.ResourceName != contact.ResourceName {
		t.Errorf("ResourceName = %s, want %s", apiPerson.ResourceName, contact.ResourceName)
	}
	if apiPerson.Etag != contact.ETag {
		t.Errorf("Etag = %s, want %s", apiPerson.Etag, contact.ETag)
	}
}

// TestDomainToApiPerson_Names tests name conversion.
func TestDomainToApiPerson_Names(t *testing.T) {
	contact := &contacts.Contact{
		Names: []contacts.Name{
			{
				GivenName:  "John",
				FamilyName: "Doe",
			},
		},
	}

	apiPerson := domainToApiPerson(contact)

	if len(apiPerson.Names) != 1 {
		t.Fatalf("expected 1 name, got %d", len(apiPerson.Names))
	}

	name := apiPerson.Names[0]
	if name.GivenName != "John" {
		t.Errorf("GivenName = %s, want John", name.GivenName)
	}
	if name.FamilyName != "Doe" {
		t.Errorf("FamilyName = %s, want Doe", name.FamilyName)
	}
}

// TestDomainToApiPerson_Emails tests email conversion.
func TestDomainToApiPerson_Emails(t *testing.T) {
	contact := &contacts.Contact{
		EmailAddresses: []contacts.Email{
			{
				Value: "john@example.com",
				Type:  "work",
			},
		},
	}

	apiPerson := domainToApiPerson(contact)

	if len(apiPerson.EmailAddresses) != 1 {
		t.Fatalf("expected 1 email, got %d", len(apiPerson.EmailAddresses))
	}

	email := apiPerson.EmailAddresses[0]
	if email.Value != "john@example.com" {
		t.Errorf("Value = %s, want john@example.com", email.Value)
	}
}

// TestApiGroupToDomain_BasicFields tests basic group mapping.
func TestApiGroupToDomain_BasicFields(t *testing.T) {
	apiGroup := &people.ContactGroup{
		ResourceName:  "contactGroups/g123",
		Etag:          "etag789",
		Name:          "Friends",
		FormattedName: "Friends",
		GroupType:     "USER_CONTACT_GROUP",
		MemberCount:   5,
	}

	group := apiGroupToDomain(apiGroup)

	if group.ResourceName != apiGroup.ResourceName {
		t.Errorf("ResourceName = %s, want %s", group.ResourceName, apiGroup.ResourceName)
	}
	if group.ETag != apiGroup.Etag {
		t.Errorf("ETag = %s, want %s", group.ETag, apiGroup.Etag)
	}
	if group.Name != apiGroup.Name {
		t.Errorf("Name = %s, want %s", group.Name, apiGroup.Name)
	}
	if group.MemberCount != int(apiGroup.MemberCount) {
		t.Errorf("MemberCount = %d, want %d", group.MemberCount, apiGroup.MemberCount)
	}
}

// TestApiGroupToDomain_Metadata tests group metadata mapping.
func TestApiGroupToDomain_Metadata(t *testing.T) {
	updateTime := "2023-01-15T10:30:00Z"
	parsedTime, _ := time.Parse(time.RFC3339, updateTime)

	apiGroup := &people.ContactGroup{
		ResourceName: "contactGroups/g123",
		Metadata: &people.ContactGroupMetadata{
			UpdateTime: updateTime,
			Deleted:    false,
		},
	}

	group := apiGroupToDomain(apiGroup)

	if group.Metadata == nil {
		t.Fatal("expected non-nil metadata")
	}
	if !group.Metadata.UpdateTime.Equal(parsedTime) {
		t.Errorf("UpdateTime = %v, want %v", group.Metadata.UpdateTime, parsedTime)
	}
	if group.Metadata.Deleted {
		t.Error("expected deleted to be false")
	}
}

// TestApiGroupToDomain_MemberResourceNames tests member resource names mapping.
func TestApiGroupToDomain_MemberResourceNames(t *testing.T) {
	apiGroup := &people.ContactGroup{
		ResourceName:        "contactGroups/g123",
		MemberResourceNames: []string{"people/c1", "people/c2", "people/c3"},
	}

	group := apiGroupToDomain(apiGroup)

	if len(group.MemberResourceNames) != 3 {
		t.Fatalf("expected 3 member resource names, got %d", len(group.MemberResourceNames))
	}
	if group.MemberResourceNames[0] != "people/c1" {
		t.Errorf("MemberResourceName[0] = %s, want people/c1", group.MemberResourceNames[0])
	}
}

// TestDomainToApiGroup_BasicFields tests basic group conversion.
func TestDomainToApiGroup_BasicFields(t *testing.T) {
	group := &contacts.ContactGroup{
		ResourceName: "contactGroups/g123",
		ETag:         "etag789",
		Name:         "Friends",
	}

	apiGroup := domainToApiGroup(group)

	if apiGroup.ResourceName != group.ResourceName {
		t.Errorf("ResourceName = %s, want %s", apiGroup.ResourceName, group.ResourceName)
	}
	if apiGroup.Etag != group.ETag {
		t.Errorf("Etag = %s, want %s", apiGroup.Etag, group.ETag)
	}
	if apiGroup.Name != group.Name {
		t.Errorf("Name = %s, want %s", apiGroup.Name, group.Name)
	}
}

// TestBidirectionalTransformation_Contact tests round-trip conversion for contacts.
func TestBidirectionalTransformation_Contact(t *testing.T) {
	original := &contacts.Contact{
		ResourceName: "people/c123",
		ETag:         "etag456",
		Names: []contacts.Name{
			{
				GivenName:  "John",
				FamilyName: "Doe",
			},
		},
		EmailAddresses: []contacts.Email{
			{
				Value: "john@example.com",
				Type:  "work",
			},
		},
	}

	apiPerson := domainToApiPerson(original)
	result := apiPersonToDomain(apiPerson)

	if result.ResourceName != original.ResourceName {
		t.Errorf("ResourceName = %s, want %s", result.ResourceName, original.ResourceName)
	}
	if len(result.Names) != len(original.Names) {
		t.Errorf("Names count = %d, want %d", len(result.Names), len(original.Names))
	}
	if len(result.EmailAddresses) != len(original.EmailAddresses) {
		t.Errorf("EmailAddresses count = %d, want %d", len(result.EmailAddresses), len(original.EmailAddresses))
	}
}

// TestBidirectionalTransformation_Group tests round-trip conversion for groups.
func TestBidirectionalTransformation_Group(t *testing.T) {
	original := &contacts.ContactGroup{
		ResourceName: "contactGroups/g123",
		ETag:         "etag789",
		Name:         "Friends",
		GroupType:    "USER_CONTACT_GROUP",
	}

	apiGroup := domainToApiGroup(original)
	result := apiGroupToDomain(apiGroup)

	if result.ResourceName != original.ResourceName {
		t.Errorf("ResourceName = %s, want %s", result.ResourceName, original.ResourceName)
	}
	if result.Name != original.Name {
		t.Errorf("Name = %s, want %s", result.Name, original.Name)
	}
}

// TestEmptyArrays tests handling of empty arrays.
func TestEmptyArrays(t *testing.T) {
	contact := &contacts.Contact{
		ResourceName:   "people/c123",
		Names:          []contacts.Name{},
		EmailAddresses: []contacts.Email{},
		PhoneNumbers:   []contacts.Phone{},
	}

	apiPerson := domainToApiPerson(contact)
	result := apiPersonToDomain(apiPerson)

	if len(result.Names) != 0 {
		t.Errorf("expected 0 names, got %d", len(result.Names))
	}
	if len(result.EmailAddresses) != 0 {
		t.Errorf("expected 0 emails, got %d", len(result.EmailAddresses))
	}
}

// TestDateConversion tests date field conversion.
func TestDateConversion(t *testing.T) {
	apiDate := &people.Date{
		Year:  1990,
		Month: 6,
		Day:   15,
	}

	domainDate := apiDateToDomain(apiDate)

	if domainDate == nil {
		t.Fatal("expected non-nil date")
	}
	if domainDate.Year != 1990 {
		t.Errorf("Year = %d, want 1990", domainDate.Year)
	}
	if domainDate.Month != 6 {
		t.Errorf("Month = %d, want 6", domainDate.Month)
	}
	if domainDate.Day != 15 {
		t.Errorf("Day = %d, want 15", domainDate.Day)
	}

	apiDateResult := domainDateToApi(domainDate)
	if apiDateResult.Year != apiDate.Year {
		t.Errorf("Year = %d, want %d", apiDateResult.Year, apiDate.Year)
	}
}

// TestNilDateConversion tests nil date handling.
func TestNilDateConversion(t *testing.T) {
	domainDate := apiDateToDomain(nil)
	if domainDate != nil {
		t.Error("expected nil date")
	}

	apiDate := domainDateToApi(nil)
	if apiDate != nil {
		t.Error("expected nil date")
	}
}
