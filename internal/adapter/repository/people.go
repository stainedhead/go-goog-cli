package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/stainedhead/go-goog-cli/internal/domain/contacts"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

// personFields defines the fields to retrieve for Person resources.
const personFields = "names,emailAddresses,phoneNumbers,addresses,organizations,birthdays,biographies,photos,urls,memberships,metadata"

// groupFields defines the fields to retrieve for ContactGroup resources.
const groupFields = "name,groupType,memberCount,memberResourceNames,metadata"

// PeopleRepository is the base repository that wraps the Google People API service.
type PeopleRepository struct {
	service     *people.Service
	maxRetries  int
	baseBackoff time.Duration
}

// PeopleContactRepository implements ContactRepository using the Google People API.
type PeopleContactRepository struct {
	*PeopleRepository
}

// PeopleGroupRepository implements ContactGroupRepository using the Google People API.
type PeopleGroupRepository struct {
	*PeopleRepository
}

// Compile-time interface compliance checks.
var (
	_ contacts.ContactRepository      = (*PeopleContactRepository)(nil)
	_ contacts.ContactGroupRepository = (*PeopleGroupRepository)(nil)
)

// NewPeopleRepository creates a new PeopleRepository with the given OAuth2 token source.
func NewPeopleRepository(ctx context.Context, tokenSource oauth2.TokenSource) (*PeopleRepository, error) {
	httpClient := oauth2.NewClient(ctx, tokenSource)

	service, err := people.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("failed to create People service: %w", err)
	}

	return &PeopleRepository{
		service:     service,
		maxRetries:  defaultMaxRetries,
		baseBackoff: defaultBaseBackoff,
	}, nil
}

// NewPeopleRepositoryWithService creates a PeopleRepository with a pre-configured service.
// This is useful for testing with mock servers.
func NewPeopleRepositoryWithService(service *people.Service) *PeopleRepository {
	return &PeopleRepository{
		service:     service,
		maxRetries:  defaultMaxRetries,
		baseBackoff: defaultBaseBackoff,
	}
}

// NewPeopleContactRepository creates a new PeopleContactRepository.
func NewPeopleContactRepository(repo *PeopleRepository) *PeopleContactRepository {
	return &PeopleContactRepository{PeopleRepository: repo}
}

// NewPeopleGroupRepository creates a new PeopleGroupRepository.
func NewPeopleGroupRepository(repo *PeopleRepository) *PeopleGroupRepository {
	return &PeopleGroupRepository{PeopleRepository: repo}
}

// =============================================================================
// ContactRepository Implementation
// =============================================================================

// List retrieves all contacts with pagination support.
func (r *PeopleContactRepository) List(ctx context.Context, opts contacts.ListOptions) (*contacts.ListResult[*contacts.Contact], error) {
	call := r.service.People.Connections.List("people/me")
	call = call.PersonFields(personFields)

	if opts.MaxResults > 0 {
		call = call.PageSize(opts.MaxResults)
	}
	if opts.PageToken != "" {
		call = call.PageToken(opts.PageToken)
	}
	if opts.SortOrder != "" {
		call = call.SortOrder(opts.SortOrder)
	}

	result, err := retryWithBackoff(ctx, r.maxRetries, defaultBaseBackoff, func() (*people.ListConnectionsResponse, error) {
		return call.Do()
	})
	if err != nil {
		return nil, mapPeopleError(err, "list contacts")
	}

	domainContacts := make([]*contacts.Contact, 0, len(result.Connections))
	for _, person := range result.Connections {
		domainContacts = append(domainContacts, apiPersonToDomain(person))
	}

	return &contacts.ListResult[*contacts.Contact]{
		Items:         domainContacts,
		NextPageToken: result.NextPageToken,
		TotalSize:     int(result.TotalItems),
	}, nil
}

// Get retrieves a specific contact by resource name.
func (r *PeopleContactRepository) Get(ctx context.Context, resourceName string) (*contacts.Contact, error) {
	call := r.service.People.Get(resourceName)
	call = call.PersonFields(personFields)

	result, err := retryWithBackoff(ctx, r.maxRetries, defaultBaseBackoff, func() (*people.Person, error) {
		return call.Do()
	})
	if err != nil {
		return nil, mapPeopleError(err, "get contact")
	}

	return apiPersonToDomain(result), nil
}

// Create creates a new contact.
func (r *PeopleContactRepository) Create(ctx context.Context, contact *contacts.Contact) (*contacts.Contact, error) {
	apiPerson := domainToApiPerson(contact)

	call := r.service.People.CreateContact(apiPerson)
	call = call.PersonFields(personFields)

	result, err := retryWithBackoff(ctx, r.maxRetries, defaultBaseBackoff, func() (*people.Person, error) {
		return call.Do()
	})
	if err != nil {
		return nil, mapPeopleError(err, "create contact")
	}

	return apiPersonToDomain(result), nil
}

