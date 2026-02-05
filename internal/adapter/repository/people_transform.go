package repository

import (
	"time"

	"github.com/stainedhead/go-goog-cli/internal/domain/contacts"
	"google.golang.org/api/people/v1"
)

// =============================================================================
// Contact Transformations
// =============================================================================

// apiPersonToDomain converts a People API Person to a domain Contact.
func apiPersonToDomain(api *people.Person) *contacts.Contact {
	if api == nil {
		return nil
	}

	contact := &contacts.Contact{
		ResourceName:   api.ResourceName,
		ETag:           api.Etag,
		Names:          make([]contacts.Name, 0),
		Nicknames:      make([]contacts.Nickname, 0),
		EmailAddresses: make([]contacts.Email, 0),
		PhoneNumbers:   make([]contacts.Phone, 0),
		Addresses:      make([]contacts.Address, 0),
		Organizations:  make([]contacts.Organization, 0),
		Birthdays:      make([]contacts.Birthday, 0),
		Biographies:    make([]contacts.Biography, 0),
		Photos:         make([]contacts.Photo, 0),
		URLs:           make([]contacts.URL, 0),
		Memberships:    make([]contacts.Membership, 0),
	}

	// Convert names
	for _, name := range api.Names {
		contact.Names = append(contact.Names, contacts.Name{
			DisplayName:        name.DisplayName,
			GivenName:          name.GivenName,
			FamilyName:         name.FamilyName,
			MiddleName:         name.MiddleName,
			HonorificPrefix:    name.HonorificPrefix,
			HonorificSuffix:    name.HonorificSuffix,
			PhoneticGivenName:  name.PhoneticGivenName,
			PhoneticFamilyName: name.PhoneticFamilyName,
			PhoneticFullName:   name.PhoneticFullName,
		})
	}

	// Convert nicknames
	for _, nickname := range api.Nicknames {
		contact.Nicknames = append(contact.Nicknames, contacts.Nickname{
			Value: nickname.Value,
			Type:  nickname.Type,
		})
	}

	// Convert email addresses
	for _, email := range api.EmailAddresses {
		isPrimary := false
		if email.Metadata != nil {
			isPrimary = email.Metadata.Primary
		}
		contact.EmailAddresses = append(contact.EmailAddresses, contacts.Email{
			Value:         email.Value,
			Type:          email.Type,
			DisplayName:   email.DisplayName,
			Primary:       isPrimary,
			FormattedType: email.FormattedType,
		})
	}

	// Convert phone numbers
	for _, phone := range api.PhoneNumbers {
		isPrimary := false
		if phone.Metadata != nil {
			isPrimary = phone.Metadata.Primary
		}
		contact.PhoneNumbers = append(contact.PhoneNumbers, contacts.Phone{
			Value:         phone.Value,
			Type:          phone.Type,
			Primary:       isPrimary,
			FormattedType: phone.FormattedType,
		})
	}

	// Convert addresses
	for _, addr := range api.Addresses {
		contact.Addresses = append(contact.Addresses, contacts.Address{
			FormattedValue:  addr.FormattedValue,
			Type:            addr.Type,
			StreetAddress:   addr.StreetAddress,
			City:            addr.City,
			Region:          addr.Region,
			PostalCode:      addr.PostalCode,
			Country:         addr.Country,
			CountryCode:     addr.CountryCode,
			ExtendedAddress: addr.ExtendedAddress,
			PoBox:           addr.PoBox,
			FormattedType:   addr.FormattedType,
		})
	}

	// Convert organizations
	for _, org := range api.Organizations {
		contact.Organizations = append(contact.Organizations, contacts.Organization{
			Name:       org.Name,
			Title:      org.Title,
			Department: org.Department,
			Symbol:     org.Symbol,
			Domain:     org.Domain,
			Location:   org.Location,
			Type:       org.Type,
			StartDate:  apiDateToDomain(org.StartDate),
			EndDate:    apiDateToDomain(org.EndDate),
			Current:    org.Current,
		})
	}

	// Convert birthdays
	for _, birthday := range api.Birthdays {
		contact.Birthdays = append(contact.Birthdays, contacts.Birthday{
			Date: apiDateToDomain(birthday.Date),
			Text: birthday.Text,
		})
	}

	// Convert biographies
	for _, bio := range api.Biographies {
		contact.Biographies = append(contact.Biographies, contacts.Biography{
			Value:       bio.Value,
			ContentType: bio.ContentType,
		})
	}

	// Convert photos
	for _, photo := range api.Photos {
		contact.Photos = append(contact.Photos, contacts.Photo{
			URL:     photo.Url,
			Default: photo.Default,
		})
	}

	// Convert URLs
	for _, url := range api.Urls {
		contact.URLs = append(contact.URLs, contacts.URL{
			Value: url.Value,
			Type:  url.Type,
		})
	}

	// Convert memberships
	for _, membership := range api.Memberships {
		domainMembership := contacts.Membership{}
		if membership.ContactGroupMembership != nil {
			domainMembership.ContactGroupMembership = &contacts.ContactGroupMembership{
				ContactGroupResourceName: membership.ContactGroupMembership.ContactGroupResourceName,
				ContactGroupID:           membership.ContactGroupMembership.ContactGroupId,
			}
		}
		if membership.DomainMembership != nil {
			domainMembership.DomainMembership = &contacts.DomainMembership{
				InViewerDomain: membership.DomainMembership.InViewerDomain,
			}
		}
		contact.Memberships = append(contact.Memberships, domainMembership)
	}

	// Convert metadata
	if api.Metadata != nil {
		contact.Metadata = &contacts.ResourceMetadata{
			Sources: make([]contacts.Source, 0),
		}
		for _, source := range api.Metadata.Sources {
			updateTime := time.Time{}
			if source.UpdateTime != "" {
				updateTime, _ = time.Parse(time.RFC3339, source.UpdateTime)
			}
			contact.Metadata.Sources = append(contact.Metadata.Sources, contacts.Source{
				Type:       source.Type,
				ID:         source.Id,
				ETag:       source.Etag,
				UpdateTime: updateTime,
			})
		}
	}

	return contact
}

