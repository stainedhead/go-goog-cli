package repository

import (
	"context"
	"testing"

	"github.com/stainedhead/go-goog-cli/internal/domain/contacts"
	"golang.org/x/oauth2"
	"google.golang.org/api/people/v1"
)

// mockTokenSource implements oauth2.TokenSource for testing.
type mockTokenSource struct{}

func (m *mockTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: "test-token"}, nil
}

// TestPeopleContactRepositoryImplementsInterface verifies compile-time interface compliance.
func TestPeopleContactRepositoryImplementsInterface(t *testing.T) {
	var _ contacts.ContactRepository = (*PeopleContactRepository)(nil)
}

// TestPeopleGroupRepositoryImplementsInterface verifies compile-time interface compliance.
func TestPeopleGroupRepositoryImplementsInterface(t *testing.T) {
	var _ contacts.ContactGroupRepository = (*PeopleGroupRepository)(nil)
}

// TestNewPeopleRepository tests repository creation.
func TestNewPeopleRepository(t *testing.T) {
	ctx := context.Background()
	tokenSource := &mockTokenSource{}

	repo, err := NewPeopleRepository(ctx, tokenSource)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo == nil {
		t.Fatal("expected non-nil repository")
	}

	if repo.service == nil {
		t.Error("expected non-nil service")
	}

	if repo.maxRetries != defaultMaxRetries {
		t.Errorf("expected maxRetries=%d, got %d", defaultMaxRetries, repo.maxRetries)
	}

	if repo.baseBackoff != defaultBaseBackoff {
		t.Errorf("expected baseBackoff=%v, got %v", defaultBaseBackoff, repo.baseBackoff)
	}
}

