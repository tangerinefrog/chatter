package files

import "github.com/google/uuid"

type File struct {
	ID         uuid.UUID
	ChatID     uuid.UUID
	UploaderID uuid.UUID
	FileKey    string
	FileName   string
	MimeType   string
	SizeBytes  int64
	UpdatedAt  string
}