// domainToApiPerson converts a domain Contact to a People API Person.
func domainToApiPerson(domain *contacts.Contact) *people.Person {
	if domain == nil {
		return nil
	}

	api := &people.Person{
		ResourceName:   domain.ResourceName,
		Etag:           domain.ETag,
		Names:          make([]*people.Name, 0),
		Nicknames:      make([]*people.Nickname, 0),
		EmailAddresses: make([]*people.EmailAddress, 0),
		PhoneNumbers:   make([]*people.PhoneNumber, 0),
		Addresses:      make([]*people.Address, 0),
		Organizations:  make([]*people.Organization, 0),
		Birthdays:      make([]*people.Birthday, 0),
		Biographies:    make([]*people.Biography, 0),
		Photos:         make([]*people.Photo, 0),
		Urls:           make([]*people.Url, 0),
		Memberships:    make([]*people.Membership, 0),
	}

	// Convert names
	for _, name := range domain.Names {
		api.Names = append(api.Names, &people.Name{
			DisplayName:        name.DisplayName,
			GivenName:          name.GivenName,
			FamilyName:         name.FamilyName,
			MiddleName:         name.MiddleName,
			HonorificPrefix:    name.HonorificPrefix,
			HonorificSuffix:    name.HonorificSuffix,
			PhoneticGivenName:  name.PhoneticGivenName,
			PhoneticFamilyName: name.PhoneticFamilyName,
			PhoneticFullName:   name.PhoneticFullName,
		})
	}

	// Convert nicknames
	for _, nickname := range domain.Nicknames {
		api.Nicknames = append(api.Nicknames, &people.Nickname{
			Value: nickname.Value,
			Type:  nickname.Type,
		})
	}

	// Convert email addresses
	for _, email := range domain.EmailAddresses {
		api.EmailAddresses = append(api.EmailAddresses, &people.EmailAddress{
			Value:         email.Value,
			Type:          email.Type,
			DisplayName:   email.DisplayName,
			FormattedType: email.FormattedType,
		})
	}

	// Convert phone numbers
	for _, phone := range domain.PhoneNumbers {
		api.PhoneNumbers = append(api.PhoneNumbers, &people.PhoneNumber{
			Value:         phone.Value,
			Type:          phone.Type,
			FormattedType: phone.FormattedType,
		})
	}

	// Convert addresses
	for _, addr := range domain.Addresses {
		api.Addresses = append(api.Addresses, &people.Address{
			FormattedValue:  addr.FormattedValue,
			Type:            addr.Type,
			StreetAddress:   addr.StreetAddress,
			City:            addr.City,
			Region:          addr.Region,
			PostalCode:      addr.PostalCode,
			Country:         addr.Country,
			CountryCode:     addr.CountryCode,
			ExtendedAddress: addr.ExtendedAddress,
			PoBox:           addr.PoBox,
			FormattedType:   addr.FormattedType,
		})
	}

	// Convert organizations
	for _, org := range domain.Organizations {
		api.Organizations = append(api.Organizations, &people.Organization{
			Name:       org.Name,
			Title:      org.Title,
			Department: org.Department,
			Symbol:     org.Symbol,
			Domain:     org.Domain,
			Location:   org.Location,
			Type:       org.Type,
			StartDate:  domainDateToApi(org.StartDate),
			EndDate:    domainDateToApi(org.EndDate),
			Current:    org.Current,
		})
	}

	// Convert birthdays
	for _, birthday := range domain.Birthdays {
		api.Birthdays = append(api.Birthdays, &people.Birthday{
			Date: domainDateToApi(birthday.Date),
			Text: birthday.Text,
		})
	}

	// Convert biographies
	for _, bio := range domain.Biographies {
		api.Biographies = append(api.Biographies, &people.Biography{
			Value:       bio.Value,
			ContentType: bio.ContentType,
		})
	}

	// Convert photos
	for _, photo := range domain.Photos {
		api.Photos = append(api.Photos, &people.Photo{
			Url:     photo.URL,
			Default: photo.Default,
		})
	}

	// Convert URLs
	for _, url := range domain.URLs {
		api.Urls = append(api.Urls, &people.Url{
			Value: url.Value,
			Type:  url.Type,
		})
	}

	// Convert memberships
	for _, membership := range domain.Memberships {
		apiMembership := &people.Membership{}
		if membership.ContactGroupMembership != nil {
			apiMembership.ContactGroupMembership = &people.ContactGroupMembership{
				ContactGroupResourceName: membership.ContactGroupMembership.ContactGroupResourceName,
				ContactGroupId:           membership.ContactGroupMembership.ContactGroupID,
			}
		}
		if membership.DomainMembership != nil {
			apiMembership.DomainMembership = &people.DomainMembership{
				InViewerDomain: membership.DomainMembership.InViewerDomain,
			}
		}
		api.Memberships = append(api.Memberships, apiMembership)
	}

	return api
}

