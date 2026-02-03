package presenter

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		wantType string
	}{
		{
			name:     "json format returns JSONPresenter",
			format:   FormatJSON,
			wantType: "*presenter.JSONPresenter",
		},
		{
			name:     "table format returns TablePresenter",
			format:   FormatTable,
			wantType: "*presenter.TablePresenter",
		},
		{
			name:     "plain format returns PlainPresenter",
			format:   FormatPlain,
			wantType: "*presenter.PlainPresenter",
		},
		{
			name:     "unknown format returns TablePresenter as default",
			format:   "unknown",
			wantType: "*presenter.TablePresenter",
		},
		{
			name:     "empty format returns TablePresenter as default",
			format:   "",
			wantType: "*presenter.TablePresenter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.format)
			gotType := getTypeName(got)
			if gotType != tt.wantType {
				t.Errorf("New(%q) = %s, want %s", tt.format, gotType, tt.wantType)
			}
		})
	}
}

func getTypeName(p Presenter) string {
	switch p.(type) {
	case *JSONPresenter:
		return "*presenter.JSONPresenter"
	case *TablePresenter:
		return "*presenter.TablePresenter"
	case *PlainPresenter:
		return "*presenter.PlainPresenter"
	default:
		return "unknown"
	}
}