// Update updates an existing contact.
func (r *PeopleContactRepository) Update(ctx context.Context, contact *contacts.Contact, updateMask []string) (*contacts.Contact, error) {
	apiPerson := domainToApiPerson(contact)

	call := r.service.People.UpdateContact(contact.ResourceName, apiPerson)
	call = call.UpdatePersonFields(joinUpdateMask(updateMask))
	call = call.PersonFields(personFields)

	result, err := retryWithBackoff(ctx, r.maxRetries, defaultBaseBackoff, func() (*people.Person, error) {
		return call.Do()
	})
	if err != nil {
		return nil, mapPeopleError(err, "update contact")
	}

	return apiPersonToDomain(result), nil
}

// Delete deletes a contact.
func (r *PeopleContactRepository) Delete(ctx context.Context, resourceName string) error {
	call := r.service.People.DeleteContact(resourceName)

	_, err := retryWithBackoff(ctx, r.maxRetries, defaultBaseBackoff, func() (*people.Empty, error) {
		return call.Do()
	})
	if err != nil {
		return mapPeopleError(err, "delete contact")
	}

	return nil
}

// Search searches for contacts by query.
func (r *PeopleContactRepository) Search(ctx context.Context, opts contacts.SearchOptions) (*contacts.ListResult[*contacts.Contact], error) {
	call := r.service.People.SearchContacts()
	call = call.Query(opts.Query)
	call = call.ReadMask(personFields)

	if opts.MaxResults > 0 {
		call = call.PageSize(opts.MaxResults)
	}

	result, err := retryWithBackoff(ctx, r.maxRetries, defaultBaseBackoff, func() (*people.SearchResponse, error) {
		return call.Do()
	})
	if err != nil {
		return nil, mapPeopleError(err, "search contacts")
	}

	domainContacts := make([]*contacts.Contact, 0, len(result.Results))
	for _, searchResult := range result.Results {
		if searchResult.Person != nil {
			domainContacts = append(domainContacts, apiPersonToDomain(searchResult.Person))
		}
	}

	return &contacts.ListResult[*contacts.Contact]{
		Items:         domainContacts,
		NextPageToken: "",
		TotalSize:     len(domainContacts),
	}, nil
}

// BatchGet retrieves multiple contacts by resource names.
func (r *PeopleContactRepository) BatchGet(ctx context.Context, resourceNames []string) ([]*contacts.Contact, error) {
	call := r.service.People.GetBatchGet()
	call = call.ResourceNames(resourceNames...)
	call = call.PersonFields(personFields)

	result, err := retryWithBackoff(ctx, r.maxRetries, defaultBaseBackoff, func() (*people.GetPeopleResponse, error) {
		return call.Do()
	})
	if err != nil {
		return nil, mapPeopleError(err, "batch get contacts")
	}

	domainContacts := make([]*contacts.Contact, 0, len(result.Responses))
	for _, response := range result.Responses {
		if response.Person != nil {
			domainContacts = append(domainContacts, apiPersonToDomain(response.Person))
		}
	}

	return domainContacts, nil
}

// =============================================================================
// ContactGroupRepository Implementation
// =============================================================================

// List retrieves all contact groups.
func (r *PeopleGroupRepository) List(ctx context.Context) ([]*contacts.ContactGroup, error) {
	call := r.service.ContactGroups.List()
	call = call.GroupFields(groupFields)

	result, err := retryWithBackoff(ctx, r.maxRetries, defaultBaseBackoff, func() (*people.ListContactGroupsResponse, error) {
		return call.Do()
	})
	if err != nil {
		return nil, mapPeopleError(err, "list contact groups")
	}

	groups := make([]*contacts.ContactGroup, 0, len(result.ContactGroups))
	for _, apiGroup := range result.ContactGroups {
		groups = append(groups, apiGroupToDomain(apiGroup))
	}

	return groups, nil
}

// Get retrieves a specific contact group by resource name.
func (r *PeopleGroupRepository) Get(ctx context.Context, resourceName string) (*contacts.ContactGroup, error) {
	call := r.service.ContactGroups.Get(resourceName)
	call = call.GroupFields(groupFields)
	call = call.MaxMembers(1000)

	result, err := retryWithBackoff(ctx, r.maxRetries, defaultBaseBackoff, func() (*people.ContactGroup, error) {
		return call.Do()
	})
	if err != nil {
		return nil, mapPeopleError(err, "get contact group")
	}

	return apiGroupToDomain(result), nil
}

