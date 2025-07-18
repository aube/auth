package entities

import (
	"time"
)

type File struct {
	ID          string    `json:"id"`
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	Path        string    `json:"path"`
	UploadedAt  time.Time `json:"uploaded_at"`
	Description string    `json:"description"`
}

func NewFile(uuid, name, contentType, path string, size int64, description string) *File {
	return &File{
		ID:          "",
		UUID:        uuid,
		Name:        name,
		Size:        size,
		ContentType: contentType,
		Path:        path,
		UploadedAt:  time.Now(),
		Description: description,
	}
}
