// Package entities defines the core domain models for the application.
package entities

import "time"

// Upload represents metadata about a user's file upload.
// Fields:
//   - ID: Database primary key
//   - UserID: Owner of the upload
//   - UUID: Server-side file identifier
//   - Name: Original client filename
//   - Category: User-defined classification
//   - Size: File size in bytes
//   - ContentType: MIME type
//   - Path: Physical storage location
//   - UploadedAt: Creation timestamp
//   - Description: User-provided description
//
// JSON tags support serialization for API responses.
type Upload struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	Path        string    `json:"path"`
	UploadedAt  time.Time `json:"uploaded_at"`
	Description string    `json:"description"`
}

// Uploads is a collection type for multiple Upload entities.
type Uploads []Upload

// NewUpload constructs an Upload from components.
// file: Associated File entity
// id: Database ID (0 for new records)
// userID: Owner identifier
// name: Original filename
// category: User classification
// contentType: MIME type
// description: User notes
// createdAt: Upload timestamp
// Returns: *Upload instance
// Note: Distinguishes between server UUID (File.Name) and original name
func NewUpload(file *File, id int64, userID int64, name string, category string, contentType string, description string, createdAt time.Time) *Upload {
	return &Upload{
		ID:          id,
		UserID:      userID,
		UUID:        file.Name, // is UUID on server filesysten
		Name:        name,      // is original name in database
		Category:    category,
		Size:        file.Size,
		ContentType: contentType,
		Path:        file.Path,
		UploadedAt:  createdAt,
		Description: description,
	}
}
