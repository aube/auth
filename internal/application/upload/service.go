// Package upload provides business logic for upload metadata operations.
package upload

import (
	"context"
	"time"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"
)

// UploadService implements core upload management functionality.
// Wraps repository operations with business logic, validation and logging.
// Fields:
//   - repo: Underlying metadata repository
//   - log: Structured logger instance
//
// NewUploadService creates a new UploadService instance.
// repo: Metadata repository implementation
// Returns: Configured *UploadService
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

// RegisterUploadedFile creates new upload metadata record:
// 1. Constructs upload entity from components
// 2. Persists via repository
// 3. Returns created metadata
//
// ctx: Context for cancellation/timeout
// userID: Upload owner
// file: Associated file entity
// name: Original filename
// category: User classification
// contentType: MIME type
// description: User description
// Returns: (*entities.Upload, error)
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

// ListByUserID retrieves paginated uploads for a user:
// Proxies to repository with pagination parameters
//
// ctx: Context for cancellation/timeout
// userID: Owner filter
// offset: Pagination start
// limit: Maximum results
// Returns: (*entities.Uploads, *dto.Pagination, error)
func (s *UploadService) ListByUserID(ctx context.Context, userID int64, offset, limit int, params map[string]any) (*entities.Uploads, *dto.Pagination, error) {
	return s.repo.ListByUserID(ctx, userID, offset, limit, params)
}

// GetByUUID retrieves upload by identifier with ownership check
// ctx: Context for cancellation/timeout
// uuid: Upload identifier
// userID: Owner verification
// Returns: (*entities.Upload, error)
func (s *UploadService) GetByUUID(ctx context.Context, uuid string, userID int64) (*entities.Upload, error) {
	return s.repo.GetByUUID(ctx, uuid, userID)
}

// GetByName retrieves upload by filename with ownership check
// ctx: Context for cancellation/timeout
// name: Original filename
// userID: Owner verification
// Returns: (*entities.Upload, error)
func (s *UploadService) GetByName(ctx context.Context, name string, userID int64) (*entities.Upload, error) {
	return s.repo.GetByName(ctx, name, userID)
}

// Delete removes upload metadata (standard operation)
// ctx: Context for cancellation/timeout
// uuid: Upload identifier
// userID: Owner verification
// Returns: error on failure
func (s *UploadService) Delete(ctx context.Context, uuid string, userID int64) error {
	return s.repo.Delete(ctx, uuid, userID)
}

// DeleteForce removes upload metadata (admin/cleanup operation)
// ctx: Context for cancellation/timeout
// uuid: Upload identifier
// userID: Owner verification
// Returns: error on failure

func (s *UploadService) DeleteForce(ctx context.Context, uuid string, userID int64) error {
	return s.repo.DeleteForce(ctx, uuid, userID)
}

// CanBeDeleted validates if upload can be deleted
// ctx: Context for cancellation/timeout
// uuid: Upload identifier
// userID: Owner verification
// Returns: error if deletion not allowed
func (s *UploadService) CanBeDeleted(ctx context.Context, uuid string, userID int64) error {
	_, err := s.repo.GetByUUID(ctx, uuid, userID)

	if err != nil {
		s.log.Debug().Err(err).Msg("Delete")
		return err
	}

	return nil
}
