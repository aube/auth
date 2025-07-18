package upload

import "github.com/aube/auth/internal/domain/entities"

type RegisterUploadDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func NewUploadedFile(file *entities.File) RegisterUploadDTO {
	return RegisterUploadDTO{
		ID:   file.ID,
		Name: file.Name,
		Size: file.Size,
	}
}
