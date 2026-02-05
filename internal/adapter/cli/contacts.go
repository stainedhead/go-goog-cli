package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stainedhead/go-goog-cli/internal/adapter/presenter"
	domaincontacts "github.com/stainedhead/go-goog-cli/internal/domain/contacts"
)

// Command flags for contacts actions.
var (
	contactsMaxResults    int64
	contactsPageToken     string
	contactsQuery         string
	contactsGivenName     string
	contactsFamilyName    string
	contactsEmail         string
	contactsEmailType     string
	contactsPhone         string
	contactsPhoneType     string
	contactsAddress       string
	contactsAddressType   string
	contactsOrganization  string
	contactsTitle         string
	contactsNotes         string
	contactsBirthday      string
	contactsURL           string
	contactsGroupName     string
	contactsDeleteConfirm bool
	contactsUpdateMask    string
)

// contactsCmd represents the contacts command group.
var contactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "Manage Google Contacts",
	Long: `Manage Google Contacts.

The contacts commands allow you to list, create, update, and manage
contacts and contact groups in your Google Contacts account.`,
}

// ================ Contact Commands ================

// contactsListCmd lists all contacts.
var contactsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all contacts",
	Long: `List all contacts in your Google Contacts account.

Use --max-results to limit the number of results.
Use --page-token to retrieve the next page of results.`,
	Example: `  # List all contacts
  goog contacts list

  # List with JSON output
  goog contacts list --format json

  # List with pagination
  goog contacts list --max-results 50`,
	Args: cobra.NoArgs,
	RunE: runContactsList,
}

// contactsGetCmd gets a specific contact.
var contactsGetCmd = &cobra.Command{
	Use:   "get <resource-name>",
	Short: "Get details of a specific contact",
	Long: `Get detailed information about a specific contact.

The resource-name should be in the format "people/c123456789".`,
	Example: `  # Get a contact
  goog contacts get people/c123456789

  # Get with JSON output
  goog contacts get people/c123456789 --format json`,
	Args: cobra.ExactArgs(1),
	RunE: runContactsGet,
}

// contactsCreateCmd creates a new contact.
var contactsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new contact",
	Long: `Create a new contact with the specified details.

Use flags to set contact information such as name, email, phone, etc.`,
	Example: `  # Create a contact with name
  goog contacts create --given-name John --family-name Doe

  # Create a contact with name and email
  goog contacts create --given-name Jane --family-name Smith --email jane@example.com

  # Create a contact with phone
  goog contacts create --given-name Bob --phone "+1234567890" --phone-type mobile`,
	Args: cobra.NoArgs,
	RunE: runContactsCreate,
}

// contactsUpdateCmd updates a contact.
var contactsUpdateCmd = &cobra.Command{
	Use:   "update <resource-name>",
	Short: "Update a contact",
	Long: `Update properties of an existing contact.

Use flags to update contact information.`,
	Example: `  # Update contact name
  goog contacts update people/c123 --given-name John --family-name Doe

  # Update email
  goog contacts update people/c123 --email newemail@example.com

  # Update phone
  goog contacts update people/c123 --phone "+1234567890"`,
	Args: cobra.ExactArgs(1),
	RunE: runContactsUpdate,
}

// contactsDeleteCmd deletes a contact.
var contactsDeleteCmd = &cobra.Command{
	Use:   "delete <resource-name>",
	Short: "Delete a contact",
	Long: `Delete a contact permanently.

WARNING: This action is irreversible. The contact will be permanently deleted.

The --confirm flag is required to prevent accidental deletion.`,
	Example: `  # Delete a contact (requires --confirm)
  goog contacts delete people/c123 --confirm`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !contactsDeleteConfirm {
			cmd.PrintErrln("Error: deletion requires --confirm flag")
			cmd.PrintErrln("Use --confirm to confirm this action")
			return fmt.Errorf("confirmation required")
		}
		return nil
	},
	RunE: runContactsDelete,
}

// contactsSearchCmd searches contacts.
var contactsSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search contacts",
	Long: `Search contacts by name, email, or other fields.

The search query will match against contact names, email addresses,
phone numbers, and other text fields.`,
	Example: `  # Search by name
  goog contacts search "John"

  # Search by email
  goog contacts search "john@example.com"

  # Search with max results
  goog contacts search "Smith" --max-results 20`,
	Args: cobra.ExactArgs(1),
	RunE: runContactsSearch,
}

// ================ Contact Group Commands ================

// contactsGroupsCmd lists all contact groups.
var contactsGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List all contact groups",
	Long: `List all contact groups in your Google Contacts account.

Contact groups allow you to organize contacts into categories.`,
	Example: `  # List all contact groups
  goog contacts groups

  # List with JSON output
  goog contacts groups --format json`,
	Args: cobra.NoArgs,
	RunE: runContactsGroups,
}

