// Package file provides business logic for file operations.
package file

import (
	"context"
	"io"

	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// FileService implements core file management functionality.
// Wraps repository operations with business logic and logging.
//
// Fields:
//   - repo: Underlying storage repository
//   - log: Structured logger instance

type FileService struct {
	repo FileRepository
	log  zerolog.Logger
}

// NewFileService creates a new FileService instance.
// repo: Storage repository implementation
// Returns: Configured *FileService
func NewFileService(repo FileRepository) *FileService {
	return &FileService{
		repo: repo,
		log:  logger.Get().With().Str("file", "service").Logger(),
	}
}

// Upload handles file upload business logic:
// 1. Generates a new UUID for the file
// 2. Creates file metadata entity
// 3. Delegates storage to repository
//
// ctx: Context for cancellation/timeout
// size: File size in bytes
// data: File content stream
// Returns: (*entities.File, error) - created file metadata
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

// Delete proxies file deletion to the repository.
// ctx: Context for cancellation/timeout
// id: File UUID to delete
// Returns: error on failure
func (s *FileService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// Download retrieves file content via repository.
// ctx: Context for cancellation/timeout
// uuid: File identifier
// Returns: (io.ReadCloser, error) - caller must close the stream
func (s *FileService) Download(ctx context.Context, uuid string) (io.ReadCloser, error) {
	return s.repo.GetFileContent(ctx, uuid)
}

// generateFileUUID creates a new UUID string for file identification.
// Returns: UUID string
// Note: Unexported as it's an internal implementation detail
func generateFileUUID() string {
	return uuid.New().String()
}
