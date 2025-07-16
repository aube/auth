package upload

import (
	"context"
	"errors"

	"github.com/aube/auth/internal/domain/entities"
)

var ErrFileNotFound = errors.New("file not found")

type UploadRepository interface {
	Create(ctx context.Context, file *entities.File) error
	ListByUserID(ctx context.Context, id string) (*[]entities.File, error)
	Delete(ctx context.Context, uuid string) error
}