// =============================================================================
// Contact Group Transformations
// =============================================================================

// apiGroupToDomain converts a People API ContactGroup to a domain ContactGroup.
func apiGroupToDomain(api *people.ContactGroup) *contacts.ContactGroup {
	if api == nil {
		return nil
	}

	group := &contacts.ContactGroup{
		ResourceName:        api.ResourceName,
		ETag:                api.Etag,
		Name:                api.Name,
		FormattedName:       api.FormattedName,
		GroupType:           api.GroupType,
		MemberCount:         int(api.MemberCount),
		MemberResourceNames: api.MemberResourceNames,
	}

	// Convert metadata
	if api.Metadata != nil {
		updateTime := time.Time{}
		if api.Metadata.UpdateTime != "" {
			updateTime, _ = time.Parse(time.RFC3339, api.Metadata.UpdateTime)
		}
		group.Metadata = &contacts.GroupMetadata{
			UpdateTime: updateTime,
			Deleted:    api.Metadata.Deleted,
		}
	}

	return group
}

// domainToApiGroup converts a domain ContactGroup to a People API ContactGroup.
func domainToApiGroup(domain *contacts.ContactGroup) *people.ContactGroup {
	if domain == nil {
		return nil
	}

	api := &people.ContactGroup{
		ResourceName:        domain.ResourceName,
		Etag:                domain.ETag,
		Name:                domain.Name,
		FormattedName:       domain.FormattedName,
		GroupType:           domain.GroupType,
		MemberCount:         int64(domain.MemberCount),
		MemberResourceNames: domain.MemberResourceNames,
	}

	return api
}

// =============================================================================
// Date Transformations
// =============================================================================

// apiDateToDomain converts a People API Date to a domain Date.
func apiDateToDomain(api *people.Date) *contacts.Date {
	if api == nil {
		return nil
	}

	return &contacts.Date{
		Year:  int(api.Year),
		Month: int(api.Month),
		Day:   int(api.Day),
	}
}

// domainDateToApi converts a domain Date to a People API Date.
func domainDateToApi(domain *contacts.Date) *people.Date {
	if domain == nil {
		return nil
	}

	return &people.Date{
		Year:  int64(domain.Year),
		Month: int64(domain.Month),
		Day:   int64(domain.Day),
	}
}
