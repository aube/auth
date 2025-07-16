package dto

import "github.com/aube/auth/internal/domain/entities"

type FileResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewFileResponse(file *entities.File) FileResponse {
	return FileResponse{
		ID:   file.ID,
		Name: file.Name,
		Size: file.Size,
	}
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error: err.Error(),
	}
}
