package upload

import (
	"context"
	"errors"

	"github.com/aube/auth/internal/domain/entities"
)

var ErrFileNotFound = errors.New("file not found")

type UploadRepository interface {
	Create(ctx context.Context, userID int64, upload *entities.Upload) error
	ListByUserID(ctx context.Context, userID int64) (*entities.Uploads, error)
	GetByUUID(ctx context.Context, uuid string, userID int64) (*entities.Upload, error)
}
