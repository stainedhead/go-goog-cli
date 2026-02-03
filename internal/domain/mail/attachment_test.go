package mail

import "testing"

func TestNewAttachment(t *testing.T) {
	att := NewAttachment("att-123", "document.pdf", "application/pdf")

	if att.ID != "att-123" {
		t.Errorf("expected ID 'att-123', got '%s'", att.ID)
	}
	if att.Filename != "document.pdf" {
		t.Errorf("expected Filename 'document.pdf', got '%s'", att.Filename)
	}
	if att.MimeType != "application/pdf" {
		t.Errorf("expected MimeType 'application/pdf', got '%s'", att.MimeType)
	}
	if att.Size != 0 {
		t.Errorf("expected Size 0, got %d", att.Size)
	}
	if att.Data != nil {
		t.Error("expected Data to be nil for new attachment")
	}
}

func TestAttachment_SetData(t *testing.T) {
	att := NewAttachment("att-123", "document.pdf", "application/pdf")
	data := []byte("test file content")

	att.SetData(data)

	if att.Size != int64(len(data)) {
		t.Errorf("expected Size %d, got %d", len(data), att.Size)
	}
	if string(att.Data) != "test file content" {
		t.Errorf("expected Data 'test file content', got '%s'", string(att.Data))
	}
}

func TestAttachment_HasData(t *testing.T) {
	att := NewAttachment("att-123", "document.pdf", "application/pdf")

	if att.HasData() {
		t.Error("expected HasData to return false for new attachment")
	}

	att.SetData([]byte("content"))
	if !att.HasData() {
		t.Error("expected HasData to return true after SetData")
	}

	att.SetData([]byte{})
	if att.HasData() {
		t.Error("expected HasData to return false for empty data")
	}
}

func TestAttachment_ClearData(t *testing.T) {
	att := NewAttachment("att-123", "document.pdf", "application/pdf")
	att.SetData([]byte("content"))

	if !att.HasData() {
		t.Error("expected attachment to have data")
	}

	att.ClearData()

	if att.Data != nil {
		t.Error("expected Data to be nil after ClearData")
	}
	// Size is intentionally not cleared to preserve metadata
}

func TestAttachment_IsImage(t *testing.T) {
	testCases := []struct {
		mimeType string
		expected bool
	}{
		{"image/jpeg", true},
		{"image/png", true},
		{"image/gif", true},
		{"image/webp", true},
		{"image/bmp", true},
		{"image/svg+xml", true},
		{"application/pdf", false},
		{"text/plain", false},
		{"application/octet-stream", false},
	}

	for _, tc := range testCases {
		t.Run(tc.mimeType, func(t *testing.T) {
			att := NewAttachment("att-123", "file", tc.mimeType)
			if att.IsImage() != tc.expected {
				t.Errorf("IsImage() for %s: expected %v, got %v", tc.mimeType, tc.expected, att.IsImage())
			}
		})
	}
}

func TestAttachment_IsPDF(t *testing.T) {
	pdfAtt := NewAttachment("att-123", "document.pdf", "application/pdf")
	if !pdfAtt.IsPDF() {
		t.Error("expected IsPDF to return true for application/pdf")
	}

	imageAtt := NewAttachment("att-456", "image.png", "image/png")
	if imageAtt.IsPDF() {
		t.Error("expected IsPDF to return false for image/png")
	}
}
