// Package dto contains data transfer objects for file upload operations.
package dto

import (
	"github.com/aube/auth/internal/domain/entities"
)

// UploadResponse represents the response structure for file upload operations.
// Fields:
//   - UUID: Unique identifier of the uploaded file.
//   - Name: Original filename.
//   - Category: User-defined file category.
//   - Size: File size in bytes.
//   - ContentType: MIME type of the file.
//   - Description: User-provided file description.
type UploadResponse struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	Description string `json:"description"`
}

// NewUploadResponse creates an UploadResponse from an entities.Upload.
// upload: Source entity containing file metadata.
// Returns: Populated UploadResponse DTO.
func NewUploadResponse(upload *entities.Upload) UploadResponse {
	return UploadResponse{
		UUID:        upload.UUID,
		Name:        upload.Name,
		Category:    upload.Category,
		Size:        upload.Size,
		ContentType: upload.ContentType,
		Description: upload.Description,
	}
}
