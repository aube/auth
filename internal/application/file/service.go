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

func (s *FileService) Upload(ctx context.Context, name, contentType string, size int64, data io.Reader, description string) (*entities.File, error) {
	file := entities.NewFile(
		generateFileID(),
		name,
		contentType,
		"", // Путь будет установлен в репозитории
		size,
		description,
	)
	if err := s.repo.Save(ctx, file, data); err != nil {
		return nil, err
	}

	return file, nil
}

func (s *FileService) GetByID(ctx context.Context, id string) (*entities.File, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *FileService) List(ctx context.Context) ([]*entities.File, error) {
	return s.repo.FindAll(ctx)
}

func (s *FileService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *FileService) Download(ctx context.Context, file *entities.File) (io.ReadCloser, error) {
	return s.repo.GetFileContent(ctx, file)
}

func generateFileID() string {
	return uuid.New().String()
}