// contactsGroupCreateCmd creates a new contact group.
var contactsGroupCreateCmd = &cobra.Command{
	Use:   "group-create <name>",
	Short: "Create a new contact group",
	Long: `Create a new contact group with the specified name.

Contact groups help organize contacts into categories.`,
	Example: `  # Create a contact group
  goog contacts group-create "Family"

  # Create with different account
  goog contacts group-create "Work Contacts" --account work`,
	Args: cobra.ExactArgs(1),
	RunE: runContactsGroupCreate,
}

// contactsGroupUpdateCmd updates a contact group.
var contactsGroupUpdateCmd = &cobra.Command{
	Use:   "group-update <resource-name>",
	Short: "Update a contact group",
	Long: `Update the name of a contact group.

Use --group-name to specify the new name.`,
	Example: `  # Update group name
  goog contacts group-update contactGroups/g123 --group-name "New Name"`,
	Args: cobra.ExactArgs(1),
	RunE: runContactsGroupUpdate,
}

// contactsGroupDeleteCmd deletes a contact group.
var contactsGroupDeleteCmd = &cobra.Command{
	Use:   "group-delete <resource-name>",
	Short: "Delete a contact group",
	Long: `Delete a contact group permanently.

WARNING: This action is irreversible. The group will be permanently deleted,
but contacts in the group will not be deleted.

The --confirm flag is required to prevent accidental deletion.`,
	Example: `  # Delete a contact group (requires --confirm)
  goog contacts group-delete contactGroups/g123 --confirm`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !contactsDeleteConfirm {
			cmd.PrintErrln("Error: deletion requires --confirm flag")
			cmd.PrintErrln("Use --confirm to confirm this action")
			return fmt.Errorf("confirmation required")
		}
		return nil
	},
	RunE: runContactsGroupDelete,
}

// contactsGroupMembersCmd lists members of a contact group.
var contactsGroupMembersCmd = &cobra.Command{
	Use:   "group-members <resource-name>",
	Short: "List members of a contact group",
	Long:  `List all contacts that are members of the specified contact group.`,
	Example: `  # List group members
  goog contacts group-members contactGroups/g123

  # List with pagination
  goog contacts group-members contactGroups/g123 --max-results 50`,
	Args: cobra.ExactArgs(1),
	RunE: runContactsGroupMembers,
}

// contactsGroupAddCmd adds contacts to a group.
var contactsGroupAddCmd = &cobra.Command{
	Use:   "group-add <group-resource-name> <contact-resource-name...>",
	Short: "Add contacts to a group",
	Long: `Add one or more contacts to a contact group.

Specify the group resource name followed by one or more contact resource names.`,
	Example: `  # Add one contact to group
  goog contacts group-add contactGroups/g123 people/c456

  # Add multiple contacts to group
  goog contacts group-add contactGroups/g123 people/c456 people/c789`,
	Args: cobra.MinimumNArgs(2),
	RunE: runContactsGroupAdd,
}

// contactsGroupRemoveCmd removes contacts from a group.
var contactsGroupRemoveCmd = &cobra.Command{
	Use:   "group-remove <group-resource-name> <contact-resource-name...>",
	Short: "Remove contacts from a group",
	Long: `Remove one or more contacts from a contact group.

Specify the group resource name followed by one or more contact resource names.`,
	Example: `  # Remove one contact from group
  goog contacts group-remove contactGroups/g123 people/c456

  # Remove multiple contacts from group
  goog contacts group-remove contactGroups/g123 people/c456 people/c789`,
	Args: cobra.MinimumNArgs(2),
	RunE: runContactsGroupRemove,
}

// ================ Command Implementations ================

func runContactsList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	repo, err := getContactRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	opts := domaincontacts.ListOptions{
		MaxResults: contactsMaxResults,
		PageToken:  contactsPageToken,
	}

	result, err := repo.List(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list contacts: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderContacts(result.Items))

	return nil
}

