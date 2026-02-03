package mail

// LabelType represents the type of a label.
const (
	LabelTypeSystem = "system"
	LabelTypeUser   = "user"
)

// Label visibility constants.
const (
	LabelVisibilityShow              = "show"
	LabelVisibilityHide              = "hide"
	LabelVisibilityShowIfUnread      = "showIfUnread"
	LabelVisibilityLabelShow         = "labelShow"
	LabelVisibilityLabelShowIfUnread = "labelShowIfUnread"
	LabelVisibilityLabelHide         = "labelHide"
)

// LabelColor represents the color configuration for a label.
type LabelColor struct {
	Background string
	Text       string
}

// Label represents an email label (folder/category).
type Label struct {
	ID                    string
	Name                  string
	Type                  string // system or user
	MessageListVisibility string
	LabelListVisibility   string
	Color                 *LabelColor
}

// NewLabel creates a new user Label with the given ID and name.
func NewLabel(id, name string) *Label {
	return &Label{
		ID:   id,
		Name: name,
		Type: LabelTypeUser,
	}
}

// NewSystemLabel creates a new system Label with the given ID and name.
func NewSystemLabel(id, name string) *Label {
	return &Label{
		ID:   id,
		Name: name,
		Type: LabelTypeSystem,
	}
}

// IsSystemLabel returns true if this is a system-defined label.
func (l *Label) IsSystemLabel() bool {
	return l.Type == LabelTypeSystem
}

// IsUserLabel returns true if this is a user-defined label.
func (l *Label) IsUserLabel() bool {
	return l.Type == LabelTypeUser
}

// SetColor sets the label's color.
func (l *Label) SetColor(background, text string) {
	l.Color = &LabelColor{
		Background: background,
		Text:       text,
	}
}

// ClearColor removes the label's color.
func (l *Label) ClearColor() {
	l.Color = nil
}

// HasColor returns true if the label has a color set.
func (l *Label) HasColor() bool {
	return l.Color != nil
}
