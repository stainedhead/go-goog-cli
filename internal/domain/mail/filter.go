package mail

// FilterCriteria defines the conditions that trigger a filter.
type FilterCriteria struct {
	From    string
	To      string
	Subject string
	Query   string
}

// FilterAction defines the actions to take when a filter matches.
type FilterAction struct {
	AddLabels    []string
	RemoveLabels []string
	Forward      string
	Archive      bool
	MarkRead     bool
	Star         bool
	Trash        bool
}

// Filter represents an email filter rule.
type Filter struct {
	ID       string
	Criteria *FilterCriteria
	Action   *FilterAction
}

// NewFilter creates a new Filter with the given ID.
func NewFilter(id string) *Filter {
	return &Filter{
		ID:       id,
		Criteria: &FilterCriteria{},
		Action: &FilterAction{
			AddLabels:    []string{},
			RemoveLabels: []string{},
		},
	}
}

// NewFilterWithCriteria creates a new Filter with the given ID and criteria.
func NewFilterWithCriteria(id string, criteria *FilterCriteria) *Filter {
	return &Filter{
		ID:       id,
		Criteria: criteria,
		Action: &FilterAction{
			AddLabels:    []string{},
			RemoveLabels: []string{},
		},
	}
}

// SetCriteria sets the filter criteria.
func (f *Filter) SetCriteria(criteria *FilterCriteria) {
	f.Criteria = criteria
}

// SetAction sets the filter action.
func (f *Filter) SetAction(action *FilterAction) {
	f.Action = action
}

// HasCriteria returns true if the filter has any criteria defined.
func (f *Filter) HasCriteria() bool {
	if f.Criteria == nil {
		return false
	}
	return f.Criteria.From != "" ||
		f.Criteria.To != "" ||
		f.Criteria.Subject != "" ||
		f.Criteria.Query != ""
}

// HasAction returns true if the filter has any actions defined.
func (f *Filter) HasAction() bool {
	if f.Action == nil {
		return false
	}
	return len(f.Action.AddLabels) > 0 ||
		len(f.Action.RemoveLabels) > 0 ||
		f.Action.Forward != "" ||
		f.Action.Archive ||
		f.Action.MarkRead ||
		f.Action.Star ||
		f.Action.Trash
}

// IsValid returns true if the filter has both criteria and actions.
func (f *Filter) IsValid() bool {
	return f.HasCriteria() && f.HasAction()
}

// AddLabelToAction adds a label to the filter's add labels action.
func (f *Filter) AddLabelToAction(label string) {
	if f.Action == nil {
		f.Action = &FilterAction{
			AddLabels:    []string{},
			RemoveLabels: []string{},
		}
	}
	f.Action.AddLabels = append(f.Action.AddLabels, label)
}

// AddRemoveLabelToAction adds a label to the filter's remove labels action.
func (f *Filter) AddRemoveLabelToAction(label string) {
	if f.Action == nil {
		f.Action = &FilterAction{
			AddLabels:    []string{},
			RemoveLabels: []string{},
		}
	}
	f.Action.RemoveLabels = append(f.Action.RemoveLabels, label)
}
