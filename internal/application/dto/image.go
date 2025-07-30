package dto

import (
	"github.com/aube/auth/internal/domain/entities"
)

type ImageResponse struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	Description string `json:"description"`
}

func NewImageResponse(upload *entities.Image) ImageResponse {
	return ImageResponse{
		UUID:        upload.UUID,
		Name:        upload.Name,
		Category:    upload.Category,
		Size:        upload.Size,
		ContentType: upload.ContentType,
		Description: upload.Description,
	}
}
