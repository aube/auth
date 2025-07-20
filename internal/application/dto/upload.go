package dto

import (
	"github.com/aube/auth/internal/domain/entities"
)

type UploadResponse struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	Description string `json:"description"`
}

func NewUploadResponse(upload *entities.Upload) UploadResponse {
	return UploadResponse{
		UUID:        upload.UUID,
		Name:        upload.Name,
		Size:        upload.Size,
		ContentType: upload.ContentType,
		Description: upload.Description,
	}
}
