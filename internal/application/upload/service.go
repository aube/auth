package upload

import (
	"context"

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

func (s *UploadService) RegisterUploadedFile(ctx context.Context, userID string, uploadedFile *entities.File) error {
	s.repo.Create(ctx, userID, uploadedFile)
	return nil
}

func (s *UploadService) ListByUserID(ctx context.Context, id string) (*[]entities.Upload, error) {
	return s.repo.ListByUserID(ctx, id)
}

// func (s *UploadService) Create(ctx context.Context, file *entities.File, description string) (*entities.Upload, error) {
// 	upload := entities.NewUpload(file, description)
// 	return s.repo.Create(ctx, upload)
// }
