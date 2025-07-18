package upload

import (
	"context"
	"errors"

	"github.com/aube/auth/internal/domain/entities"
)

var ErrFileNotFound = errors.New("file not found")

type UploadRepository interface {
	Create(ctx context.Context, userID string, file *entities.File) error
	ListByUserID(ctx context.Context, userID string) (*[]entities.Upload, error)
	Delete(ctx context.Context, uuid string) error
}