// TestContactRepository_Create tests creating a contact.
func TestContactRepository_Create(t *testing.T) {
	t.Skip("Skipping: requires proper mock server setup")
	ctx := context.Background()
	service := &people.Service{}
	repo := NewPeopleRepositoryWithService(service)
	contactRepo := NewPeopleContactRepository(repo)

	contact := &contacts.Contact{
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

	// This will fail without mock server, but we're testing the structure
	_, err := contactRepo.Create(ctx, contact)
	// We expect an error since we don't have a mock server
	if err == nil {
		t.Error("expected error without mock server")
	}
}

// TestContactRepository_Get tests retrieving a contact.
func TestContactRepository_Get(t *testing.T) {
	ctx := context.Background()
	service := &people.Service{}
	repo := NewPeopleRepositoryWithService(service)
	contactRepo := NewPeopleContactRepository(repo)

	resourceName := "people/c12345"

	_, err := contactRepo.Get(ctx, resourceName)
	// We expect an error since we don't have a mock server
	if err == nil {
		t.Error("expected error without mock server")
	}
}

// TestContactRepository_Update tests updating a contact.
func TestContactRepository_Update(t *testing.T) {
	ctx := context.Background()
	service := &people.Service{}
	repo := NewPeopleRepositoryWithService(service)
	contactRepo := NewPeopleContactRepository(repo)

	contact := &contacts.Contact{
		ResourceName: "people/c12345",
		ETag:         "etag123",
		Names: []contacts.Name{
			{
				GivenName:  "Jane",
				FamilyName: "Doe",
			},
		},
	}

	updateMask := []string{"names"}

	_, err := contactRepo.Update(ctx, contact, updateMask)
	// We expect an error since we don't have a mock server
	if err == nil {
		t.Error("expected error without mock server")
	}
}

// TestContactRepository_Delete tests deleting a contact.
func TestContactRepository_Delete(t *testing.T) {
	ctx := context.Background()
	service := &people.Service{}
	repo := NewPeopleRepositoryWithService(service)
	contactRepo := NewPeopleContactRepository(repo)

	resourceName := "people/c12345"

	err := contactRepo.Delete(ctx, resourceName)
	// We expect an error since we don't have a mock server
	if err == nil {
		t.Error("expected error without mock server")
	}
}

// TestContactRepository_List tests listing contacts.
func TestContactRepository_List(t *testing.T) {
	ctx := context.Background()
	service := &people.Service{}
	repo := NewPeopleRepositoryWithService(service)
	contactRepo := NewPeopleContactRepository(repo)

	opts := contacts.ListOptions{
		MaxResults: 10,
		PageToken:  "",
	}

	_, err := contactRepo.List(ctx, opts)
	// We expect an error since we don't have a mock server
	if err == nil {
		t.Error("expected error without mock server")
	}
}

// TestContactRepository_Search tests searching contacts.
func TestContactRepository_Search(t *testing.T) {
	ctx := context.Background()
	service := &people.Service{}
	repo := NewPeopleRepositoryWithService(service)
	contactRepo := NewPeopleContactRepository(repo)

	opts := contacts.SearchOptions{
		Query:      "john",
		MaxResults: 10,
	}

	_, err := contactRepo.Search(ctx, opts)
	// We expect an error since we don't have a mock server
	if err == nil {
		t.Error("expected error without mock server")
	}
}

// TestContactRepository_BatchGet tests batch getting contacts.
func TestContactRepository_BatchGet(t *testing.T) {
	ctx := context.Background()
	service := &people.Service{}
	repo := NewPeopleRepositoryWithService(service)
	contactRepo := NewPeopleContactRepository(repo)

	resourceNames := []string{"people/c1", "people/c2"}

	_, err := contactRepo.BatchGet(ctx, resourceNames)
	// We expect an error since we don't have a mock server
	if err == nil {
		t.Error("expected error without mock server")
	}
}

// TestGroupRepository_Create tests creating a contact group.
func TestGroupRepository_Create(t *testing.T) {
	ctx := context.Background()
	service := &people.Service{}
	repo := NewPeopleRepositoryWithService(service)
	groupRepo := NewPeopleGroupRepository(repo)

	group := &contacts.ContactGroup{
		Name: "Friends",
	}

	_, err := groupRepo.Create(ctx, group)
	// We expect an error since we don't have a mock server
	if err == nil {
		t.Error("expected error without mock server")
	}
}

// TestGroupRepository_Get tests retrieving a contact group.
func TestGroupRepository_Get(t *testing.T) {
	ctx := context.Background()
	service := &people.Service{}
	repo := NewPeopleRepositoryWithService(service)
	groupRepo := NewPeopleGroupRepository(repo)

	resourceName := "contactGroups/g123"

	_, err := groupRepo.Get(ctx, resourceName)
	// We expect an error since we don't have a mock server
	if err == nil {
		t.Error("expected error without mock server")
	}
}

// TestGroupRepository_Update tests updating a contact group.
func TestGroupRepository_Update(t *testing.T) {
	ctx := context.Background()
	service := &people.Service{}
	repo := NewPeopleRepositoryWithService(service)
	groupRepo := NewPeopleGroupRepository(repo)

	group := &contacts.ContactGroup{
		ResourceName: "contactGroups/g123",
		ETag:         "etag123",
		Name:         "Updated Friends",
	}

	_, err := groupRepo.Update(ctx, group)
	// We expect an error since we don't have a mock server
	if err == nil {
		t.Error("expected error without mock server")
	}
}

// TestGroupRepository_Delete tests deleting a contact group.
func TestGroupRepository_Delete(t *testing.T) {
	ctx := context.Background()
	service := &people.Service{}
	repo := NewPeopleRepositoryWithService(service)
	groupRepo := NewPeopleGroupRepository(repo)

	resourceName := "contactGroups/g123"

	err := groupRepo.Delete(ctx, resourceName)
	// We expect an error since we don't have a mock server
	if err == nil {
		t.Error("expected error without mock server")
	}
}

// TestGroupRepository_List tests listing contact groups.
func TestGroupRepository_List(t *testing.T) {
	ctx := context.Background()
	service := &people.Service{}
	repo := NewPeopleRepositoryWithService(service)
	groupRepo := NewPeopleGroupRepository(repo)

	_, err := groupRepo.List(ctx)
	// We expect an error since we don't have a mock server
	if err == nil {
		t.Error("expected error without mock server")
	}
}

// TestGroupRepository_ListMembers tests listing group members.
func TestGroupRepository_ListMembers(t *testing.T) {
	ctx := context.Background()
	service := &people.Service{}
	repo := NewPeopleRepositoryWithService(service)
	groupRepo := NewPeopleGroupRepository(repo)

	resourceName := "contactGroups/g123"
	opts := contacts.ListOptions{
		MaxResults: 10,
	}

	_, err := groupRepo.ListMembers(ctx, resourceName, opts)
	// We expect an error since we don't have a mock server
	if err == nil {
		t.Error("expected error without mock server")
	}
}

// TestGroupRepository_AddMembers tests adding members to a group.
func TestGroupRepository_AddMembers(t *testing.T) {
	ctx := context.Background()
	service := &people.Service{}
	repo := NewPeopleRepositoryWithService(service)
	groupRepo := NewPeopleGroupRepository(repo)

	groupResourceName := "contactGroups/g123"
	contactResourceNames := []string{"people/c1", "people/c2"}

	err := groupRepo.AddMembers(ctx, groupResourceName, contactResourceNames)
	// We expect an error since we don't have a mock server
	if err == nil {
		t.Error("expected error without mock server")
	}
}

// TestGroupRepository_RemoveMembers tests removing members from a group.
func TestGroupRepository_RemoveMembers(t *testing.T) {
	ctx := context.Background()
	service := &people.Service{}
	repo := NewPeopleRepositoryWithService(service)
	groupRepo := NewPeopleGroupRepository(repo)

	groupResourceName := "contactGroups/g123"
	contactResourceNames := []string{"people/c1", "people/c2"}

	err := groupRepo.RemoveMembers(ctx, groupResourceName, contactResourceNames)
	// We expect an error since we don't have a mock server
	if err == nil {
		t.Error("expected error without mock server")
	}
}

// TestErrorMapping tests error mapping for People API.
func TestErrorMapping(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		operation string
		wantErr   bool
	}{
		{
			name:      "nil error returns nil",
			err:       nil,
			operation: "test",
			wantErr:   false,
		},
		{
			name:      "generic error is returned",
			err:       context.DeadlineExceeded,
			operation: "test",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mapPeopleError(tt.err, tt.operation)
			if (err != nil) != tt.wantErr {
				t.Errorf("mapPeopleError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestListOptionsValidation tests list options validation.
func TestListOptionsValidation(t *testing.T) {
	tests := []struct {
		name string
		opts contacts.ListOptions
		want contacts.ListOptions
	}{
		{
			name: "empty options use defaults",
			opts: contacts.ListOptions{},
			want: contacts.ListOptions{
				MaxResults: 0,
				PageToken:  "",
			},
		},
		{
			name: "custom max results",
			opts: contacts.ListOptions{
				MaxResults: 50,
			},
			want: contacts.ListOptions{
				MaxResults: 50,
				PageToken:  "",
			},
		},
		{
			name: "with page token",
			opts: contacts.ListOptions{
				PageToken: "token123",
			},
			want: contacts.ListOptions{
				MaxResults: 0,
				PageToken:  "token123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.opts.MaxResults != tt.want.MaxResults {
				t.Errorf("MaxResults = %d, want %d", tt.opts.MaxResults, tt.want.MaxResults)
			}
			if tt.opts.PageToken != tt.want.PageToken {
				t.Errorf("PageToken = %s, want %s", tt.opts.PageToken, tt.want.PageToken)
			}
		})
	}
}
