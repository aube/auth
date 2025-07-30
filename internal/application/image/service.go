package image

import (
	"context"
	"time"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"
)

type ImageService struct {
	repo ImageRepository
	log  zerolog.Logger
}

func NewImageService(repo ImageRepository) *ImageService {
	return &ImageService{
		repo: repo,
		log:  logger.Get().With().Str("image", "service").Logger(),
	}
}

func (s *ImageService) RegisterImageedFile(
	ctx context.Context,
	userID int64,
	file *entities.File,
	name,
	category,
	contentType,
	description string,
) (*entities.Image, error) {

	image := entities.NewImage(file, 0, userID, name, category, contentType, description, time.Now())

	err := s.repo.Create(ctx, userID, image)
	if err != nil {
		s.log.Debug().Err(err).Msg("RegisterUploadedImage")
		return nil, err
	}

	return image, nil
}

func (s *ImageService) ListByUserID(ctx context.Context, userID int64, offset, limit int) (*entities.Images, *dto.Pagination, error) {
	return s.repo.ListByUserID(ctx, userID, offset, limit)
}

func (s *ImageService) GetByUUID(ctx context.Context, uuid string, userID int64) (*entities.Image, error) {
	return s.repo.GetByUUID(ctx, uuid, userID)
}

func (s *ImageService) GetByName(ctx context.Context, name string, userID int64) (*entities.Image, error) {
	return s.repo.GetByName(ctx, name, userID)
}

func (s *ImageService) Delete(ctx context.Context, uuid string, userID int64) error {
	return s.repo.Delete(ctx, uuid, userID)
}

func (s *ImageService) DeleteForce(ctx context.Context, uuid string, userID int64) error {
	return s.repo.DeleteForce(ctx, uuid, userID)
}

func (s *ImageService) CanBeDeleted(ctx context.Context, uuid string, userID int64) error {
	_, err := s.repo.GetByUUID(ctx, uuid, userID)

	if err != nil {
		s.log.Debug().Err(err).Msg("Delete")
		return err
	}

	return nil
}
