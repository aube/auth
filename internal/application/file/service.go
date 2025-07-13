package file

import (
	"context"
	"io"

	"github.com/aube/auth/internal/domain/entities"
	"github.com/google/uuid"
)

// FileService предоставляет бизнес-логику для работы с файлами
type FileService struct {
	repo FileRepository
}

// NewFileService создает новый экземпляр FileService
func NewFileService(repo FileRepository) *FileService {
	return &FileService{repo: repo}
}

// Upload загружает новый файл
func (s *FileService) Upload(ctx context.Context, name, contentType string, size int64, data io.Reader) (*entities.File, error) {
	file := entities.NewFile(
		generateFileID(),
		name,
		contentType,
		"", // Путь будет установлен в репозитории
		size,
	)

	if err := s.repo.Save(ctx, file, data); err != nil {
		return nil, err
	}

	return file, nil
}

// GetByID возвращает файл по ID
func (s *FileService) GetByID(ctx context.Context, id string) (*entities.File, error) {
	return s.repo.FindByID(ctx, id)
}

// List возвращает список всех файлов
func (s *FileService) List(ctx context.Context) ([]*entities.File, error) {
	return s.repo.FindAll(ctx)
}

// Delete удаляет файл
func (s *FileService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// Download возвращает содержимое файла
func (s *FileService) Download(ctx context.Context, file *entities.File) (io.ReadCloser, error) {
	return s.repo.GetFileContent(ctx, file)
}

// generateFileID генерирует уникальный ID для файла
func generateFileID() string {
	return uuid.New().String()
}
