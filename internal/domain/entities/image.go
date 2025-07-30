package entities

import "time"

type Image struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	Path        string    `json:"path"`
	ImageedAt   time.Time `json:"uploaded_at"`
	Description string    `json:"description"`
}

type Images []Image

func NewImage(file *File, id int64, userID int64, name string, category string, contentType string, description string, createdAt time.Time) *Image {
	return &Image{
		ID:          id,
		UserID:      userID,
		UUID:        file.Name, // is UUID on server filesysten
		Name:        name,      // is original name in database
		Category:    category,
		Size:        file.Size,
		ContentType: contentType,
		Path:        file.Path,
		ImageedAt:   createdAt,
		Description: description,
	}
}
