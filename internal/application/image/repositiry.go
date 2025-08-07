package image

import (
	"context"
	"errors"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
)

var ErrFileNotFound = errors.New("file not found")

type ImageRepository interface {
	Create(ctx context.Context, userID int64, image *entities.Image) error
	ListByUserID(ctx context.Context, userID int64, offset, limit int, params map[string]any) (*entities.Images, *dto.Pagination, error)
	GetByUUID(ctx context.Context, uuid string, userID int64) (*entities.Image, error)
	GetByName(ctx context.Context, name string, userID int64) (*entities.Image, error)
	Delete(ctx context.Context, uuid string, userID int64) error
	DeleteForce(ctx context.Context, uuid string, userID int64) error
}