func runContactsGet(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	resourceName := args[0]

	repo, err := getContactRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	contact, err := repo.Get(ctx, resourceName)
	if err != nil {
		return fmt.Errorf("failed to get contact: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderContact(contact))

	return nil
}

func runContactsCreate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	repo, err := getContactRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	contact := domaincontacts.NewContact()

	if contactsGivenName != "" || contactsFamilyName != "" {
		contact.Names = []domaincontacts.Name{
			{
				GivenName:  contactsGivenName,
				FamilyName: contactsFamilyName,
			},
		}
	}

	if contactsEmail != "" {
		err = contact.AddEmail(contactsEmail, contactsEmailType, true)
		if err != nil {
			return fmt.Errorf("failed to add email: %w", err)
		}
	}

	if contactsPhone != "" {
		err = contact.AddPhone(contactsPhone, contactsPhoneType, true)
		if err != nil {
			return fmt.Errorf("failed to add phone: %w", err)
		}
	}

	if contactsNotes != "" {
		contact.Biographies = []domaincontacts.Biography{
			{Value: contactsNotes},
		}
	}

	created, err := repo.Create(ctx, contact)
	if err != nil {
		return fmt.Errorf("failed to create contact: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderContact(created))

	return nil
}

func runContactsUpdate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	resourceName := args[0]

	repo, err := getContactRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	contact, err := repo.Get(ctx, resourceName)
	if err != nil {
		return fmt.Errorf("failed to get contact: %w", err)
	}

	if contactsGivenName != "" || contactsFamilyName != "" {
		if len(contact.Names) == 0 {
			contact.Names = []domaincontacts.Name{{}}
		}
		if contactsGivenName != "" {
			contact.Names[0].GivenName = contactsGivenName
		}
		if contactsFamilyName != "" {
			contact.Names[0].FamilyName = contactsFamilyName
		}
	}

	if contactsEmail != "" {
		err = contact.AddEmail(contactsEmail, contactsEmailType, true)
		if err != nil {
			return fmt.Errorf("failed to add email: %w", err)
		}
	}

	if contactsPhone != "" {
		err = contact.AddPhone(contactsPhone, contactsPhoneType, true)
		if err != nil {
			return fmt.Errorf("failed to add phone: %w", err)
		}
	}

	updateMask := []string{}
	if contactsUpdateMask != "" {
		updateMask = []string{contactsUpdateMask}
	}

	updated, err := repo.Update(ctx, contact, updateMask)
	if err != nil {
		return fmt.Errorf("failed to update contact: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderContact(updated))

	return nil
}

func runContactsDelete(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	resourceName := args[0]

	repo, err := getContactRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	err = repo.Delete(ctx, resourceName)
	if err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderSuccess(fmt.Sprintf("Contact '%s' deleted", resourceName)))

	return nil
}

func runContactsSearch(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	query := args[0]

	repo, err := getContactRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	opts := domaincontacts.SearchOptions{
		Query:      query,
		MaxResults: contactsMaxResults,
		PageToken:  contactsPageToken,
	}

	result, err := repo.Search(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to search contacts: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderContacts(result.Items))

	return nil
}

func runContactsGroups(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	repo, err := getContactGroupRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	groups, err := repo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list contact groups: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderContactGroups(groups))

	return nil
}

func runContactsGroupCreate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	name := args[0]

	repo, err := getContactGroupRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	group, err := domaincontacts.NewContactGroup(name)
	if err != nil {
		return fmt.Errorf("invalid contact group: %w", err)
	}

	created, err := repo.Create(ctx, group)
	if err != nil {
		return fmt.Errorf("failed to create contact group: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderContactGroup(created))

	return nil
}

func runContactsGroupUpdate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	resourceName := args[0]

	repo, err := getContactGroupRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	group, err := repo.Get(ctx, resourceName)
	if err != nil {
		return fmt.Errorf("failed to get contact group: %w", err)
	}

	if contactsGroupName != "" {
		group.Name = contactsGroupName
	}

	updated, err := repo.Update(ctx, group)
	if err != nil {
		return fmt.Errorf("failed to update contact group: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderContactGroup(updated))

	return nil
}

func runContactsGroupDelete(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	resourceName := args[0]

	repo, err := getContactGroupRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	err = repo.Delete(ctx, resourceName)
	if err != nil {
		return fmt.Errorf("failed to delete contact group: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderSuccess(fmt.Sprintf("Contact group '%s' deleted", resourceName)))

	return nil
}

func runContactsGroupMembers(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	resourceName := args[0]

	repo, err := getContactGroupRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	opts := domaincontacts.ListOptions{
		MaxResults: contactsMaxResults,
		PageToken:  contactsPageToken,
	}

	result, err := repo.ListMembers(ctx, resourceName, opts)
	if err != nil {
		return fmt.Errorf("failed to list group members: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderContacts(result.Items))

	return nil
}

func runContactsGroupAdd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	groupResourceName := args[0]
	contactResourceNames := args[1:]

	repo, err := getContactGroupRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	err = repo.AddMembers(ctx, groupResourceName, contactResourceNames)
	if err != nil {
		return fmt.Errorf("failed to add members to group: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderSuccess(fmt.Sprintf("Added %d contact(s) to group", len(contactResourceNames))))

	return nil
}

func runContactsGroupRemove(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	groupResourceName := args[0]
	contactResourceNames := args[1:]

	repo, err := getContactGroupRepositoryFromDeps(ctx)
	if err != nil {
		return err
	}

	err = repo.RemoveMembers(ctx, groupResourceName, contactResourceNames)
	if err != nil {
		return fmt.Errorf("failed to remove members from group: %w", err)
	}

	p := presenter.New(formatFlag)
	cmd.Println(p.RenderSuccess(fmt.Sprintf("Removed %d contact(s) from group", len(contactResourceNames))))

	return nil
}

// init registers all contacts commands.
func init() {
	rootCmd.AddCommand(contactsCmd)

	// Contact commands
	contactsCmd.AddCommand(contactsListCmd)
	contactsCmd.AddCommand(contactsGetCmd)
	contactsCmd.AddCommand(contactsCreateCmd)
	contactsCmd.AddCommand(contactsUpdateCmd)
	contactsCmd.AddCommand(contactsDeleteCmd)
	contactsCmd.AddCommand(contactsSearchCmd)

	// Group commands
	contactsCmd.AddCommand(contactsGroupsCmd)
	contactsCmd.AddCommand(contactsGroupCreateCmd)
	contactsCmd.AddCommand(contactsGroupUpdateCmd)
	contactsCmd.AddCommand(contactsGroupDeleteCmd)
	contactsCmd.AddCommand(contactsGroupMembersCmd)
	contactsCmd.AddCommand(contactsGroupAddCmd)
	contactsCmd.AddCommand(contactsGroupRemoveCmd)

	// Flags for list command
	contactsListCmd.Flags().Int64Var(&contactsMaxResults, "max-results", 100, "maximum number of contacts to return")
	contactsListCmd.Flags().StringVar(&contactsPageToken, "page-token", "", "token for pagination")

	// Flags for create command
	contactsCreateCmd.Flags().StringVar(&contactsGivenName, "given-name", "", "contact's given name (first name)")
	contactsCreateCmd.Flags().StringVar(&contactsFamilyName, "family-name", "", "contact's family name (last name)")
	contactsCreateCmd.Flags().StringVar(&contactsEmail, "email", "", "contact's email address")
	contactsCreateCmd.Flags().StringVar(&contactsEmailType, "email-type", "home", "email type (home, work, other)")
	contactsCreateCmd.Flags().StringVar(&contactsPhone, "phone", "", "contact's phone number")
	contactsCreateCmd.Flags().StringVar(&contactsPhoneType, "phone-type", "mobile", "phone type (mobile, home, work, other)")
	contactsCreateCmd.Flags().StringVar(&contactsAddress, "address", "", "contact's address")
	contactsCreateCmd.Flags().StringVar(&contactsAddressType, "address-type", "home", "address type (home, work, other)")
	contactsCreateCmd.Flags().StringVar(&contactsOrganization, "organization", "", "contact's organization")
	contactsCreateCmd.Flags().StringVar(&contactsTitle, "title", "", "contact's job title")
	contactsCreateCmd.Flags().StringVar(&contactsNotes, "notes", "", "notes about the contact")
	contactsCreateCmd.Flags().StringVar(&contactsBirthday, "birthday", "", "contact's birthday (YYYY-MM-DD)")
	contactsCreateCmd.Flags().StringVar(&contactsURL, "url", "", "contact's URL")

	// Flags for update command
	contactsUpdateCmd.Flags().StringVar(&contactsGivenName, "given-name", "", "contact's given name (first name)")
	contactsUpdateCmd.Flags().StringVar(&contactsFamilyName, "family-name", "", "contact's family name (last name)")
	contactsUpdateCmd.Flags().StringVar(&contactsEmail, "email", "", "contact's email address")
	contactsUpdateCmd.Flags().StringVar(&contactsEmailType, "email-type", "home", "email type (home, work, other)")
	contactsUpdateCmd.Flags().StringVar(&contactsPhone, "phone", "", "contact's phone number")
	contactsUpdateCmd.Flags().StringVar(&contactsPhoneType, "phone-type", "mobile", "phone type (mobile, home, work, other)")
	contactsUpdateCmd.Flags().StringVar(&contactsUpdateMask, "update-mask", "", "fields to update")

	// Flags for delete command
	contactsDeleteCmd.Flags().BoolVar(&contactsDeleteConfirm, "confirm", false, "confirm deletion")

	// Flags for search command
	contactsSearchCmd.Flags().Int64Var(&contactsMaxResults, "max-results", 100, "maximum number of results")
	contactsSearchCmd.Flags().StringVar(&contactsPageToken, "page-token", "", "token for pagination")

	// Flags for group update command
	contactsGroupUpdateCmd.Flags().StringVar(&contactsGroupName, "group-name", "", "new name for the contact group")

	// Flags for group delete command
	contactsGroupDeleteCmd.Flags().BoolVar(&contactsDeleteConfirm, "confirm", false, "confirm deletion")

	// Flags for group members command
	contactsGroupMembersCmd.Flags().Int64Var(&contactsMaxResults, "max-results", 100, "maximum number of results")
	contactsGroupMembersCmd.Flags().StringVar(&contactsPageToken, "page-token", "", "token for pagination")
}
