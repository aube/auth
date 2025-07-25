package upload

import (
	"context"
	"time"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"
)

type UploadService struct {
	repo UploadRepository
	log  zerolog.Logger
}

func NewUploadService(repo UploadRepository) *UploadService {
	return &UploadService{
		repo: repo,
		log:  logger.Get().With().Str("upload", "service").Logger(),
	}
}

func (s *UploadService) RegisterUploadedFile(
	ctx context.Context,
	userID int64,
	file *entities.File,
	name,
	category,
	contentType,
	description string,
) (*entities.Upload, error) {

	upload := entities.NewUpload(file, 0, userID, name, category, contentType, description, time.Now())

	err := s.repo.Create(ctx, userID, upload)
	if err != nil {
		s.log.Debug().Err(err).Msg("RegisterUploadedFile")
		return nil, err
	}

	return upload, nil
}

func (s *UploadService) ListByUserID(ctx context.Context, userID int64, offset, limit int) (*entities.Uploads, *dto.Pagination, error) {
	return s.repo.ListByUserID(ctx, userID, offset, limit)
}

func (s *UploadService) GetByUUID(ctx context.Context, uuid string, userID int64) (*entities.Upload, error) {
	return s.repo.GetByUUID(ctx, uuid, userID)
}

func (s *UploadService) GetByName(ctx context.Context, name string, userID int64) (*entities.Upload, error) {
	return s.repo.GetByName(ctx, name, userID)
}

func (s *UploadService) Delete(ctx context.Context, uuid string, userID int64) error {
	return s.repo.Delete(ctx, uuid, userID)
}

func (s *UploadService) DeleteForce(ctx context.Context, uuid string, userID int64) error {
	return s.repo.DeleteForce(ctx, uuid, userID)
}

func (s *UploadService) CanBeDeleted(ctx context.Context, uuid string, userID int64) error {
	_, err := s.repo.GetByUUID(ctx, uuid, userID)

	if err != nil {
		s.log.Debug().Err(err).Msg("Delete")
		return err
	}

	return nil
}

// func (s *UploadService) Create(ctx context.Context, file *entities.File, description string) (*entities.Upload, error) {
// 	upload := entities.NewUpload(file, description)
// 	return s.repo.Create(ctx, upload)
// }
