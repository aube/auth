package entities

import "time"

type Dataset struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	Path        string    `json:"path"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

func NewDataset(id, name, contentType, path string, size int64) *Dataset {
	return &Dataset{
		ID:          id,
		Name:        name,
		Size:        size,
		ContentType: contentType,
		Path:        path,
		UploadedAt:  time.Now(),
	}
}
