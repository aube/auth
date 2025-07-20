package file

import (
	"context"
	"io"

	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type FileService struct {
	repo FileRepository
	log  zerolog.Logger
}

func NewFileService(repo FileRepository) *FileService {
	return &FileService{
		repo: repo,
		log:  logger.Get().With().Str("file", "service").Logger(),
	}
}

func (s *FileService) Upload(ctx context.Context, size int64, data io.Reader) (*entities.File, error) {
	file := entities.NewFile(
		generateFileUUID(),
		"", // Путь будет установлен в репозитории
		size,
	)
	if err := s.repo.Save(ctx, file, data); err != nil {
		return nil, err
	}

	return file, nil
}

func (s *FileService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *FileService) Download(ctx context.Context, uuid string) (io.ReadCloser, error) {
	return s.repo.GetFileContent(ctx, uuid)
}

func generateFileUUID() string {
	return uuid.New().String()
}
