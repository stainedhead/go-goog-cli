package contacts

import (
	"testing"
)

func TestNewContactGroup(t *testing.T) {
	tests := []struct {
		name       string
		groupName  string
		wantErr    bool
		errMessage string
	}{
		{
			name:      "creates group successfully",
			groupName: "Friends",
			wantErr:   false,
		},
		{
			name:       "rejects empty name",
			groupName:  "",
			wantErr:    true,
			errMessage: "group name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group, err := NewContactGroup(tt.groupName)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewContactGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMessage {
				t.Errorf("NewContactGroup() error message = %v, want %v", err.Error(), tt.errMessage)
			}

			if !tt.wantErr {
				if group == nil {
					t.Error("NewContactGroup() returned nil group")
				}
				if group.Name != tt.groupName {
					t.Errorf("NewContactGroup() name = %v, want %v", group.Name, tt.groupName)
				}
			}
		})
	}
}

func TestContactGroup_IsSystemGroup(t *testing.T) {
	tests := []struct {
		name      string
		groupType string
		want      bool
	}{
		{
			name:      "system group",
			groupType: GroupTypeSystem,
			want:      true,
		},
		{
			name:      "user group",
			groupType: GroupTypeUserContactGroup,
			want:      false,
		},
		{
			name:      "empty type is not system",
			groupType: "",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &ContactGroup{GroupType: tt.groupType}
			got := group.IsSystemGroup()

			if got != tt.want {
				t.Errorf("IsSystemGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactGroup_CanModify(t *testing.T) {
	tests := []struct {
		name      string
		groupType string
		want      bool
	}{
		{
			name:      "user group can be modified",
			groupType: GroupTypeUserContactGroup,
			want:      true,
		},
		{
			name:      "system group cannot be modified",
			groupType: GroupTypeSystem,
			want:      false,
		},
		{
			name:      "empty type can be modified",
			groupType: "",
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &ContactGroup{GroupType: tt.groupType}
			got := group.CanModify()

			if got != tt.want {
				t.Errorf("CanModify() = %v, want %v", got, tt.want)
			}
		})
	}
}
