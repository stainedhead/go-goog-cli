package calendar

import "testing"

func TestNewCalendar(t *testing.T) {
	title := "Work Calendar"
	cal := NewCalendar(title)

	if cal.Title != title {
		t.Errorf("expected title %q, got %q", title, cal.Title)
	}

	if cal.AccessRole != AccessRoleOwner {
		t.Errorf("expected access role %q, got %q", AccessRoleOwner, cal.AccessRole)
	}

	if cal.Primary {
		t.Error("expected Primary to be false by default")
	}

	if cal.Selected {
		t.Error("expected Selected to be false by default")
	}
}

func TestIsValidAccessRole(t *testing.T) {
	tests := []struct {
		role  string
		valid bool
	}{
		{AccessRoleOwner, true},
		{AccessRoleWriter, true},
		{AccessRoleReader, true},
		{AccessRoleFreeBusyReader, true},
		{"owner", true},
		{"writer", true},
		{"reader", true},
		{"freeBusyReader", true},
		{"admin", false},
		{"", false},
		{"OWNER", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			got := IsValidAccessRole(tt.role)
			if got != tt.valid {
				t.Errorf("IsValidAccessRole(%q) = %v, want %v", tt.role, got, tt.valid)
			}
		})
	}
}

func TestCalendarCanWrite(t *testing.T) {
	tests := []struct {
		role     string
		canWrite bool
	}{
		{AccessRoleOwner, true},
		{AccessRoleWriter, true},
		{AccessRoleReader, false},
		{AccessRoleFreeBusyReader, false},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			cal := &Calendar{AccessRole: tt.role}
			got := cal.CanWrite()
			if got != tt.canWrite {
				t.Errorf("Calendar{AccessRole: %q}.CanWrite() = %v, want %v", tt.role, got, tt.canWrite)
			}
		})
	}
}

func TestCalendarCanRead(t *testing.T) {
	tests := []struct {
		role    string
		canRead bool
	}{
		{AccessRoleOwner, true},
		{AccessRoleWriter, true},
		{AccessRoleReader, true},
		{AccessRoleFreeBusyReader, false},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			cal := &Calendar{AccessRole: tt.role}
			got := cal.CanRead()
			if got != tt.canRead {
				t.Errorf("Calendar{AccessRole: %q}.CanRead() = %v, want %v", tt.role, got, tt.canRead)
			}
		})
	}
}

func TestCalendarIsOwner(t *testing.T) {
	tests := []struct {
		role    string
		isOwner bool
	}{
		{AccessRoleOwner, true},
		{AccessRoleWriter, false},
		{AccessRoleReader, false},
		{AccessRoleFreeBusyReader, false},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			cal := &Calendar{AccessRole: tt.role}
			got := cal.IsOwner()
			if got != tt.isOwner {
				t.Errorf("Calendar{AccessRole: %q}.IsOwner() = %v, want %v", tt.role, got, tt.isOwner)
			}
		})
	}
}

func TestCalendarFields(t *testing.T) {
	cal := &Calendar{
		ID:          "cal123",
		Title:       "Personal",
		Description: "My personal calendar",
		TimeZone:    "America/New_York",
		ColorID:     "9",
		Primary:     true,
		Selected:    true,
		AccessRole:  AccessRoleOwner,
	}

	if cal.ID != "cal123" {
		t.Errorf("unexpected ID: %s", cal.ID)
	}

	if cal.Description != "My personal calendar" {
		t.Errorf("unexpected Description: %s", cal.Description)
	}

	if cal.TimeZone != "America/New_York" {
		t.Errorf("unexpected TimeZone: %s", cal.TimeZone)
	}

	if cal.ColorID != "9" {
		t.Errorf("unexpected ColorID: %s", cal.ColorID)
	}

	if !cal.Primary {
		t.Error("expected Primary to be true")
	}

	if !cal.Selected {
		t.Error("expected Selected to be true")
	}
}
