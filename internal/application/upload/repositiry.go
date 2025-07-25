package upload

import (
	"context"
	"errors"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
)

var ErrFileNotFound = errors.New("file not found")

type UploadRepository interface {
	Create(ctx context.Context, userID int64, upload *entities.Upload) error
	ListByUserID(ctx context.Context, userID int64, offset, limit int) (*entities.Uploads, *dto.Pagination, error)
	GetByUUID(ctx context.Context, uuid string, userID int64) (*entities.Upload, error)
	GetByName(ctx context.Context, name string, userID int64) (*entities.Upload, error)
	Delete(ctx context.Context, uuid string, userID int64) error
	DeleteForce(ctx context.Context, uuid string, userID int64) error
}
