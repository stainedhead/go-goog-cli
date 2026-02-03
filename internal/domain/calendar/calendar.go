package calendar

// Calendar represents a Google Calendar.
type Calendar struct {
	// ID is the unique identifier for the calendar.
	ID string
	// Title is the calendar's title/summary.
	Title string
	// Description is the calendar's description.
	Description string
	// TimeZone is the calendar's time zone (e.g., "America/New_York").
	TimeZone string
	// ColorID is the color ID for the calendar.
	ColorID string
	// Primary indicates whether this is the user's primary calendar.
	Primary bool
	// Selected indicates whether the calendar is selected in the UI.
	Selected bool
	// AccessRole is the user's access role: owner, writer, reader, freeBusyReader.
	AccessRole string
}

// Access role constants.
const (
	AccessRoleOwner          = "owner"
	AccessRoleWriter         = "writer"
	AccessRoleReader         = "reader"
	AccessRoleFreeBusyReader = "freeBusyReader"
)

// NewCalendar creates a new Calendar with the given title.
func NewCalendar(title string) *Calendar {
	return &Calendar{
		Title:      title,
		AccessRole: AccessRoleOwner,
	}
}

// IsValidAccessRole checks if the given role is a valid access role.
func IsValidAccessRole(role string) bool {
	switch role {
	case AccessRoleOwner, AccessRoleWriter, AccessRoleReader, AccessRoleFreeBusyReader:
		return true
	default:
		return false
	}
}

// CanWrite returns true if the access role allows writing to the calendar.
func (c *Calendar) CanWrite() bool {
	return c.AccessRole == AccessRoleOwner || c.AccessRole == AccessRoleWriter
}

// CanRead returns true if the access role allows reading the calendar.
func (c *Calendar) CanRead() bool {
	switch c.AccessRole {
	case AccessRoleOwner, AccessRoleWriter, AccessRoleReader:
		return true
	default:
		return false
	}
}

// IsOwner returns true if the user owns this calendar.
func (c *Calendar) IsOwner() bool {
	return c.AccessRole == AccessRoleOwner
}