// Create creates a new contact group.
func (r *PeopleGroupRepository) Create(ctx context.Context, group *contacts.ContactGroup) (*contacts.ContactGroup, error) {
	request := &people.CreateContactGroupRequest{
		ContactGroup: domainToApiGroup(group),
	}

	call := r.service.ContactGroups.Create(request)

	result, err := retryWithBackoff(ctx, r.maxRetries, defaultBaseBackoff, func() (*people.ContactGroup, error) {
		return call.Do()
	})
	if err != nil {
		return nil, mapPeopleError(err, "create contact group")
	}

	return apiGroupToDomain(result), nil
}

// Update updates an existing contact group.
func (r *PeopleGroupRepository) Update(ctx context.Context, group *contacts.ContactGroup) (*contacts.ContactGroup, error) {
	apiGroup := domainToApiGroup(group)

	request := &people.UpdateContactGroupRequest{
		ContactGroup: apiGroup,
	}

	call := r.service.ContactGroups.Update(group.ResourceName, request)

	result, err := retryWithBackoff(ctx, r.maxRetries, defaultBaseBackoff, func() (*people.ContactGroup, error) {
		return call.Do()
	})
	if err != nil {
		return nil, mapPeopleError(err, "update contact group")
	}

	return apiGroupToDomain(result), nil
}

// Delete deletes a contact group.
func (r *PeopleGroupRepository) Delete(ctx context.Context, resourceName string) error {
	call := r.service.ContactGroups.Delete(resourceName)
	call = call.DeleteContacts(false)

	_, err := retryWithBackoff(ctx, r.maxRetries, defaultBaseBackoff, func() (*people.Empty, error) {
		return call.Do()
	})
	if err != nil {
		return mapPeopleError(err, "delete contact group")
	}

	return nil
}

// ListMembers retrieves all members of a contact group.
func (r *PeopleGroupRepository) ListMembers(ctx context.Context, resourceName string, opts contacts.ListOptions) (*contacts.ListResult[*contacts.Contact], error) {
	call := r.service.ContactGroups.Get(resourceName)
	call = call.GroupFields("memberResourceNames")
	call = call.MaxMembers(1000)

	groupResult, err := retryWithBackoff(ctx, r.maxRetries, defaultBaseBackoff, func() (*people.ContactGroup, error) {
		return call.Do()
	})
	if err != nil {
		return nil, mapPeopleError(err, "list group members")
	}

	if len(groupResult.MemberResourceNames) == 0 {
		return &contacts.ListResult[*contacts.Contact]{
			Items:         []*contacts.Contact{},
			NextPageToken: "",
			TotalSize:     0,
		}, nil
	}

	// Batch get the member contacts
	contactRepo := NewPeopleContactRepository(r.PeopleRepository)
	members, err := contactRepo.BatchGet(ctx, groupResult.MemberResourceNames)
	if err != nil {
		return nil, err
	}

	return &contacts.ListResult[*contacts.Contact]{
		Items:         members,
		NextPageToken: "",
		TotalSize:     len(members),
	}, nil
}

// AddMembers adds members to a contact group.
func (r *PeopleGroupRepository) AddMembers(ctx context.Context, groupResourceName string, contactResourceNames []string) error {
	request := &people.ModifyContactGroupMembersRequest{
		ResourceNamesToAdd: contactResourceNames,
	}

	call := r.service.ContactGroups.Members.Modify(groupResourceName, request)

	_, err := retryWithBackoff(ctx, r.maxRetries, defaultBaseBackoff, func() (*people.ModifyContactGroupMembersResponse, error) {
		return call.Do()
	})
	if err != nil {
		return mapPeopleError(err, "add group members")
	}

	return nil
}

// RemoveMembers removes members from a contact group.
func (r *PeopleGroupRepository) RemoveMembers(ctx context.Context, groupResourceName string, contactResourceNames []string) error {
	request := &people.ModifyContactGroupMembersRequest{
		ResourceNamesToRemove: contactResourceNames,
	}

	call := r.service.ContactGroups.Members.Modify(groupResourceName, request)

	_, err := retryWithBackoff(ctx, r.maxRetries, defaultBaseBackoff, func() (*people.ModifyContactGroupMembersResponse, error) {
		return call.Do()
	})
	if err != nil {
		return mapPeopleError(err, "remove group members")
	}

	return nil
}

// =============================================================================
// Helper Functions
// =============================================================================

// joinUpdateMask joins update mask fields with commas.
func joinUpdateMask(fields []string) string {
	return strings.Join(fields, ",")
}

// mapPeopleError maps Google People API errors to domain errors.
func mapPeopleError(err error, operation string) error {
	if err == nil {
		return nil
	}

	// Use the existing mapAPIError function from errors.go
	return mapAPIError(err, operation)
}
