package mail

import "testing"

func TestNewLabel(t *testing.T) {
	label := NewLabel("label-123", "My Label")

	if label.ID != "label-123" {
		t.Errorf("expected ID 'label-123', got '%s'", label.ID)
	}
	if label.Name != "My Label" {
		t.Errorf("expected Name 'My Label', got '%s'", label.Name)
	}
	if label.Type != LabelTypeUser {
		t.Errorf("expected Type '%s', got '%s'", LabelTypeUser, label.Type)
	}
}

func TestNewSystemLabel(t *testing.T) {
	label := NewSystemLabel("INBOX", "Inbox")

	if label.ID != "INBOX" {
		t.Errorf("expected ID 'INBOX', got '%s'", label.ID)
	}
	if label.Name != "Inbox" {
		t.Errorf("expected Name 'Inbox', got '%s'", label.Name)
	}
	if label.Type != LabelTypeSystem {
		t.Errorf("expected Type '%s', got '%s'", LabelTypeSystem, label.Type)
	}
}

func TestLabel_IsSystemLabel(t *testing.T) {
	userLabel := NewLabel("label-123", "My Label")
	systemLabel := NewSystemLabel("INBOX", "Inbox")

	if userLabel.IsSystemLabel() {
		t.Error("expected user label not to be system label")
	}
	if !systemLabel.IsSystemLabel() {
		t.Error("expected system label to be system label")
	}
}

func TestLabel_IsUserLabel(t *testing.T) {
	userLabel := NewLabel("label-123", "My Label")
	systemLabel := NewSystemLabel("INBOX", "Inbox")

	if !userLabel.IsUserLabel() {
		t.Error("expected user label to be user label")
	}
	if systemLabel.IsUserLabel() {
		t.Error("expected system label not to be user label")
	}
}

func TestLabel_SetColor(t *testing.T) {
	label := NewLabel("label-123", "My Label")

	if label.HasColor() {
		t.Error("expected new label not to have color")
	}

	label.SetColor("#ff0000", "#ffffff")

	if !label.HasColor() {
		t.Error("expected label to have color after SetColor")
	}
	if label.Color.Background != "#ff0000" {
		t.Errorf("expected Background '#ff0000', got '%s'", label.Color.Background)
	}
	if label.Color.Text != "#ffffff" {
		t.Errorf("expected Text '#ffffff', got '%s'", label.Color.Text)
	}
}

func TestLabel_ClearColor(t *testing.T) {
	label := NewLabel("label-123", "My Label")
	label.SetColor("#ff0000", "#ffffff")

	if !label.HasColor() {
		t.Error("expected label to have color")
	}

	label.ClearColor()

	if label.HasColor() {
		t.Error("expected label not to have color after ClearColor")
	}
	if label.Color != nil {
		t.Error("expected Color to be nil after ClearColor")
	}
}

func TestLabel_HasColor(t *testing.T) {
	label := NewLabel("label-123", "My Label")

	if label.HasColor() {
		t.Error("expected HasColor to return false for new label")
	}

	label.SetColor("#ff0000", "#ffffff")
	if !label.HasColor() {
		t.Error("expected HasColor to return true after SetColor")
	}

	label.ClearColor()
	if label.HasColor() {
		t.Error("expected HasColor to return false after ClearColor")
	}
}

func TestLabelColor(t *testing.T) {
	color := &LabelColor{
		Background: "#0000ff",
		Text:       "#ffffff",
	}

	if color.Background != "#0000ff" {
		t.Errorf("expected Background '#0000ff', got '%s'", color.Background)
	}
	if color.Text != "#ffffff" {
		t.Errorf("expected Text '#ffffff', got '%s'", color.Text)
	}
}
