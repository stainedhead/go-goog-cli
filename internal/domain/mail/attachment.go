package mail

// Attachment represents an email attachment.
type Attachment struct {
	ID       string
	Filename string
	MimeType string
	Size     int64
	Data     []byte
}

// NewAttachment creates a new Attachment with the given parameters.
func NewAttachment(id, filename, mimeType string) *Attachment {
	return &Attachment{
		ID:       id,
		Filename: filename,
		MimeType: mimeType,
	}
}

// SetData sets the attachment data and updates the size.
func (a *Attachment) SetData(data []byte) {
	a.Data = data
	a.Size = int64(len(data))
}

// HasData returns true if the attachment has data loaded.
func (a *Attachment) HasData() bool {
	return len(a.Data) > 0
}

// ClearData clears the attachment data to free memory.
func (a *Attachment) ClearData() {
	a.Data = nil
}

// IsImage returns true if the attachment appears to be an image.
func (a *Attachment) IsImage() bool {
	switch a.MimeType {
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/bmp", "image/svg+xml":
		return true
	default:
		return false
	}
}

// IsPDF returns true if the attachment is a PDF document.
func (a *Attachment) IsPDF() bool {
	return a.MimeType == "application/pdf"
}
